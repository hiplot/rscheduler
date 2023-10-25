package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rscheduler/global"
)

type baseAPI struct {
}

var BaseAPI = &baseAPI{}

func (b baseAPI) Info(c *gin.Context) {
	c.JSON(http.StatusOK, BaseInfoResponse{
		BaseResponse: NewBaseSuccessResponse(),
		BaseInfo:     baseInfo{Version: global.VERSION},
	})
}
