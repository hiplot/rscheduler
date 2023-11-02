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
	BaseInfo baseInfo `json:"base_info"`
}

type baseInfo struct {
	Version string `json:"version"`
}

type ProcessorInfoResp struct {
	BaseResp
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

type TaskInfoResp struct {
	BaseResp
	TaskInfo []taskInfo `json:"task_info"`
}

type taskInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProcessorID string `json:"processor_id"`
	Memory      uint64 `json:"memory"`
	CPU         uint64 `json:"cpu"`
	CreateAt    int64  `json:"create_at"`
	StartAt     int64  `json:"start_at"`
	RunTime     int64  `json:"run_time"`
}
