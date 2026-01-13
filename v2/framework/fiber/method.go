package fiber

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ghinknet/payutils/v2/model"
	"github.com/ghinknet/payutils/v2/payment/alipay"
	"github.com/ghinknet/payutils/v2/payment/wechat"
	"github.com/ghinknet/payutils/v2/utils"
	"github.com/go-pay/gopay"
	"github.com/gofiber/fiber/v3"
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

// Close an order
func (c *Client) Close(orderID string) error {
	// Try to close order in WeChat Pay
	if c.Payment.WeChat != nil {
		// Close order in WeChat-Pay
		wxRsp, err := c.Payment.WeChat.V3TransactionCloseOrder(context.Background(), orderID)
		if err != nil {
			return err
		}

		// Check return status
		if wxRsp.Code != 0 && wxRsp.Code != 404 {
			return wechat.ErrWeChatPayRespCodeInvalid
		}
	}

	// Try to close order in Alipay
	if c.Payment.Alipay != nil {
		// Prepare params
		bm := make(gopay.BodyMap)
		bm.Set("out_trade_no", orderID)

		// Close order in Alipay
		aliRsp, err := c.Payment.Alipay.TradeClose(context.Background(), bm)
		if err != nil {
			return err
		}

		// Check return status
		if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
			return alipay.ErrAlipayRespCodeInvalid
		}
	}

	return nil
}

// Refund an order
func (c *Client) Refund(
	orderID string, method model.TradeMethod,
	currency string, totalAmount int64,
	refundID string, refundAmount int64, reason string,
) error {
	switch method {
	case model.TradeMethodAlipay:
		if c.Payment.Alipay != nil {
			// Check currency
			if currency != "CNY" {
				return model.ErrUnsupportedCurrency
			}

			// Prepare params
			bm := make(gopay.BodyMap)
			bm.Set("out_trade_no", orderID)
			bm.Set("refund_amount", utils.CentsToYuan(refundAmount))
			bm.Set("out_request_no", refundID)

			// Set reason
			if reason != "" {
				bm.Set("refund_reason", reason)
			}

			// Refund order in Alipay
			aliRsp, err := c.Payment.Alipay.TradeRefund(context.Background(), bm)
			if err != nil {
				return err
			}

			// Check return status
			if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
				return alipay.ErrAlipayRespCodeInvalid
			}

			return nil
		}
		return model.ErrPaymentMethodDisabled
	case model.TradeMethodWeChatPay:
		if c.Payment.WeChat != nil {
			// Prepare params
			bm := make(gopay.BodyMap)
			bm.Set("out_trade_no", orderID).
				Set("out_refund_no", refundID).
				Set("reason", reason).
				Set("notify_url", fmt.Sprintf(
					"%s%s/wechat/callback", c.Config.Basic.Endpoint, c.Config.Fiber.(*fiber.Group).Prefix,
				)).
				SetBodyMap("refundAmount", func(bm gopay.BodyMap) {
					bm.Set("total", totalAmount).
						Set("refund", refundAmount).
						Set("currency", currency)
				})

			// Refund order in WeChat-Pay
			wxRsp, err := c.Payment.WeChat.V3Refund(context.Background(), bm)
			if err != nil {
				return err
			}

			// Check return status
			if wxRsp.Code != 0 && wxRsp.Code != 404 {
				return wechat.ErrWeChatPayRespCodeInvalid
			}

			return nil
		}
		return model.ErrPaymentMethodDisabled
	default:
		return model.ErrUnsupportedMethod
	}
}
