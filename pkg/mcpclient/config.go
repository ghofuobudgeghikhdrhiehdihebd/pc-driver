package mcpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// NewMcpClient loads AirLock's Remote.json and Local.json files.
func NewMcpClient() (*McpClient, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	base := filepath.Join(home, "Application Support", "AirLock", "MCP")
	return NewMcpClientFromFiles(filepath.Join(base, "Remote.json"), filepath.Join(base, "Local.json"))
}

// NewMcpClientFromFiles creates a client from remote and local JSON files.
func NewMcpClientFromFiles(remotePath, localPath string) (*McpClient, error) {
	all := map[string]Server{}
	found := false
	for _, path := range []string{remotePath, localPath} {
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read MCP config %s: %w", path, err)
		}
		var config fileConfig
		if err = json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("parse MCP config %s: %w", path, err)
		}
		found = true
		for name, server := range config.MCPServers {
			all[name] = server
		}
	}
	if !found {
		return nil, fmt.Errorf("no MCP configuration files found")
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("MCP configuration contains no servers")
	}
	client := &McpClient{servers: map[string]*connection{}}
	for name, server := range all {
		client.servers[name] = &connection{config: server, client: &http.Client{}}
	}
	return client, nil
}

// Servers returns configured server names.
func (c *McpClient) Servers() []string {
	names := make([]string, 0, len(c.servers))
	for name := range c.servers {
		names = append(names, name)
	}
	return names
}
