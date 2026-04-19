package model

import (
	"context"
	"net/http"
	"time"
)

type Contract interface {
	DetailProvider(ctx context.Context, r *http.Request, upstream string, tradeID string) (TradeDetail, error)
	StatusUpdater(
		ctx context.Context, r *http.Request, upstream string, tradeID string, status TradeState, time time.Time,
	) error
}
