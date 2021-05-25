package model

type DownloadRecord struct {
	RecordId     int64  `xorm:"pk autoincr" json:"record_id"`
	UserId       int64  `xorm:"int" json:"user_id"`
	Phone        string `xorm:"varchar(11)" json:"phone"`
	DownloadTime string `xorm:"bigint" json:"download_time"`
}
