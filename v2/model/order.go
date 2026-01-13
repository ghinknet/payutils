package model

// Platform provides enum type for platforms
type Platform string

const (
	PlatformPC     Platform = "PC"
	PlatformMobile Platform = "Mobile"
	PlatformWeChat Platform = "WeChat"
)

// OrderDetail provides order's detail
type OrderDetail struct {
	Subject  string
	Price    int64
	Currency string
	Expiry   int64
}

// TradeMethod provides enum type for trade payment method
type TradeMethod string

const (
	TradeMethodAlipay    TradeMethod = "Alipay"
	TradeMethodWeChatPay TradeMethod = "WeChatPay"
	// TradeMethodUnknown trade method unknown due to system error
	TradeMethodUnknown TradeMethod = "Unknown"
)

// TradeState provides enum type for trade state const
type TradeState string

const (
	// TradeStatePending waiting for user to pay
	TradeStatePending TradeState = "PENDING"
	// TradeStateSuccess trade finished successfully
	TradeStateSuccess TradeState = "SUCCESS"
	// TradeStateClosed trade closed due to timed out or closed manually
	TradeStateClosed TradeState = "CLOSED"
	// TradeStateFinished trade finished and cannot be refund
	TradeStateFinished TradeState = "FINISHED"
	// TradeStateUnknown trade status unknown due to system error
	TradeStateUnknown TradeState = "UNKNOWN"
)

// OrderRequest provides basic requests params bind
type OrderRequest struct {
	OrderID  string   `json:"orderID" binding:"required" validate:"required"`
	Platform Platform `json:"platform" binding:"required" validate:"required"`
	OpenID   string   `json:"openID"`
}
