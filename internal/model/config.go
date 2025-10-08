package model

import "github.com/gin-gonic/gin"

type AlipayConfig struct {
	AppID             string
	AppCertPrivateKey string
	AppCert           string
	RootCert          string
	PublicCert        string
	IsProd            bool
}

type WechatConfig struct {
	AppID                    string
	AppSecret                string
	MerchantID               string
	MerchantAPIv3Key         string
	MerchantCertSerialNumber string
	MerchantPrivateKey       string
}

type Config struct {
	Alipay       *AlipayConfig
	Wechat       *WechatConfig
	Gin          *gin.RouterGroup
	Debug        bool
	ErrorHandler func(c *gin.Context, err error)
	OrderInfo    func(orderID string, authorization string) (OrderInfo, error)
	OrderStatus  func(orderID string, authorization string, status TradeStatus) error
}
