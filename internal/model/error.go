package model

import (
	"errors"
)

var ErrMissAllowedOrigin = errors.New("miss allowedorigin")
var ErrMissEndpoint = errors.New("miss endpoint")
var ErrMissOrderHandler = errors.New("miss order handler")
var ErrNoEnoughTimeToPay = errors.New("no enough time to pay")
var ErrWeChatPayRespCodeInvalid = errors.New("wechat pay resp code invalid")
var ErrAlipayRespCodeInvalid = errors.New("alipay resp code invalid")
var ErrWeChatRedirectURIMismatch = errors.New("wechat redirect_uri mismatch")
var ErrOpenIDIsRequired = errors.New("open id is required")
