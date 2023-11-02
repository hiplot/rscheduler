package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, obj)
}
