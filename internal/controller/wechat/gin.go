package wechat

import (
	"net/http"

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

}

func (g *GinController) BasicInfo(c *gin.Context) {

}
