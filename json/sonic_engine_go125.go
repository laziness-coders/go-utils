//go:build go1.25

package json

import "io"

type sonicEngine struct {
	fallback Engine
}

func NewSonicEngine() Engine {
	return &sonicEngine{fallback: NewStdEngine()}
}

func (e *sonicEngine) Marshal(v interface{}) ([]byte, error) {
	return e.fallback.Marshal(v)
}

func (e *sonicEngine) Unmarshal(data []byte, v interface{}) error {
	return e.fallback.Unmarshal(data, v)
}

func (e *sonicEngine) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return e.fallback.MarshalIndent(v, prefix, indent)
}

func (e *sonicEngine) NewEncoder(w io.Writer) Encoder {
	return e.fallback.NewEncoder(w)
}

func (e *sonicEngine) NewDecoder(r io.Reader) Decoder {
	return e.fallback.NewDecoder(r)
}
