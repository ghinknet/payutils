package alipay

import (
	"github.com/ghinknet/payutils/v2/model"
)

func ErrAlipayRespCodeInvalid(
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

var ErrAlipayNotifyVerifyFailed = model.New("alipay notify verify failed")
