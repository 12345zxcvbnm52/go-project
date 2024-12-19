package global

import "github.com/go-playground/validator/v10"

var ValidateStructConfig map[string]func(validator.StructLevel) = make(map[string]func(validator.StructLevel))

var ValidateFiledConfig map[string]func(validator.FieldLevel) bool = make(map[string]func(validator.FieldLevel) bool)

var TranslateConfig map[string]func(string, string) validator.RegisterTranslationsFunc = map[string]func(string, string) validator.RegisterTranslationsFunc{}

var TranslateMsg map[string]string = map[string]string{}

func AddValidator(tag string, f func(validator.FieldLevel) bool) {
	ValidateFiledConfig[tag] = f
}

func AddTranslate(tag string, f func(string, string) validator.RegisterTranslationsFunc) {
	TranslateConfig[tag] = f
}
