package model

type User struct {
	Id          int64  `xorm:"pk autoincr" json:"id"`
	Phone       string `xorm:"varchar(11)" json:"phone"`
	Password    string `xorm:"varchar(255)" json:"password"`
	FingerPrint string `xorm:"varchar(20)" json:"finger_print"`
	IsSuperUser int8   `xorm:"tinyint" json:"is_super_user"`
	DateJoined  string  `xorm:"bigint" json:"date_joined"`
}
