package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"sync"
)

var translationMutex sync.Mutex

// i18n 国际化中间件
func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			// RegisterDefaultTranslations 中有数据静态的 bug
			translationMutex.Lock()
			// RegisterDefaultTranslations 将验证器和对应语言类型的 Translator 注册进来
			switch locale {
			case "zh":
				_ = zhTranslations.RegisterDefaultTranslations(v, trans)
			case "en":
				_ = enTranslations.RegisterDefaultTranslations(v, trans)
			default:
				_ = zhTranslations.RegisterDefaultTranslations(v, trans)
			}
			translationMutex.Unlock()
			// 供后续翻译使用
			c.Set("trans", trans)
		}
		c.Next()
	}
}
