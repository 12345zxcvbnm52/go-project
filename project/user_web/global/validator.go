package global

import "github.com/go-playground/validator"

var ValidateItemConfig map[string]func(validator.StructLevel)

var ValidateStructConfig map[string]func(validator.FieldLevel) bool

var TranslateConfig map[string]func(string, string) validator.RegisterTranslationsFunc
