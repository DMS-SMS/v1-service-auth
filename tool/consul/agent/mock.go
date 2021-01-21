// add package in v.1.1.6
// clone from tool/consul/agent in club

package agent

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/stretchr/testify/mock"
)

type _mock struct {
	mock *mock.Mock
}

func Mock(mock *mock.Mock) _mock {
	return _mock{mock: mock}
}

func (m _mock) ChangeAllServiceNodes() { // add in v.1.1.6
	m.mock.Called()
}

func (m _mock) GetNextServiceNode(service string) (*registry.Node, error) {
	args := m.mock.Called(service)
	return args.Get(0).(*registry.Node), args.Error(1)
}
