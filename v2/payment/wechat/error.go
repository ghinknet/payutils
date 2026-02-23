package wechat

import (
	"github.com/ghinknet/payutils/v2/model"
)

func ErrWeChatPayRespCodeInvalid(
	upstreamCode int,
	upstreamResponse string,
	upstreamMessage string,
) error {
	return model.New(
		"wechat pay resp code invalid",
		model.WithUpstreamCode(upstreamCode),
		model.WithUpstreamResponse(upstreamResponse),
		model.WithUpstreamMessage(upstreamMessage),
	)
}

var ErrWeChatRedirectURIMismatch = model.New("wechat redirect uri mismatch")
var ErrWeChatOpenIDIsRequired = model.New("wechat open id is required")
