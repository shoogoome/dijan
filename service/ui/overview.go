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
	keyNum := make(chan int)
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
			var nodeInfo []string
			json.Unmarshal(bode, &nodeInfo)
			keyNum <- len(nodeInfo)
		}(node.Name)
	}
	for i := 0; i < cluster.Member.NumMembers(); i++ {
		keys += <- keyNum
	}
	ctx.HTML(templatesUtils.Render("overview.html", pongo2.Context{
		"verification": true,
		"login": true,
		"nodes": nodeUtils.GetNodeInfo(),
		"keys": keys,
		"node_number": cluster.Member.NumMembers(),
		"circular_number": utils.GlobalSystemConfig.Cluster.CircleNumberOfNode,
	}))
}

