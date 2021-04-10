package util

import (
	"github.com/mojocn/base64Captcha"

	"github.com/gin-gonic/gin"
)

var store = base64Captcha.DefaultMemStore

type Captcha struct {
	Id          string `json:"id"`
	Base64Blob  string `json:"base_64_blob"`
	VerifyValue string `json:"verify_value"`
}

func GenerateCaptcha(ctx *gin.Context) map[string]interface{} {
	var driver base64Captcha.Driver
	driver = &base64Captcha.DriverDigit{Height: 40, Width: 180, Length: 4, MaxSkew: 0.4, DotCount: 80}
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	body := map[string]interface{}{
		"code":      1,
		"data":      b64s,
		"captchaId": id,
		"msg":       "success",
	}
	if err != nil {
		body = map[string]interface{}{"code": 0, "msg": err.Error()}
	}
	return body
}

func VerifyCaptcha(id string, value string) bool {
	if store.Verify(id, value, true) {
		return true
	}
	return false
}
