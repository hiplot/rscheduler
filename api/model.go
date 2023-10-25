package api

const (
	BaseSuccessCode = 0
	BaseFailCode    = 1
)

type BaseResponse struct {
	StatusCode    int64  `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

func NewBaseSuccessResponse() BaseResponse {
	return BaseResponse{
		StatusCode:    BaseSuccessCode,
		StatusMessage: "Success",
	}
}

func NewBaseFailResponse() BaseResponse {
	return BaseResponse{
		StatusCode:    BaseFailCode,
		StatusMessage: "Fail",
	}
}

type BaseInfoResponse struct {
	BaseResponse
	BaseInfo baseInfo `json:"base_info"`
}

type baseInfo struct {
	Version string `json:"version"`
}

type ProcessorInfoResponse struct {
	BaseResponse
	ProcessorInfo []processorInfo `json:"processor_info"`
}

type processorInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PID          int    `json:"pid"`
	TotalTaskNum uint32 `json:"total_task_num"`
	Running      bool   `json:"running"`
	TaskID       string `json:"task_id"`
}
