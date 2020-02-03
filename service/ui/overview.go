package ui

import (
	"dijan/core/cluster"
	"dijan/core/view/auth"
	"dijan/utils"
	"dijan/utils/node"
	templatesUtils "dijan/utils/templates"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/kataras/iris"
	"io/ioutil"
)

// 总览页
func Overview(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	if !auth.IsLogin() {
		ctx.Redirect("/")
	}
	// 获取总存储数据数
	keys := 0
	size := 0
	keyNum := make(chan []int)
	for _, node := range cluster.Member.Members() {
		go func(name string) {
			defer func() {
				recover()
			}()
			response, err := utils.Requests("GET", fmt.Sprintf("http://%s%s/api/node", name, utils.GlobalSystemConfig.Server.HttpListenPort), nil)
			if err != nil {
				panic(err)
			}
			bode, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			var nodeInfo []map[string]interface{}
			json.Unmarshal(bode, &nodeInfo)

			size := 0
			for _, i := range nodeInfo {
				size += int(i["size"].(float64))
			}

			keyNum <- []int{len(nodeInfo), size}
		}(node.Name)
	}
	for i := 0; i < cluster.Member.NumMembers(); i++ {
		s := <- keyNum
		keys += s[0]
		size += s[1]
	}
	ctx.HTML(templatesUtils.Render("overview.html", pongo2.Context{
		"verification": true,
		"login": true,
		"size": fmt.Sprintf("%.3f", float64(size) / 1024 / 1024),
		"nodes": nodeUtils.GetNodeInfo(),
		"keys": keys,
		"node_number": cluster.Member.NumMembers(),
		"circular_number": utils.GlobalSystemConfig.Cluster.CircleNumberOfNode,
	}))
}

