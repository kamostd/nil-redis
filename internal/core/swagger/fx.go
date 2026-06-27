package swagger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"swagger docs",

		fx.Invoke(func(router *gin.RouterGroup) {
			Routing(router)
		},
		),
	)
}
