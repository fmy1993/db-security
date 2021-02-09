package param

type RegisterParam struct {
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	CheckCode string `json:"checkCode"`
}
