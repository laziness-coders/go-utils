package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewEnvelope(t *testing.T) {
	data := map[string]string{"id": "1", "name": "Alice"}
	env := NewEnvelope(data)

	if !env.Success {
		t.Error("Expected Success to be true")
	}
	if env.Data == nil {
		t.Error("Expected Data to be set")
	}
	if env.Error != nil {
		t.Error("Expected Error to be nil")
	}
}

func TestNewEnvelopeWithMeta(t *testing.T) {
	data := map[string]string{"id": "1"}
	meta := map[string]int{"page": 1, "size": 10}
	env := NewEnvelopeWithMeta(data, meta)

	if !env.Success {
		t.Error("Expected Success to be true")
	}
	if env.Meta == nil {
		t.Error("Expected Meta to be set")
	}
}

func TestNewEnvelopeError(t *testing.T) {
	err := NewAppError("NOT_FOUND", "Resource not found")
	env := NewEnvelopeError(err)

	if env.Success {
		t.Error("Expected Success to be false")
	}
	if env.Error == nil {
		t.Error("Expected Error to be set")
	}
}

func TestNewFlat(t *testing.T) {
	data := map[string]string{"id": "1"}
	flat := NewFlat("Fetched successfully", data)

	if flat.Status != "ok" {
		t.Errorf("Expected Status to be 'ok', got '%s'", flat.Status)
	}
	if flat.Message != "Fetched successfully" {
		t.Error("Expected Message to be set")
	}
	if flat.Data == nil {
		t.Error("Expected Data to be set")
	}
}

func TestNewFlatError(t *testing.T) {
	flat := NewFlatError("Something went wrong")

	if flat.Status != "error" {
		t.Errorf("Expected Status to be 'error', got '%s'", flat.Status)
	}
	if flat.Message != "Something went wrong" {
		t.Error("Expected Message to be set")
	}
}

func TestNewLean(t *testing.T) {
	data := map[string]string{"id": "1"}
	lean := NewLean(data)

	if lean.Data == nil {
		t.Error("Expected Data to be set")
	}
}

func TestNewLeanError(t *testing.T) {
	err := NewAppError("ERROR", "Error occurred")
	lean := NewLeanError(err)

	if lean.Error == nil {
		t.Error("Expected Error to be set")
	}
}

func TestNewHAL(t *testing.T) {
	data := map[string]string{"id": "1", "name": "Alice"}
	links := map[string]Link{
		"self": {Href: "/v1/users/1"},
		"list": {Href: "/v1/users"},
	}
	hal := NewHAL(data, links)

	if hal.Links == nil {
		t.Error("Expected Links to be set")
	}
	if len(hal.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(hal.Links))
	}
	if hal.Data == nil {
		t.Error("Expected Data to be set")
	}
}

func TestNewHALWithSelfLink(t *testing.T) {
	data := map[string]string{"id": "1"}
	hal := NewHALWithSelfLink(data, "/v1/users/1")

	if hal.Links == nil {
		t.Error("Expected Links to be set")
	}
	if _, ok := hal.Links["self"]; !ok {
		t.Error("Expected self link to exist")
	}
	if hal.Links["self"].Href != "/v1/users/1" {
		t.Error("Expected correct self link href")
	}
}

func TestNewValidation(t *testing.T) {
	errors := []FieldError{
		{Field: "email", Message: "must be valid email"},
		{Field: "password", Message: "min length 6"},
	}
	validation := NewValidation("INVALID_INPUT", errors)

	if validation.Code != "INVALID_INPUT" {
		t.Error("Expected Code to be 'INVALID_INPUT'")
	}
	if len(validation.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(validation.Errors))
	}
}

func TestValidationAddError(t *testing.T) {
	validation := Validation{Code: "INVALID_INPUT"}
	validation.AddError("username", "required")
	validation.AddError("email", "invalid format")

	if len(validation.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(validation.Errors))
	}
	if validation.Errors[0].Field != "username" {
		t.Error("Expected first error field to be 'username'")
	}
}

func TestWriteResponseEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"id": "1", "name": "Alice"}

	err := WriteResponse(w, http.StatusOK, data, nil, nil, FormatEnvelope)
	if err != nil {
		t.Fatalf("WriteResponse failed: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response Envelope[map[string]string]
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}
	if response.Data["id"] != "1" {
		t.Error("Expected data to be preserved")
	}
}

func TestWriteResponseFlat(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"id": "1"}

	err := WriteResponse(w, http.StatusOK, data, nil, nil, FormatFlat)
	if err != nil {
		t.Fatalf("WriteResponse failed: %v", err)
	}

	var response Flat[map[string]string]
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}
}

func TestWriteResponseLean(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"id": "1"}

	err := WriteResponse(w, http.StatusOK, data, nil, nil, FormatLean)
	if err != nil {
		t.Fatalf("WriteResponse failed: %v", err)
	}

	var response OnlyData[map[string]string]
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Data["id"] != "1" {
		t.Error("Expected data to be preserved")
	}
}

func TestWriteResponseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	appErr := &AppError{Code: "NOT_FOUND", Message: "Resource not found"}

	err := WriteResponse[any](w, http.StatusNotFound, nil, appErr, nil, FormatEnvelope)
	if err != nil {
		t.Fatalf("WriteResponse failed: %v", err)
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response Envelope[any]
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected Success to be false")
	}
	if response.Error == nil {
		t.Error("Expected Error to be set")
	}
}

func TestWriteValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	errors := []FieldError{
		{Field: "email", Message: "must be valid email"},
	}

	err := WriteValidationError(w, http.StatusBadRequest, errors, FormatEnvelope)
	if err != nil {
		t.Fatalf("WriteValidationError failed: %v", err)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected ResponseFormat
	}{
		{"envelope", FormatEnvelope},
		{"flat", FormatFlat},
		{"lean", FormatLean},
		{"hal", FormatHAL},
		{"unknown", FormatEnvelope},
		{"", FormatEnvelope},
	}

	for _, tt := range tests {
		result := ParseFormat(tt.input)
		if result != tt.expected {
			t.Errorf("ParseFormat(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestNewAppError(t *testing.T) {
	err := NewAppError("BAD_REQUEST", "Invalid request")

	if err.Code != "BAD_REQUEST" {
		t.Error("Expected Code to be 'BAD_REQUEST'")
	}
	if err.Message != "Invalid request" {
		t.Error("Expected Message to be set")
	}
}

func TestNewAppErrorWithDetails(t *testing.T) {
	details := map[string]string{"field": "email"}
	err := NewAppErrorWithDetails("VALIDATION_ERROR", "Validation failed", details)

	if err.Details == nil {
		t.Error("Expected Details to be set")
	}
}
