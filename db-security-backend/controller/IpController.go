package controller

import (
	"db-security-backend/middleware"
	"db-security-backend/model"
	"db-security-backend/service"
	"db-security-backend/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

type IpController struct {
	is service.BannedIpService
}

func (ic *IpController) Router(engine *gin.Engine) {
	engine.DELETE("/free_ip/:ipId", middleware.JWTAuth(), middleware.AdminCheck(), ic.freeIp)
	engine.GET("/all_ip", middleware.JWTAuth(), middleware.AdminCheck(), ic.getAllIp)
	engine.POST("/freeze_ip", middleware.JWTAuth(), middleware.AdminCheck(), ic.freezeIp)
}

//解封ip
func (ic *IpController) freeIp(ctx *gin.Context) {
	ipId, _ := strconv.Atoi(ctx.Param("ipId"))
	err := ic.is.FreeIp(int64(ipId))
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "解封成功")
}

//获取所有被封禁的ip
func (ic *IpController) getAllIp(ctx *gin.Context) {
	ips, err := ic.is.GetAllIp()
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, ips)
}

//封禁ip
func (ic *IpController) freezeIp(ctx *gin.Context) {
	var ip model.BannedIp
	err := util.Decode(ctx.Request.Body, &ip)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	ic.is.BanIp(ip.Ip)
	util.Success(ctx, "封禁成功")
}
