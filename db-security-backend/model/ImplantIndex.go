package model

type ImplantIndex struct {
	Id        int64  `xorm:"pk autoincr" json:"id"`
	IndexName string `xorm:"varchar(20)" json:"index_name"`
	Phone     string `xorm:"varchar(11)" json:"phone"`
	Datetime  string `xorm:"datetime" json:"datetime"`
}
