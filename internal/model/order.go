package model

// TradeStatus provides enum type for trade status const
type TradeStatus string

const (
	TradePending  TradeStatus = "PENDING"
	TradeSuccess  TradeStatus = "SUCCESS"
	TradeClosed   TradeStatus = "CLOSED"
	TradeFinished TradeStatus = "FINISHED"
)

// Platform provides enum type for platforms
type Platform string

const (
	PlatformPC     Platform = "PC"
	PlatformMobile Platform = "Mobile"
	PlatformWeChat Platform = "WeChat"
)

// OrderRequest provides basic requests params bind
type OrderRequest struct {
	OrderID  string   `json:"orderID" binding:"required"`
	Platform Platform `json:"platform" binding:"required"`
	OpenID   string   `json:"openID"`
}

// OrderInfo provides basic order info
type OrderInfo struct {
	Subject string
	Price   int64
}
