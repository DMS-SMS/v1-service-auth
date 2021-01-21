// add package in v.1.1.6
// clone from tool/consul/agent in club

package consul

import "github.com/micro/go-micro/v2/registry"

type ServiceName string

type Agent interface {
	// method to refresh all service node list
	ChangeAllServiceNodes() error         // add in v.1.1.6
	// method to refresh specific service node list
	ChangeServiceNodes(ServiceName) error // add in v.1.1.6
	GetNextServiceNode(ServiceName) (*registry.Node, error)
}
