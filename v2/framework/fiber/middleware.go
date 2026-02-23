package fiber

import (
	"github.com/ghinknet/payutils/v2/model"
	"github.com/gofiber/fiber/v3"
)

// CheckPaymentActivate provides a middleware to check a payment method activate status
func CheckPaymentActivate(paymentMethod model.TradeMethod, client *Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		if client == nil {
			return RespInternalServerError(c, model.ErrPaymentMethodDisabled)
		}
		switch paymentMethod {
		case model.TradeMethodAlipay:
			if client.Config.Alipay == nil {
				return client.Config.ErrorHandler(c, model.ErrPaymentMethodDisabled)
			}
		case model.TradeMethodWeChatPay:
			if client.Config.WeChatPay == nil {
				return client.Config.ErrorHandler(c, model.ErrPaymentMethodDisabled)
			}
		default:
			return client.Config.ErrorHandler(c, model.ErrPaymentMethodDisabled)
		}
		return c.Next()
	}
}
