package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Resp(ctx *gin.Context, code int, data any, msg string) {
	type response struct {
		Code int    `json:"code"`
		Data any    `json:"data"`
		Msg  string `json:"msg"`
	}

	ctx.JSON(http.StatusOK, response{code, data, msg})
}

func RespSuccess(c *gin.Context, data any) {
	Resp(c, http.StatusOK, data, "success")
}
