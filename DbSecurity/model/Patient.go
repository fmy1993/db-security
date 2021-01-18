package model

import "github.com/shopspring/decimal"

type Patient struct {
	Id      int64           `xorm:"pk autoincr" json:"id"`
	Name    string          `xorm:"varchar(10)" json:"name"`
	Weight  decimal.Decimal `xorm:"decimal(3,1)" json:"weight"`
	High    int             `xorm:"int" json:"high"`
	Age     int             `xorm:"int" json:"age"`
	IdCard  string          `xorm:"varchar(18)" json:"id_card"`
	Phone   string          `xorm:"varchar(11)" json:"phone"`
	Address string          `xorm:"varchar(50)" json:"address"`
	Bill    decimal.Decimal `xorm:"decimal(6,2)" json:"bill"`
}

type PatientCopy struct {
	Id      int64           `xorm:"pk autoincr" json:"id"`
	Weight  decimal.Decimal `xorm:"decimal(3,1)" json:"weight"`
	High    int             `xorm:"int" json:"high"`
	Age     int             `xorm:"int" json:"age"`
	Phone   string          `xorm:"varchar(11)" json:"phone"`
	Address string          `xorm:"varchar(50)" json:"address"`
	Bill    decimal.Decimal `xorm:"decimal(6,2)" json:"bill"`
}
