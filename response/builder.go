package response

import (
	"encoding/json"
	"net/http"
)

// AppError represents a common application error structure
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// NewAppError creates a new application error
func NewAppError(code, message string) AppError {
	return AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithDetails creates a new application error with details
func NewAppErrorWithDetails(code, message string, details any) AppError {
	return AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ResponseFormat represents the supported response formats
type ResponseFormat string

const (
	FormatEnvelope ResponseFormat = "envelope"
	FormatFlat     ResponseFormat = "flat"
	FormatLean     ResponseFormat = "lean"
	FormatHAL      ResponseFormat = "hal"
)

// WriteResponse writes a response in the specified format
func WriteResponse[T any](w http.ResponseWriter, statusCode int, data T, appErr *AppError, meta any, format ResponseFormat) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var response interface{}

	switch format {
	case FormatFlat:
		if appErr != nil {
			response = NewFlatError(appErr.Message)
		} else {
			response = NewFlat("Success", data)
		}
	case FormatLean:
		if appErr != nil {
			response = NewLeanError(appErr)
		} else {
			response = NewLean(data)
		}
	case FormatHAL:
		// For HAL format, expect data to contain links information
		// This is a simplified version; in real use, you'd extract links from data or meta
		if appErr != nil {
			response = NewEnvelopeError(appErr)
		} else {
			response = NewHAL(data, make(map[string]Link))
		}
	default: // FormatEnvelope
		if appErr != nil {
			response = NewEnvelopeError(appErr)
		} else {
			if meta != nil {
				response = NewEnvelopeWithMeta(data, meta)
			} else {
				response = NewEnvelope(data)
			}
		}
	}

	return json.NewEncoder(w).Encode(response)
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, statusCode int, errors []FieldError, format ResponseFormat) error {
	validation := NewValidation("INVALID_INPUT", errors)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var response interface{}

	switch format {
	case FormatFlat:
		response = Flat[any]{
			Status:  "error",
			Message: "Validation failed",
			Data:    validation,
		}
	case FormatLean:
		response = NewLeanError(validation)
	default: // FormatEnvelope or FormatHAL
		response = NewEnvelopeError(validation)
	}

	return json.NewEncoder(w).Encode(response)
}

// ParseFormat parses a format string and returns the corresponding ResponseFormat
func ParseFormat(format string) ResponseFormat {
	switch format {
	case "flat":
		return FormatFlat
	case "lean":
		return FormatLean
	case "hal":
		return FormatHAL
	default:
		return FormatEnvelope
	}
}
