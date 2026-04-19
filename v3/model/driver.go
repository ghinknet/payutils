package model

import (
	"context"
	"net/http"
	"time"
)

type PayDriver interface {
	NewClient(params PayDriverClientParam) (PayClient, error)
}

type PayDriverClientParam struct {
	// Pay client credential
	Credential map[string]string
	// Contract
	StatusUpdater func(
		ctx context.Context, r *http.Request, upstream string, tradeID string, status TradeState, time time.Time,
	) error
	ErrorHandler func(ctx context.Context, r *http.Request, err error) error
	// JSON
	Unmarshal func(data []byte, v any) error
	Marshal   func(v any) ([]byte, error)
}

type HttpDriver interface {
	NewInstance(instance any) (HttpInstance, error)
}
