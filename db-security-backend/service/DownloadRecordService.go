package service

import (
	"db-security-backend/dao"
	"db-security-backend/model"
	"time"
)

type DownloadRecordService struct {
}

// InsertDownloadRecord 插入下载记录
func (drs *DownloadRecordService) InsertDownloadRecord(userId int64, phone string) int64 {
	downloadRecordDao := dao.NewDownloadRecordDao()
	downloadRecord := model.DownloadRecord{}
	downloadRecord.UserId = userId
	downloadRecord.Phone = phone
	downloadRecord.DownloadTime = time.Now().Format("2006-01-02 15:04:05")
	return downloadRecordDao.InsertDownloadRecord(downloadRecord)
}

// GetAllRecords 得到downloadRecord表的所有数据
func (drs *DownloadRecordService) GetAllRecords() (*[]model.DownloadRecord, error) {
	var downloadRecordDao = dao.NewDownloadRecordDao()
	return downloadRecordDao.GetAllRecord()
}
