package fiber

import (
	"context"
	"net/http"

	"github.com/ghinknet/payutils/v2/model"
	"github.com/ghinknet/payutils/v2/payment/alipay"
	"github.com/ghinknet/payutils/v2/payment/wechat"
	"github.com/go-pay/gopay"
)

// Status returns the status of the order
func (c *Client) Status(orderID string) (model.TradeState, model.TradeMethod, error) {
	// Try to check in WeChat Pay
	if c.Payment.WeChat != nil {
		// Check WeChat-Pay
		wxRsp, err := c.Payment.WeChat.V3TransactionQueryOrder(context.Background(), 2, orderID)
		if err != nil {
			return model.TradeStateUnknown, model.TradeMethodUnknown, err
		}

		// Check return status
		if wxRsp.Code != 0 && wxRsp.Code != 404 {
			return model.TradeStateUnknown, model.TradeMethodUnknown, wechat.ErrWeChatPayRespCodeInvalid
		}

		return wechat.MapState(wxRsp.Response.TradeState), model.TradeMethodWeChatPay, nil
	}

	// Try to check in Alipay
	if c.Payment.Alipay != nil {
		// Prepare params
		bm := make(gopay.BodyMap)
		bm.Set("out_trade_no", orderID)

		// Check Alipay
		aliRsp, err := c.Payment.Alipay.TradeQuery(context.Background(), bm)
		if err != nil {
			return model.TradeStateUnknown, model.TradeMethodUnknown, err
		}

		// Check return status
		if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
			return model.TradeStateUnknown, model.TradeMethodUnknown, alipay.ErrAlipayRespCodeInvalid
		}

		return alipay.MapState(aliRsp.TradeStatus), model.TradeMethodAlipay, nil
	}

	return model.TradeStateUnknown, model.TradeMethodUnknown, nil
}
