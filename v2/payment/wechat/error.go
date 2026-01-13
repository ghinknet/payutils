package wechat

import "errors"

var ErrWeChatPayRespCodeInvalid = errors.New("wechat pay resp code invalid")
var ErrWeChatRedirectURIMismatch = errors.New("wechat redirect uri mismatch")
var ErrWeChatOpenIDIsRequired = errors.New("wechat open id is required")
