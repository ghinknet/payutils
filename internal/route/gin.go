package route

import (
	"git.ghink.net/ghink/payutils/v2/internal/client"
	"git.ghink.net/ghink/payutils/v2/internal/controller/alipay"
	"git.ghink.net/ghink/payutils/v2/internal/controller/wechat"
	"git.ghink.net/ghink/payutils/v2/internal/model"
	"github.com/gin-gonic/gin"
)

// GinRegister registers gin router
func GinRegister(r *gin.RouterGroup, client *client.Client, config model.Config) {
	{
		alipayRoute := r.Group("/alipay")
		alipayGinController := alipay.GinController{Client: client, Config: config}
		alipayRoute.POST("/create", alipayGinController.Create)
		alipayRoute.POST("/callback", alipayGinController.Callback)
	}
	{
		wechatRoute := r.Group("/wechat")
		wechatGinController := wechat.GinController{Client: client, Config: config}
		wechatRoute.POST("/create", wechatGinController.Create)
		wechatRoute.POST("/callback", wechatGinController.Callback)
		wechatRoute.POST("/openIDCallback", wechatGinController.OpenIDCallback)
		wechatRoute.POST("/authorizeLink", wechatGinController.AuthorizeLinkGen)
	}
}
