package alipay

import (
	"github.com/ghinknet/payutils/v2/model"
)

func ErrAlipayRespCodeInvalidBuilder(
	upstreamCode int,
	upstreamResponse string,
	upstreamMessage string,
) error {
	return model.New(
		"alipay resp code invalid",
		model.WithUpstreamCode(upstreamCode),
		model.WithUpstreamResponse(upstreamResponse),
		model.WithUpstreamMessage(upstreamMessage),
	)
}

var ErrAlipayRespCodeInvalid = model.New("alipay resp code invalid")
var ErrAlipayNotifyVerifyFailed = model.New("alipay notify verify failed")
