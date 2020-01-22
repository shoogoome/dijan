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
	utils.InitGlobalSystemConfig()
	c := cache.New("rocksdb", 9999999)

	service.InitHttpService()
	recoveryUtils.MemoryRecovery()
	n, e := cluster.New(utils.GlobalSystemConfig.Server.HostName, utils.GlobalSystemConfig.Server.Address, utils.GlobalSystemConfig.Server.ClusterAddress)
	if e != nil {
		panic(e)
	}
	tcp.New(c, n).Listen()
}
