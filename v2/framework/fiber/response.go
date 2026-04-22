package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.gh.ink/payutils/v2/model"
)

func Resp(c fiber.Ctx, code int, data any, msg string) error {
	return c.Status(http.StatusOK).JSON(model.Response{Code: code, Data: data, Msg: msg})
}

func RespSuccess(c fiber.Ctx, data any) error {
	return Resp(c, http.StatusOK, data, "success")
}

func RespInternalServerError(c fiber.Ctx, err error) error {
	return Resp(c, http.StatusInternalServerError, nil, "internal server error")
}
