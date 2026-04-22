package wechat

import (
	"go.gh.ink/payutils/v2/model"
)

func ErrWeChatPayRespCodeInvalidBuilder(
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

var ErrWeChatPayRespCodeInvalid = model.New("wechat pay resp code invalid")
var ErrWeChatRedirectURIMismatch = model.New("wechat redirect uri mismatch")
var ErrWeChatOpenIDIsRequired = model.New("wechat open id is required")
