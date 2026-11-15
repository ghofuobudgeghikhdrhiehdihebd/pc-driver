package ollamaish

import "net/http"

func RegisterAll(mux *http.ServeMux) {
	mux.HandleFunc("/api/chat", ChatHandler)
	mux.HandleFunc("/health", HealthHandler)
}
