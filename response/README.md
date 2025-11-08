# 🌐 Go REST API Response Pattern — Flexible and Framework-Agnostic

A clean, consistent, and extensible **Result/Response Pattern** for REST APIs in Go.  
Supports multiple formats (Envelope, Flat, Lean, HAL, Validation) and works across frameworks like **Gin**, **Fiber**, **Echo**, etc.

---

## 🚀 Features

✅ Consistent and reusable response structure  
✅ Supports multiple response formats:
- **Envelope** – enterprise-style, predictable format  
- **Flat** – lightweight and human-readable  
- **Lean** – minimal and efficient  
- **HAL** – hypermedia (HATEOAS) style  
- **Validation** – rich field-level validation errors  

✅ Internationalization (**i18n**) with multi-language support (e.g., English, Vietnamese)  
✅ Compatible with **any HTTP framework** (Gin, Fiber, Echo, Chi, net/http)  
✅ Easy to extend with custom error codes or metadata  

---

## 🧱 Folder Structure

```
response/
│── envelope.go      # Envelope format
│── flat.go          # Flat format
│── lean.go          # Lean format
│── hal.go           # HAL format
│── validation.go    # Validation format
│── builder.go       # Common writer and helpers
│── i18n.go          # Validator + translations
```

---

## 📦 Install Dependencies

```bash
go get github.com/go-playground/validator/v10
go get github.com/go-playground/locales
go get github.com/go-playground/universal-translator
go get github.com/go-playground/validator/v10/translations/en
go get github.com/go-playground/validator/v10/translations/vi
```

---

## 🧩 Core Response Models

### 1. Envelope Format (Default)

```go
type Envelope[T any] struct {
    Success bool   `json:"success"`
    Data    T      `json:"data,omitempty"`
    Error   any    `json:"error,omitempty"`
    Meta    any    `json:"meta,omitempty"`
}
```

**Example**

```json
{
  "success": true,
  "data": { "id": 1, "name": "Alice" },
  "meta": { "page": 1, "size": 10, "total": 30 }
}
```

---

### 2. Flat Format

```go
type Flat[T any] struct {
    Status  string `json:"status"`   // "ok" | "error"
    Message string `json:"message"`
    Data    T      `json:"data,omitempty"`
}
```

**Example**

```json
{ 
  "status": "ok", 
  "message": "Fetched successfully", 
  "data": { "id": 1 } 
}
```

---

### 3. Lean Format

```go
type OnlyData[T any] struct { 
    Data T `json:"data"` 
}

type OnlyError struct { 
    Error any `json:"error"` 
}
```

**Example**

```json
{ "data": { "id": 1, "name": "Alice" } }
```

---

### 4. HAL Format (Hypermedia)

```go
type Link struct { 
    Href string `json:"href"` 
}

type HAL[T any] struct {
    Links map[string]Link `json:"_links"`
    Data  T               `json:"data"`
}
```

**Example**

```json
{
  "_links": {
    "self": { "href": "/v1/users/1" },
    "list": { "href": "/v1/users" }
  },
  "data": { "id": 1, "name": "Alice" }
}
```

---

### 5. Validation Format

```go
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type Validation struct {
    Code   string       `json:"code"` // e.g. INVALID_INPUT
    Errors []FieldError `json:"errors"`
}
```

**Example**

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "errors": [
      { "field": "email", "message": "must be valid email" },
      { "field": "password", "message": "min length 6" }
    ]
  }
}
```

---

## 🌍 Multi-Language Validation Setup

### Step 1: Initialize the validator with translators

```go
import (
    "github.com/laziness-coders/go-utils/response"
)

func init() {
    response.InitValidator()
}
```

### Step 2: Detect and register language per request

```go
func GetTranslator(lang string) ut.Translator {
    return response.GetTranslator(lang)
}
```

### Step 3: Translate validation errors

```go
type SignupRequest struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=6"`
}

func ValidateRequest(req SignupRequest, lang string) []response.FieldError {
    return response.ValidateStruct(req, lang)
}
```

---

## 🧰 Common Writer / Response Builder

### Using WriteResponse

```go
import (
    "net/http"
    "github.com/laziness-coders/go-utils/response"
)

