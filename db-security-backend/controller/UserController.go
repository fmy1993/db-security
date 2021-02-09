package controller

import (
	"log"
	"strconv"

	"db-security-backend/middleware"
	"db-security-backend/model"
	"db-security-backend/param"
	"db-security-backend/service"
	"db-security-backend/tool"
	"db-security-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type UserController struct {
	us service.UserService
	bis service.BannedIpService
	pic util.Picture
	csvOperation util.CsvOperation
}

func (uc *UserController) Router(engine *gin.Engine) {
	userGroup := engine.Group("/api/user")
	{
		userGroup.GET("/captcha", uc.captcha)
		userGroup.GET("/salt", uc.salt)
		userGroup.POST("/register", uc.register)
		userGroup.POST("/login", uc.login)
		userGroup.POST("/logout", uc.logout)
		userGroup.POST("/revise", middleware.JWTAuth(), uc.revise)
		userGroup.POST("/admin/freeze", middleware.JWTAuth(), middleware.AdminCheck(), uc.freezeUser)
		userGroup.POST("/admin/free", middleware.JWTAuth(), middleware.AdminCheck(), uc.freeUser)
		userGroup.POST("/admin/users", middleware.JWTAuth(), middleware.AdminCheck(), uc.getAllUsers)
	}
}

//生成验证码
func (uc *UserController) captcha(ctx *gin.Context) {
	tool.Success(ctx, tool.GenerateCaptcha(ctx))
}

//生成盐值
func (uc *UserController) salt(ctx *gin.Context) {
	tool.Success(ctx, uc.us.Salt())
}

//用户注册
func (uc *UserController) register(ctx *gin.Context) {
	xForwardFor := ctx.Request.Header.Values("X-Forwarded-For")
	var ip string
	if len(xForwardFor) != 0 {
		ip = xForwardFor[0]
	} else {
		ip = ctx.Request.Header.Values("Remote_Addr")[0]
	}
	if uc.bis.IpExist(ip) {
		tool.Failed(ctx, "此ip已经被封禁,请联系管理员解封")
		ctx.Abort()
		return
	}
	var registerParam param.RegisterParam
	err := tool.Decode(ctx.Request.Body, &registerParam)
	if err != nil {
		tool.Failed(ctx, "注册失败,请检查注册字段")
		ctx.Abort()
		return
	}
	if registerParam.CheckCode != tool.Md5(tool.GetConfig().Register.CheckCode) {
		conn := tool.NewRedisPool().Get()
		res, _ := redis.String(conn.Do("hget", "bannedIps", ip))
		re, _ := strconv.Atoi(res)
		if re > 2 {
			uc.bis.BanIp(ip)
			tool.Failed(ctx, "你的ip已经被封禁,请联系管理员解封")
			ctx.Abort()
			return
		}
		_, err = conn.Do("hincrby", "bannedIps", ip, 1)
		if err != nil {
			tool.Failed(ctx, err)
			ctx.Abort()
			return
		}
		tool.Failed(ctx, "校验码错误")
		ctx.Abort()
		return
	}
	password, err := tool.Base64Decode(registerParam.Password)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	registerParam.Password = string(tool.RSADecrypt(password, "config/ssl.pem"))[6:]
	_, exist := uc.us.Register(registerParam)
	if exist == 0 {
		tool.Success(ctx, "注册成功")
		return
	} else {
		tool.Failed(ctx, "已经存在该用户")
		ctx.Abort()
		return
	}
}

