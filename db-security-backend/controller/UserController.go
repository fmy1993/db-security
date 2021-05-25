package controller

import (
	"db-security-backend/config"
	"db-security-backend/param/userParam"
	"log"
	"strconv"

	"db-security-backend/middleware"
	"db-security-backend/model"
	"db-security-backend/service"
	"db-security-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type UserController struct {
	us  service.UserService
	bis service.BannedIpService
	cs  service.CsvService
	wms service.WaterMarkService
	cfs service.ConfigService
}

func (uc *UserController) Router(engine *gin.Engine) {
	engine.GET("/captcha", uc.captcha)
	engine.GET("/salt", uc.salt)
	engine.POST("/register", uc.register)
	engine.POST("/login", uc.login)
	engine.POST("/logout", uc.logout)
	engine.PUT("/revise", middleware.JWTAuth(), uc.revise)
	engine.GET("/user/info", middleware.JWTAuth(), uc.getUserInfo)
	engine.POST("/freeze/:userId", middleware.JWTAuth(), middleware.AdminCheck(), uc.freezeUser)
	engine.POST("/free/:userId", middleware.JWTAuth(), middleware.AdminCheck(), uc.freeUser)
	engine.GET("/all_users", middleware.JWTAuth(), middleware.AdminCheck(), uc.getAllUsers)
}

//生成验证码
func (uc *UserController) captcha(ctx *gin.Context) {
	util.Success(ctx, util.GenerateCaptcha(ctx))
}

//获取自身信息
func (uc *UserController) getUserInfo(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*service.CustomClaims)
	var user = uc.us.GetUser(claims.Phone)
	var res = make(map[string]interface{})
	res["id"] = user.Id
	res["phone"] = user.Phone
	res["is_super_user"] = user.IsSuperUser
	util.Success(ctx, res)
}

//生成盐值
func (uc *UserController) salt(ctx *gin.Context) {
	util.Success(ctx, uc.us.Salt())
}

//用户注册
func (uc *UserController) register(ctx *gin.Context) {
	// ip检测部分
	xForwardFor := ctx.Request.Header.Values("X-Forwarded-For")
	var ip string
	if len(xForwardFor) != 0 {
		ip = xForwardFor[0]
	} else {
		ip = ctx.Request.Header.Values("Remote_Addr")[0]
	}
	if uc.bis.IpExist(ip) {
		util.Failed(ctx, "此ip已经被封禁,请联系管理员解封")
		ctx.Abort()
		return
	}
	// 注册部分
	var registerParam userParam.RegisterParam
	err := util.Decode(ctx.Request.Body, &registerParam)
	if err != nil {
		util.Failed(ctx, "注册失败,请检查注册字段")
		ctx.Abort()
		return
	}
	if registerParam.CheckCode != util.Md5(config.GetConfig().Register.CheckCode) {
		conn := util.NewRedisPool().Get()
		res, _ := redis.String(conn.Do("hget", "bannedIps", ip))
		re, _ := strconv.Atoi(res)
		if re > 2 {
			uc.bis.BanIp(ip)
			util.Failed(ctx, "你的ip已经被封禁,请联系管理员解封")
			ctx.Abort()
			return
		}
		_, err = conn.Do("hincrby", "bannedIps", ip, 1)
		if err != nil {
			util.Failed(ctx, err)
			ctx.Abort()
			return
		}
		util.Failed(ctx, "校验码错误")
		ctx.Abort()
		return
	}
	// 校验密码部分
	password, err := util.Base64Decode(registerParam.Password)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	registerParam.Password = string(util.RSADecrypt(password, "config/ssl.pem"))[6:]
	_, exist := uc.us.Register(registerParam)
	if exist == 0 {
		util.Success(ctx, "注册成功")
		return
	} else {
		util.Failed(ctx, "已经存在该用户")
		ctx.Abort()
		return
	}
}

