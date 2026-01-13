package fiber

import (
	"github.com/gofiber/fiber/v3"
)

// register the framework route
func register(r fiber.Router, c *Client) {
	{
		alipayRoute := r.Group("/alipay")
		alipayFiberController := Alipay{Client: c}
		alipayRoute.Post("/create", alipayFiberController.Create)
		alipayRoute.Post("/callback", alipayFiberController.Callback)
	}
	{
		wechatRoute := r.Group("/wechat")
		wechatFiberController := WeChatPay{Client: c}
		wechatRoute.Post("/create", wechatFiberController.Create)
		wechatRoute.Post("/callback", wechatFiberController.Callback)
		wechatRoute.Post("/openIDCallback", wechatFiberController.OpenIDCallback)
		wechatRoute.Post("/authorizeLink", wechatFiberController.AuthorizeLinkGen)
	}
}
