package schemas

type Message struct {
	Type     string         `msgpack:"type"`
	String   string         `msgpack:"string,omitempty"`
	Request  McpCallRequest `msgpack:"request,omitempty"`
	Response McpCallResult  `msgpack:"response,omitempty"`
}

// McpServer describes the connection information needed by mcpclient to
// reach one MCP server. A server is either remote (URL) or local (Command).
// Media-specific configuration is intentionally not represented here.
type McpServer struct {
	Name    string            `msgpack:"name"`
	Type    string            `msgpack:"type,omitempty"`
	URL     string            `msgpack:"url,omitempty"`
	Command string            `msgpack:"command,omitempty"`
	Args    []string          `msgpack:"args,omitempty"`
	Env     map[string]string `msgpack:"env,omitempty"`
	Headers map[string]string `msgpack:"headers,omitempty"`
}

// McpToolCall contains the complete input required for McpClient.CallTool.
// Arguments must contain JSON/MessagePack-compatible scalar, array, and map
// values accepted by the selected MCP tool.
type McpToolCall struct {
	Server    string         `msgpack:"server"`
	Tool      string         `msgpack:"tool"`
	Arguments map[string]any `msgpack:"arguments,omitempty"`
}

// McpCallRequest is a transport-level request that can be converted directly
// into mcpclient.McpClient.CallTool arguments.
type McpCallRequest struct {
	Server string         `msgpack:"server"`
	Call   McpToolCall    `msgpack:"call"`
	Params map[string]any `msgpack:"params,omitempty"`
}

// McpCallResult contains the tool result returned by an MCP server. Result is
// kept generic because each MCP tool defines its own result shape.
type McpCallResult struct {
	Server string `msgpack:"server"`
	Tool   string `msgpack:"tool"`
	Result any    `msgpack:"result,omitempty"`
	Error  string `msgpack:"error,omitempty"`
}
