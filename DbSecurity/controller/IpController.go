package controller

import (
	"DbSecurity/middleware"
	"DbSecurity/model"
	"DbSecurity/service"
	"DbSecurity/tool"
	"github.com/gin-gonic/gin"
)

type IpController struct {
}

func (is *IpController) Router(engine *gin.Engine) {
	ipGroup := engine.Group("/api/ip/admin")
	ipGroup.Use(middleware.JWTAuth())
	ipGroup.Use(middleware.AdminCheck())
	{
		ipGroup.POST("/free", is.freeIp)
		ipGroup.POST("/ips", is.getAllIp)
		ipGroup.POST("/add", is.freezeIp)
	}
}

//解封ip
func (is *IpController) freeIp(ctx *gin.Context) {
	var ip model.BannedIp
	err := tool.Decode(ctx.Request.Body, &ip)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	var ipService service.BannedIpService
	err = ipService.FreeIp(&ip)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "解封成功")
}

//获取所有被封禁的ip
func (is *IpController) getAllIp(ctx *gin.Context) {
	var ipService service.BannedIpService
	ips, err := ipService.GetAllIp()
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, ips)
}

//封禁ip
func (is *IpController) freezeIp(ctx *gin.Context) {
	var ipService service.BannedIpService
	var ip model.BannedIp
	err := tool.Decode(ctx.Request.Body, &ip)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	ipService.BanIp(ip.Ip)
	tool.Success(ctx, "封禁成功")
}
