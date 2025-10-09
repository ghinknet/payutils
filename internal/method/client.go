package method

import (
	"git.ghink.net/ghink/payutils/internal/model"
	"git.ghink.net/ghink/payutils/internal/route"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/gopay/wechat/v3"
)

// CreateClient creates a unified client and register route
func CreateClient(config model.Config) (*model.Client, error) {
	// Create payutils client
	client := &model.Client{}

	var err error

	// Debug switch
	var debugOption gopay.DebugSwitch
	if config.Debug {
		debugOption = gopay.DebugOn
	} else {
		debugOption = gopay.DebugOff
	}

	if config.Alipay != nil {
		// Create alipay client
		client.Alipay, err = alipay.NewClientV3(
			config.Alipay.AppID,
			config.Alipay.AppCertPrivateKey,
			config.Alipay.IsProd,
		)
		if err != nil {
			return nil, err
		}
		// Set alipay cert
		err = client.Alipay.SetCert(
			[]byte(config.Alipay.AppCert),
			[]byte(config.Alipay.RootCert),
			[]byte(config.Alipay.PublicCert),
		)
		if err != nil {
			return nil, err
		}
		// Debug switch
		client.Alipay.DebugSwitch = debugOption
	}

	if config.WeChatPay != nil {
		// Create wechat-pay client
		client.WeChat, err = wechat.NewClientV3(
			config.WeChatPay.MerchantID,
			config.WeChatPay.MerchantCertSerialNumber,
			config.WeChatPay.MerchantAPIv3Key,
			config.WeChatPay.MerchantPrivateKey,
		)
		if err != nil {
			return nil, err
		}
		// Auto verify sign by public key
		err = client.WeChat.AutoVerifySignByPublicKey(
			[]byte(config.WeChatPay.PublicKey),
			config.WeChatPay.PublicKeyID,
		)
		if err != nil {
			return nil, err
		}
		// Debug switch
		client.WeChat.DebugSwitch = debugOption
	}

	// Register gin route
	if config.Gin != nil {
		route.GinRegister(config.Gin, client, config)
	}

	return client, nil
}
