package json

import (
	"bytes"
	"testing"
)

type LargeStruct struct {
	ID      int
	Name    string
	Values  []int
	Nested  []NestedStruct
	MapData map[string]NestedStruct
}

type NestedStruct struct {
	Field1 string
	Field2 float64
	Field3 []byte
}

func newLargeStruct() LargeStruct {
	nested := make([]NestedStruct, 100)
	for i := range nested {
		nested[i] = NestedStruct{
			Field1: "nested field",
			Field2: float64(i) * 1.23,
			Field3: make([]byte, 256),
		}
	}
	mapData := make(map[string]NestedStruct, 100)
	for i := 0; i < 100; i++ {
		mapData[string(rune('A'+i))] = nested[i]
	}
	return LargeStruct{
		ID:      12345,
		Name:    "Large Test Struct",
		Values:  make([]int, 1000),
		Nested:  nested,
		MapData: mapData,
	}
}

func BenchmarkStdEngineMarshal(b *testing.B) {
	e := NewStdEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterEngineMarshal(b *testing.B) {
	e := NewJsoniterEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoJsonEngineMarshal(b *testing.B) {
	e := NewGoJsonEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSonicEngineMarshal(b *testing.B) {
	SetMarshalEngine(NewSonicEngine())
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStdEngineUnmarshal(b *testing.B) {
	e := NewStdEngine()
	SetUnmarshalEngine(e)
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out LargeStruct
		if err := Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterEngineUnmarshal(b *testing.B) {
	e := NewJsoniterEngine()
	SetUnmarshalEngine(e)
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out LargeStruct
		if err := Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoJsonEngineUnmarshal(b *testing.B) {
	e := NewGoJsonEngine()
	SetUnmarshalEngine(e)
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out LargeStruct
		if err := Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSonicEngineUnmarshal(b *testing.B) {
	SetUnmarshalEngine(NewSonicEngine())
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out LargeStruct
		if err := Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

// MarshalIndent Benchmarks
func BenchmarkStdEngineMarshalIndent(b *testing.B) {
	e := NewStdEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalIndent(obj, "", "  ")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterEngineMarshalIndent(b *testing.B) {
	e := NewJsoniterEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalIndent(obj, "", "  ")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoJsonEngineMarshalIndent(b *testing.B) {
	e := NewGoJsonEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalIndent(obj, "", "  ")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSonicEngineMarshalIndent(b *testing.B) {
	SetMarshalEngine(NewSonicEngine())
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalIndent(obj, "", "  ")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Encoder Benchmarks
func BenchmarkStdEngineEncoder(b *testing.B) {
	e := NewStdEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		enc := NewEncoder(buf)
		if err := enc.Encode(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterEngineEncoder(b *testing.B) {
	e := NewJsoniterEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		enc := NewEncoder(buf)
		if err := enc.Encode(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoJsonEngineEncoder(b *testing.B) {
	e := NewGoJsonEngine()
	SetMarshalEngine(e)
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		enc := NewEncoder(buf)
		if err := enc.Encode(obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSonicEngineEncoder(b *testing.B) {
	SetMarshalEngine(NewSonicEngine())
	obj := newLargeStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		enc := NewEncoder(buf)
		if err := enc.Encode(obj); err != nil {
			b.Fatal(err)
		}
	}
}

// Decoder Benchmarks
func BenchmarkStdEngineDecoder(b *testing.B) {
	e := NewStdEngine()
	SetUnmarshalEngine(e)
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data)
		dec := NewDecoder(buf)
		var out LargeStruct
		if err := dec.Decode(&out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterEngineDecoder(b *testing.B) {
	e := NewJsoniterEngine()
	SetUnmarshalEngine(e)
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data)
		dec := NewDecoder(buf)
		var out LargeStruct
		if err := dec.Decode(&out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoJsonEngineDecoder(b *testing.B) {
	e := NewGoJsonEngine()
	SetUnmarshalEngine(e)
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data)
		dec := NewDecoder(buf)
		var out LargeStruct
		if err := dec.Decode(&out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSonicEngineDecoder(b *testing.B) {
	SetUnmarshalEngine(NewSonicEngine())
	obj := newLargeStruct()
	data, err := Marshal(obj)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data)
		dec := NewDecoder(buf)
		var out LargeStruct
		if err := dec.Decode(&out); err != nil {
			b.Fatal(err)
		}
	}
}
