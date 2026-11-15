package mcpclient

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func (x *connection) callStdio(ctx context.Context, id int64, method string, params any) (json.RawMessage, error) {
	if x.cmd == nil {
		if x.config.Command == "" {
			return nil, fmt.Errorf("MCP server has neither URL nor command")
		}
		x.cmd = exec.Command(x.config.Command, x.config.Args...)
		x.cmd.Stderr = os.Stderr
		var err error
		x.in, err = x.cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		stdout, err := x.cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		x.out = bufio.NewReader(stdout)
		x.cmd.Env = os.Environ()
		for key, value := range x.config.Env {
			x.cmd.Env = append(x.cmd.Env, key+"="+value)
		}
		if err = x.cmd.Start(); err != nil {
			return nil, err
		}
	}
	if method == "tools/call" && x.session == "" {
		if _, err := x.exchange(ctx, id, "initialize", map[string]any{"protocolVersion": "2024-11-05", "capabilities": map[string]any{}, "clientInfo": map[string]string{"name": "pc-usbc-int", "version": "0.1.0"}}); err != nil {
			return nil, err
		}
		message, _ := json.Marshal(rpcRequest{"2.0", 0, "notifications/initialized", nil})
		message = append(message, '\n')
		if _, err := x.in.Write(message); err != nil {
			return nil, err
		}
		x.session = "initialized"
	}
	return x.exchange(ctx, id, method, params)
}
func (x *connection) exchange(ctx context.Context, id int64, method string, params any) (json.RawMessage, error) {
	message, err := json.Marshal(rpcRequest{"2.0", id, method, params})
	if err != nil {
		return nil, err
	}
	message = append(message, '\n')
	if _, err = x.in.Write(message); err != nil {
		return nil, err
	}
	type response struct {
		value json.RawMessage
		err   error
	}
	result := make(chan response, 1)
	go func() { value, err := decodeResponse(x.out); result <- response{value, err} }()
	select {
	case response := <-result:
		return response.value, response.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
