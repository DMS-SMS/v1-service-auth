// add package in v.1.1.6
// clone from tool/consul/agent in club

package agent

import "errors"

var (
	ErrAvailableNodeNotFound = errors.New("there is no currently available services")
)
