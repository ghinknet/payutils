package model

import "time"

type TradeState string

const (
	// TradeStatePending waiting for user to pay
	TradeStatePending TradeState = "PENDING"
	// TradeStateSuccess trade finished successfully
	TradeStateSuccess TradeState = "SUCCESS"
	// TradeStateClosed trade closed due to timed out or closed manually
	TradeStateClosed TradeState = "CLOSED"
	// TradeStateFinished trade finished and cannot be updated (like refund)
	TradeStateFinished TradeState = "FINISHED"
	// TradeStateUnknown trade status unknown due to system error
	TradeStateUnknown TradeState = "UNKNOWN"
)

type TradeDetail struct {
	Subject  string
	Price    int64
	Currency string
	Expiry   time.Time
}
