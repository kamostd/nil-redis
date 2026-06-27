package transport

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"transport",

		fx.Provide(NewHTTPHandler),

		fx.Invoke(func(handler *HTTPHandler, group *gin.RouterGroup) {
			handler.Routing(group)
		}),

		fx.Invoke(func(logger *zap.Logger) *zap.Logger {
			return logger.With(zap.String("module", "transport"))
		}),
	)
}
