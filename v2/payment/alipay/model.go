package alipay

import "github.com/ghinknet/payutils/v2/model"

// Config provides needed config options for Alipay
type Config struct {
	AppID             string
	AppCertPrivateKey string
	AppCert           string
	RootCert          string
	PublicCert        string
	IsProd            bool
}

const (
	// TradeStateWaitBuyerPay Trade created and waiting for buyer to pay
	// -> TradeStatePending
	TradeStateWaitBuyerPay string = "WAIT_BUYER_PAY"
	// TradeStateSuccess Trade successes
	// -> TradeStateSuccess
	TradeStateSuccess string = "TRADE_SUCCESS"
	// TradeStateClosed Trade closed due to time out or refunded after pay
	// -> TradeStateClosed
	TradeStateClosed string = "TRADE_CLOSED"
	// TradeStateFinished Trade finished and cannot be refund
	// ->  TradeStateFinished
	TradeStateFinished string = "TRADE_FINISHED"
)

// TradeStateMap provides a map refer to internal trade state
var TradeStateMap = map[string]model.TradeState{
	TradeStateWaitBuyerPay: model.TradeStatePending,
	TradeStateSuccess:      model.TradeStateSuccess,
	TradeStateClosed:       model.TradeStateClosed,
	TradeStateFinished:     model.TradeStateFinished,
}

// MapState provides a method to map trade state (string) to internal trade state
func MapState(state string) model.TradeState {
	// Try to find in map
	internalState, ok := TradeStateMap[state]
	if !ok {
		return model.TradeStateUnknown
	}
	return internalState
}
