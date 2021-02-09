package model

type BannedIp struct {
	Id int64  `xorm:"pk autoincr" json:"id"`
	Ip string `xorm:"varchar(20)" json:"ip"`
}
