package payutils

import (
	"git.ghink.net/ghink/payutils/internal/client"
	"git.ghink.net/ghink/payutils/internal/model"
	"git.ghink.net/ghink/payutils/internal/route"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/gopay/wechat/v3"
)

type AlipayConfig = model.AlipayConfig
type WeChatPayConfig = model.WeChatPayConfig
type Config = model.Config

type Client = client.Client

type TradeStatus = model.TradeStatus
type OrderInfo = model.OrderInfo

var ErrMissEndpoint = model.ErrMissEndpoint
var ErrMissOrderHandler = model.ErrMissOrderHandler
var ErrWeChatPayRespCodeInvalid = model.ErrWeChatPayRespCodeInvalid
var ErrWeChatRedirectURIMismatch = model.ErrWeChatRedirectURIMismatch
var ErrOpenIDIsRequired = model.ErrOpenIDIsRequired

// CreateClient creates a unified client and register route
func CreateClient(config model.Config) (*client.Client, error) {
	// Create payutils client
	c := &client.Client{}

	var err error

	// Debug switch
	var debugOption gopay.DebugSwitch
	if config.Debug {
		debugOption = gopay.DebugOn
	} else {
		debugOption = gopay.DebugOff
	}

	// Check endpoint
	if config.Endpoint == "" {
		return nil, model.ErrMissEndpoint
	}

	// Check order handler
	if config.OrderInfo == nil || config.OrderStatus == nil {
		return nil, model.ErrMissOrderHandler
	}

	// Check error handler
	if config.ErrorHandler == nil {
		config.ErrorHandler = model.RespInternalServerError
	}

	if config.Alipay != nil {
		// Create alipay client
		c.Alipay, err = alipay.NewClientV3(
			config.Alipay.AppID,
			config.Alipay.AppCertPrivateKey,
			config.Alipay.IsProd,
		)
		if err != nil {
			return nil, err
		}
		// Set alipay cert
		err = c.Alipay.SetCert(
			[]byte(config.Alipay.AppCert),
			[]byte(config.Alipay.RootCert),
			[]byte(config.Alipay.PublicCert),
		)
		if err != nil {
			return nil, err
		}
		// Debug switch
		c.Alipay.DebugSwitch = debugOption
	}

	if config.WeChatPay != nil {
		// Create wechat-pay client
		c.WeChat, err = wechat.NewClientV3(
			config.WeChatPay.MerchantID,
			config.WeChatPay.MerchantCertSerialNumber,
			config.WeChatPay.MerchantAPIv3Key,
			config.WeChatPay.MerchantPrivateKey,
		)
		if err != nil {
			return nil, err
		}
		// Auto verify sign by public key
		err = c.WeChat.AutoVerifySignByPublicKey(
			[]byte(config.WeChatPay.PublicKey),
			config.WeChatPay.PublicKeyID,
		)
		if err != nil {
			return nil, err
		}
		// Debug switch
		c.WeChat.DebugSwitch = debugOption
	}

	// Register gin route
	if config.Gin != nil {
		route.GinRegister(config.Gin, c, config)
	}

	return c, nil
}
