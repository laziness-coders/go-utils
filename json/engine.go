package json

import (
	"github.com/spf13/viper"
	"io"
	"log"
)

// Set the best engines based on benchmarks:
// - Marshal: Sonic (fastest)
// - Unmarshal: Sonic (fastest)
var marshalEngine Engine = NewSonicEngine()
var unmarshalEngine Engine = NewSonicEngine()

func SetEngineFromEnv() {
	marshalEngineNam := viper.GetString("JSON_MARSHAL_ENGINE")
	unmarshalEngineNam := viper.GetString("JSON_UNMARSHAL_ENGINE")
	if marshalEngineNam != "" {
		SetMarshalEngine(GetEngineByName(marshalEngineNam))
	}

	if unmarshalEngineNam != "" {
		SetUnmarshalEngine(GetEngineByName(unmarshalEngineNam))
	}

	log.Printf("Using json marshal engine: %T", marshalEngine)
	log.Printf("Using json unmarshal engine: %T", unmarshalEngine)
}

func SetEngine(marshal, unmarshal Engine) {
	SetMarshalEngine(marshal)
	SetUnmarshalEngine(unmarshal)
}

func SetMarshalEngine(engine Engine) {
	if engine != nil {
		marshalEngine = engine
	}
}

func SetUnmarshalEngine(engine Engine) {
	if engine != nil {
		unmarshalEngine = engine
	}
}

func GetEngineByName(name string) Engine {
	switch name {
	case EngineStandard:
		return NewStdEngine()
	case EngineGoJson:
		return NewGoJsonEngine()
	case EngineJsoniter:
		return NewJsoniterEngine()
	case EngineSonic:
		return NewSonicEngine()
	default:
		return NewStdEngine()
	}
}

// Marshal serializes the input using the marshal engine.
func Marshal(v interface{}) ([]byte, error) {
	return marshalEngine.Marshal(v)
}

// Unmarshal deserializes the input using the unmarshal engine.
func Unmarshal(data []byte, v interface{}) error {
	return unmarshalEngine.Unmarshal(data, v)
}

// MarshalIndent serializes the input with indentation using the marshal engine.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return marshalEngine.MarshalIndent(v, prefix, indent)
}

// NewEncoder returns a new encoder using the marshal engine.
func NewEncoder(w io.Writer) Encoder {
	return marshalEngine.NewEncoder(w)
}

// NewDecoder returns a new decoder using the unmarshal engine.
func NewDecoder(r io.Reader) Decoder {
	return unmarshalEngine.NewDecoder(r)
}
