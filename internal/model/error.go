package model

import "errors"

var ErrWeChatPayRespCodeInvalid = errors.New("wechat pay resp code invalid")
var ErrWeChatRedirectURIMismatch = errors.New("wechat redirect_uri mismatch")
