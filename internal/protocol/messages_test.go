package protocol

import (
	"encoding/json"
	"testing"
)

func TestAuthorizationRequestJSON(t *testing.T) {
	request := AuthorizationRequest{
		Type:      TypeAuthorizationRequest,
		RequestID: "req-001",
		Action:    "read",
		Resource:  "report.pdf",
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	var decoded AuthorizationRequest

	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded != request {
		t.Fatalf("expected %+v, got %+v", request, decoded)
	}
}

func TestAuthorizationResponseJSON(t *testing.T) {
	response := AuthorizationResponse{
		Type:      TypeAuthorizationResponse,
		RequestID: "req-001",
		Approved:  true,
		Message:   "Access granted",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	var decoded AuthorizationResponse

	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded != response {
		t.Fatalf("expected %+v, got %+v", response, decoded)
	}
}

func TestErrorResponseJSON(t *testing.T) {
	response := ErrorResponse{
		Type:      TypeError,
		RequestID: "req-001",
		Message:   "Unsupported action",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	var decoded ErrorResponse

	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded != response {
		t.Fatalf("expected %+v, got %+v", response, decoded)
	}
}

func TestEncodeMessage(t *testing.T) {
	request := AuthorizationRequest{
		Type:      TypeAuthorizationRequest,
		RequestID: "req-001",
		Action:    "read",
		Resource:  "report.pdf",
	}

	data, err := EncodeMessage(request)
	if err != nil {
		t.Fatal(err)
	}

	expected := "{\"type\":\"authorization_request\",\"request_id\":\"req-001\",\"action\":\"read\",\"resource\":\"report.pdf\"}\n"

	if string(data) != expected {
		t.Fatalf("expected %q, got %q", expected, string(data))
	}
}

func TestDecodeMessage(t *testing.T) {
	data := []byte(`{
		"type": "authorization_response",
		"request_id": "req-001",
		"approved": true,
		"message": "Access granted"
	}`)

	var response AuthorizationResponse

	err := DecodeMessage(data, &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.RequestID != "req-001" {
		t.Fatalf("unexpected request ID: %q", response.RequestID)
	}

	if !response.Approved {
		t.Fatal("expected access to be approved")
	}

	if response.Message != "Access granted" {
		t.Fatalf("unexpected message: %q", response.Message)
	}
}
