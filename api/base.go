package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"rscheduler/global"
)

type baseAPI struct {
}

var BaseAPI = &baseAPI{}

func (b baseAPI) Info(c *gin.Context) {
	c.JSON(http.StatusOK, BaseInfoResp{
		BaseResp: NewBaseSuccessResp(),
		BaseInfo: baseInfo{Version: getVersion()},
	})
}

func getVersion() string {
	// 读取VERSION文件
	content, err := os.ReadFile("VERSION")
	if err != nil {
		global.Logger.Errorln(err)
		return "Missing VERSION file"
	}
	return string(content)
}
