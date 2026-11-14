// in Go, the package protocol must always come first
package protocol

import "encoding/json"

const (
	TypeAuthorizationRequest  = "authorization_request"
	TypeAuthorizationResponse = "authorization_response"
	TypeError                 = "error"
)

// AuthorizationRequest represents a request from the MacBook to AirLock.
type AuthorizationRequest struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
}

// ErrorResponse represents an error returned by AirLock.
type ErrorResponse struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id,omitempty"`
	Message   string `json:"message"`
}

// AuthorizationResponse represents AirLock's decision.
type AuthorizationResponse struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
	Approved  bool   `json:"approved"`
	Message   string `json:"message,omitempty"`
}

// EncodeMessage serializes a protocol message into newline-terminated JSON.
func EncodeMessage(message any) ([]byte, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return append(data, '\n'), nil
}

// DecodeMessage deserializes JSON data into a protocol message.
func DecodeMessage(data []byte, message any) error {
	return json.Unmarshal(data, message)
}
