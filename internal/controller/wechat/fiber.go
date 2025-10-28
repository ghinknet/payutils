package wechat

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.ghink.net/ghink/payutils/v2/internal/client"
	"git.ghink.net/ghink/payutils/v2/internal/model"
	"github.com/bytedance/sonic"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
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
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	bm.Set("appid", f.Config.WeChatPay.AppID).
		Set("mchid", f.Config.WeChatPay.MerchantID).
		Set("description", orderInfo.Subject).
		Set("out_trade_no", req.OrderID).
		Set("time_expire", expire).
		Set("notify_url", fmt.Sprintf(
			"%s%s/wechat/callback", f.Config.Endpoint, f.Config.Fiber.(*fiber.Group).Prefix,
		)).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", orderInfo.Price).
				Set("currency", "CNY")
		})

	switch req.Platform {
	case model.PlatformPC:
		// Create a native transaction
		wxRsp, err := f.Client.WeChat.V3TransactionNative(context.Background(), bm)
		if err != nil {
			return f.Config.ErrorHandler(c, err)
		}
		if wxRsp.Code != 0 {
			return f.Config.ErrorHandler(c, model.ErrWeChatPayRespCodeInvalid)
		}

		return model.FiberRespSuccess(c, map[string]string{
			"payUrl": wxRsp.Response.CodeUrl,
		})
	case model.PlatformMobile:
		fallthrough
	case model.PlatformWeChat:
		if req.OpenID == "" {
			return f.Config.ErrorHandler(c, model.ErrOpenIDIsRequired)
		}
		bm.SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", req.OpenID)
		})

		// Create a jsapi transaction
		wxRsp, err := f.Client.WeChat.V3TransactionJsapi(context.Background(), bm)
		if err != nil {
			return f.Config.ErrorHandler(c, err)
		}
		if wxRsp.Code != 0 {
			return f.Config.ErrorHandler(c, model.ErrWeChatPayRespCodeInvalid)
		}

		// Get jsapi sign
		jsapi, err := f.Client.WeChat.PaySignOfJSAPI(
			f.Config.WeChatPay.AppID,
			wxRsp.Response.PrepayId,
		)
		if err != nil {
			return f.Config.ErrorHandler(c, err)
		}

		return model.FiberRespSuccess(c, jsapi)
	}
	return nil
}

func (f *FiberController) Callback(c fiber.Ctx) error {
	// Convert request
	httpReq, err := adaptor.ConvertRequest(c, false)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Parse notify params
	notifyReq, err := wechat.V3ParseNotify(httpReq)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Get public key
	certMap := f.Client.WeChat.WxPublicKeyMap()
	// Verify sign
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Decrypt message from WeChat Pay
	wechatPayCallback := &model.WeChatPayCallback{}
	err = notifyReq.DecryptCipherTextToStruct(
		f.Config.WeChatPay.MerchantAPIv3Key, wechatPayCallback)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	var status model.TradeStatus
	switch wechatPayCallback.TradeState {
	case model.WeChatTradeStateSuccess:
		status = model.TradeSuccess
	case model.WeChatTradeStateClosed:
		status = model.TradeClosed
	case model.WeChatTradeStateNotPay:
		status = model.TradePending
	case model.WeChatTradeStateRefund:
		status = model.TradeClosed
	}
	// Return status
	err = f.Config.OrderStatus(
		wechatPayCallback.OutTradeNo,
		status,
	)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	return c.Status(http.StatusOK).JSON(&wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
}

func (f *FiberController) OpenIDCallback(c fiber.Ctx) error {
	// Read request params
	var req model.OpenIDCallbackRequest
	if err := c.Bind().JSON(&req); err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Request URI
	URL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		f.Config.WeChatPay.AppID,
		f.Config.WeChatPay.AppSecret,
		req.Code,
	)

	// Send Request
	resp, err := http.Get(URL)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Parse JSON Data
	var result model.AccessTokenResponse
	err = sonic.Unmarshal(body, &result)
	if err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	return model.FiberRespSuccess(c, map[string]string{
		"openID": result.OpenID,
	})
}

func (f *FiberController) AuthorizeLinkGen(c fiber.Ctx) error {
	// Read request params
	var req model.AuthorizeLinkRequest
	if err := c.Bind().JSON(&req); err != nil {
		return f.Config.ErrorHandler(c, err)
	}

	// Check same-site origin(?)
	if !strings.HasPrefix(req.RedirectURI, f.Config.Endpoint) {
		return f.Config.ErrorHandler(c, model.ErrWeChatRedirectURIMismatch)
	}

	// Encode redirect_uri
	req.RedirectURI = url.QueryEscape(req.RedirectURI)

	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",
		f.Config.WeChatPay.AppID,
		req.RedirectURI,
		req.State,
	)

	// Return authorize link
	return model.FiberRespSuccess(c, map[string]string{
		"url": authURL,
	})
}
