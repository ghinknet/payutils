package payutils

import (
	"git.ghink.net/ghink/payutils/internal/method"
	"git.ghink.net/ghink/payutils/internal/model"
)

type AlipayConfig = model.AlipayConfig
type WeChatPayConfig = model.WeChatPayConfig
type Config = model.Config

type Client = model.Client

type TradeStatus = model.TradeStatus
type OrderInfo = model.OrderInfo

var ErrMissEndpoint = model.ErrMissEndpoint
var ErrMissOrderHandler = model.ErrMissOrderHandler
var ErrWeChatPayRespCodeInvalid = model.ErrWeChatPayRespCodeInvalid
var ErrWeChatRedirectURIMismatch = model.ErrWeChatRedirectURIMismatch
var ErrOpenIDIsRequired = model.ErrOpenIDIsRequired

var CreateClient = method.CreateClient
