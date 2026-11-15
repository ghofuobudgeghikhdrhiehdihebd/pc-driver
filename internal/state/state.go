package statepkg

import "net/http"

type State struct {
	Mux       *http.ServeMux
	McpUdsLoc string
}
