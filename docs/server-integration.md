# Server changes for environment-configured agents

The client is implemented in this repository. The server is external and must implement the contract below before registration and automatic subdomains work end to end.

## Client setup

Copy `.env.example` to `.env`, fill in the server-issued token and server URL, then run:

```sh
go run . agent --target http://localhost:4000
```

Required settings: `PORTUNE_SERVER_URL`, `PORTUNE_AGENT_TOKEN`, `PORTUNE_AGENT_NAME_PREFIX`. There are no code defaults for these values. `.env` is read from the current working directory. Existing process environment takes precedence (including empty values, which fail validation). Files support single-line KEY=value, optional export, quotes and comments; shell expansion, escape processing and multiline values are not supported. Real .env files are ignored by Git.

The CLI accepts only --target (plus help). --server, --token, --subdomain and --name have been removed. Tunnel display name is the target hostname: http://localhost:4000 becomes localhost. It is not a unique identifier.

## 1. Authenticate and register/update the agent

Continue accepting Socket.IO with WebSocket transport at the configured URL (the existing /tunnel namespace). Handshake auth now contains:

```json
{
  "token": "<server-issued-token>",
  "agentIp": "192.168.1.10",
  "agentName": "192.168.1.10-portune-agent"
}
```

- Validate token before registering any metadata. Never log or return tokens.
- Associate registration with the identity/owner authorized by the token; never trust the supplied name or IP as authentication or a unique database key.
- Store/update agentName, reported local IP, socket ID, online state and last-seen time. Validate metadata types and lengths.
- Client name format is <local-IP>-<PORTUNE_AGENT_NAME_PREFIX>-agent. The prefix is configurable. IPv6 names may contain colons: this is a display name, not a DNS label.
- The client prefers a private IPv4 on an active non-loopback interface, otherwise another unicast IP. Multiple interfaces/VPNs can affect selection. No public-IP lookup is performed. Store server-observed remote IP separately if required; NAT and changing local IPs make names unsuitable as stable identities.
- Reconnect registration must be idempotent for the authenticated agent identity. If tokens currently represent accounts rather than individual agents, define persistent per-device identity on the backend before supporting multiple devices per token.
- On successful registration emit authenticated exactly once per connection, with the following payload. The client requires boolean ok=true before requesting a tunnel.

```json
{"ok": true, "agentId": 123, "agentName": "192.168.1.10-portune-agent"}
```

Reject invalid credentials through the Socket.IO connection error mechanism. Do not emit a successful authenticated event on failure. No new agent:register event or registration REST endpoint is called by this client.

## 2. Generate the subdomain on the server

After authenticated, the client emits tunnel:create with an acknowledgment callback:

```json
{"name": "localhost"}
```

- requestedSubdomain is deliberately absent. Generate an available DNS-safe random subdomain on the server.
- Enforce uniqueness with a database constraint/atomic reservation; retry collisions. A pre-check alone is insufficient under concurrent requests.
- Associate the tunnel with the authenticated agent and socket. Enforce ownership and tunnel limits.
- Construct the public URL using server configuration for public domain and scheme. Do not derive it from the local target or agent display name.
- Return the acknowledgment within the client's 10-second timeout:

```json
{"ok": true, "tunnelId": 456, "subdomain": "generated-value", "url": "https://generated-value.<configured-public-domain>"}
```

Failure acknowledgment:

```json
{"ok": false, "error": "Unable to allocate a tunnel"}
```

The agent prints the returned URL as-is. It does not generate or reserve domains. Custom subdomains are deferred to a future authenticated UI; no custom-domain UI is included here.

## 3. Preserve the HTTP forwarding protocol

Server sends http:request with requestId, subdomain, method, path (including query string), headers (object of string values) and body (string). Route each request only to its owning agent.

Agent emits http:response with requestId, statusCode, headers, body and isBase64Encoded. Decode Base64 when true before returning bytes to the visitor. Correlate requests, apply server-side response deadlines, and fail pending requests when agents disconnect.

## 4. Connection lifecycle

Mark agents offline and remove stale socket routing on disconnect. Define whether reconnect reuses a tunnel or replaces it; make tunnel creation idempotent for that policy so reconnects do not leak reservations. Do not use display name localhost to distinguish tunnels on different ports/devices. The current client requests a tunnel each time it receives authenticated.

## Acceptance checklist

- Valid token registers metadata; invalid token creates neither agent nor tunnel.
- Starting with only --target succeeds after .env is configured.
- Prefix changes affect the registered display name; target hostname determines tunnel name.
- Missing required settings fail locally without exposing credentials.
- Two concurrent allocations cannot reserve the same subdomain; collisions are retried.
- Acknowledgment includes the server-generated public URL; HTTP traffic reaches the local target.
- Reconnects/disconnects update registration and routing without duplicate stale reservations.
- No token appears in logs or committed files.
