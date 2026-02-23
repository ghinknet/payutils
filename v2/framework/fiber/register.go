package fiber

import (
	"github.com/ghinknet/payutils/v2/model"
	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
)

// register the framework route
func register(r fiber.Router, c *Client) {
	{
		alipayRoute := r.Group("/alipay")
		alipayRoute.Use(CheckPaymentActivate(model.TradeMethodAlipay, c))
		alipayRoute.Use(recoverer.New())
		alipayFiberController := Alipay{Client: c}
		alipayRoute.Post("/create", alipayFiberController.Create)
		alipayRoute.Post("/callback", alipayFiberController.Callback)
	}
	{
		wechatRoute := r.Group("/wechat")
		wechatRoute.Use(CheckPaymentActivate(model.TradeMethodWeChatPay, c))
		wechatRoute.Use(recoverer.New())
		wechatFiberController := WeChatPay{Client: c}
		wechatRoute.Post("/create", wechatFiberController.Create)
		wechatRoute.Post("/callback", wechatFiberController.Callback)
		wechatRoute.Post("/openIDCallback", wechatFiberController.OpenIDCallback)
		wechatRoute.Post("/authorizeLink", wechatFiberController.AuthorizeLinkGen)
	}
}
