package cacheException

import "dijan/models"

func RocksdbSetFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5305,
		Message: "存储失败",
	}
}

func RocksdbDeleteFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5306,
		Message: "删除失败",
	}
}

func RocksdbGetFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5307,
		Message: "获取失败",
	}
}
