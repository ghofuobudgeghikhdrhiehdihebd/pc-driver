package mcpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// CallTool invokes a tool on a configured server and returns its raw MCP result.
func (c *McpClient) CallTool(ctx context.Context, server, tool string, arguments any) (json.RawMessage, error) {
	conn, ok := c.servers[server]
	if !ok {
		return nil, fmt.Errorf("unknown MCP server %q", server)
	}
	conn.mu.Lock()
	defer conn.mu.Unlock()
	params := map[string]any{"name": tool, "arguments": arguments}
	if conn.config.URL != "" || strings.EqualFold(conn.config.Type, "http") || strings.EqualFold(conn.config.Type, "streamable-http") {
		return conn.callHTTP(ctx, c.next(), "tools/call", params)
	}
	return conn.callStdio(ctx, c.next(), "tools/call", params)
}
func (c *McpClient) next() int64 { return c.seq.Add(1) }
