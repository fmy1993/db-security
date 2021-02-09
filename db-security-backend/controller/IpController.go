package controller

import (
	"db-security-backend/middleware"
	"db-security-backend/model"
	"db-security-backend/service"
	"db-security-backend/tool"
	"github.com/gin-gonic/gin"
)

type IpController struct {
	is service.BannedIpService
}

func (ic *IpController) Router(engine *gin.Engine) {
	ipGroup := engine.Group("/api/ip/admin")
	ipGroup.Use(middleware.JWTAuth())
	ipGroup.Use(middleware.AdminCheck())
	{
		ipGroup.POST("/free", ic.freeIp)
		ipGroup.POST("/ips", ic.getAllIp)
		ipGroup.POST("/add", ic.freezeIp)
	}
}

//解封ip
func (ic *IpController) freeIp(ctx *gin.Context) {
	var ip model.BannedIp
	err := tool.Decode(ctx.Request.Body, &ip)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	err = ic.is.FreeIp(&ip)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "解封成功")
}

//获取所有被封禁的ip
func (ic *IpController) getAllIp(ctx *gin.Context) {
	ips, err := ic.is.GetAllIp()
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, ips)
}

//封禁ip
func (ic *IpController) freezeIp(ctx *gin.Context) {
	var ip model.BannedIp
	err := tool.Decode(ctx.Request.Body, &ip)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	ic.is.BanIp(ip.Ip)
	tool.Success(ctx, "封禁成功")
}
