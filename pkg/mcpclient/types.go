package mcpclient

import (
	"bufio"
	"io"
	"net/http"
	"os/exec"
	"sync"
	"sync/atomic"
)

type Server struct {
	Type    string            `json:"type"`
	URL     string            `json:"url"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
	Headers map[string]string `json:"headers"`
}
type fileConfig struct {
	MCPServers map[string]Server `json:"mcpServers"`
}
type McpClient struct {
	servers map[string]*connection
	seq     atomic.Int64
}
type connection struct {
	config  Server
	mu      sync.Mutex
	cmd     *exec.Cmd
	in      io.WriteCloser
	out     *bufio.Reader
	client  *http.Client
	session string
}
