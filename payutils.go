package payutils

import (
	"git.ghink.net/ghink/payutils/internal/method"
	"git.ghink.net/ghink/payutils/internal/model"
)

type AlipayConfig = model.AlipayConfig
type WechatConfig = model.WechatConfig
type Config = model.Config

type TradeStatus = model.TradeStatus
type OrderInfo = model.OrderInfo

var CreateClient = method.CreateClient
