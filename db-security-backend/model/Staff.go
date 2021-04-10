package model

import "github.com/shopspring/decimal"

type Staff struct {
	StaffId       int64           `xorm:"pk autoincr" json:"staff_id"`
	StaffName     string          `xorm:"varchar(30)" json:"staff_name"`
	Height        decimal.Decimal `xorm:"decimal(10,1)" json:"height"`
	Weight        decimal.Decimal `xorm:"decimal(10,2)" json:"weight"`
	Qualification string          `xorm:"varchar(30)" json:"qualification"`
	IdCard        string          `xorm:"char(18)" json:"id_card"`
	Salary        decimal.Decimal `xorm:"decimal(10)" json:"salary"`
}

type StaffCopy struct {
	StaffId       int64           `xorm:"pk autoincr" json:"staff_id"`
	Height        decimal.Decimal `xorm:"decimal(10,1)" json:"height"`
	Weight        decimal.Decimal `xorm:"decimal(10,2)" json:"weight"`
	Qualification string          `xorm:"varchar(30)" json:"qualification"`
	Salary        decimal.Decimal `xorm:"decimal(10)" json:"salary"`
}
