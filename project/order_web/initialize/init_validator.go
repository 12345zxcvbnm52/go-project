package initialize

import (
	"errors"
	"regexp"

	gb "order_web/global"

	"github.com/dlclark/regexp2"
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

func ValidateMobile() {
	gb.AddValidator("mobile",
		func(fl validator.FieldLevel) bool {
			//fl.Field可以获得正在认证的结构体字段(binding标签对应的字段)
			mobile := fl.Field().String()
			ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
			return ok
		},
	)
}

func TranslateMobile() {
	gb.AddTranslate("mobile", func(tag string, msg string) validator.RegisterTranslationsFunc {
		return func(trans ut.Translator) error {
			if err := trans.Add(tag, msg, false); err != nil {
				zap.S().Errorw("对手机号的检测消息的翻译失败", "msg", err.Error())
				return err
			}
			return nil
		}
	})
	gb.TranslateMsg["mobile"] = "{0}格式错误"
}

func ValidatePassword() {
	gb.AddValidator("password",
		func(fl validator.FieldLevel) bool {
			password := fl.Field().String()
			reg, _ := regexp2.Compile(`^(?![a-zA-Z!@#$%^&*]+$)(?![0-9!@#$%^&*]+$)[0-9A-Za-z!@#$%^&*]{8,16}$`, 0)
			ok, _ := reg.MatchString(password)
			return ok
		},
	)
}

func TranslatePassword() {
	gb.AddTranslate("password", func(tag string, msg string) validator.RegisterTranslationsFunc {
		return func(trans ut.Translator) error {
			if err := trans.Add(tag, msg, false); err != nil {
				zap.S().Errorw("对密码的检测消息的翻译失败", "msg", err.Error())
				return err
			}
			return nil
		}
	})
	gb.TranslateMsg["password"] = "{0}格式出错,密码必须在7-18位之内且必须由字母和数字组成"
}

func ValidateUserName() {
	gb.AddValidator("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		return len(username) <= 20
	})
}

func TranslateUserName() {
	gb.AddTranslate("username", func(tag, msg string) validator.RegisterTranslationsFunc {
		return func(trans ut.Translator) error {
			if err := trans.Add(tag, msg, false); err != nil {
				zap.S().Errorw("对用户名的检测消息的翻译失败", "msg", err.Error())
				return err
			}
			return nil
		}
	})
	gb.TranslateMsg["username"] = "{0}格式出错,用户名字数应不大于10个字符"
}

func InitValidator() {
	ValidateMobile()
	TranslateMobile()
	TranslatePassword()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for k, f := range gb.ValidateFiledConfig {
			v.RegisterValidation(k, f)
		}
		for k, f := range gb.TranslateConfig {
			v.RegisterTranslation(k, gb.Translator, f(k, gb.TranslateMsg[k]), translate)
		}
	} else {
		zap.S().Errorw("Validate转化失败")
		panic(errors.New("Validate转化失败"))
	}
}
