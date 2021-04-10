package middleware

import (
	"db-security-backend/service"
	"db-security-backend/util"
	"github.com/gin-gonic/gin"
)

func AdminCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		j := service.NewJwt()
		claims, _ := j.ParseToken(token)
		var userService service.UserService
		if userService.GetUser(claims.Phone).IsSuperUser == int8(1) {
			ctx.Next()
		} else {
			util.Failed(ctx, "权限不足")
			ctx.Abort()
			return
		}
	}
}
