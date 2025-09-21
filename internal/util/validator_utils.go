package util

import (
	"github.com/gin-gonic/gin/binding"
	zh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// 启用中文翻译器

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhCn := zh.New()
		uni := ut.New(zhCn, zhCn)
		trans, _ = uni.GetTranslator("zh")
		_ = zhtranslations.RegisterDefaultTranslations(v, trans)
	}
}
