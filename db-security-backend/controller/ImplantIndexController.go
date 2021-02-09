package controller

import (
	"db-security-backend/middleware"
	"db-security-backend/service"
	"db-security-backend/tool"
	"github.com/gin-gonic/gin"
)

type ImplantIndexController struct {
	iiss service.ImplantIndexService
}

func (iis *ImplantIndexController) Router(engine *gin.Engine) {
	iiGroup := engine.Group("/api/ii/admin")
	iiGroup.Use(middleware.JWTAuth())
	iiGroup.Use(middleware.AdminCheck())
	{
		iiGroup.POST("/iis", iis.getAllIp)
	}
}

//获取所有记录
func (iis *ImplantIndexController) getAllIp(ctx *gin.Context) {
	iiss, err := iis.iiss.GetAllII()
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, iiss)
}
