package api

import (
	"dijan/core/cache"
	"dijan/core/cluster"
	"dijan/core/view/auth"
	"dijan/utils"
	"fmt"
	"github.com/kataras/iris"
)

// 清空数据库
func Empty(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.CheckLogin()
	s := cache.Conn.NewScanner()
	defer s.Close()

	for s.Scan() {
		k := s.Key()
		cache.Conn.Del(k)
	}
	ctx.JSON(iris.Map {
		"status": true,
	})
}

// 跳转清空
func EmptySearch(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.CheckLogin()
	hostname := ctx.URLParam("hostname")
	utils.Requests(
		"GET",
		fmt.Sprintf("http://%s%s/api/empty", hostname, utils.GlobalSystemConfig.Server.HttpListenPort), nil)
	ctx.JSON(iris.Map {
		"status": true,
	})
}

// 清空全部节点数据库
func EmptyAll(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.CheckLogin()
	lock := make(chan int)
	for _, node := range cluster.Member.Members() {
		go func(name string) {
			utils.Requests(
				"GET",
				fmt.Sprintf("http://%s%s/api/empty", name, utils.GlobalSystemConfig.Server.HttpListenPort), nil)
			lock <- 1
		}(node.Name)
	}
	for i := 0; i < cluster.Member.NumMembers(); i++ {
		<- lock
	}
	ctx.JSON(iris.Map {
		"status": true,
	})
}