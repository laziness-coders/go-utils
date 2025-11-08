package response

// Envelope represents the enterprise-style response format with predictable structure
type Envelope[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
	Error   any  `json:"error,omitempty"`
	Meta    any  `json:"meta,omitempty"`
}

// NewEnvelope creates a new Envelope response with success status and data
func NewEnvelope[T any](data T) Envelope[T] {
	return Envelope[T]{
		Success: true,
		Data:    data,
	}
}

// NewEnvelopeWithMeta creates a new Envelope response with success status, data, and metadata
func NewEnvelopeWithMeta[T any](data T, meta any) Envelope[T] {
	return Envelope[T]{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

// NewEnvelopeError creates a new Envelope response with error
func NewEnvelopeError(err any) Envelope[any] {
	return Envelope[any]{
		Success: false,
		Error:   err,
	}
}
