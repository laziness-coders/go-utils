# Custom JSON Package

This package provides a unified interface for JSON serialization and deserialization in Go, with pluggable backends for maximum performance and flexibility. Supported engines:

- Standard Library (`encoding/json`)
- [jsoniter](https://github.com/json-iterator/go)
- [go-json](https://github.com/goccy/go-json)
- [sonic](https://github.com/bytedance/sonic)

## Features
- Drop-in replacement for `encoding/json`
- Switch engines at runtime for optimal performance
- Supports Marshal, Unmarshal, MarshalIndent, Encoder, and Decoder

## Benchmark Results

| Operation                | Std      | Jsoniter | Go-Json  | Sonic    | Fastest Engine      |
|--------------------------|----------|----------|----------|----------|---------------------|
| Marshal                  | 172,145  | 199,846  | 147,158  | **52,686**   | **Sonic**           |
| Unmarshal                | 863,716  | 366,048  | 266,201  | **124,791**  | **Sonic**           |
| MarshalIndent            | 665,004  | 214,448  | **177,620**  | 474,772  | **Go-Json**         |
| Encoder (Encode)         | 156,342  | 264,314  | 137,496  | **90,464**   | **Sonic**           |
| Decoder (Decode)         | 829,561  | 485,304  | 312,691  | **134,663**  | **Sonic**           |

*All times are in ns/op (lower is better). Benchmarks run on Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz, goos: darwin, goarch: amd64.*

## Benchmark Chart

```mermaid
graph LR
    Marshal[Marshal]
    Unmarshal[Unmarshal]
    MarshalIndent[MarshalIndent]
    Encoder[Encoder]
    Decoder[Decoder]
    Std[Std]
    Jsoniter[Jsoniter]
    GoJson[Go-Json]
    Sonic[Sonic]

    Marshal -- 172145 --> Std
    Marshal -- 199846 --> Jsoniter
    Marshal -- 147158 --> GoJson
    Marshal -- 52686 --> Sonic
    Unmarshal -- 863716 --> Std
    Unmarshal -- 366048 --> Jsoniter
    Unmarshal -- 266201 --> GoJson
    Unmarshal -- 124791 --> Sonic
    MarshalIndent -- 665004 --> Std
    MarshalIndent -- 214448 --> Jsoniter
    MarshalIndent -- 177620 --> GoJson
    MarshalIndent -- 474772 --> Sonic
    Encoder -- 156342 --> Std
    Encoder -- 264314 --> Jsoniter
    Encoder -- 137496 --> GoJson
    Encoder -- 90464 --> Sonic
    Decoder -- 829561 --> Std
    Decoder -- 485304 --> Jsoniter
    Decoder -- 312691 --> GoJson
    Decoder -- 134663 --> Sonic
```

## Usage Example

```go
import myjson "github.com/laziness-coders/go-utils/json"

// Use Sonic for best performance
myjson.SetMarshalEngine(myjson.NewSonicEngine())
myjson.SetUnmarshalEngine(myjson.NewSonicEngine())

// For pretty-printing, use Go-Json
myjson.SetMarshalEngine(myjson.NewGoJsonEngine())

// Marshal
data, err := myjson.Marshal(obj)

// Unmarshal
err = myjson.Unmarshal(data, &obj)

// MarshalIndent
pretty, err := myjson.MarshalIndent(obj, "", "  ")
```

## License
MIT 