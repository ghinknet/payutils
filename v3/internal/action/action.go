package action

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ghinknet/payutils/v3/model"
)

func StatusUpdaterConstructor(conf model.Config) func(
	ctx context.Context, r *http.Request, upstream string, tradeID string, status model.TradeState, tm time.Time,
) error {
	return func(
		ctx context.Context, r *http.Request, upstream string, tradeID string, status model.TradeState, tm time.Time,
	) error {
		tradeID = strings.TrimSuffix(strings.TrimPrefix(tradeID, conf.TradeIDPrefix), conf.TradeIDSuffix)
		return conf.Contract.StatusUpdater(ctx, r, upstream, tradeID, status, tm)
	}
}
