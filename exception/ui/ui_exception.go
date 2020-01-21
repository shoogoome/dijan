package uiException

import (
	"dijan/models"
)

func UIFileNoFoundException() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5100,
		Message: "查找静态资源失败",
	}
}