//用户登录
func (uc *UserController) login(ctx *gin.Context) {
	var loginParam userParam.LoginParam
	err := util.Decode(ctx.Request.Body, &loginParam)
	if err != nil {
		util.Failed(ctx, "登录失败, 请检查登录字段")
		ctx.Abort()
		return
	}
	// -------------------------------------------------------------------
	//// 对比验证码
	//if !util.VerifyCaptcha(loginParam.CaptchaId, loginParam.CaptchaValue) {
	//	util.Failed(ctx, "验证码有误")
	//	ctx.Abort()
	//	return
	//}
	// -------------------------------------------------------------------
	password, err := util.Base64Decode(loginParam.Password)
	if err != nil {
		util.Failed(ctx, err.Error())
		ctx.Abort()
		return
	}
	// 对比密码
	loginParam.Password = string(util.RSADecrypt(password, "config/ssl.pem"))[6:]
	var user *model.User
	user = uc.us.GetUser(loginParam.Phone)
	truePwd := user.Password
	if truePwd == "" {
		util.Failed(ctx, "不存在该用户")
		ctx.Abort()
		return
	}
	if truePwd[7:] != util.EncoderSha256(truePwd[:6]+loginParam.Password) {
		conn := util.NewRedisPool().Get()
		res, _ := redis.String(conn.Do("hget", "user_"+strconv.FormatInt(user.Id, 10), "FCount"))
		re, _ := strconv.Atoi(res)
		if re > 3 {
			util.Failed(ctx, "此账号已经被冻结")
			ctx.Abort()
			return
		}
		_, err = conn.Do("hincrby", "user_"+strconv.FormatInt(user.Id, 10), "FCount", 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		util.Failed(ctx, "密码错误")
		ctx.Abort()
		return
	}
	if user.IsSuperUser == 1 {
		util.Success(ctx, map[string]interface{}{
			"token": service.GenerateToken(user),
		})
		ctx.Abort()
		return
	}
	// 身份校验成功
	uc.wms.BlendFingerPrintToPic(user.FingerPrint, user.Phone)
	atoi, err := strconv.Atoi(uc.cfs.GetConfigValueByKey("arnold_key"))
	if err != nil {
		return
	}
	uc.wms.Arnold(user.Phone, atoi)
	err = uc.wms.EmbedWatermarkToData(user.Phone)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	length, err := uc.cs.GenerateCsv(user.Phone)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(
		ctx, map[string]interface{}{
			"token": service.GenerateToken(user),
			"total": length,
		},
	)
}

//修改密码
func (uc *UserController) revise(ctx *gin.Context) {
	var reviseParam userParam.ReviseParam
	err := util.Decode(ctx.Request.Body, &reviseParam)
	if err != nil {
		util.Failed(ctx, "注册失败,请检查注册字段")
		ctx.Abort()
		return
	}
	oldPwd, err := util.Base64Decode(reviseParam.OldPassword)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	reviseParam.OldPassword = string(util.RSADecrypt(oldPwd, "config/ssl.pem"))
	newPwd, err := util.Base64Decode(reviseParam.NewPassword)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	reviseParam.NewPassword = string(util.RSADecrypt(newPwd, "config/ssl.pem"))
	var user *model.User
	claims := ctx.MustGet("claims").(*service.CustomClaims)
	user = uc.us.GetUser(claims.Phone)
	if user.Password == "" {
		util.Failed(ctx, "不存在该用户")
		ctx.Abort()
		return
	}
	if user.Password[7:] != util.EncoderSha256(user.Password[:6]+reviseParam.OldPassword) {
		util.Failed(ctx, "密码错误")
		ctx.Abort()
		return
	}
	salt := uc.us.Salt()
	user.Password = salt + ":" + util.EncoderSha256(salt+reviseParam.NewPassword)
	err = uc.us.RevisePassword(claims.Id, user)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "修改成功")
}

//退出登录
func (uc *UserController) logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	j := service.NewJwt()
	claims, err := j.ParseToken(token)
	if err != nil {
		util.Success(ctx, "授权过期")
		ctx.Abort()
		return
	}
	if token == "" {
		util.Failed(ctx, "没有登录")
		ctx.Abort()
		return
	}
	conn := util.NewRedisPool().Get()
	_, err = conn.Do("SET", "expireJwt:"+strconv.Itoa(int(claims.Id)), token)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	_, err = conn.Do("EXPIRE", "expireJwt:"+strconv.Itoa(int(claims.Id)), 3600)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "退出成功")
}

//封禁用户
func (uc *UserController) freezeUser(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))
	conn := util.NewRedisPool().Get()
	_, err := conn.Do("hset", "user_"+strconv.FormatInt(int64(userId), 10), "FCount", 99)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "封禁成功")
}

//解冻用户
func (uc *UserController) freeUser(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))
	conn := util.NewRedisPool().Get()
	_, err := conn.Do("hset", "user_"+strconv.FormatInt(int64(userId), 10), "FCount", 0)
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, "解冻成功")
}

//获取所有用户数据
func (uc *UserController) getAllUsers(ctx *gin.Context) {
	users, err := uc.us.GetUsers()
	if err != nil {
		util.Failed(ctx, err)
		ctx.Abort()
		return
	}
	util.Success(ctx, users)
}
