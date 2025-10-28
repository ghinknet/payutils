package alipay

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"git.ghink.net/ghink/payutils/internal/client"
	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

type FiberController struct {
	Client *client.Client
	Config model.Config
}

func (f *FiberController) Create(c fiber.Ctx) error {
	// Read request params
	var req model.OrderRequest
	if err := c.Bind().JSON(&req); err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Get order
	orderInfo, err := f.Config.OrderInfo(
		req.OrderID,
		c.Get("Authorization"),
	)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Prepare params
	bm := make(gopay.BodyMap)
	bm.Set("subject", orderInfo.Subject).
		Set("out_trade_no", req.OrderID).
		Set("total_amount", centsToYuan(orderInfo.Price)).
		Set("notify_url", fmt.Sprintf(
			"%s%s/alipay/callback", f.Config.Endpoint, f.Config.Fiber.(*fiber.Group).Prefix,
		))

	// Create order
	var url string
	switch req.Platform {
	case model.PlatformPC:
		url, err = f.Client.Alipay.TradePagePay(context.Background(), bm)
		if err != nil {
			return f.Config.ErrorHandler(c, err)
		}
	case model.PlatformWeChat:
		fallthrough
	case model.PlatformMobile:
		url, err = f.Client.Alipay.TradeWapPay(context.Background(), bm)
		if err != nil {
			return f.Config.ErrorHandler(c, err)
		}
	}
	return model.FiberRespSuccess(c, map[string]string{
		"payUrl": url,
	})
}

func (f *FiberController) Callback(c fiber.Ctx) error {
	// Convert request
	httpReq, err := adaptor.ConvertRequest(c, false)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Parse notify params
	notifyReq, err := alipay.ParseNotifyToBodyMap(httpReq)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Verify sign by alipay public cert
	ok, err := alipay.VerifySignWithCert([]byte(f.Config.Alipay.PublicCert), notifyReq)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}
	if !ok {
		return f.Config.ErrorHandler(c, errors.New("failed to verify"))
	}

	// Parse data
	// Docs: https://opendocs.alipay.com/open/203/105286
	notifyRequest := &model.NotifyRequest{}
	err = notifyReq.Unmarshal(notifyRequest)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}
	var status model.TradeStatus
	switch notifyRequest.TradeStatus {
	case model.AlipayTradeClosed:
		status = model.TradeClosed
	case model.AlipayTradeSuccess:
		status = model.TradeSuccess
	case model.AlipayTradeFinished:
		status = model.TradeFinished
	}

	// Return status
	err = f.Config.OrderStatus(
		notifyRequest.OutTradeNo,
		status,
	)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	return c.Status(http.StatusOK).SendString("success")
}
