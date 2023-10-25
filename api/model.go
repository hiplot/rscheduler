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
