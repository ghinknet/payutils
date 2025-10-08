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
	Alipay    *AlipayConfig
	Wechat    *WechatConfig
	Gin       *gin.RouterGroup
	Debug     bool
	OrderInfo func(orderID string, authorization string) string
}
