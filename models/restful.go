package models

type RestfulAPIResult struct {
	Status bool `json:"status"`
	ErrCode int `json:"err_code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

