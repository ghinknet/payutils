package model

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
)

type response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func DefaultErrorHandler(c any, err error) error {
	switch {
	case reflect.TypeOf(c).String() != "*gin.Context":
		GinRespInternalServerError(c.(*gin.Context), err)
	case reflect.TypeOf(err).String() != "fiber.Ctx":
		return FiberRespInternalServerError(c.(fiber.Ctx), err)
	}
	return nil
}

func GinResp(ctx *gin.Context, code int, data any, msg string) {
	ctx.PureJSON(http.StatusOK, response{code, data, msg})
}

func GinRespSuccess(c *gin.Context, data any) {
	GinResp(c, http.StatusOK, data, "success")
}

func GinRespInternalServerError(c *gin.Context, err error) {
	GinResp(c, http.StatusInternalServerError, map[string]any{
		"error": err,
	}, "internal server error")
}

func FiberResp(c fiber.Ctx, code int, data any, msg string) error {
	return c.Status(http.StatusOK).JSON(response{code, data, msg})
}

func FiberRespSuccess(c fiber.Ctx, data any) error {
	return FiberResp(c, http.StatusOK, data, "success")
}

func FiberRespInternalServerError(c fiber.Ctx, err error) error {
	return FiberResp(c, http.StatusInternalServerError, nil, "internal server error")
}
