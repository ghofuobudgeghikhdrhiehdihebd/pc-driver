package ollamaish

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/pkg/httphelpers"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/pkg/schemas"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httphelpers.WriteError(w, http.StatusMethodNotAllowed, "use POST /api/chat")
		return
	}

	var request schemas.ChatHandlerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := validateChatRequest(request); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response := pseudoHandleChat(request) // REPLACE THIS LINE AND IMPLEMENT THE DRIVER COMPATIBLE WITH OLLAMA!!!
	httphelpers.WriteJSON(w, http.StatusOK, response)
}

func pseudoHandleChat(request schemas.ChatHandlerRequest) schemas.ChatHandlerResponse {
	for _, tool := range request.Tools {
		if tool.Name == "get_current_time" && strings.Contains(strings.ToLower(request.CurrentPrompt), "time") {
			return schemas.ChatHandlerResponse{
				ToolCalls: []schemas.ToolCall{{
					Name:      tool.Name,
					Arguments: map[string]string{},
				}},
			}
		}
	}

	return schemas.ChatHandlerResponse{
		Content: "[Pseudo response] " + request.CurrentPrompt,
	}
}

func validateChatRequest(request schemas.ChatHandlerRequest) error {
	if len(request.Messages) == 0 {
		return errors.New("messages must contain at least one message")
	}
	if strings.TrimSpace(request.CurrentPrompt) == "" {
		return errors.New("currentPrompt must not be empty")
	}

	for index, message := range request.Messages {
		if strings.TrimSpace(message.Sender) == "" {
			return errors.New("messages[" + strconv.Itoa(index) + "].sender must not be empty")
		}
		if strings.TrimSpace(message.Content) == "" {
			return errors.New("messages[" + strconv.Itoa(index) + "].content must not be empty")
		}
	}

	return nil
}
