package paramsException

import (
	"dijan/models"
	"fmt"
)

func UnmarshalBodyJsonFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5201,
		Message: "读取json数据失败",
	}
}

func UnmarshalTextJsonFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5205,
		Message: "解析json失败",
	}
}

func LackParams(mes string) models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5202,
		Message: fmt.Sprintf("%s 参数为必填参数", mes),
	}
}

func ParamsIsNotStandard(key string, ty string) models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5203,
		Message: fmt.Sprintf("%s 参数不规范, 应为 %s 类型", key, ty),
	}
}

func DataUrlParserFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5204,
		Message: "durl解析失败",
	}
}
