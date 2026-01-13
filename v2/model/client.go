package model

import (
	"github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/gopay/wechat/v3"
)

// Client is the public template of client
type Client struct {
	Alipay *alipay.ClientV3
	WeChat *wechat.ClientV3
}
