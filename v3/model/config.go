package model

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	// Payutils
	Debug               bool
	Endpoint            string
	TradeIDPrefix       string
	TradeIDSuffix       string
	NoNewPaymentWindows *time.Duration
	SafetyMargin        *time.Duration
	// Http
	Instances map[string]any
	// Pay
	Credentials C
	// Contract
	Contract     Contract
	ErrorHandler func(ctx context.Context, r *http.Request, err error) error
	// JSON
	Unmarshal func(data []byte, v any) error
	Marshal   func(v any) ([]byte, error)
}
