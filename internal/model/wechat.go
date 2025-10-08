package model

// WechatTradeState provides enum type for WeChat trade status const
type WechatTradeState string

const (
	WechatTradeStateSuccess  WechatTradeState = "SUCCESS"    // 支付成功
	WechatTradeStateRefund   WechatTradeState = "REFUND"     // 转入退款
	WechatTradeStateNotPay   WechatTradeState = "NOTPAY"     // 未支付
	WechatTradeStateClosed   WechatTradeState = "CLOSED"     // 已关闭
	WechatTradeStateRevoked  WechatTradeState = "REVOKED"    // 已撤销（仅付款码支付会返回）
	WechatTradeStatePaying   WechatTradeState = "USERPAYING" // 用户支付中（仅付款码支付会返回）
	WechatTradeStatePayError WechatTradeState = "PAYERROR"   // 支付失败（仅付款码支付会返回）
)

// WechatTradeType provides enum type for WeChat trade type const
type WechatTradeType string

const (
	WechatTradeTypeApp      WechatTradeType = "APP"      // APP支付
	WechatTradeTypeJSAPI    WechatTradeType = "JSAPI"    // JSAPI支付
	WechatTradeTypeNative   WechatTradeType = "NATIVE"   // Native支付
	WechatTradeTypeH5       WechatTradeType = "MWEB"     // H5支付
	WechatTradeTypeMicropay WechatTradeType = "MICROPAY" // 付款码支付
	WechatTradeTypeFacepay  WechatTradeType = "FACEPAY"  // 刷脸支付
)

type WechatPayCallback struct {
	TransactionID   string            `json:"transaction_id"`
	Amount          AmountInfo        `json:"amount"`
	MchID           string            `json:"mchid"`
	TradeState      WechatTradeState  `json:"trade_state"`
	BankType        string            `json:"bank_type"`
	PromotionDetail []PromotionDetail `json:"promotion_detail,omitempty"`
	SuccessTime     string            `json:"success_time"`
	Payer           PayerInfo         `json:"payer"`
	OutTradeNo      string            `json:"out_trade_no"`
	AppID           string            `json:"appid"`
	TradeStateDesc  string            `json:"trade_state_desc"`
	TradeType       WechatTradeType   `json:"trade_type"`
	Attach          string            `json:"attach,omitempty"`
	SceneInfo       SceneInfo         `json:"scene_info,omitempty"`
}

type AmountInfo struct {
	PayerTotal    int    `json:"payer_total"`
	Total         int    `json:"total"`
	Currency      string `json:"currency"`
	PayerCurrency string `json:"payer_currency"`
}

type PromotionDetail struct {
	Amount              int           `json:"amount"`
	WechatpayContribute int           `json:"wechatpay_contribute"`
	CouponID            string        `json:"coupon_id"`
	Scope               string        `json:"scope"`
	MerchantContribute  int           `json:"merchant_contribute"`
	Name                string        `json:"name"`
	OtherContribute     int           `json:"other_contribute"`
	Currency            string        `json:"currency"`
	StockID             string        `json:"stock_id"`
	GoodsDetail         []GoodsDetail `json:"goods_detail,omitempty"`
}

type GoodsDetail struct {
	GoodsRemark    string `json:"goods_remark"`
	Quantity       int    `json:"quantity"`
	DiscountAmount int    `json:"discount_amount"`
	GoodsID        string `json:"goods_id"`
	UnitPrice      int    `json:"unit_price"`
}

type PayerInfo struct {
	OpenID string `json:"openid"`
}

type SceneInfo struct {
	DeviceID string `json:"device_id"`
}
