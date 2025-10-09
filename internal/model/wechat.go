package model

// WeChatTradeState provides enum type for WeChat trade status const
type WeChatTradeState string

const (
	WeChatTradeStateSuccess WeChatTradeState = "SUCCESS" // Pay success
	WeChatTradeStateRefund  WeChatTradeState = "REFUND"  // Pay refund
	WeChatTradeStateNotPay  WeChatTradeState = "NOTPAY"  // Unpay
	WeChatTradeStateClosed  WeChatTradeState = "CLOSED"  // Closed
)

// WeChatTradeType provides enum type for WeChat trade type const
type WeChatTradeType string

const (
	WeChatTradeTypeApp    WeChatTradeType = "APP"    // APP pay
	WeChatTradeTypeJSAPI  WeChatTradeType = "JSAPI"  // JSAPI pay
	WeChatTradeTypeNative WeChatTradeType = "NATIVE" // Native pay
	WeChatTradeTypeH5     WeChatTradeType = "MWEB"   // H5 pay
)

type WeChatPayCallback struct {
	TransactionID   string            `json:"transaction_id"`
	Amount          AmountInfo        `json:"amount"`
	MchID           string            `json:"mchid"`
	TradeState      WeChatTradeState  `json:"trade_state"`
	BankType        string            `json:"bank_type"`
	PromotionDetail []PromotionDetail `json:"promotion_detail,omitempty"`
	SuccessTime     string            `json:"success_time"`
	Payer           PayerInfo         `json:"payer"`
	OutTradeNo      string            `json:"out_trade_no"`
	AppID           string            `json:"appid"`
	TradeStateDesc  string            `json:"trade_state_desc"`
	TradeType       WeChatTradeType   `json:"trade_type"`
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
	WeChatPayContribute int           `json:"wechatpay_contribute"`
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

// AccessTokenResponse provides WeChat access token response struct
type AccessTokenResponse struct {
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	OpenID         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotUser int    `json:"is_snapshotuser"`
	UnionID        string `json:"unionid"`
}

// OpenIDCallbackRequest provides basic WeChat OpenID callback requests params bind
type OpenIDCallbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

// AuthorizeLinkRequest provides basic WeChat OpenID authorize link request params bind
type AuthorizeLinkRequest struct {
	RedirectURI string `json:"redirect_uri" binding:"required"`
	State       string `json:"state" binding:"required"`
}
