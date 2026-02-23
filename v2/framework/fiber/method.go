package fiber

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ghinknet/payutils/v2/model"
	"github.com/ghinknet/payutils/v2/payment/alipay"
	"github.com/ghinknet/payutils/v2/payment/wechat"
	"github.com/ghinknet/payutils/v2/utils"
	"github.com/go-pay/gopay"
	"github.com/gofiber/fiber/v3"
)

// Status returns the status of the order
func (c *Client) Status(orderID string) (model.TradeState, model.TradeMethod, time.Time, error) {
	// Try to check in WeChat Pay
	if c.Payment.WeChat != nil {
		// Check WeChat-Pay
		wxRsp, err := c.Payment.WeChat.V3TransactionQueryOrder(
			context.Background(), 2, fmt.Sprintf(
				"%s%s%s", c.Config.Basic.Prefix, orderID, c.Config.Basic.Suffix,
			),
		)
		if err != nil {
			return model.TradeStateUnknown, model.TradeMethodUnknown, time.Time{}, err
		}

		// Check return status
		if wxRsp.Code != 0 && wxRsp.Code != 404 {
			return model.TradeStateUnknown, model.TradeMethodUnknown,
				time.Time{}, wechat.ErrWeChatPayRespCodeInvalidBuilder(
					wxRsp.Code, wxRsp.ErrResponse.Code, wxRsp.ErrResponse.Message,
				)
		}

		// Success
		if wxRsp.Code == 0 {
			return wechat.MapState(wxRsp.Response.TradeState), model.TradeMethodWeChatPay,
				wechat.FormatTime(wxRsp.Response.SuccessTime), nil
		}
	}

	// Try to check in Alipay
	if c.Payment.Alipay != nil {
		// Prepare params
		bm := make(gopay.BodyMap)
		bm.Set("out_trade_no", fmt.Sprintf("%s%s%s", c.Config.Basic.Prefix, orderID, c.Config.Basic.Suffix)).
			Set("query_options", []string{"send_pay_date"})

		// Check Alipay
		aliRsp, err := c.Payment.Alipay.TradeQuery(context.Background(), bm)
		if err != nil {
			return model.TradeStateUnknown, model.TradeMethodUnknown, time.Time{}, err
		}

		// Check return status
		if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
			return model.TradeStateUnknown, model.TradeMethodUnknown,
				time.Time{}, alipay.ErrAlipayRespCodeInvalidBuilder(
					aliRsp.StatusCode, aliRsp.ErrResponse.Code, aliRsp.ErrResponse.Message,
				)
		}

		// Success
		if aliRsp.StatusCode == 200 {
			return alipay.MapState(aliRsp.TradeStatus), model.TradeMethodAlipay, alipay.FormatTime(aliRsp.SendPayDate), nil
		}
	}

	return model.TradeStateUnknown, model.TradeMethodUnknown, time.Time{}, nil
}

// Close an order
func (c *Client) Close(orderID string) error {
	// Try to close order in WeChat Pay
	if c.Payment.WeChat != nil {
		// Close order in WeChat-Pay
		wxRsp, err := c.Payment.WeChat.V3TransactionCloseOrder(
			context.Background(), fmt.Sprintf("%s%s%s", c.Config.Basic.Prefix, orderID, c.Config.Basic.Suffix),
		)
		if err != nil {
			return err
		}

		// Check return status
		if wxRsp.Code != 0 && wxRsp.Code != 404 {
			return wechat.ErrWeChatPayRespCodeInvalidBuilder(
				wxRsp.Code, wxRsp.ErrResponse.Code, wxRsp.ErrResponse.Message,
			)
		}
	}

	// Try to close order in Alipay
	if c.Payment.Alipay != nil {
		// Prepare params
		bm := make(gopay.BodyMap)
		bm.Set("out_trade_no", fmt.Sprintf("%s%s%s", c.Config.Basic.Prefix, orderID, c.Config.Basic.Suffix))

		// Close order in Alipay
		aliRsp, err := c.Payment.Alipay.TradeClose(context.Background(), bm)
		if err != nil {
			return err
		}

		// Check return status
		if aliRsp.StatusCode != http.StatusOK && aliRsp.ErrResponse.Code != "ACQ.TRADE_NOT_EXIST" {
			return alipay.ErrAlipayRespCodeInvalidBuilder(
				aliRsp.StatusCode, aliRsp.ErrResponse.Code, aliRsp.ErrResponse.Message,
			)
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
			bm.Set("out_trade_no",
				fmt.Sprintf("%s%s%s", c.Config.Basic.Prefix, orderID, c.Config.Basic.Suffix),
			)
			bm.Set("refund_amount", utils.CentsToYuan(refundAmount))
			bm.Set("out_request_no",
				fmt.Sprintf("%s%s%s", c.Config.Basic.Prefix, refundID, c.Config.Basic.Suffix),
			)

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
				return alipay.ErrAlipayRespCodeInvalidBuilder(
					aliRsp.StatusCode, aliRsp.ErrResponse.Code, aliRsp.ErrResponse.Message,
				)
			}

			return nil
		}
		return model.ErrPaymentMethodDisabled
	case model.TradeMethodWeChatPay:
		if c.Payment.WeChat != nil {
			// Prepare params
			bm := make(gopay.BodyMap)
			bm.Set("out_trade_no",
				fmt.Sprintf("%s%s%s", c.Config.Basic.Prefix, orderID, c.Config.Basic.Suffix),
			).
				Set("out_refund_no",
					fmt.Sprintf("%s%s%s", c.Config.Basic.Prefix, refundID, c.Config.Basic.Suffix),
				).
				Set("reason", reason).
				Set("notify_url", fmt.Sprintf(
					"%s%s/wechat/callback", c.Config.Basic.Endpoint, c.Config.Fiber.(*fiber.Group).Prefix,
				)).
				SetBodyMap("amount", func(bm gopay.BodyMap) {
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
				return wechat.ErrWeChatPayRespCodeInvalidBuilder(
					wxRsp.Code, wxRsp.ErrResponse.Code, wxRsp.ErrResponse.Message,
				)
			}

			return nil
		}
		return model.ErrPaymentMethodDisabled
	default:
		return model.ErrUnsupportedMethod
	}
}

// internalStatusUpdater preprocess orderID then call external updater
func (c *Client) internalStatusUpdater(
	ctx fiber.Ctx, orderID string, status model.TradeState, method model.TradeMethod, tm time.Time,
) error {
	// Process orderID
	orderID = strings.TrimPrefix(orderID, c.Config.Basic.Prefix)
	orderID = strings.TrimSuffix(orderID, c.Config.Basic.Suffix)

	return c.Config.StatusUpdater(ctx, orderID, status, method, tm)
}
