package controller

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"DbSecurity/middleware"
	"DbSecurity/model"
	"DbSecurity/param"
	"DbSecurity/service"
	"DbSecurity/tool"
	"DbSecurity/util"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type PatientController struct {
}

func (pc *PatientController) Router(engine *gin.Engine) {
	patientGroup := engine.Group("/api/patient")
	patientGroup.Use(middleware.JWTAuth())
	{
		patientGroup.POST("/data", pc.userPatientCopyData)
		patientGroup.GET("/download", pc.download)
		patientGroup.GET("/search", pc.search)
		patientGroup.POST("/analysis", pc.analysis)
	}
	adminGroup := engine.Group("/api/patient/admin")
	adminGroup.Use(middleware.JWTAuth())
	adminGroup.Use(middleware.AdminCheck())
	{
		adminGroup.POST("/arnold", pc.generateSecretOldPic)
		adminGroup.POST("/dp", pc.differentialPrivacy)
		adminGroup.POST("/patient", pc.patientData)
		adminGroup.GET("/search", pc.searchOrigin)
		adminGroup.POST("/add", pc.addPatient)
		adminGroup.POST("/delete", pc.deletePatient)
		adminGroup.POST("/update", pc.updatePatient)
		adminGroup.POST("/track", pc.track)
	}
}

//差分隐私
func (pc *PatientController) differentialPrivacy(ctx *gin.Context) {
	var patientService service.PatientService
	success, err := patientService.CreateNewPatientCopyTable()
	if err != nil {
		tool.Failed(ctx, "failed")
		ctx.Abort()
		return
	}
	if !success {
		tool.Failed(ctx, "failed")
		ctx.Abort()
		return
	}
	patientsCopy, err := patientService.GetPatientCopyData()
	if err != nil {
		tool.Failed(ctx, "failed")
		ctx.Abort()
		return
	}
	var dp util.DifferentialPrivacy
	dp.DifferPrivacy(patientsCopy)
	_, err = util.EmbedOldPic()
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "success")
}

//生成变换后图片
func (pc *PatientController) generateSecretOldPic(ctx *gin.Context) {
	var pic util.Picture
	pic.Arnold("", 3)
}

//获取用户数据
func (pc *PatientController) userPatientCopyData(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*tool.CustomClaims)
	page, _ := strconv.Atoi(ctx.Query("page"))
	var patientService service.PatientService
	patientData, err := patientService.GetPageUserPatientCopyData(claims.Phone, 50, page)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, patientData)
}

//获取原始表数据
func (pc *PatientController) patientData(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	var patientService service.PatientService
	patientData, err := patientService.GetPagePatientData(50, page)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, patientData)
}

//增加原始表数据
func (pc *PatientController) addPatient(ctx *gin.Context) {
	var patient model.Patient
	err := tool.Decode(ctx.Request.Body, &patient)
	if err != nil {
		tool.Failed(ctx, "添加失败")
		ctx.Abort()
		return
	}
	var patientService service.PatientService
	err = patientService.InsertPatient(patient)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "添加成功")
}

//删除原始表数据
func (pc *PatientController) deletePatient(ctx *gin.Context) {
	var patient model.Patient
	err := tool.Decode(ctx.Request.Body, &patient)
	if err != nil {
		tool.Failed(ctx, "删除失败")
		ctx.Abort()
		return
	}
	var patientService service.PatientService
	err = patientService.DeletePatient(patient)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "删除成功")
}

//修改原始表数据
func (pc *PatientController) updatePatient(ctx *gin.Context) {
	var patient model.Patient
	err := tool.Decode(ctx.Request.Body, &patient)
	if err != nil {
		tool.Failed(ctx, "修改失败")
		ctx.Abort()
		return
	}
	var patientService service.PatientService
	err = patientService.UpdatePatient(patient)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "修改成功")
}

//文件下载
func (pc *PatientController) download(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*tool.CustomClaims)
	var cfg = tool.GetConfig()
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "病人表.csv"))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(cfg.Path.Form + "patient" + claims.Phone + ".csv")
}

//搜索用户表
func (pc *PatientController) search(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*tool.CustomClaims)
	page, _ := strconv.Atoi(ctx.Query("page"))
	pattern := ctx.Query("pattern")
	var patientService service.PatientService
	patientData, err := patientService.GetPageUserPatientCopyDataByPattern(claims.Phone, 50, page, pattern)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, patientData)
}

//搜索原始表
func (pc *PatientController) searchOrigin(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pattern := ctx.Query("pattern")
	var patientService service.PatientService
	patientData, err := patientService.GetPagePatientDataByPattern(50, page, pattern)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, patientData)
}

