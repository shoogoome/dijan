package api

import (
	"bytes"
	"dijan/core/cache"
	"dijan/core/cluster"
	"dijan/core/view/auth"
	"dijan/utils"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"io/ioutil"
)

// 节点再平衡
func Rebalance(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.CheckLogin()
	go func() {
		s := cache.Conn.NewScanner()
		defer s.Close()
		for s.Scan() {
			k := s.Key()
			n, ok := cluster.NodeObject.ShouldProcess(k)
			if !ok {
				data := iris.Map {
					"key": k,
					"value": s.Value(),
				}
				body, err := json.Marshal(data); if err == nil {
					response, err := utils.Requests(
						"POST",
						fmt.Sprintf("http://%s%s/api/storage/set", n, utils.GlobalSystemConfig.Server.HttpListenPort),
						bytes.NewBuffer(body)); if err == nil {
						body, err := ioutil.ReadAll(response.Body); if err == nil {
							var resData map[string]interface{}
							err = json.Unmarshal(body, &resData); if err == nil && resData["status"].(bool) {
								cache.Conn.Del(k)
							}
						}
						response.Body.Close()
					}
				}
			}
		}
	}()
	ctx.JSON(iris.Map {
		"status": "已在后台启动节点再平衡",
	})
}
