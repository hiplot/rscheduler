package api

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	g := gin.Default()
	initRouter(g)
	_ = g.Run(":8080")
}

func initRouter(g *gin.Engine) {
	g.GET("/completed", TaskAPI.TaskCompleteHandler)
	g.GET("/base/info", BaseAPI.Info)
}
