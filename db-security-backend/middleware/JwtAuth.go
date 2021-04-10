package middleware

import (
	"db-security-backend/util"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"

	"db-security-backend/service"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			util.Failed(ctx, "未携带token,非法访问")
			ctx.Abort()
			return
		}
		j := service.NewJwt()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == service.TokenExpired {
				util.Failed(ctx, "授权过期")
				ctx.Abort()
				return
			}
			util.Failed(ctx, err)
			ctx.Abort()
			return
		}
		conn := util.NewRedisPool().Get()
		res, err := redis.String(conn.Do("GET", "expireJwt:"+strconv.Itoa(int(claims.Id)), token))
		if res == "1" {
			util.Failed(ctx, "授权过期")
			ctx.Abort()
			return
		}
		if claims.ExpiresAt-time.Now().Unix() < 450 {
			var userService service.UserService
			newToken := service.GenerateToken(userService.GetUser(claims.Phone))
			ctx.Set("refreshToken", newToken)
		}
		ctx.Set("claims", claims)
	}
}
