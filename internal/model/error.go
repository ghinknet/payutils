package model

import (
	"errors"
)

var ErrMissEndpoint = errors.New("miss endpoint")
var ErrMissOrderHandler = errors.New("miss order handler")
var ErrWeChatPayRespCodeInvalid = errors.New("wechat pay resp code invalid")
var ErrWeChatRedirectURIMismatch = errors.New("wechat redirect_uri mismatch")
var ErrOpenIDIsRequired = errors.New("open id is required")
