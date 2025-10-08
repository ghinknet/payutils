package wechat

import (
	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/gin-gonic/gin"
)

type GinController struct {
	Client *model.Client
	Config model.Config
}

func (g *GinController) Create(c *gin.Context) {

}

func (g *GinController) Callback(c *gin.Context) {

}

func (g *GinController) OpenIDCallback(c *gin.Context) {

}

func (g *GinController) BasicInfo(c *gin.Context) {

}
