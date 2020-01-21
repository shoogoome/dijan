package api

import (
	"bytes"
	"dijan/core/cache"
	"dijan/core/cluster"
	"dijan/core/view/auth"
	"dijan/exception/cache"
	"dijan/utils"
	"dijan/utils/params"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
)


// set value
func ModifyData(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.CheckLogin()
	reqBody := paramsUtils.RequestJsonInterface(ctx)
	params := paramsUtils.NewParamsParser(reqBody)

	// 其他节点存储
	if addr, ok := cluster.NodeObject.ShouldProcess(params.Str("key", "存储键")); !ok {
		body, _ := json.Marshal(reqBody)
		response, err := utils.Requests(
			"POST",
			fmt.Sprintf("http://%s%s/api/storage/set", addr, utils.GlobalSystemConfig.Server.HttpListenPort),
			bytes.NewBuffer(body))
		if err != nil || response.StatusCode != 200 {
			panic(cacheException.RocksdbSetFail())
		}
	// 本机节点存储
	} else {
		if err := cache.Conn.Set(params.Str("key", "存储键"), []byte(params.Str("value", "存储值")), params.Int("ttl", "过期时间", -1)); err != nil {
			panic(cacheException.RocksdbSetFail())
		}
	}
	ctx.JSON(iris.Map {
		"status": true,
	})
}

// 删除
func DeleteData(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.CheckLogin()
	reqBody := paramsUtils.RequestJsonInterface(ctx)
	params := paramsUtils.NewParamsParser(reqBody)

	// 其他节点存储
	if addr, ok := cluster.NodeObject.ShouldProcess(params.Str("key", "存储键")); !ok {
		body, _ := json.Marshal(reqBody)
		response, err := utils.Requests(
			"DELETE",
			fmt.Sprintf("http://%s%s/api/storage/delete", addr, utils.GlobalSystemConfig.Server.HttpListenPort),
			bytes.NewBuffer(body))
		if err != nil || response.StatusCode != 200 {
			panic(cacheException.RocksdbDeleteFail())
		}
		// 本机节点存储
	} else {
		if err := cache.Conn.Del(params.Str("key", "存储键")); err != nil {
			panic(cacheException.RocksdbDeleteFail())
		}
	}
	ctx.JSON(iris.Map {
		"status": true,
	})
}

// 获取数据
func GetData(ctx iris.Context, auth auth.DijanAuthAuthorization, key string) {
	auth.CheckLogin()

	value, _ := cache.Conn.Get(key)
	ctx.JSON(iris.Map {
		"value": string(value),
	})
}