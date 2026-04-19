package model

import "net/http"

type PayClient interface {

	// Http Handler

	Create(http.ResponseWriter, *http.Request)
	Callback(http.ResponseWriter, *http.Request)

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
