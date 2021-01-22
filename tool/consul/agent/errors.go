// add package in v.1.1.6
// clone from tool/consul/agent in club

package agent

import "errors"

var (
	ErrAvailableNodeNotFound = errors.New("there is no currently available service node")
	ErrUnavailableService = errors.New("unavailable service, please put in agent.Services if you want to use")
)
