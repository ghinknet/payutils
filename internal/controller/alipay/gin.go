package alipay

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

type GinController struct {
	Client *model.Client
	Config model.Config
}

// centsToYuan transfer cents to yuan
func centsToYuan(cents int64) string {
	yuan := cents / 100
	remainder := cents % 100

	if cents < 0 {
		yuan = -yuan
		remainder = -remainder
		return fmt.Sprintf("-%d.%02d", yuan, remainder)
	}

	return fmt.Sprintf("%d.%02d", yuan, remainder)
}

func (g *GinController) Create(c *gin.Context) {
	// Read request params
	var req model.OrderRequest
	err := c.ShouldBind(&req)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	// Get order
	orderInfo, err := g.Config.OrderInfo(
		req.OrderID,
		c.Request.Header.Get("Authorization"),
	)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	// Prepare params
	bm := make(gopay.BodyMap)
	bm.Set("subject", orderInfo.Subject).
		Set("out_trade_no", req.OrderID).
		Set("total_amount", centsToYuan(orderInfo.Price)).
		Set("notify_url", fmt.Sprintf(
			"%s%s/alipay/callback", g.Config.Endpoint, g.Config.Gin.BasePath(),
		))

	// Create order
	var url string
	switch req.Platform {
	case model.PlatformPC:
		url, err = g.Client.Alipay.TradePagePay(context.Background(), bm)
		if err != nil {
			g.Config.ErrorHandler(c, err)
			return
		}
	case model.PlatformWeChat:
		fallthrough
	case model.PlatformMobile:
		url, err = g.Client.Alipay.TradeWapPay(context.Background(), bm)
		if err != nil {
			g.Config.ErrorHandler(c, err)
			return
		}
	}
	model.RespSuccess(c, map[string]string{
		"payUrl": url,
	})
}

func (g *GinController) Callback(c *gin.Context) {
	// Parse notify params
	notifyReq, err := alipay.ParseNotifyToBodyMap(c.Request)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	// Verify sign by alipay public cert
	ok, err := alipay.VerifySignWithCert([]byte(g.Config.Alipay.PublicCert), notifyReq)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}
	if !ok {
		g.Config.ErrorHandler(c, errors.New("failed to verify"))
		return
	}

	// Parse data
	// Docs: https://opendocs.alipay.com/open/203/105286
	notifyRequest := &model.NotifyRequest{}
	err = notifyReq.Unmarshal(notifyRequest)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
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
	err = g.Config.OrderStatus(
		notifyRequest.OutTradeNo,
		status,
	)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	c.String(http.StatusOK, "%s", "success")
}
