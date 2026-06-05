package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	return gin.Default()
}

func NewGroupRouter(engine *gin.Engine) *gin.RouterGroup {
	router := engine.Group("/api")

	return router
}

func NewHttpServer() *http.Server {
	return &http.Server{
		Addr: "localhost:8080",
	}
}
