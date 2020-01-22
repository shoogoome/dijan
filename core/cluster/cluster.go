package cluster

import (
	"dijan/utils"
	"fmt"
	"github.com/hashicorp/memberlist"
	"os"
	"stathat.com/c/consistent"
	"time"
)

var Member *memberlist.Memberlist
var NodeObject Node

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

func New(hostname, addr, cluster string) (Node, error) {
	fmt.Println("flag", hostname, addr, cluster)
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = os.Stdout
	l, e := memberlist.Create(conf)
	Member = l
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
Conn:
	_, e = l.Join(clu)
	if e != nil {
		fmt.Println("[!] 集群连接失败，5s后尝试重连...", e)
		time.Sleep(time.Second * 5)
		goto Conn
	}
	circle := consistent.New()
	circle.NumberOfReplicas = utils.GlobalSystemConfig.Cluster.CircleNumberOfNode
	go func() {
		for {
			m := l.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	NodeObject = &node{circle, addr}
	return NodeObject, nil
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}
