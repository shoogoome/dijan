package api

import (
	"dijan/core/cache"
	"dijan/core/view/auth"
	"dijan/models"
	"encoding/json"
	"github.com/kataras/iris"
	"strconv"
	"time"
)

// 获取节点存储信息
func NodeStorageInfo(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.CheckLogin()

	index := 1
	var storageValue models.StorageValue
	var records []map[string]interface{}
	scanner := cache.Conn.NewScanner()
	for scanner.Scan() {

		key := scanner.Key()
		value := scanner.Value()
		info := map[string]interface{}{
			"key": key,
			"index": strconv.Itoa(index),
			"size": len(value),
		}
		err := json.Unmarshal(value, &storageValue)
		if err == nil {
			if storageValue.TTL != -1 && time.Now().Unix() >= storageValue.TTL {
				cache.Conn.Del(key)
				continue
			} else {
				info["value"] = string(storageValue.Value)
				info["ttl"] = storageValue.TTL
			}
		} else {
			info["value"] = string(value)
			info["ttl"] = -1
		}
		records = append(records, info)
		index += 1
	}
	defer scanner.Close()
	ctx.JSON(records)
}