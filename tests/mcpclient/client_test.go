package mcpclient_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/pkg/mcpclient"
)

func writeConfig(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestNewMcpClientFromFilesMergesRemoteAndLocal(t *testing.T) {
	dir := t.TempDir()
	remote := filepath.Join(dir, "Remote.json")
	local := filepath.Join(dir, "Local.json")
	writeConfig(t, remote, map[string]any{"mcpServers": map[string]any{
		"remote": map[string]any{"type": "http", "url": "https://example.test/mcp"},
		"shared": map[string]any{"type": "http", "url": "https://remote.test/mcp"},
	}})
	writeConfig(t, local, map[string]any{"mcpServers": map[string]any{
		"local":  map[string]any{"command": "server"},
		"shared": map[string]any{"type": "http", "url": "https://local.test/mcp"},
	}})

	client, err := mcpclient.NewMcpClientFromFiles(remote, local)
	if err != nil {
		t.Fatal(err)
	}
	got := client.Servers()
	want := []string{"local", "remote", "shared"}
	// Servers is intentionally unordered.
	for _, name := range want {
		found := false
		for _, candidate := range got {
			if candidate == name {
				found = true
			}
		}
		if !found {
			t.Errorf("Servers() missing %q: %v", name, got)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("Servers() = %v, want %v", got, want)
	}

	// Local.json is loaded second and overrides duplicate names.
	if _, err := client.CallTool(context.Background(), "shared", "unused", nil); err == nil || !contains(err.Error(), "local.test") {
		t.Fatalf("expected the local override to be used, got %v", err)
	}
}

func TestNewMcpClientFromFilesErrors(t *testing.T) {
	dir := t.TempDir()
	_, err := mcpclient.NewMcpClientFromFiles(filepath.Join(dir, "missing-remote"), filepath.Join(dir, "missing-local"))
	if err == nil {
		t.Fatal("expected missing configuration error")
	}

	bad := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(bad, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = mcpclient.NewMcpClientFromFiles(bad, filepath.Join(dir, "missing-local"))
	if err == nil {
		t.Fatal("expected invalid JSON error")
	}

	empty := filepath.Join(dir, "empty.json")
	writeConfig(t, empty, map[string]any{"mcpServers": map[string]any{}})
	_, err = mcpclient.NewMcpClientFromFiles(empty, filepath.Join(dir, "missing-local"))
	if err == nil {
		t.Fatal("expected empty configuration error")
	}
}

func TestCallToolUsesConfiguredHTTPServer(t *testing.T) {
	dir := t.TempDir()
	remote := filepath.Join(dir, "Remote.json")
	writeConfig(t, remote, map[string]any{"mcpServers": map[string]any{"demo": map[string]any{"type": "http", "url": "https://configured.example/mcp"}}})
	client, err := mcpclient.NewMcpClientFromFiles(remote, filepath.Join(dir, "missing-local"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.CallTool(context.Background(), "demo", "greet", map[string]any{"name": "Ada"})
	if err == nil || !contains(err.Error(), "configured.example") {
		t.Fatalf("expected configured URL error, got %v", err)
	}
}

func contains(value, substring string) bool {
	for i := 0; i+len(substring) <= len(value); i++ {
		if value[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}
