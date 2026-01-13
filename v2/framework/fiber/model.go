package fiber

import (
	"github.com/ghinknet/payutils/v2/model"
	"github.com/ghinknet/payutils/v2/payment/alipay"
	"github.com/ghinknet/payutils/v2/payment/wechat"
	"github.com/gofiber/fiber/v3"
)

// Config provides specific config options for framework
type Config struct {
	// Public options
	model.Config
	// Framework options
	Fiber          fiber.Router
	ErrorHandler   func(c fiber.Ctx, err error) error
	DetailProvider func(c fiber.Ctx, orderID string, method model.TradeMethod) (model.OrderDetail, error)
	StatusUpdater  func(c fiber.Ctx, orderID string, status model.TradeState) error
	// Payment options
	Alipay    *alipay.Config
	WeChatPay *wechat.Config
}

// Client provides specific config options in client for framework
type Client struct {
	Config  *Config
	Payment *model.Client
}
