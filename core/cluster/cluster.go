package cluster

import (
	"dijan/utils"
	"fmt"
	"github.com/hashicorp/memberlist"
	"os"
	"stathat.com/c/consistent"
	"sync"
	"time"
)

var Member *memberlist.Memberlist
var NodeObject Node
var lock sync.RWMutex

type Node interface {
	ShouldProcess(key string) (string, bool)
	Members() []string
	Addr() string
}

type node struct {
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}


func New(addr, cluster string, ipChan chan string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.LogOutput = os.Stderr
	if addr != "" {
		conf.AdvertiseAddr = addr
	} else {
		conf.BindAddr = utils.GlobalSystemConfig.Server.HostName
		conf.Name = utils.GlobalSystemConfig.Server.HostName
	}
	l, e := memberlist.Create(conf)
	Member = l
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
	lock.Lock()
	joinCluster(clu)
	lock.Unlock()

	circle := consistent.New()
	circle.NumberOfReplicas = utils.GlobalSystemConfig.Cluster.CircleNumberOfNode
	// 更新集群信息
	go func() {
		t := time.NewTimer(time.Second * 1)
		for {
			select {
			case <-t.C:
				m := l.Members()
				nodes := make([]string, len(m))
				for i, n := range m {
					nodes[i] = n.Name
				}
				circle.Set(nodes)
				t.Reset(time.Second * 1)
			case ip := <-ipChan:
				lock.Lock()
				joinCluster([]string{ip})
				utils.GlobalSystemConfig.Server.ClusterAddress = ip
				lock.Unlock()
				fmt.Println("[*] 集群重搭建成功")
				t.Reset(time.Second * 1)
			}
		}
	}()
	NodeObject = &node{circle, addr}
	return NodeObject, nil
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}

func joinCluster(cIp []string) {
	_, err := Member.Join(cIp)
	if err != nil {
		fmt.Println("[!] 集群连接失败，5s后尝试重连...", err)
		time.Sleep(time.Second * 5)
		joinCluster(cIp)
	}
}