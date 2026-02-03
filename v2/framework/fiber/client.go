package fiber

import (
	"github.com/ghinknet/json"
	"github.com/ghinknet/payutils/v2/model"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/gopay/wechat/v3"
)

// CreateClient creates a unified client and register route
func CreateClient(config Config) (*Client, error) {
	// Create a new client
	client := &Client{
		Config:  &config,
		Payment: new(model.Client),
	}

	// Debug switch
	var debugOption gopay.DebugSwitch
	if client.Config.Basic.Debug {
		debugOption = gopay.DebugOn
	} else {
		debugOption = gopay.DebugOff
	}

	// Check AllowOrigins
	if client.Config.Basic.AllowOrigins == nil {
		return nil, model.ErrMissAllowedOrigin
	}

	// Check endpoint
	if client.Config.Basic.Endpoint == "" {
		return nil, model.ErrMissEndpoint
	}

	// Check handler
	if client.Config.DetailProvider == nil || client.Config.StatusUpdater == nil {
		return nil, model.ErrMissHandler
	}

	// Check error handler
	if client.Config.ErrorHandler == nil {
		client.Config.ErrorHandler = RespInternalServerError
	}

	// Check methods
	if client.Config.Marshal == nil {
		client.Config.Marshal = json.Marshal
	}
	if client.Config.Unmarshal == nil {
		client.Config.Unmarshal = json.Unmarshal
	}

	// Init payment
	var err error
	// Alipay
	if client.Config.Alipay != nil {
		// Create alipay client
		client.Payment.Alipay, err = alipay.NewClientV3(
			client.Config.Alipay.AppID,
			client.Config.Alipay.AppCertPrivateKey,
			client.Config.Alipay.IsProd,
		)
		if err != nil {
			return nil, err
		}
		// Set alipay cert
		if err = client.Payment.Alipay.SetCert(
			[]byte(client.Config.Alipay.AppCert),
			[]byte(client.Config.Alipay.RootCert),
			[]byte(client.Config.Alipay.PublicCert),
		); err != nil {
			return nil, err
		}
		// Debug switch
		client.Payment.Alipay.DebugSwitch = debugOption
	}
	// WeChatPay
	if client.Config.WeChatPay != nil {
		// Create wechat-pay client
		client.Payment.WeChat, err = wechat.NewClientV3(
			client.Config.WeChatPay.MerchantID,
			client.Config.WeChatPay.MerchantCertSerialNumber,
			client.Config.WeChatPay.MerchantAPIv3Key,
			client.Config.WeChatPay.MerchantPrivateKey,
		)
		if err != nil {
			return nil, err
		}
		// Auto verify sign by public key
		if err = client.Payment.WeChat.AutoVerifySignByPublicKey(
			[]byte(client.Config.WeChatPay.PublicKey),
			client.Config.WeChatPay.PublicKeyID,
		); err != nil {
			return nil, err
		}
		// Debug switch
		client.Payment.WeChat.DebugSwitch = debugOption
	}

	// Register framework route
	register(client.Config.Fiber, client)

	return client, nil
}
