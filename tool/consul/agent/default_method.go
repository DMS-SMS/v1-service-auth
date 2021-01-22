// Add file in v.1.0.6
// default_method.go is file to declare method of default struct

package agent

import (
	"auth/tool/consul"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/registry"
	"reflect"
)

const StatusMustBePassing = "Status==passing"


// call changeServiceNode once with mutex.Lock & Unlock
func (d *_default) ChangeServiceNodes(service consul.ServiceName) (err error) {
	d.nodeMutex.Lock()
	defer d.nodeMutex.Unlock()

	err = d.changeServiceNodes(service)
	return
}

// private method to handle business logic of changing specific service node list
func (d *_default) changeServiceNodes(service consul.ServiceName) error {
	checks, _, err := d.client.Health().Checks(string(service), &api.QueryOptions{Filter: StatusMustBePassing})
	if err != nil {
		return errors.New(fmt.Sprintf("unable to query health checkes, err: %v", err))
	}

	var nodes []*registry.Node
	for _, check := range checks {
		as, _, err := d.client.Agent().Service(check.ServiceID, nil)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to query service, err: %v", err))
		}
		var md = map[string]string{"CheckID": check.CheckID}
		node := &registry.Node{Id: as.ID, Address: fmt.Sprintf("%s:%d", as.Address, as.Port), Metadata: md}
		nodes = append(nodes, node)
	}

	if !reflect.DeepEqual(d.nodes[service], nodes) {
		d.nodes[service] = nodes
		d.next[service] = d.Strategy([]*registry.Service{{Nodes: nodes}})
	}

	return nil
}

// move from agent/default.go to agent/default_method.go
func (d *_default) GetNextServiceNode(service string) (*registry.Node, error) {
	checks, _, err := d.client.Health().Checks(service, &api.QueryOptions{Filter: StatusMustBePassing})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to query health checkes, err: %v", err))
	}

	if len(checks) == 0 {
		return nil, ErrAvailableNodeNotFound
	}

	var nodes []*registry.Node
	for _, check := range checks {
		as, _, err := d.client.Agent().Service(check.ServiceID, nil)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to query service, err: %v", err))
		}
		var md = map[string]string{"CheckID": check.CheckID}
		node := &registry.Node{Id: as.ID, Address: fmt.Sprintf("%s:%d", as.Address, as.Port), Metadata: md}
		nodes = append(nodes, node)
	}

	if !reflect.DeepEqual(d.nodes, nodes) {
		d.nodes = nodes
		d.next = d.Strategy([]*registry.Service{{Nodes: nodes}})
	}

	selectedNode, err := d.next()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to select node in selector, err: %v", err))
	}

	return selectedNode, nil
}
