package client

import (
	"context"
	"net/http"

	"git.ghink.net/ghink/payutils/v2/internal/model"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/gopay/wechat/v3"
)

type Client struct {
	Alipay *alipay.ClientV3
	WeChat *wechat.ClientV3
}

// CheckStatus checks whether an order has been paid
func (c *Client) CheckStatus(orderID string) (bool, error) {
	if c.WeChat != nil {
		// Check WeChat-Pay
		wxRsp, err := c.WeChat.V3TransactionQueryOrder(context.Background(), 2, orderID)
		if err != nil {
			return false, err
		}
		if wxRsp.Code != 0 && wxRsp.Code != 404 {
			return false, model.ErrWeChatPayRespCodeInvalid
		}
		if wxRsp.Response.TradeState == "SUCCESS" {
			return true, nil
		}
	}

	if c.Alipay != nil {
		// Prepare params
		bm := make(gopay.BodyMap)
		bm.Set("out_trade_no", orderID)

		// Check Alipay
		aliRsp, err := c.Alipay.TradeQuery(context.Background(), bm)
		if err != nil {
			return false, err
		}
		if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
			return false, model.ErrAlipayRespCodeInvalid
		}
		if aliRsp.TradeStatus == "TRADE_SUCCESS" {
			return true, nil
		}
	}

	return false, nil
}
