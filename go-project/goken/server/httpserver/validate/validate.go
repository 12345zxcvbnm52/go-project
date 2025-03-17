package validate

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type Validator struct {
	Trans ut.Translator
	//存储对结构体中自定义数量字段验证的函数
	ValidateStructFuncs map[string]func(validator.StructLevel)
	//存储对结构体中某一字段验证的函数
	ValidateFleidFuncs map[string]func(validator.FieldLevel) bool
	//存储对结构体中某一字段验证信息进行翻译的函数
	TranslateFleidFuncs map[string]func(string, string) validator.RegisterTranslationsFunc
	//认证失败的原始信息,可按需求翻译
	ValidateFleidFailedMsg map[string]string
	//翻译器函数,用于翻译错误信息
	TranslateFunc func(ut.Translator, validator.FieldError) string
}

func NewValidator(locale string) (*Validator, error) {
	v := &Validator{
		ValidateStructFuncs:    make(map[string]func(validator.StructLevel)),
		ValidateFleidFuncs:     make(map[string]func(validator.FieldLevel) bool),
		TranslateFleidFuncs:    make(map[string]func(string, string) validator.RegisterTranslationsFunc),
		ValidateFleidFailedMsg: make(map[string]string),
	}
	v.TranslateFunc = v.translate
	trans, err := NewTranslator(locale)
	if err != nil {
		return nil, err
	}
	v.Trans = trans
	return v, nil
}

func NewTranslator(locale string) (ut.Translator, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//将struct字段转为tag的json字段
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)
		trans, ok := uni.GetTranslator(locale)
		if !ok {
			err := errors.New("获取本地翻译器失败")
			return nil, err
		}
		switch locale {
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
		case "en":
			en_translations.RegisterDefaultTranslations(v, trans)
		default:
			zh_translations.RegisterDefaultTranslations(v, trans)
		}
		return trans, nil
	}
	return nil, errors.New("无法获取支持gin的翻译器")
}

func (v *Validator) translate(trans ut.Translator, fe validator.FieldError) string {
	msg, _ := trans.T(fe.Tag(), fe.Field())
	return msg
}

func (v *Validator) AddFleidValidator(tag string, errMsg string, f func(validator.FieldLevel) bool) {
	v.ValidateFleidFailedMsg[tag] = errMsg
	v.ValidateFleidFuncs[tag] = f
}

func (v *Validator) AddFleidTranslator(tag string, f func(string, string) validator.RegisterTranslationsFunc) {
	v.TranslateFleidFuncs[tag] = f
}

func (v *Validator) Excute() error {
	if va, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for k, f := range v.ValidateFleidFuncs {
			va.RegisterValidation(k, f)
		}
		for k, f := range v.TranslateFleidFuncs {
			va.RegisterTranslation(k, v.Trans, f(k, v.ValidateFleidFailedMsg[k]), v.TranslateFunc)
		}
		return nil
	} else {
		return errors.New("Validate出现未知错误,无法得到gin内部Validator")
	}
}
