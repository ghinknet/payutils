package alipay

import (
	"git.ghink.net/ghink/payutils/internal/model"
	"github.com/gin-gonic/gin"
)

type GinController struct {
	Client *model.Client
}

func (g *GinController) Create(c *gin.Context) {
}

func (g *GinController) Callback(c *gin.Context) {

}
