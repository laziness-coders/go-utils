package json

import (
	"encoding/json"
	"io"
)

type stdEngine struct{}

func NewStdEngine() Engine {
	return &stdEngine{}
}

func (e *stdEngine) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (e *stdEngine) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (e *stdEngine) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (e *stdEngine) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (e *stdEngine) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}
