package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("# comment\nexport A=hello # comment\nB='token#literal'\nC=\"quoted value\" # comment\n"), 0600); err != nil {
		t.Fatal(err)
	}
	values, err := readEnv(path)
	if err != nil {
		t.Fatal(err)
	}
	if values["A"] != "hello" || values["B"] != "token#literal" || values["C"] != "quoted value" {
		t.Fatalf("unexpected values: %v", values)
	}
	if err := os.WriteFile(path, []byte("TOKEN='private-unclosed"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err = readEnv(path)
	if err == nil || strings.Contains(err.Error(), "private") {
		t.Fatal("expected redacted parsing error")
	}
}

func TestConfigValidationAndEnvironmentPrecedence(t *testing.T) {
	t.Setenv("PORTUNE_SERVER_URL", "https://example.test/tunnel")
	t.Setenv("PORTUNE_AGENT_TOKEN", "")
	t.Setenv("PORTUNE_AGENT_NAME_PREFIX", "demo")
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("PORTUNE_AGENT_TOKEN=file-token\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadConfig(path, "http://localhost:4000")
	if err == nil || !strings.Contains(err.Error(), "PORTUNE_AGENT_TOKEN") {
		t.Fatal("empty environment must override file token")
	}
	t.Setenv("PORTUNE_AGENT_TOKEN", "test-token")
	for _, target := range []string{"", "localhost:4000", "ftp://localhost", "http://user:password@localhost", "http://localhost/#fragment", "http://localhost/?q=1"} {
		if _, err := LoadConfig(path, target); err == nil {
			t.Errorf("accepted invalid target %q", target)
		}
	}
}
