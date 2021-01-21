// add package in v.1.1.6
// clone from tool/consul/agent in club

package consul

import "github.com/micro/go-micro/v2/registry"

type Agent interface {
	// method to refresh all service node list
	ChangeAllServiceNodes()  // add in v.1.1.6
	GetNextServiceNode(service string) (*registry.Node, error)
}
