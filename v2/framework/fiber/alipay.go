package fiber

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ghinknet/payutils/v2/model"
	internalAlipay "github.com/ghinknet/payutils/v2/payment/alipay"
	"github.com/ghinknet/payutils/v2/utils"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

type Alipay struct {
	Client *Client
}

func (a *Alipay) Create(c fiber.Ctx) error {
	// Read request params
	var req model.OrderRequest
	if err := c.Bind().JSON(&req); err != nil {
		return a.Client.Config.ErrorHandler(c, err)
	}

	// Get order
	orderInfo, err := a.Client.Config.DetailProvider(
		c,
		req.OrderID,
		model.TradeMethodAlipay,
	)
	if err != nil {
		return a.Client.Config.ErrorHandler(c, err)
	}

	// Check currency
	if orderInfo.Currency != "CNY" {
		return a.Client.Config.ErrorHandler(c, model.ErrUnsupportedCurrency)
	}

	// Check time
	if orderInfo.Expiry-time.Now().Unix() < 60 {
		return a.Client.Config.ErrorHandler(c, model.ErrNoEnoughTimeToPay)
	}

	// Prepare params
	expire := time.Unix(orderInfo.Expiry, 0).Add(-5 * time.Second).Format("2006-01-02 15:04:05")
	// Prepare params
	bm := make(gopay.BodyMap)
	bm.Set("subject", orderInfo.Subject).
		Set("time_expire", expire).
		Set("out_trade_no", fmt.Sprintf("%s%s%s", a.Client.Config.Basic.Prefix, req.OrderID, a.Client.Config.Basic.Suffix)).
		Set("total_amount", utils.CentsToYuan(orderInfo.Price)).
		Set("notify_url", fmt.Sprintf(
			"%s%s/alipay/callback", a.Client.Config.Basic.Endpoint, a.Client.Config.Fiber.(*fiber.Group).Prefix,
		))

	// Create order
	var url string
	switch req.Platform {
	case model.PlatformPC:
		url, err = a.Client.Payment.Alipay.TradePagePay(context.Background(), bm)
		if err != nil {
			return a.Client.Config.ErrorHandler(c, err)
		}
	case model.PlatformWeChat:
		fallthrough
	case model.PlatformMobile:
		url, err = a.Client.Payment.Alipay.TradeWapPay(context.Background(), bm)
		if err != nil {
			return a.Client.Config.ErrorHandler(c, err)
		}
	}
	return RespSuccess(c, map[string]string{
		"payUrl": url,
	})
}

func (a *Alipay) Callback(c fiber.Ctx) error {
	// Convert request
	httpReq, err := adaptor.ConvertRequest(c, false)
	if err != nil {
		return a.Client.Config.ErrorHandler(c, err)
	}

	// Parse notify params
	notifyReq, err := alipay.ParseNotifyToBodyMap(httpReq)
	if err != nil {
		return a.Client.Config.ErrorHandler(c, err)
	}

	// Verify sign by alipay public cert
	ok, err := alipay.VerifySignWithCert([]byte(a.Client.Config.Alipay.PublicCert), notifyReq)
	if err != nil {
		return a.Client.Config.ErrorHandler(c, err)
	}
	if !ok {
		return a.Client.Config.ErrorHandler(c, internalAlipay.ErrAlipayNotifyVerifyFailed)
	}

	// Parse data
	// Docs: https://opendocs.alipay.com/open/203/105286
	notifyRequest := &alipay.NotifyRequest{}
	err = notifyReq.Unmarshal(notifyRequest)
	if err != nil {
		return a.Client.Config.ErrorHandler(c, err)
	}

	// Return status
	err = a.Client.internalStatusUpdater(
		c,
		notifyRequest.OutTradeNo,
		internalAlipay.MapState(notifyRequest.TradeStatus),
	)
	if err != nil {
		return a.Client.Config.ErrorHandler(c, err)
	}

	return c.Status(http.StatusOK).SendString("success")
}
