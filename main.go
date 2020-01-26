package main

import (
	"dijan/core/cache"
	"dijan/core/cluster"
	"dijan/core/tcp"
	"dijan/service"
	"dijan/utils"
	"dijan/utils/recovery"
)

func main() {
	ipChan := utils.InitGlobalSystemConfig()
	c := cache.New()

	service.InitHttpService()
	recoveryUtils.MemoryRecovery()
	n, e := cluster.New(utils.GlobalSystemConfig.Server.Address, utils.GlobalSystemConfig.Server.ClusterAddress, ipChan)
	if e != nil {
		panic(e)
	}
	tcp.New(c, n).Listen()
}
