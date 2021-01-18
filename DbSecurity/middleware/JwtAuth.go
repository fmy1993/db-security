package middleware

import (
	"time"

	"DbSecurity/service"
	"DbSecurity/tool"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			tool.Failed(ctx, "未携带token,非法访问")
			ctx.Abort()
			return
		}
		conn := tool.NewRedisPool().Get()
		res, err := redis.String(conn.Do("SISMEMBER", "expireJwt", token))
		if res == "1" {
			tool.Failed(ctx, "授权过期")
			ctx.Abort()
			return
		}
		j := tool.NewJwt()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == tool.TokenExpired {
				tool.Failed(ctx, "授权过期")
				ctx.Abort()
				return
			}
			tool.Failed(ctx, err)
			ctx.Abort()
			return
		}
		if claims.ExpiresAt-time.Now().Unix() < 450{
			var userService service.UserService
			newToken := tool.GenerateToken(userService.GetUser(claims.Phone))
			ctx.Set("refreshToken", newToken)
		}
		ctx.Set("claims", claims)
	}
}
