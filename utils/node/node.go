package nodeUtils

import (
	"dijan/core/cluster"
	"sort"
)

func GetNodeInfo() []string {
	nodes := make([]string, cluster.Member.NumMembers())
	for i, node := range cluster.Member.Members() {
		nodes[i] = node.Name
	}
	sort.Sort(sort.StringSlice(nodes))
	return nodes
}
