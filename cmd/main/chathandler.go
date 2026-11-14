package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type chatHandlerMessage struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type chatHandlerTool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  map[string]string `json:"parameters"`
}

type chatHandlerRequest struct {
	Model         string               `json:"model"`
	Messages      []chatHandlerMessage `json:"messages"`
	CurrentPrompt string               `json:"currentPrompt"`
	Tools         []chatHandlerTool    `json:"tools"`
}

type chatHandlerResponse struct {
	Content   string           `json:"content,omitempty"`
	Message   *responseMessage `json:"message,omitempty"`
	ToolCalls []pseudoToolCall `json:"toolCalls,omitempty"`
}

type responseMessage struct {
	Content string `json:"content"`
}

type pseudoToolCall struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/chat", chatHandler)
	mux.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr:              ":11434",
		Handler:           withCORS(withLogging(mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("local chat backend listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST /api/chat")
		return
	}

	var request chatHandlerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := validateChatRequest(request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	response := pseudoHandleChat(request) // REPLACE THIS LINE AND IMPLEMENT THE DRIVER COMPATIBLE WITH OLLAMA!!!
	writeJSON(w, http.StatusOK, response)
}

func validateChatRequest(request chatHandlerRequest) error {
	if len(request.Messages) == 0 {
		return errors.New("messages must contain at least one message")
	}
	if strings.TrimSpace(request.CurrentPrompt) == "" {
		return errors.New("currentPrompt must not be empty")
	}

	for index, message := range request.Messages {
		if strings.TrimSpace(message.Sender) == "" {
			return errors.New("messages[" + itoa(index) + "].sender must not be empty")
		}
		if strings.TrimSpace(message.Content) == "" {
			return errors.New("messages[" + itoa(index) + "].content must not be empty")
		}
	}

	return nil
}

func pseudoHandleChat(request chatHandlerRequest) chatHandlerResponse {
	for _, tool := range request.Tools {
		if tool.Name == "get_current_time" && strings.Contains(strings.ToLower(request.CurrentPrompt), "time") {
			return chatHandlerResponse{
				ToolCalls: []pseudoToolCall{{
					Name:      tool.Name,
					Arguments: map[string]string{},
				}},
			}
		}
	}

	return chatHandlerResponse{
		Content: "[Pseudo response] " + request.CurrentPrompt,
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(started).Round(time.Millisecond))
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	index := len(digits)
	for value > 0 {
		index--
		digits[index] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[index:])
}
