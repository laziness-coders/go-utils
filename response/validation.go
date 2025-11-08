package response

// FieldError represents a field-level validation error
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validation represents a validation error response with field-level details
type Validation struct {
	Code   string       `json:"code"` // e.g., INVALID_INPUT
	Errors []FieldError `json:"errors"`
}

// NewValidation creates a new validation error response
func NewValidation(code string, errors []FieldError) Validation {
	return Validation{
		Code:   code,
		Errors: errors,
	}
}

// AddError adds a field error to the validation response
func (v *Validation) AddError(field, message string) {
	v.Errors = append(v.Errors, FieldError{
		Field:   field,
		Message: message,
	})
}
