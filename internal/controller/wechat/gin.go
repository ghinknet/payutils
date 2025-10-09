package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

type GinController struct {
	Client *model.Client
	Config model.Config
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
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	bm.Set("appid", g.Config.WeChatPay.AppID).
		Set("mchid", g.Config.WeChatPay.MerchantID).
		Set("description", orderInfo.Subject).
		Set("out_trade_no", req.OrderID).
		Set("time_expire", expire).
		Set("notify_url", fmt.Sprintf(
			"%s%s/wechat/callback", g.Config.Endpoint, g.Config.Gin.BasePath(),
		)).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", orderInfo.Price).
				Set("currency", "CNY")
		})

	switch req.Platform {
	case model.PlatformPC:
		// Create a native transaction
		wxRsp, err := g.Client.WeChat.V3TransactionNative(context.Background(), bm)
		if err != nil {
			g.Config.ErrorHandler(c, err)
			return
		}
		if wxRsp.Code != 0 {
			g.Config.ErrorHandler(c, model.ErrWeChatPayRespCodeInvalid)
			return
		}

		model.RespSuccess(c, map[string]string{
			"payUrl": wxRsp.Response.CodeUrl,
		})
	case model.PlatformMobile:
		fallthrough
	case model.PlatformWeChat:
		if req.OrderID == "" {
			g.Config.ErrorHandler(c, model.ErrOrderIDIsRequired)
			return
		}
		bm.SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", req.OrderID)
		})

		// Create a jsapi transaction
		wxRsp, err := g.Client.WeChat.V3TransactionJsapi(context.Background(), bm)
		if err != nil {
			g.Config.ErrorHandler(c, err)
			return
		}
		if wxRsp.Code != 0 {
			g.Config.ErrorHandler(c, model.ErrWeChatPayRespCodeInvalid)
			return
		}

		// Get jsapi sign
		jsapi, err := g.Client.WeChat.PaySignOfJSAPI(
			g.Config.WeChatPay.AppID,
			wxRsp.Response.PrepayId,
		)

		model.RespSuccess(c, jsapi)
	}
}

func (g *GinController) Callback(c *gin.Context) {
	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		return
	}

	//TODO:这玩意到底能不能跑？
	// 获取微信平台证书
	certMap := g.Client.WeChat.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return
	}

	// 微信消息解密
	wechatPayCallback := &model.WeChatPayCallback{}
	err = notifyReq.DecryptCipherTextToStruct(
		g.Config.WeChatPay.MerchantAPIv3Key, *wechatPayCallback)
	if err != nil {
		return
	}

	//TODO:这四个状态是不是要修改model
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
	err = g.Config.OrderStatus(
		wechatPayCallback.OutTradeNo,
		status,
	)
	if err != nil {
		return
	}

	// ====↓↓↓====异步通知应答====↓↓↓====
	// 退款通知http应答码为200且返回状态码为SUCCESS才会当做商户接收成功，否则会重试。
	// 注意：重试过多会导致微信支付端积压过多通知而堵塞，影响其他正常通知。

	// 此写法是 gin 框架返回微信的写法
	c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
}

func (g *GinController) OpenIDCallback(c *gin.Context) {
	// Read request params
	var req model.OpenIDCallbackRequest
	err := c.ShouldBind(&req)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	// Request URI
	URL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		g.Config.WeChatPay.AppID,
		g.Config.WeChatPay.AppSecret,
		req.Code,
	)

	// Send Request
	resp, err := http.Get(URL)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	// Parse JSON Data
	var result model.AccessTokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	model.RespSuccess(c, map[string]string{
		"openID": result.OpenID,
	})
}

func (g *GinController) AuthorizeLinkGen(c *gin.Context) {
	// Read request params
	var req model.AuthorizeLinkRequest
	err := c.ShouldBind(&req)
	if err != nil {
		g.Config.ErrorHandler(c, err)
		return
	}

	// Check same-site origin(?)
	if !strings.HasPrefix(req.RedirectURI, g.Config.Endpoint) {
		g.Config.ErrorHandler(c, model.ErrWeChatRedirectURIMismatch)
	}

	// Encode redirect_uri
	req.RedirectURI = url.QueryEscape(req.RedirectURI)

	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",
		g.Config.WeChatPay.AppID,
		req.RedirectURI,
		req.State,
	)

	// Return authorize link
	model.RespSuccess(c, map[string]string{
		"url": authURL,
	})
}
