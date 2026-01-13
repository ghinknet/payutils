package alipay

import "errors"

var ErrAlipayRespCodeInvalid = errors.New("alipay resp code invalid")
var ErrAlipayNotifyVerifyFailed = errors.New("alipay notify verify failed")
