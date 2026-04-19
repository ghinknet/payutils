package model

import "time"

type ReturnStatus struct {
	TradeStatus TradeState
	Upstream    string
	Time        time.Time
}
