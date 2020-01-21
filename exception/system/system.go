package systemException

import (
	"dijan/models"
)

func SystemException() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5000,
		Message: "系统错误",
	}
}

func SystemApiTokenVerificationFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5001,
		Message: "token验证失败",
	}
}

func SystemCommunicationFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5002,
		Message: "通讯失败",
	}
}




