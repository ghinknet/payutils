package action

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go.gh.ink/payutils/v3/model"
)

type recordingContract struct {
	called   bool
	tradeID  string
	upstream string
	status   model.TradeState
}

func (c *recordingContract) StatusUpdater(
	_ context.Context, _ *http.Request, upstream string, tradeID string, status model.TradeState, _ time.Time,
) error {
	c.called = true
	c.tradeID = tradeID
	c.upstream = upstream
	c.status = status
	return nil
}

func TestStatusUpdaterConstructor_NilContract(t *testing.T) {
	// Without a contract there is nowhere to push status; must be a safe no-op.
	updater := StatusUpdaterConstructor(model.Config{Contract: nil})
	if err := updater(context.Background(), nil, "alipay", "T123", model.TradeStateSuccess, time.Now()); err != nil {
		t.Fatalf("updater with nil contract returned error: %v", err)
	}
}

func TestStatusUpdaterConstructor_TrimsPrefixSuffix(t *testing.T) {
	rc := &recordingContract{}
	conf := model.Config{
		TradeIDPrefix: "pre_",
		TradeIDSuffix: "_suf",
		Contract:      rc,
	}
	updater := StatusUpdaterConstructor(conf)

	err := updater(context.Background(), nil, "wechat", "pre_T123_suf", model.TradeStateSuccess, time.Now())
	if err != nil {
		t.Fatalf("updater returned error: %v", err)
	}
	if !rc.called {
		t.Fatal("contract StatusUpdater was not called")
	}
	if rc.tradeID != "T123" {
		t.Errorf("tradeID = %q, want T123 (prefix/suffix should be trimmed)", rc.tradeID)
	}
	if rc.upstream != "wechat" {
		t.Errorf("upstream = %q, want wechat", rc.upstream)
	}
	if rc.status != model.TradeStateSuccess {
		t.Errorf("status = %q, want SUCCESS", rc.status)
	}
}