func handler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"message": "Success"}
    format := response.ParseFormat(r.URL.Query().Get("format"))
    
    response.WriteResponse(w, http.StatusOK, data, nil, nil, format)
}
```

### Error Response

```go
func errorHandler(w http.ResponseWriter, r *http.Request) {
    appErr := &response.AppError{
        Code:    "NOT_FOUND",
        Message: "Resource not found",
    }
    format := response.ParseFormat(r.URL.Query().Get("format"))
    
    response.WriteResponse[any](w, http.StatusNotFound, nil, appErr, nil, format)
}
```

### Validation Error Response

```go
func validationHandler(w http.ResponseWriter, r *http.Request) {
    errors := []response.FieldError{
        {Field: "email", Message: "must be valid email"},
        {Field: "password", Message: "min length 6"},
    }
    format := response.ParseFormat(r.URL.Query().Get("format"))
    
    response.WriteValidationError(w, http.StatusBadRequest, errors, format)
}
```

---

## ⚡ Framework Integration Examples

### Gin

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/laziness-coders/go-utils/response"
)

type SignupRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func Signup(c *gin.Context) {
    lang := c.DefaultQuery("lang", "en")
    format := response.ParseFormat(c.DefaultQuery("format", "envelope"))
    
    var req SignupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        appErr := response.NewAppError("BAD_REQUEST", "invalid JSON")
        response.WriteResponse(c.Writer, 400, nil, &appErr, nil, format)
        return
    }

    if errors := response.ValidateStruct(req, lang); len(errors) > 0 {
        response.WriteValidationError(c.Writer, 400, errors, format)
        return
    }

    data := gin.H{"message": "Registered successfully"}
    response.WriteResponse(c.Writer, 200, data, nil, nil, format)
}

func main() {
    response.InitValidator()
    
    r := gin.Default()
    r.POST("/signup", Signup)
    r.Run(":8080")
}
```

### Fiber

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/laziness-coders/go-utils/response"
)

type SignupRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func Signup(c *fiber.Ctx) error {
    lang := c.Query("lang", "en")
    format := response.ParseFormat(c.Query("format", "envelope"))
    
    var req SignupRequest
    if err := c.BodyParser(&req); err != nil {
        appErr := response.NewAppError("BAD_REQUEST", "invalid JSON")
        c.Status(400)
        return response.WriteResponse(c, 400, nil, &appErr, nil, format)
    }

    if errors := response.ValidateStruct(req, lang); len(errors) > 0 {
        c.Status(400)
        return response.WriteValidationError(c, 400, errors, format)
    }

    data := map[string]string{"message": "Registered successfully"}
    return response.WriteResponse(c, 200, data, nil, nil, format)
}

func main() {
    response.InitValidator()
    
    app := fiber.New()
    app.Post("/signup", Signup)
    app.Listen(":8080")
}
```

### Standard net/http

```go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/laziness-coders/go-utils/response"
)

type SignupRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
    lang := r.URL.Query().Get("lang")
    if lang == "" {
        lang = "en"
    }
    format := response.ParseFormat(r.URL.Query().Get("format"))
    
    var req SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        appErr := response.NewAppError("BAD_REQUEST", "invalid JSON")
        response.WriteResponse(w, 400, nil, &appErr, nil, format)
        return
    }

    if errors := response.ValidateStruct(req, lang); len(errors) > 0 {
        response.WriteValidationError(w, 400, errors, format)
        return
    }

    data := map[string]string{"message": "Registered successfully"}
    response.WriteResponse(w, 200, data, nil, nil, format)
}

func main() {
    response.InitValidator()
    
    http.HandleFunc("/signup", Signup)
    http.ListenAndServe(":8080", nil)
}
```

---

## 🧪 Query Parameters

| Param    | Description                           | Example                      |
|----------|---------------------------------------|------------------------------|
| `format` | Choose output format                  | envelope, flat, lean, hal    |
| `lang`   | Choose language for validation msgs   | en, vi                       |

**Example Requests:**

```bash
# Envelope format with English validation
POST /signup?format=envelope&lang=en

# Flat format with Vietnamese validation
POST /signup?format=flat&lang=vi

# Lean format
POST /signup?format=lean

