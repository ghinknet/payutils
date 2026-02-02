package fiber

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ghinknet/json"
	"github.com/ghinknet/payutils/v2/model"
	internalWeChat "github.com/ghinknet/payutils/v2/payment/wechat"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

type WeChatPay struct {
	Client *Client
}

func (w *WeChatPay) Create(c fiber.Ctx) error {
	// Read request params
	var req model.OrderRequest
	if err := c.Bind().JSON(&req); err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Get order
	orderInfo, err := w.Client.Config.DetailProvider(
		c,
		req.OrderID,
		model.TradeMethodWeChatPay,
	)
	if err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Check time
	if orderInfo.Expiry-time.Now().Unix() < 60 {
		return w.Client.Config.ErrorHandler(c, model.ErrNoEnoughTimeToPay)
	}

	// Prepare params
	expire := time.Unix(orderInfo.Expiry, 0).Add(-5 * time.Second).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	bm.Set("appid", w.Client.Config.WeChatPay.AppID).
		Set("mchid", w.Client.Config.WeChatPay.MerchantID).
		Set("description", orderInfo.Subject).
		Set("out_trade_no", fmt.Sprintf("%s%s%s", w.Client.Config.Basic.Prefix, req.OrderID, w.Client.Config.Basic.Suffix)).
		Set("time_expire", expire).
		Set("notify_url", fmt.Sprintf(
			"%s%s/wechat/callback", w.Client.Config.Basic.Endpoint, w.Client.Config.Fiber.(*fiber.Group).Prefix,
		)).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", orderInfo.Price).
				Set("currency", orderInfo.Currency)
		})

	switch req.Platform {
	case model.PlatformPC:
		// Create a native transaction
		wxRsp, err := w.Client.Payment.WeChat.V3TransactionNative(context.Background(), bm)
		if err != nil {
			return w.Client.Config.ErrorHandler(c, err)
		}
		if wxRsp.Code != 0 {
			return w.Client.Config.ErrorHandler(c, internalWeChat.ErrWeChatPayRespCodeInvalid)
		}

		return RespSuccess(c, map[string]string{
			"payUrl": wxRsp.Response.CodeUrl,
		})
	case model.PlatformMobile:
		fallthrough
	case model.PlatformWeChat:
		if req.OpenID == "" {
			return w.Client.Config.ErrorHandler(c, internalWeChat.ErrWeChatOpenIDIsRequired)
		}
		bm.SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", req.OpenID)
		})

		// Create a jsapi transaction
		wxRsp, err := w.Client.Payment.WeChat.V3TransactionJsapi(context.Background(), bm)
		if err != nil {
			return w.Client.Config.ErrorHandler(c, err)
		}
		if wxRsp.Code != 0 {
			return w.Client.Config.ErrorHandler(c, internalWeChat.ErrWeChatPayRespCodeInvalid)
		}

		// Get jsapi sign
		jsapi, err := w.Client.Payment.WeChat.PaySignOfJSAPI(
			w.Client.Config.WeChatPay.AppID,
			wxRsp.Response.PrepayId,
		)
		if err != nil {
			return w.Client.Config.ErrorHandler(c, err)
		}

		return RespSuccess(c, jsapi)
	}
	return nil
}

func (w *WeChatPay) Callback(c fiber.Ctx) error {
	// Convert request
	httpReq, err := adaptor.ConvertRequest(c, false)
	if err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Parse notify params
	notifyReq, err := wechat.V3ParseNotify(httpReq)
	if err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Get public key
	certMap := w.Client.Payment.WeChat.WxPublicKeyMap()
	// Verify sign
	if err = notifyReq.VerifySignByPKMap(certMap); err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Decrypt message from WeChat Pay
	wechatPayCallback := new(internalWeChat.WeChatPayCallback)
	if err = notifyReq.DecryptCipherTextToStruct(
		w.Client.Config.WeChatPay.MerchantAPIv3Key, wechatPayCallback); err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Return status
	if err = w.Client.internalStatusUpdater(
		c,
		wechatPayCallback.OutTradeNo,
		internalWeChat.MapState(wechatPayCallback.TradeState),
		model.TradeMethodWeChatPay,
		internalWeChat.FormatTime(wechatPayCallback.SuccessTime),
	); err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	return c.Status(http.StatusOK).JSON(&wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
}

func (w *WeChatPay) OpenIDCallback(c fiber.Ctx) error {
	// Read request params
	var req internalWeChat.OpenIDCallbackRequest
	if err := c.Bind().JSON(&req); err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Request URI
	URL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		w.Client.Config.WeChatPay.AppID,
		w.Client.Config.WeChatPay.AppSecret,
		req.Code,
	)

	// Send Request
	resp, err := http.Get(URL)
	if err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Parse JSON Data
	var result internalWeChat.AccessTokenResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	return RespSuccess(c, map[string]string{
		"openID": result.OpenID,
	})
}

func (w *WeChatPay) AuthorizeLinkGen(c fiber.Ctx) error {
	// Read request params
	var req internalWeChat.AuthorizeLinkRequest
	if err := c.Bind().JSON(&req); err != nil {
		return w.Client.Config.ErrorHandler(c, err)
	}

	// Check allowed origins
	allowed := false
	for _, origin := range w.Client.Config.Basic.AllowOrigins {
		if strings.HasPrefix(req.RedirectURI, origin) {
			allowed = true
			break
		}
	}
	if !allowed {
		return w.Client.Config.ErrorHandler(c, internalWeChat.ErrWeChatRedirectURIMismatch)
	}

	// Encode redirect_uri
	req.RedirectURI = url.QueryEscape(req.RedirectURI)

	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",
		w.Client.Config.WeChatPay.AppID,
		req.RedirectURI,
		req.State,
	)

	// Return authorize link
	return RespSuccess(c, map[string]string{
		"url": authURL,
	})
}
