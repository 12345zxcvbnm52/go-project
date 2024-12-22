package initialize

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
	"go.uber.org/zap"

	gb "goods_web/global"
)

func InitTranslator(locale string) (err error) {
	//更换为gin的验证引擎
	//本质是gin内部封装了Validator包
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
		gb.Translator, ok = uni.GetTranslator(locale)
		if !ok {
			err := errors.New("获取本地翻译器失败")
			zap.S().Errorw(err.Error())
			return err
		}
		switch locale {
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, gb.Translator)
		case "en":
			en_translations.RegisterDefaultTranslations(v, gb.Translator)
		default:
			zh_translations.RegisterDefaultTranslations(v, gb.Translator)
		}
		return nil
	}
	panic(errors.New("无法转化为validator.Validate"))
}
