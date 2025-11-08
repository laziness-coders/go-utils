package response_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/laziness-coders/go-utils/response"
)

// User represents a simple user model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SignupRequest represents a signup request with validation rules
type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"required,min=18"`
}

// ExampleEnvelope demonstrates the envelope response format
func Example_envelope() {
	w := httptest.NewRecorder()
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	response.WriteResponse(w, http.StatusOK, user, nil, nil, response.FormatEnvelope)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	fmt.Printf("Success: %v\n", result["success"])
	// Output: Success: true
}

// ExampleFlat demonstrates the flat response format
func Example_flat() {
	w := httptest.NewRecorder()
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	response.WriteResponse(w, http.StatusOK, user, nil, nil, response.FormatFlat)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	fmt.Printf("Status: %v\n", result["status"])
	// Output: Status: ok
}

// ExampleLean demonstrates the lean response format
func Example_lean() {
	w := httptest.NewRecorder()
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	response.WriteResponse(w, http.StatusOK, user, nil, nil, response.FormatLean)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	fmt.Printf("User ID: %.0f\n", data["id"])
	// Output: User ID: 1
}

// ExampleHAL demonstrates the HAL response format with hypermedia links
func Example_hal() {
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	links := map[string]response.Link{
		"self": {Href: "/v1/users/1"},
		"list": {Href: "/v1/users"},
	}
	halResp := response.NewHAL(user, links)

	fmt.Printf("Self link: %s\n", halResp.Links["self"].Href)
	// Output: Self link: /v1/users/1
}

// ExampleValidation demonstrates validation with English translations
func Example_validation() {
	response.InitValidator()

	req := SignupRequest{
		Email:    "invalid-email",
		Password: "123",
		Age:      15,
	}

	errors := response.ValidateStruct(req, "en")
	fmt.Printf("Number of errors: %d\n", len(errors))
	fmt.Printf("First error field: %s\n", errors[0].Field)
	// Output:
	// Number of errors: 3
	// First error field: Email
}

// ExampleValidation_vietnamese demonstrates validation with Vietnamese translations
func Example_validation_vietnamese() {
	response.InitValidator()

	req := SignupRequest{
		Email:    "",
		Password: "",
		Age:      0,
	}

	errors := response.ValidateStruct(req, "vi")
	fmt.Printf("Number of errors: %d\n", len(errors))
	// Output: Number of errors: 3
}

// ExampleAppError demonstrates error responses
func Example_appError() {
	w := httptest.NewRecorder()
	appErr := &response.AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
	}

	response.WriteResponse[any](w, http.StatusNotFound, nil, appErr, nil, response.FormatEnvelope)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	fmt.Printf("Success: %v\n", result["success"])
	// Output: Success: false
}

// ExampleParseFormat demonstrates format string parsing
func Example_parseFormat() {
	formats := []string{"envelope", "flat", "lean", "hal", "unknown"}

	for _, f := range formats {
		parsed := response.ParseFormat(f)
		fmt.Printf("%s -> %s\n", f, parsed)
	}
	// Output:
	// envelope -> envelope
	// flat -> flat
	// lean -> lean
	// hal -> hal
	// unknown -> envelope
}
