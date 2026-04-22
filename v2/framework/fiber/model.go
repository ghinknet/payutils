package fiber

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.gh.ink/payutils/v2/model"
	"go.gh.ink/payutils/v2/payment/alipay"
	"go.gh.ink/payutils/v2/payment/wechat"
)

// Config provides specific config options for framework
type Config struct {
	// Public options
	Basic model.Config
	// Framework options
	Fiber          fiber.Router
	ErrorHandler   func(c fiber.Ctx, err error) error
	DetailProvider func(c fiber.Ctx, orderID string, method model.TradeMethod) (model.OrderDetail, error)
	StatusUpdater  func(c fiber.Ctx, orderID string, status model.TradeState, method model.TradeMethod, tm time.Time) error
	// Method
	Unmarshal func(data []byte, v any) error
	Marshal   func(v any) ([]byte, error)
	// Payment options
	Alipay    *alipay.Config
	WeChatPay *wechat.Config
}

// Client provides specific config options in client for framework
type Client struct {
	Config  *Config
	Payment *model.Client
}
