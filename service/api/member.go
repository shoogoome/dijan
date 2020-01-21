package api

import (
	"dijan/core/cluster"
	"dijan/core/view/auth"
	"dijan/utils"
	"github.com/kataras/iris"
)

func GetMember(ctx iris.Context, auth auth.DijanAuthAuthorization) {

	nodes := make([]string, cluster.Member.NumMembers())
	for index, node := range cluster.Member.Members() {
		nodes[index] = node.Name + utils.GlobalSystemConfig.Server.TcpListenPort
	}
	ctx.JSON(nodes)
}