package mcpclient

import (
	"encoding/json"
	"fmt"
	"io"
)

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}
type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
}
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func decodeResponse(r io.Reader) (json.RawMessage, error) {
	var p rpcResponse
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return nil, err
	}
	if p.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", p.Error.Code, p.Error.Message)
	}
	return p.Result, nil
}
