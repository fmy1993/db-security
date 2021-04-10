package userParam

type LoginParam struct {
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	CaptchaId    string `json:"captchaId"`
	CaptchaValue string `json:"captchaValue"`
}
