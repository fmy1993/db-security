package controller

import (
	"db-security-backend/config"
	"db-security-backend/middleware"
	"db-security-backend/model"
	"db-security-backend/param"
	"db-security-backend/param/staffParam"
	"db-security-backend/service"
	"db-security-backend/util"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"io"
	"strconv"
)

type StaffController struct {
	dps service.DifferentialPrivacyService
	wms service.WaterMarkService
	us  service.UserService
	scs service.StaffCopyService
	ss  service.StaffService
}

func (sc *StaffController) Router(engine *gin.Engine) {
	engine.GET("/idcard/most", sc.getMostIdCard)
	engine.GET("/download/staff", middleware.JWTAuth(), sc.downloadCsv)
	engine.POST("/user_staff", middleware.JWTAuth(), sc.getStaffCopysData)
	engine.POST("/staff", middleware.JWTAuth(), sc.getStaffCopysData)
	engine.POST("/analysis", middleware.JWTAuth(), sc.analysis)
	engine.POST("/track", middleware.JWTAuth(), middleware.AdminCheck(), sc.track)
	engine.POST("/ori_staff", middleware.JWTAuth(), middleware.AdminCheck(), sc.getStaffsData)
	engine.POST("/add_staff", middleware.JWTAuth(), middleware.AdminCheck(), sc.addStaff)
	engine.DELETE("/staff/:staffId", middleware.JWTAuth(), middleware.AdminCheck(), sc.deleteStaff)
	engine.PUT("/staff/:staffId", middleware.JWTAuth(), middleware.AdminCheck(), sc.updateStaff)
	engine.POST("/dp", middleware.JWTAuth(), middleware.AdminCheck(), sc.differentialPrivacy)
}

// 获取通过指数机制扰动过的身份证号地区
func (sc *StaffController) getMostIdCard(ctx *gin.Context) {
	util.Success(ctx, sc.dps.ExpMechanism(decimal.NewFromInt(1)))
}

// 对原始数据通过laplace机制添加噪音
func (sc *StaffController) differentialPrivacy(ctx *gin.Context) {
	sc.dps.DifferentialPrivacy()
	// 这里缺少预嵌入部分
	util.Success(ctx, "")
}

// 盗版追踪
func (sc *StaffController) track(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	fs, _ := file.Open()
	c := csv.NewReader(fs)
	var staffCopy = make(map[int64]decimal.Decimal)
	for {
		row, err := c.Read()
		if err != nil && err != io.EOF {
			util.Failed(ctx, "文件有误")
			ctx.Abort()
			return
		}
		if err == io.EOF {
			break
		}
		salary, _ := decimal.NewFromString(row[4])
		tmp, _ := strconv.Atoi(row[0])
		staffCopy[int64(tmp)] = salary
	}
	sc.wms.PickWatermarkByCsv(staffCopy)
	fingerPrint := sc.wms.PickFingerPrintByPic()
	var user, distance = sc.us.GetUsersByFingerprint(fingerPrint)
	var res = make(map[string]interface{})
	res["fingerprint"] = fingerPrint
	res["distance"] = distance
	res["user"] = user.Phone
	res["pic"] = sc.wms.PicDecode()
	util.Success(ctx, res)
}

// csv文件下载
func (sc *StaffController) downloadCsv(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*service.CustomClaims)
	var cfg = config.GetConfig()
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "staff.csv"))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(cfg.Path.Form + "staff" + claims.Phone + ".csv")
}

// 查询克隆表数据
func (sc *StaffController) getStaffCopysData(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*service.CustomClaims)
	var selectStaffParam staffParam.SelectStaffParam
	err := util.Decode(ctx.Request.Body, &selectStaffParam)
	if err != nil {
		util.Failed(ctx, err.Error())
		ctx.Abort()
		return
	}
	staffData, total, err := sc.scs.GetPageStaffCopyByPhone(claims.Phone, selectStaffParam)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	var res = make(map[string]interface{})
	res["staff"] = staffData
	res["total"] = total
	util.Success(ctx, res)
}

// 查询原始表数据
func (sc *StaffController) getStaffsData(ctx *gin.Context) {
	var selectStaffParam staffParam.SelectStaffParam
	err := util.Decode(ctx.Request.Body, &selectStaffParam)
	if err != nil {
		util.Failed(ctx, err.Error())
		ctx.Abort()
		return
	}
	staffData, total, err := sc.ss.GetPageStaff(selectStaffParam)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	var res = make(map[string]interface{})
	res["staffData"] = staffData
	res["total"] = total
	util.Success(ctx, res)
}

// 增加原始表数据
func (sc *StaffController) addStaff(ctx *gin.Context) {
	var staff model.Staff
	err := util.Decode(ctx.Request.Body, &staff)
	if err != nil {
		util.Failed(ctx, "添加失败")
		ctx.Abort()
		return
	}
	err = sc.ss.InsertStaff(staff)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "添加成功")
}

//删除原始表数据
func (sc *StaffController) deleteStaff(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("staffId"))
	err := sc.ss.DeleteStaff(int64(id))
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "删除成功")
}

// 修改原始表数据
func (sc *StaffController) updateStaff(ctx *gin.Context) {
	var staff model.Staff
	err := util.Decode(ctx.Request.Body, &staff)
	if err != nil {
		util.Failed(ctx, "修改失败")
		ctx.Abort()
		return
	}
	err = sc.ss.UpdateStaff(staff)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "修改成功")
}

//数据可视化支撑
func (sc *StaffController) analysis(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*service.CustomClaims)
	var analysisParam param.AnalysisParam
	err := util.Decode(ctx.Request.Body, &analysisParam)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	data, err := sc.scs.GetAllStaffCopyDataByPhone(claims.Phone)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, sc.scs.AnalysisStaff(data, analysisParam))
}
