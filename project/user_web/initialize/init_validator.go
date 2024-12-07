package initialize

import (
	"errors"
	"regexp"

	gb "user_web/global"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 用于辅助字段检测的翻译
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		zap.S().Errorw("翻译器翻译失败", "msg", err.Error())
		return "内部错误,请再次尝试"
	}
	return msg
}

func ValidateMobile(fl validator.FieldLevel) bool {
	//fl.Field可以获得正在认证的结构体字段(binding标签对应的字段)
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return ok
}

func TranslateMobile(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			zap.S().Errorw("对手机号的检测消息的翻译失败", "msg", err.Error())
			return err
		}
		return nil
	}
}

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", ValidateMobile)
		v.RegisterTranslation("mobile", gb.Translator, TranslateMobile("mobile", "{0}格式错误"), translate)
	} else {
		zap.S().Errorw("Validate转化失败")
		panic(errors.New("Validate转化失败"))
	}
}
