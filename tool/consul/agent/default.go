// add package in v.1.1.6
// clone from tool/consul/agent in club

package agent

import (
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"sync"
)

type _default struct {
	Strategy  selector.Strategy
	client    *api.Client
//  next      selector.Next                    // before v.1.1.6
//  nodes     []*registry.Node                 // before v.1.1.6
	next      map[ServiceName]selector.Next    // change in v.1.1.6
	nodes     map[ServiceName][]*registry.Node // change in v.1.1.6
	services  []ServiceName                    // add in v.1.1.6
	nodeMutex sync.RWMutex                     // add in v.1.1.6
}

func Default(setters ...FieldSetter) *_default {
	return newDefault(setters...)
}

func newDefault(setters ...FieldSetter) (h *_default) {
	h = new(_default)
	for _, setter := range setters {
		setter(h)
	}
	h.nodes = map[ServiceName][]*registry.Node{}
	h.nodeMutex = sync.RWMutex{}
	return
}

type FieldSetter func(*_default)

func Client(c *api.Client) FieldSetter {
	return func(d *_default) {
		d.client = c
	}
}

func Strategy(s selector.Strategy) FieldSetter {
	return func(d *_default) {
		d.Strategy = s
	}
}

func Services(s []ServiceName) FieldSetter {
	return func(d *_default) {
		d.services = s
	}
}
