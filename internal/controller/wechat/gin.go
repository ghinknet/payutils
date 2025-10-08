package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"io"
	"net/http"
	"time"
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
		// Create native transaction
		wxRsp, err := g.Client.Wechat.V3TransactionNative(context.Background(), bm)
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

	}
}

func (g *GinController) Callback(c *gin.Context) {
	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		return
	}

	//TODO:这玩意到底能不能跑？
	// 获取微信平台证书
	certMap := g.Client.Wechat.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return
	}

	// 微信消息解密
	wechatPayCallback := &model.WechatPayCallback{}
	err = notifyReq.DecryptCipherTextToStruct(
		g.Config.WeChatPay.MerchantAPIv3Key, *wechatPayCallback)
	if err != nil {
		return
	}

	//TODO:这四个状态是不是要修改model
	var status model.TradeStatus
	switch wechatPayCallback.TradeState {
	case model.WechatTradeStateSuccess:
		status = model.TradeSuccess
	case model.WechatTradeStateClosed:
		status = model.TradeClosed
	case model.WechatTradeStateNotPay:
		status = model.TradePending
	case model.WechatTradeStateRefund:
		status = model.TradeClosed
	}
	// Return status
	err = g.Config.OrderStatus(
		wechatPayCallback.OutTradeNo,
		c.Request.Header.Get("Authorization"),
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
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		g.Config.WeChatPay.AppID,
		g.Config.WeChatPay.AppSecret,
		req.Code,
	)

	// Send Request
	resp, err := http.Get(url)
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

func (g *GinController) BasicInfo(c *gin.Context) {

}
