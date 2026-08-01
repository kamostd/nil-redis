package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGinEngine() *gin.Engine {
	return gin.Default()
}

func NewRouterGroup(engine *gin.Engine) *gin.RouterGroup {
	return engine.Group("/api")
}

func NewHTTPServer(engine *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    "localhost:8080",
		Handler: engine,
	}
}
