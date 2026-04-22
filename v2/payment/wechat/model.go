package wechat

import (
	"time"

	"go.gh.ink/payutils/v2/model"
)

// Config provides needed config options for WeChat Pay
type Config struct {
	AppID                    string
	AppSecret                string
	MerchantID               string
	MerchantAPIv3Key         string
	MerchantCertSerialNumber string
	MerchantPrivateKey       string
	PublicKey                string
	PublicKeyID              string
}

const (
	// TradeStateNotPay Trade created and waiting for buyer to pay
	// -> TradeStatePending
	TradeStateNotPay string = "NOTPAY"
	// TradeStateSuccess Trade successes
	// -> TradeStateSuccess
	TradeStateSuccess string = "SUCCESS"
	// TradeStateRefund Trade refunded
	// -> TradeStateClosed
	TradeStateRefund string = "REFUND"
	// TradeStateClosed Trade closed due to time out
	// -> TradeStateClosed
	TradeStateClosed string = "CLOSED"
)

// TradeStateMap provides a map refer to internal trade state
var TradeStateMap = map[string]model.TradeState{
	TradeStateNotPay:  model.TradeStatePending,
	TradeStateSuccess: model.TradeStateSuccess,
	TradeStateRefund:  model.TradeStateClosed,
	TradeStateClosed:  model.TradeStateClosed,
}

// MapState provides a method to map trade state (string) to internal trade state
func MapState(state string) model.TradeState {
	// Try to find in map
	internalState, ok := TradeStateMap[state]
	if !ok {
		return model.TradeStateUnknown
	}
	return internalState
}

// FormatTime provides a method to map trade time (string) to internal time
func FormatTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}

	// Format time
	timeObj, err := time.ParseInLocation(time.RFC3339, timeStr, time.Local)
	if err != nil {
		return time.Time{}
	}

	return timeObj
}

// TradeType provides enum type for WeChat trade type const
type TradeType string

const (
	TradeTypeApp    TradeType = "APP"    // APP pay
	TradeTypeJSAPI  TradeType = "JSAPI"  // JSAPI pay
	TradeTypeNative TradeType = "NATIVE" // Native pay
	TradeTypeH5     TradeType = "MWEB"   // H5 pay
)

type WeChatPayCallback struct {
	TransactionID   string             `json:"transaction_id"`
	Amount          AmountInfo         `json:"amount"`
	MchID           string             `json:"mchid"`
	TradeState      string             `json:"trade_state"`
	BankType        string             `json:"bank_type"`
	PromotionDetail []*PromotionDetail `json:"promotion_detail,omitempty"`
	SuccessTime     string             `json:"success_time"`
	Payer           PayerInfo          `json:"payer"`
	OutTradeNo      string             `json:"out_trade_no"`
	AppID           string             `json:"appid"`
	TradeStateDesc  string             `json:"trade_state_desc"`
	TradeType       TradeType          `json:"trade_type"`
	Attach          string             `json:"attach,omitempty"`
	SceneInfo       SceneInfo          `json:"scene_info,omitempty"`
}

type AmountInfo struct {
	PayerTotal    int    `json:"payer_total"`
	Total         int    `json:"total"`
	Currency      string `json:"currency"`
	PayerCurrency string `json:"payer_currency"`
}

type PromotionDetail struct {
	Amount              int            `json:"amount"`
	WeChatPayContribute int            `json:"wechatpay_contribute"`
	CouponID            string         `json:"coupon_id"`
	Scope               string         `json:"scope"`
	MerchantContribute  int            `json:"merchant_contribute"`
	Name                string         `json:"name"`
	OtherContribute     int            `json:"other_contribute"`
	Currency            string         `json:"currency"`
	StockID             string         `json:"stock_id"`
	GoodsDetail         []*GoodsDetail `json:"goods_detail,omitempty"`
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

// AuthorizeLinkRequest provides basic WeChat OpenID authorise link request params bind
type AuthorizeLinkRequest struct {
	RedirectURI string `json:"redirect_uri" binding:"required"`
	State       string `json:"state" binding:"required"`
}
