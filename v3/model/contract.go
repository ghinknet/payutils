package model

import (
	"context"
	"net/http"
	"time"
)

type Contract interface {
	StatusUpdater(
		ctx context.Context, r *http.Request, upstream string, tradeID string, status TradeState, time time.Time,
	) error
}
