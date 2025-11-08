package response

// OnlyData represents a minimal response format with only data
type OnlyData[T any] struct {
	Data T `json:"data"`
}

// OnlyError represents a minimal response format with only error
type OnlyError struct {
	Error any `json:"error"`
}

// NewLean creates a new minimal data-only response
func NewLean[T any](data T) OnlyData[T] {
	return OnlyData[T]{
		Data: data,
	}
}

// NewLeanError creates a new minimal error-only response
func NewLeanError(err any) OnlyError {
	return OnlyError{
		Error: err,
	}
}
