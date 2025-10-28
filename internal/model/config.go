package model

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
)

type AlipayConfig struct {
	AppID             string
	AppCertPrivateKey string
	AppCert           string
	RootCert          string
	PublicCert        string
	IsProd            bool
}

type WeChatPayConfig struct {
	AppID                    string
	AppSecret                string
	MerchantID               string
	MerchantAPIv3Key         string
	MerchantCertSerialNumber string
	MerchantPrivateKey       string
	PublicKey                string
	PublicKeyID              string
}

type Config struct {
	Alipay        *AlipayConfig
	WeChatPay     *WeChatPayConfig
	Gin           *gin.RouterGroup
	Fiber         fiber.Router
	Debug         bool
	AllowedOrigin []string
	Endpoint      string
	ErrorHandler  func(c any, err error) error
	OrderInfo     func(orderID string, authorization string) (OrderInfo, error)
	OrderStatus   func(orderID string, status TradeStatus) error
}
