package payutils

import (
	"git.ghink.net/ghink/payutils/internal/method"
	"git.ghink.net/ghink/payutils/internal/model"
)

type AlipayConfig = model.AlipayConfig
type WechatConfig = model.WechatConfig
type Config = model.Config

var CreateClient = method.CreateClient
