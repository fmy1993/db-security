package controller

import (
	"db-security-backend/middleware"
	"db-security-backend/service"
	"db-security-backend/util"
	"github.com/gin-gonic/gin"
)

type DownloadRecordController struct {
	drs service.DownloadRecordService
}

func (drc *DownloadRecordController) Router(engine *gin.Engine) {
	engine.GET("/download_record", middleware.JWTAuth(), middleware.AdminCheck(), drc.GetAllRecords)
}

// GetAllRecords 获取所有下载记录
func (drc *DownloadRecordController) GetAllRecords(ctx *gin.Context) {
	records, err := drc.drs.GetAllRecords()
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, records)
}
