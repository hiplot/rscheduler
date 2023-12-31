package api

const (
	BaseSuccessCode = 0
	BaseFailCode    = 1
)

type BaseResp struct {
	StatusCode    int64  `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

func NewBaseSuccessResp() BaseResp {
	return BaseResp{
		StatusCode:    BaseSuccessCode,
		StatusMessage: "Success",
	}
}

func NewBaseFailResp() BaseResp {
	return BaseResp{
		StatusCode:    BaseFailCode,
		StatusMessage: "Fail",
	}
}

type BaseInfoResp struct {
	BaseResp
	BaseInfo BaseInfo `json:"base_info"`
}

type BaseInfo struct {
	Version string `json:"version"`
}

type ProcessorInfoResp struct {
	BaseResp
	ProcessorInfo []ProcessorInfo `json:"processor_info"`
}

type ProcessorInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PID          int    `json:"pid"`
	TotalTaskNum uint32 `json:"total_task_num"`
	Running      bool   `json:"running"`
	TaskID       string `json:"task_id"`
}

type TaskInfoResp struct {
	BaseResp
	TaskInfo []TaskInfo `json:"task_info"`
}

type TaskInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProcessorID string `json:"processor_id"`
	Memory      uint64 `json:"memory"`
	CPU         uint64 `json:"cpu"`
	CreateAt    int64  `json:"create_at"`
	StartAt     int64  `json:"start_at"`
	RunTime     int64  `json:"run_time"`
}

type ProcessorDeleteReq struct {
	ID    string `form:"id" binding:"required"`
	Force bool   `form:"force"`
}

type ProcessorDeleteResp struct {
	BaseResp
	Success bool   `json:"success"`
	Info    string `json:"info"`
}

type TaskDeleteReq struct {
	ID string `form:"id" binding:"required"`
}

type TaskDeleteResp struct {
	BaseResp
	Success bool   `json:"success"`
	Info    string `json:"info"`
}
