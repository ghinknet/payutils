package route

import (
	"git.ghink.net/ghink/payutils/internal/client"
	"git.ghink.net/ghink/payutils/internal/controller/alipay"
	"git.ghink.net/ghink/payutils/internal/controller/wechat"
	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/gofiber/fiber/v3"
)

// FiberRegister registers gin router
func FiberRegister(r fiber.Router, client *client.Client, config model.Config) {
	{
		alipayRoute := r.Group("/alipay")
		alipayFiberController := alipay.FiberController{Client: client, Config: config}
		alipayRoute.Post("/create", alipayFiberController.Create)
		alipayRoute.Post("/callback", alipayFiberController.Callback)
	}
	{
		wechatRoute := r.Group("/wechat")
		wechatFiberController := wechat.FiberController{Client: client, Config: config}
		wechatRoute.Post("/create", wechatFiberController.Create)
		wechatRoute.Post("/callback", wechatFiberController.Callback)
		wechatRoute.Post("/openIDCallback", wechatFiberController.OpenIDCallback)
		wechatRoute.Post("/authorizeLink", wechatFiberController.AuthorizeLinkGen)
	}
}
