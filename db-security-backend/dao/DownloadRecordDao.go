package dao

import (
	"db-security-backend/config"
	"db-security-backend/model"
)

type DownloadRecordDao struct {
	*config.Orm
}

func NewDownloadRecordDao() *DownloadRecordDao {
	return &DownloadRecordDao{config.DbEngine}
}

// InsertDownloadRecord 添加下载记录信息
func (drd *DownloadRecordDao) InsertDownloadRecord(record model.DownloadRecord) int64 {
	result, err := drd.InsertOne(record)
	if err != nil {
		return 0
	}
	return result
}

// GetAllRecord 获取所有下载记录
func (drd *DownloadRecordDao) GetAllRecord() (*[]model.DownloadRecord, error) {
	var records []model.DownloadRecord
	err := drd.Find(&records)
	if err != nil {
		return nil, err
	}
	return &records, nil
}
