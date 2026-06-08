package action

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.gh.ink/payutils/v3/model"
)

func StatusUpdaterConstructor(conf model.Config) func(
	ctx context.Context, r *http.Request, upstream string, tradeID string, status model.TradeState, tm time.Time,
) error {
	return func(
		ctx context.Context, r *http.Request, upstream string, tradeID string, status model.TradeState, tm time.Time,
	) error {
		// Contract is optional; without it there is nowhere to push status.
		if conf.Contract == nil {
			return nil
		}
		tradeID = strings.TrimSuffix(strings.TrimPrefix(tradeID, conf.TradeIDPrefix), conf.TradeIDSuffix)
		return conf.Contract.StatusUpdater(ctx, r, upstream, tradeID, status, tm)
	}
}
