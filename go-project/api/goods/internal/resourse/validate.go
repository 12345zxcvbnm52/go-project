package resourse

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func InitValidate() {
	GoodsServer.Server.Validator.AddFleidValidator("mobile",
		"{0}不是一个有效的手机号码",
		func(fl validator.FieldLevel) bool {
			//fl.Field可以获得正在认证的结构体字段(binding标签对应的字段)
			mobile := fl.Field().String()
			ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
			return ok
		},
	)
	GoodsServer.Server.Validator.AddFleidTranslator("mobile", func(tag string, msg string) validator.RegisterTranslationsFunc {
		return func(trans ut.Translator) error {
			if err := trans.Add(tag, msg, false); err != nil {
				return err
			}
			return nil
		}
	})

}
