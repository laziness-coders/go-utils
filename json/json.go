package json

import (
	"io"
)

// Engine defines the base interface for JSON operations.
type Engine interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}
