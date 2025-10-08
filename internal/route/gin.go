package route

import (
	"git.ghink.net/ghink/payutils/internal/controller/alipay"
	"git.ghink.net/ghink/payutils/internal/controller/wechat"
	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/gin-gonic/gin"
)

// GinRegister registers gin router
func GinRegister(r *gin.RouterGroup, client *model.Client) {
	{
		alipayRoute := r.Group("/alipay")
		alipayGinController := alipay.GinController{Client: client}
		alipayRoute.POST("/create", alipayGinController.Create)
		alipayRoute.POST("/callback", alipayGinController.Callback)
	}
	{
		wechatRoute := r.Group("/wechat")
		wechatGinController := wechat.GinController{Client: client}
		wechatRoute.POST("/create", wechatGinController.Create)
		wechatRoute.POST("/callback", wechatGinController.Callback)
		wechatRoute.POST("/openIDCallback", wechatGinController.OpenIDCallback)
		wechatRoute.POST("/basicInfo", wechatGinController.BasicInfo)
	}
}
