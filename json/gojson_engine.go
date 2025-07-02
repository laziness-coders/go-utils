package json

import (
	"io"

	gojson "github.com/goccy/go-json"
)

type gojsonEngine struct{}

func NewGoJsonEngine() Engine {
	return &gojsonEngine{}
}

func (e *gojsonEngine) Marshal(v interface{}) ([]byte, error) {
	return gojson.Marshal(v)
}

func (e *gojsonEngine) Unmarshal(data []byte, v interface{}) error {
	return gojson.Unmarshal(data, v)
}

func (e *gojsonEngine) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return gojson.MarshalIndent(v, prefix, indent)
}

func (e *gojsonEngine) NewEncoder(w io.Writer) Encoder {
	return gojson.NewEncoder(w)
}

func (e *gojsonEngine) NewDecoder(r io.Reader) Decoder {
	return gojson.NewDecoder(r)
}