//数据可视化支撑
func (pc *PatientController) analysis(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*tool.CustomClaims)
	var analysisParam param.AnalysisParam
	err := tool.Decode(ctx.Request.Body, &analysisParam)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	var patientService service.PatientService
	data, err := patientService.GetUserPatientCopyData(claims.Phone)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	length := len(*data)
	var response = make(map[string]int)
	if analysisParam.Bill != 0 && analysisParam.Weight == 0 && analysisParam.High == 0 {
		var (
			billAbove = 0
			billBelow = 0
		)
		for i := 0; i < length; i++ {
			if (*data)[i].Bill.LessThan(decimal.NewFromInt(int64(analysisParam.Bill))) {
				billBelow += 1
			} else {
				billAbove += 1
			}
		}
		response["data_length"] = length
		response["bill_above"] = billAbove
		response["bill_below"] = billBelow
	} else if analysisParam.Weight != 0 && analysisParam.Bill == 0 && analysisParam.High == 0 {
		var (
			weightAbove = 0
			weightBelow = 0
		)
		for i := 0; i < length; i++ {
			if (*data)[i].Weight.LessThan(decimal.NewFromInt(int64(analysisParam.Weight))) {
				weightBelow += 1
			} else {
				weightAbove += 1
			}
		}
		response["data_length"] = length
		response["weight_above"] = weightAbove
		response["weight_below"] = weightBelow
	} else if analysisParam.Weight != 0 && analysisParam.High != 0 && analysisParam.Bill == 0 {
		var (
			weightAboveHighAbove = 0
			weightBelowHighAbove = 0
			weightAboveHighBelow = 0
			weightBelowHighBelow = 0
		)
		for i := 0; i < length; i++ {
			if decimal.NewFromInt(int64(analysisParam.Weight)).LessThan((*data)[i].Weight) && (*data)[i].
				High >= analysisParam.High {
				weightAboveHighAbove += 1
			} else if decimal.NewFromInt(int64(analysisParam.Weight)).LessThan((*data)[i].Weight) && (*data)[i].
				High < analysisParam.High {
				weightAboveHighBelow += 1
			} else if (*data)[i].Weight.LessThan(decimal.NewFromInt(int64(analysisParam.Weight))) && (*data)[i].
				High >= analysisParam.High {
				weightBelowHighAbove += 1
			} else {
				weightBelowHighBelow += 1
			}
		}
		response["data_length"] = length
		response["weight_above_high_above"] = weightAboveHighAbove
		response["weight_below_high_above"] = weightBelowHighAbove
		response["weight_above_high_below"] = weightAboveHighBelow
		response["weight_below_high_below"] = weightBelowHighBelow
	} else if analysisParam.Weight == 0 && analysisParam.Bill == 0 && analysisParam.High != 0 {
		var (
			highAbove = 0
			highBelow = 0
		)
		for i := 0; i < length; i++ {
			if (*data)[i].High > analysisParam.High {
				highAbove += 1
			} else {
				highBelow += 1
			}
		}
		response["data_length"] = length
		response["high_above"] = highAbove
		response["high_below"] = highBelow
	}
	tool.Success(ctx, response)
}

//盗版追踪
func (pc *PatientController) track(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	fs, _ := file.Open()
	c := csv.NewReader(fs)
	var patientCopy = make(map[string]struct {
		Weight decimal.Decimal
		Bill   decimal.Decimal
	})
	for {
		row, err := c.Read()
		if err != nil && err != io.EOF {
			tool.Failed(ctx, "文件有误")
			ctx.Abort()
			return
		}
		if err == io.EOF {
			break
		}
		weight, _ := decimal.NewFromString(row[2])
		bill, _ := decimal.NewFromString(row[6])
		patientCopy[row[0]] = struct {
			Weight decimal.Decimal
			Bill   decimal.Decimal
		}{Weight: weight, Bill: bill}
	}
	var pic util.Picture
	fp := pic.PickWatermark(patientCopy)
	fmt.Println(fp)
	var user service.UserService
	users, _ := user.GetUsers()
	var min = 20
	var phone = (*users)[0].Phone
	for i := 0; i < len(*users); i++ {
		var count = 0
		for j := 0; j < 20; j++ {
			if (*users)[i].FingerPrint[j] != fp[j] {
				count += 1
			}
		}
		if count < min {
			min = count
			phone = (*users)[i].Phone
		}
	}
	var p util.Picture
	var dis = make(map[string]interface{})
	dis["distance"] = min
	dis["phone"] = phone
	dis["fingerPrint"] = fp
	dis["picSrc"] = p.PicDecode()
	tool.Success(ctx, dis)
}