# HAL format
POST /signup?format=hal
```

---

## 🧠 Recommended Usage

| Scenario                    | Format       | Why                              |
|-----------------------------|--------------|----------------------------------|
| Public APIs (mobile, web)   | Envelope     | Standardized and safe            |
| Internal microservices      | Lean         | Lightweight, fast                |
| Admin dashboards            | Flat         | Easy to read and debug           |
| API discovery               | HAL          | Link-rich navigation             |
| Form input validation       | Validation   | Field-level feedback             |

---

## 🧩 Usage Examples

### Success Response with Metadata

```go
data := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
meta := map[string]int{"page": 1, "size": 10, "total": 100}
response.WriteResponse(w, 200, data, nil, meta, response.FormatEnvelope)
```

**Response:**

```json
{
  "success": true,
  "data": [
    {"id": 1, "name": "Alice"},
    {"id": 2, "name": "Bob"}
  ],
  "meta": {
    "page": 1,
    "size": 10,
    "total": 100
  }
}
```

### HAL Response with Links

```go
data := User{ID: 1, Name: "Alice"}
links := map[string]response.Link{
    "self": {Href: "/v1/users/1"},
    "list": {Href: "/v1/users"},
    "edit": {Href: "/v1/users/1/edit"},
}
halResponse := response.NewHAL(data, links)
```

**Response:**

```json
{
  "_links": {
    "self": {"href": "/v1/users/1"},
    "list": {"href": "/v1/users"},
    "edit": {"href": "/v1/users/1/edit"}
  },
  "data": {
    "id": 1,
    "name": "Alice"
  }
}
```

---

## 📚 API Reference

### Response Formats

#### `NewEnvelope[T any](data T) Envelope[T]`
Creates a success envelope response.

#### `NewEnvelopeWithMeta[T any](data T, meta any) Envelope[T]`
Creates a success envelope response with metadata.

#### `NewEnvelopeError(err any) Envelope[any]`
Creates an error envelope response.

#### `NewFlat[T any](message string, data T) Flat[T]`
Creates a flat success response.

#### `NewFlatError(message string) Flat[any]`
Creates a flat error response.

#### `NewLean[T any](data T) OnlyData[T]`
Creates a lean data-only response.

#### `NewLeanError(err any) OnlyError`
Creates a lean error-only response.

#### `NewHAL[T any](data T, links map[string]Link) HAL[T]`
Creates a HAL response with hypermedia links.

#### `NewValidation(code string, errors []FieldError) Validation`
Creates a validation error response.

### Validation & i18n

#### `InitValidator()`
Initializes the validator with multi-language support.

#### `GetValidator() *validator.Validate`
Returns the validator instance (lazy initialization).

#### `GetTranslator(lang string) ut.Translator`
Returns a translator for the specified language.

#### `ValidateStruct(s interface{}, lang string) []FieldError`
Validates a struct and returns translated field errors.

#### `TranslateErrors(err error, trans ut.Translator) []FieldError`
Translates validator errors to FieldError slice.

### Response Builders

#### `WriteResponse[T any](w http.ResponseWriter, statusCode int, data T, appErr *AppError, meta any, format ResponseFormat) error`
Writes a response in the specified format.

#### `WriteValidationError(w http.ResponseWriter, statusCode int, errors []FieldError, format ResponseFormat) error`
Writes a validation error response.

#### `ParseFormat(format string) ResponseFormat`
Parses a format string and returns the corresponding ResponseFormat.

---

## 🔧 Advanced Usage

### Custom Error Codes

```go
const (
    ErrCodeNotFound     = "NOT_FOUND"
    ErrCodeUnauthorized = "UNAUTHORIZED"
    ErrCodeValidation   = "VALIDATION_ERROR"
    ErrCodeInternal     = "INTERNAL_ERROR"
)

func handleError(w http.ResponseWriter, code string, message string) {
    appErr := response.NewAppError(code, message)
    response.WriteResponse[any](w, getStatusCode(code), nil, &appErr, nil, response.FormatEnvelope)
}

func getStatusCode(code string) int {
    switch code {
    case ErrCodeNotFound:
        return http.StatusNotFound
    case ErrCodeUnauthorized:
        return http.StatusUnauthorized
    case ErrCodeValidation:
        return http.StatusBadRequest
    default:
        return http.StatusInternalServerError
    }
}
```

### Middleware for Format Detection

```go
func FormatMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        format := r.URL.Query().Get("format")
        if format == "" {
            // Check Accept header
            accept := r.Header.Get("Accept")
            if strings.Contains(accept, "application/hal+json") {
                format = "hal"
            } else {
                format = "envelope"
            }
        }
        ctx := context.WithValue(r.Context(), "format", format)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Language Detection from Accept-Language Header

```go
func LanguageMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        lang := r.URL.Query().Get("lang")
        if lang == "" {
            // Parse Accept-Language header
            acceptLang := r.Header.Get("Accept-Language")
            if strings.HasPrefix(acceptLang, "vi") {
                lang = "vi"
            } else {
                lang = "en"
            }
        }
        ctx := context.WithValue(r.Context(), "lang", lang)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 🧪 Testing

Run the tests:

```bash
go test ./response/... -v
```

Example test output:

```
=== RUN   TestNewEnvelope
--- PASS: TestNewEnvelope (0.00s)
=== RUN   TestValidateStructSuccess
--- PASS: TestValidateStructSuccess (0.00s)
=== RUN   TestWriteResponseEnvelope
--- PASS: TestWriteResponseEnvelope (0.00s)
PASS
ok      github.com/laziness-coders/go-utils/response    0.005s
```

---

## 🧩 Next Steps

- [ ] Add custom error codes for business logic
- [ ] Extend language support (e.g., French, Japanese, Spanish)
- [ ] Add middleware for automatic detection of Accept-Language and format
- [ ] Add support for JSON:API specification
- [ ] Add support for Problem Details (RFC 7807)
- [ ] Add response caching headers support
- [ ] Add metrics and observability

---

## 📘 License

MIT © 2025 — Flexible REST API Response Pattern for Go

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📞 Support

For issues and questions, please open an issue on the GitHub repository.
