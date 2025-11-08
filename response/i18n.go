package response

import (
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/vi"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	vi_trans "github.com/go-playground/validator/v10/translations/vi"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
)

// InitValidator initializes the validator with multi-language support
func InitValidator() {
	validate = validator.New()
	uni = ut.New(en.New(), en.New(), vi.New())
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}

// GetTranslator detects and returns the appropriate translator based on language code
// Supports "en" (English) and "vi" (Vietnamese)
func GetTranslator(lang string) ut.Translator {
	if validate == nil || uni == nil {
		InitValidator()
	}

	base := "en"
	if strings.HasPrefix(strings.ToLower(lang), "vi") {
		base = "vi"
	}

	trans, _ := uni.GetTranslator(base)
	if base == "vi" {
		vi_trans.RegisterDefaultTranslations(validate, trans)
	} else {
		en_trans.RegisterDefaultTranslations(validate, trans)
	}
	return trans
}

// TranslateErrors converts validator errors to FieldError slice with translations
func TranslateErrors(err error, trans ut.Translator) []FieldError {
	var out []FieldError
	if ves, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ves {
			out = append(out, FieldError{
				Field:   e.Field(),
				Message: e.Translate(trans),
			})
		}
	}
	return out
}

// ValidateStruct validates a struct and returns translated field errors
func ValidateStruct(s interface{}, lang string) []FieldError {
	v := GetValidator()
	trans := GetTranslator(lang)

	err := v.Struct(s)
	if err == nil {
		return nil
	}

	return TranslateErrors(err, trans)
}
