// add package in v.1.1.6
// clone from tool/consul/agent in club

package consul

import "github.com/micro/go-micro/v2/registry"

type Agent interface {
	GetNextServiceNode(service string) (*registry.Node, error)
}
