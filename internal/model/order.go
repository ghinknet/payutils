package model

// TradeStatus provides enum type for trade status const
type TradeStatus string

const (
	TradePending  TradeStatus = "PENDING"
	TradeSuccess  TradeStatus = "SUCCESS"
	TradeClosed   TradeStatus = "CLOSED"
	TradeFinished TradeStatus = "FINISHED"
)

// OrderRequest provides basic requests params bind
type OrderRequest struct {
	OrderID string `json:"orderID" binding:"required"`
}

// OrderInfo provides basic order info
type OrderInfo struct {
	Subject string
	Price   int64
}
