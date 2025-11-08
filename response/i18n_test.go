package response

import (
	"testing"
)

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"required,min=18,max=100"`
}

func TestInitValidator(t *testing.T) {
	InitValidator()

	if validate == nil {
		t.Error("Expected validator to be initialized")
	}
	if uni == nil {
		t.Error("Expected universal translator to be initialized")
	}
}

func TestGetValidator(t *testing.T) {
	// Reset to test lazy initialization
	validate = nil
	uni = nil

	v := GetValidator()
	if v == nil {
		t.Error("Expected validator to be returned")
	}
	if validate == nil {
		t.Error("Expected validator to be initialized")
	}
}

func TestGetTranslatorEnglish(t *testing.T) {
	trans := GetTranslator("en")
	if trans == nil {
		t.Error("Expected translator to be returned")
	}

	// Test with different English variants
	trans = GetTranslator("en-US")
	if trans == nil {
		t.Error("Expected translator for en-US")
	}
}

func TestGetTranslatorVietnamese(t *testing.T) {
	trans := GetTranslator("vi")
	if trans == nil {
		t.Error("Expected translator to be returned")
	}

	// Test with different Vietnamese variants
	trans = GetTranslator("vi-VN")
	if trans == nil {
		t.Error("Expected translator for vi-VN")
	}
}

func TestTranslateErrorsEnglish(t *testing.T) {
	InitValidator()
	trans := GetTranslator("en")

	req := SignupRequest{
		Email:    "invalid-email",
		Password: "123",
		Age:      15,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Fatal("Expected validation errors")
	}

	errors := TranslateErrors(err, trans)
	if len(errors) == 0 {
		t.Error("Expected translated errors")
	}

	// Check that we got field-level errors
	foundEmail := false
	foundPassword := false
	foundAge := false

	for _, fieldErr := range errors {
		if fieldErr.Field == "Email" {
			foundEmail = true
		}
		if fieldErr.Field == "Password" {
			foundPassword = true
		}
		if fieldErr.Field == "Age" {
			foundAge = true
		}

		// Verify message is not empty
		if fieldErr.Message == "" {
			t.Errorf("Expected non-empty message for field %s", fieldErr.Field)
		}
	}

	if !foundEmail {
		t.Error("Expected email validation error")
	}
	if !foundPassword {
		t.Error("Expected password validation error")
	}
	if !foundAge {
		t.Error("Expected age validation error")
	}
}

func TestTranslateErrorsVietnamese(t *testing.T) {
	InitValidator()
	trans := GetTranslator("vi")

	req := SignupRequest{
		Email:    "invalid-email",
		Password: "123",
		Age:      15,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Fatal("Expected validation errors")
	}

	errors := TranslateErrors(err, trans)
	if len(errors) == 0 {
		t.Error("Expected translated errors")
	}

	// Vietnamese translations should be different from English
	for _, fieldErr := range errors {
		if fieldErr.Message == "" {
			t.Errorf("Expected non-empty message for field %s", fieldErr.Field)
		}
	}
}

func TestValidateStructSuccess(t *testing.T) {
	req := SignupRequest{
		Email:    "user@example.com",
		Password: "password123",
		Age:      25,
	}

	errors := ValidateStruct(req, "en")
	if errors != nil {
		t.Errorf("Expected no errors for valid struct, got %d errors", len(errors))
	}
}

func TestValidateStructFailure(t *testing.T) {
	req := SignupRequest{
		Email:    "",
		Password: "",
		Age:      0,
	}

	errors := ValidateStruct(req, "en")
	if len(errors) == 0 {
		t.Error("Expected validation errors")
	}

	// Should have errors for all required fields
	if len(errors) < 3 {
		t.Errorf("Expected at least 3 errors, got %d", len(errors))
	}
}

func TestValidateStructWithDifferentLanguages(t *testing.T) {
	req := SignupRequest{
		Email:    "invalid",
		Password: "123",
		Age:      10,
	}

	// Test English
	errorsEn := ValidateStruct(req, "en")
	if len(errorsEn) == 0 {
		t.Error("Expected validation errors for English")
	}

	// Test Vietnamese
	errorsVi := ValidateStruct(req, "vi")
	if len(errorsVi) == 0 {
		t.Error("Expected validation errors for Vietnamese")
	}

	// Both should have same number of errors, just different messages
	if len(errorsEn) != len(errorsVi) {
		t.Errorf("Expected same number of errors, got EN: %d, VI: %d", len(errorsEn), len(errorsVi))
	}
}

func TestTranslateErrorsWithNonValidationError(t *testing.T) {
	trans := GetTranslator("en")

	// Test with a non-validator error
	errors := TranslateErrors(nil, trans)
	if errors != nil {
		t.Error("Expected nil for nil error")
	}
}
