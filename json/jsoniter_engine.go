package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

type jsoniterEngine struct {
	json jsoniter.API
}

func NewJsoniterEngine() Engine {
	return &jsoniterEngine{json: jsoniter.ConfigCompatibleWithStandardLibrary}
}

func (e *jsoniterEngine) Marshal(v interface{}) ([]byte, error) {
	return e.json.Marshal(v)
}

func (e *jsoniterEngine) Unmarshal(data []byte, v interface{}) error {
	return e.json.Unmarshal(data, v)
}

func (e *jsoniterEngine) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return e.json.MarshalIndent(v, prefix, indent)
}

func (e *jsoniterEngine) NewEncoder(w io.Writer) Encoder {
	return e.json.NewEncoder(w)
}

func (e *jsoniterEngine) NewDecoder(r io.Reader) Decoder {
	return e.json.NewDecoder(r)
}
