package middleware

import (
	"DbSecurity/service"
	"DbSecurity/tool"
	"github.com/gin-gonic/gin"
)

func AdminCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		j := tool.NewJwt()
		claims, _ := j.ParseToken(token)
		var userService service.UserService
		if userService.GetUser(claims.Phone).IsSuperUser == int8(1) {
			ctx.Next()
		} else {
			tool.Failed(ctx, "权限不足")
			ctx.Abort()
			return
		}
	}
}
