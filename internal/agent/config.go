package agent

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

func LoadConfig(path, target string) (Config, error) {
	values, err := readEnv(path)
	if err != nil {
		return Config{}, err
	}
	get := func(key string) string {
		if value, ok := os.LookupEnv(key); ok {
			return strings.TrimSpace(value)
		}
		return strings.TrimSpace(values[key])
	}
	config := Config{ServerURL: get("PORTUNE_SERVER_URL"), Token: get("PORTUNE_AGENT_TOKEN"), LocalTarget: strings.TrimSpace(target)}
	prefix := get("PORTUNE_AGENT_NAME_PREFIX")
	for _, setting := range []struct{ key, value string }{
		{"PORTUNE_SERVER_URL", config.ServerURL}, {"PORTUNE_AGENT_TOKEN", config.Token}, {"PORTUNE_AGENT_NAME_PREFIX", prefix},
	} {
		if setting.value == "" {
			return Config{}, fmt.Errorf("%s is required in .env or environment", setting.key)
		}
	}
	if _, err := validateURL(config.ServerURL); err != nil {
		return Config{}, fmt.Errorf("PORTUNE_SERVER_URL must be an absolute HTTP(S) URL without credentials or fragment")
	}
	u, err := validateURL(config.LocalTarget)
	if err != nil {
		return Config{}, fmt.Errorf("--target must be an absolute HTTP(S) URL without credentials or fragment")
	}
	if u.RawQuery != "" || u.ForceQuery {
		return Config{}, fmt.Errorf("--target must not contain a query string")
	}
	for _, c := range prefix {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
			return Config{}, fmt.Errorf("PORTUNE_AGENT_NAME_PREFIX must contain only letters, digits or hyphens")
		}
	}
	config.AgentIP, err = localIP()
	if err != nil {
		return Config{}, err
	}
	config.AgentName = config.AgentIP + "-" + prefix + "-agent"
	config.TunnelName = u.Hostname()
	return config, nil
}

func validateURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return nil, fmt.Errorf("invalid URL")
	}
	return u, nil
}

func localIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("detect local IP: %w", err)
	}
	var fallback string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, address := range addresses {
			ip, _, err := net.ParseCIDR(address.String())
			if err != nil || !ip.IsGlobalUnicast() {
				continue
			}
			if ip.To4() != nil && ip.IsPrivate() {
				return ip.String(), nil
			}
			if fallback == "" || ip.To4() != nil {
				fallback = ip.String()
			}
		}
	}
	if fallback != "" {
		return fallback, nil
	}
	return "", fmt.Errorf("cannot detect a non-loopback local IP; connect to a network before starting the agent")
}

func readEnv(path string) (map[string]string, error) {
	values := make(map[string]string)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return values, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		text = strings.TrimSpace(strings.TrimPrefix(text, "export "))
		key, value, ok := strings.Cut(text, "=")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid .env entry at line %d", line)
		}
		if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") {
			quote := value[0]
			end := strings.IndexByte(value[1:], quote)
			if end < 0 {
				return nil, fmt.Errorf("unclosed .env quote at line %d", line)
			}
			end++
			tail := strings.TrimSpace(value[end+1:])
			if tail != "" && !strings.HasPrefix(tail, "#") {
				return nil, fmt.Errorf("invalid .env value at line %d", line)
			}
			value = value[1:end]
		} else if index := strings.Index(value, " #"); index >= 0 {
			value = strings.TrimSpace(value[:index])
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}
	return values, nil
}
