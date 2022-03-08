package rsp

import (
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/isc"
	"net/http"
)

func Success(ctx *gin.Context, object any) {
	ctx.JSON(http.StatusOK, isc.ObjectToData(object))
}

func SuccessOfStandard(ctx *gin.Context, v any) {
	ctx.JSON(http.StatusOK, map[string]any{
		"code":    "success",
		"message": "成功",
		"data":    isc.ObjectToData(v),
	})
}

func FailedOfStandard(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, map[string]any{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

func FailedWithDataOfStandard(ctx *gin.Context, code string, message string, v any) {
	ctx.JSON(http.StatusOK, map[string]any{
		"code":    code,
		"message": message,
		"data":    isc.ObjectToData(v),
	})
}