//用户登录
func (uc *UserController) login(ctx *gin.Context) {
	var loginParam param.LoginParam
	err := tool.Decode(ctx.Request.Body, &loginParam)
	if err != nil {
		tool.Failed(ctx, "登录失败, 请检查登录字段")
		ctx.Abort()
		return
	}
	if !tool.VerifyCaptcha(loginParam.CaptchaId, loginParam.CaptchaValue) {
		tool.Failed(ctx, "验证码有误")
		ctx.Abort()
		return
	}
	password, err := tool.Base64Decode(loginParam.Password)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	loginParam.Password = string(tool.RSADecrypt(password, "config/ssl.pem"))[6:]
	var user *model.User
	user = uc.us.GetUser(loginParam.Phone)
	truePwd := user.Password
	if truePwd == "" {
		tool.Failed(ctx, "不存在该用户")
		ctx.Abort()
		return
	}
	if truePwd[7:] != tool.EncoderSha256(truePwd[:6]+loginParam.Password) {
		conn := tool.NewRedisPool().Get()
		res, _ := redis.String(conn.Do("hget", "user_"+strconv.FormatInt(user.Id, 10), "FCount"))
		re, _ := strconv.Atoi(res)
		if re > 3 {
			tool.Failed(ctx, "此账号已经被冻结")
			ctx.Abort()
			return
		}
		_, err = conn.Do("hincrby", "user_"+strconv.FormatInt(user.Id, 10), "FCount", 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		tool.Failed(ctx, "密码错误")
		ctx.Abort()
		return
	}
	uc.pic.BlendFP(user.FingerPrint, user.Phone)
	uc.pic.Arnold(user.Phone, 3)
	_, err = util.EmbedNewPic(user.Phone)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	length, err := uc.csvOperation.GenerateCsv(user.Phone)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, map[string]interface{}{
		"token":     tool.GenerateToken(user),
		"total":     length,
		"admin":     user.IsSuperUser,
		"csrfToken": "csrfToken",
	})
}

//修改密码
func (uc *UserController) revise(ctx *gin.Context) {
	var reviseParam param.ReviseParam
	err := tool.Decode(ctx.Request.Body, &reviseParam)
	if err != nil {
		tool.Failed(ctx, "注册失败,请检查注册字段")
		ctx.Abort()
		return
	}
	oldPwd, err := tool.Base64Decode(reviseParam.OldPassword)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	reviseParam.OldPassword = string(tool.RSADecrypt(oldPwd, "config/ssl.pem"))
	newPwd, err := tool.Base64Decode(reviseParam.NewPassword)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	reviseParam.NewPassword = string(tool.RSADecrypt(newPwd, "config/ssl.pem"))
	var user *model.User
	claims := ctx.MustGet("claims").(*tool.CustomClaims)
	user = uc.us.GetUser(claims.Phone)
	if user.Password == "" {
		tool.Failed(ctx, "不存在该用户")
		ctx.Abort()
		return
	}
	if user.Password[7:] != tool.EncoderSha256(user.Password[:6]+reviseParam.OldPassword) {
		tool.Failed(ctx, "密码错误")
		ctx.Abort()
		return
	}
	salt := uc.us.Salt()
	user.Password = salt + ":" + tool.EncoderSha256(salt+reviseParam.NewPassword)
	err = uc.us.RevisePassword(claims.Id, user)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "修改成功")
}

//退出登录
func (uc *UserController) logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		tool.Failed(ctx, "没有登录")
		ctx.Abort()
		return
	}
	conn := tool.NewRedisPool().Get()
	_, err := conn.Do("SADD", "expireJwt", token)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	_, err = conn.Do("EXPIRE", "expireJwt", token, 3600)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "退出成功")
}

//封禁用户
func (uc *UserController) freezeUser(ctx *gin.Context) {
	var user model.User
	err := tool.Decode(ctx.Request.Body, &user)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	conn := tool.NewRedisPool().Get()
	_, err = conn.Do("hset", "user_"+strconv.FormatInt(user.Id, 10), "FCount", 99)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "封禁成功")
}

//解冻用户
func (uc *UserController) freeUser(ctx *gin.Context) {
	var user model.User
	err := tool.Decode(ctx.Request.Body, &user)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	conn := tool.NewRedisPool().Get()
	_, err = conn.Do("hset", "user_"+strconv.FormatInt(user.Id, 10), "FCount", 0)
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, "解冻成功")
}

//获取所有用户数据
func (uc *UserController) getAllUsers(ctx *gin.Context) {
	users, err := uc.us.GetUsers()
	if err != nil {
		tool.Failed(ctx, err)
		ctx.Abort()
		return
	}
	tool.Success(ctx, users)
}
