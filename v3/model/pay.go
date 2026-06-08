package model

import (
	"context"
	"net/http"
)

type PayClient interface {

	// Create attempts to create a trade with the given param.
	//
	// Each driver defines its own param struct. The driver MUST assert whether
	// it claims the given param: if the concrete type does not belong to this
	// driver, it returns claimed=false so the dispatcher can try the next one.
	// When claimed=true, result holds the driver-specific payload (ready to be
	// marshalled, e.g. a pay url or a jsapi sign object) and err the outcome.
	Create(ctx context.Context, param any) (claimed bool, result any, err error)

	// Callback handles an upstream asynchronous notification.
	Callback(w http.ResponseWriter, r *http.Request)

	// Actions

	Status(tradeID string) (ReturnStatus, error)
	Close(tradeID string) error
	Refund(
		tradeID string,
		currency string, refundID string,
		totalAmount int64, refundAmount int64,
		reason string,
	) error
}
