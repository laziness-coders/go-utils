package response

// Flat represents a lightweight and human-readable response format
type Flat[T any] struct {
	Status  string `json:"status"` // "ok" | "error"
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// NewFlat creates a new Flat response with success status
func NewFlat[T any](message string, data T) Flat[T] {
	return Flat[T]{
		Status:  "ok",
		Message: message,
		Data:    data,
	}
}

// NewFlatError creates a new Flat response with error status
func NewFlatError(message string) Flat[any] {
	return Flat[any]{
		Status:  "error",
		Message: message,
	}
}
