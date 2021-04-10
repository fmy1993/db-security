package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS int = 0
	FAILED  int = 1
)

func Success(ctx *gin.Context, v interface{}) {
	refresh, _ := ctx.Get("refreshToken")
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    SUCCESS,
		"msg":     "success",
		"payload": v,
		"token":   refresh,
	})
}

func Failed(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    FAILED,
		"msg":     "failed",
		"payload": v,
	})
}
