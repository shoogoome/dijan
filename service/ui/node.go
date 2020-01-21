package ui

import (
	"dijan/core/view/auth"
	"dijan/exception/system"
	"dijan/utils"
	"dijan/utils/node"
	"dijan/utils/templates"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/kataras/iris"
	"io/ioutil"
)

// 节点页
func NodeUI(ctx iris.Context, auth auth.DijanAuthAuthorization, node string) {
	if !auth.IsLogin() {
		ctx.Redirect("/")
	}

	var data []map[string]interface{}

	response, err := utils.Requests("GET", fmt.Sprintf(
		"http://%s%s/api/node",
		node,
		utils.GlobalSystemConfig.Server.HttpListenPort), nil)

	if err != nil {
		panic(systemException.SystemCommunicationFail())
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(systemException.SystemCommunicationFail())
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(systemException.SystemCommunicationFail())
	}
	ctx.HTML(templatesUtils.Render("node.html", pongo2.Context{
		"records": data,
		"nodes": nodeUtils.GetNodeInfo(),
		"count": len(data),
		"hostname": node,
	}))
}