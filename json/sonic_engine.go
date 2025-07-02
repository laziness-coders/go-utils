package json

import (
	"io"

	"github.com/bytedance/sonic"
)

type sonicEngine struct{}

func NewSonicEngine() Engine {
	// Set the default configuration for Sonic to conform to the standard JSON format.
	sonic.ConfigDefault = sonic.ConfigStd
	return &sonicEngine{}
}

func (e *sonicEngine) Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func (e *sonicEngine) Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}

func (e *sonicEngine) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return sonic.MarshalIndent(v, prefix, indent)
}

func (e *sonicEngine) NewEncoder(w io.Writer) Encoder {
	return sonic.ConfigStd.NewEncoder(w)
}

func (e *sonicEngine) NewDecoder(r io.Reader) Decoder {
	return sonic.ConfigStd.NewDecoder(r)
}
