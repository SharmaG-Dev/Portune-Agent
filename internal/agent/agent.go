package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	socketio "github.com/zishang520/socket.io/clients/socket/v3"
	"github.com/zishang520/socket.io/v3/pkg/types"
)

type Agent struct {
	config Config

	socket *socketio.Socket

	httpClient *http.Client
}

func New(config Config) *Agent {
	return &Agent{
		config: config,

		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (a *Agent) Run(ctx context.Context) error {
	fmt.Println()
	fmt.Println("Portune Agent")
	fmt.Println("================================")

	fmt.Println("Server :", a.config.ServerURL)
	fmt.Println("Target :", a.config.LocalTarget)

	if a.config.Subdomain != "" {
		fmt.Println("Subdomain :", a.config.Subdomain)
	}

	fmt.Println()

	options := socketio.DefaultOptions()

	// Force WebSocket transport.
	options.SetTransports(
		types.NewSet(
			socketio.WebSocket,
		),
	)

	// Token gets sent as:
	//
	// socket.handshake.auth.token
	options.SetAuth(
		map[string]any{
			"token": a.config.Token,
		},
	)

	socket, err := socketio.Connect(
		a.config.ServerURL,
		options,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to connect to Portune server: %w",
			err,
		)
	}

	a.socket = socket

	a.registerEvents()

	<-ctx.Done()

	fmt.Println()
	fmt.Println("Stopping Portune Agent...")

	a.socket.Disconnect()

	fmt.Println("Portune Agent stopped.")

	return nil
}

func (a *Agent) registerEvents() {

	_ = a.socket.On(
		"connect",
		func(args ...any) {

			fmt.Println("✓ Connected to Portune Server")

			fmt.Println(
				"  Socket ID:",
				a.socket.Id(),
			)
		},
	)

	_ = a.socket.On(
		"authenticated",
		func(args ...any) {

			fmt.Println("✓ Agent authenticated")

			if len(args) > 0 {

				if data, ok :=
					args[0].(map[string]any); ok {

					if agentName, ok :=
						data["agentName"].(string); ok {

						fmt.Println(
							"  Agent Name:",
							agentName,
						)
					}

					if agentID,
						ok := data["agentId"].(float64); ok {

						fmt.Println(
							"  Agent ID:",
							int(agentID),
						)
					}
				}
			}

			// Only create tunnel after successful auth.
			a.createTunnel()
		},
	)

	_ = a.socket.On(
		"connect_error",
		func(args ...any) {

			fmt.Println(
				"✗ Connection error:",
				args,
			)
		},
	)

	_ = a.socket.On(
		"error",
		func(args ...any) {

			fmt.Println(
				"✗ Server error:",
				args,
			)
		},
	)

	_ = a.socket.On(
		"disconnect",
		func(args ...any) {

			fmt.Println(
				"⚠ Disconnected:",
				args,
			)
		},
	)

	_ = a.socket.On(
		"http:request",
		func(args ...any) {

			if len(args) == 0 {
				fmt.Println(
					"Invalid http:request: empty payload",
				)

				return
			}

			payload, ok :=
				args[0].(map[string]any)

			if !ok {

				fmt.Printf(
					"Invalid http:request payload: %#v\n",
					args[0],
				)

				return
			}

			// Each HTTP request can be processed concurrently.
			go a.handleHTTPRequest(payload)
		},
	)
}

func (a *Agent) createTunnel() {

	payload := TunnelCreateRequest{
		RequestedSubdomain: a.config.Subdomain,
		Name:               a.config.TunnelName,
	}

	fmt.Println()
	fmt.Println("Creating tunnel...")

	a.socket.
		Timeout(10*time.Second).
		EmitWithAck(
			"tunnel:create",
			payload,
		)(
		func(args []any, err error) {

			if err != nil {

				fmt.Println(
					"✗ Tunnel creation failed:",
					err,
				)

				return
			}

			if len(args) == 0 {

				fmt.Println(
					"✗ Tunnel server returned empty response",
				)

				return
			}

			response,
				ok := args[0].(map[string]any)

			if !ok {

				fmt.Printf(
					"✗ Invalid tunnel response: %#v\n",
					args[0],
				)

				return
			}

			success,
				_ := response["ok"].(bool)

			if !success {

				fmt.Println(
					"✗ Tunnel creation failed:",
					response["error"],
				)

				return
			}

			fmt.Println()
			fmt.Println("✓ Tunnel Established!")
			fmt.Println("================================")

			fmt.Println(
				"Tunnel ID :",
				response["tunnelId"],
			)

			fmt.Println(
				"Subdomain :",
				response["subdomain"],
			)

			fmt.Println(
				"Public URL:",
				response["url"],
			)

			fmt.Println(
				"Forwarding:",
				a.config.LocalTarget,
			)

			fmt.Println("================================")
			fmt.Println()
		},
	)
}
