package route

import (
	"git.ghink.net/ghink/payutils/v2/internal/client"
	"git.ghink.net/ghink/payutils/v2/internal/controller/alipay"
	"git.ghink.net/ghink/payutils/v2/internal/controller/wechat"
	"git.ghink.net/ghink/payutils/v2/internal/model"
	"github.com/gofiber/fiber/v3"
)

// FiberRegister registers gin router
func FiberRegister(r fiber.Router, client *client.Client, config model.Config) {
	{
		alipayRoute := r.Group("/alipay")
		alipayGinController := alipay.FiberController{Client: client, Config: config}
		alipayRoute.Post("/create", alipayGinController.Create)
		alipayRoute.Post("/callback", alipayGinController.Callback)
	}
	{
		wechatRoute := r.Group("/wechat")
		wechatGinController := wechat.FiberController{Client: client, Config: config}
		wechatRoute.Post("/create", wechatGinController.Create)
		wechatRoute.Post("/callback", wechatGinController.Callback)
		wechatRoute.Post("/openIDCallback", wechatGinController.OpenIDCallback)
		wechatRoute.Post("/authorizeLink", wechatGinController.AuthorizeLinkGen)
	}
}
