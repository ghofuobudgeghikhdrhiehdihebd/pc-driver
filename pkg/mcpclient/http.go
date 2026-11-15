package mcpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (x *connection) callHTTP(ctx context.Context, id int64, method string, params any) (json.RawMessage, error) {
	if x.config.URL == "" {
		return nil, fmt.Errorf("MCP HTTP server has no URL")
	}
	if method == "tools/call" && x.session == "" {
		if _, err := x.callHTTP(ctx, id, "initialize", map[string]any{"protocolVersion": "2024-11-05", "capabilities": map[string]any{}, "clientInfo": map[string]string{"name": "pc-usbc-int", "version": "0.1.0"}}); err != nil {
			return nil, err
		}
		x.session = "initialized"
	}
	body, err := json.Marshal(rpcRequest{"2.0", id, method, params})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, x.config.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range x.config.Headers {
		req.Header.Set(key, value)
	}
	if x.session != "" {
		req.Header.Set("Mcp-Session-Id", x.session)
	}
	response, err := x.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("MCP server returned HTTP %s", response.Status)
	}
	if value := response.Header.Get("Mcp-Session-Id"); value != "" {
		x.session = value
	}
	return decodeResponse(response.Body)
}
