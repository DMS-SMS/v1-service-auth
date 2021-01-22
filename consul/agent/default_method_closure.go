// Add file in v.1.1.6
// default_method_closure.go is file for declaring method to return closure of default struct

package agent

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/log"
	"net"
	"strconv"
	"strings"
)

// move from /tool/closure/consul.go in v.1.1.6
func (d *_default) ServiceNodeRegistry(s server.Server) func() error {
	return func() (err error) {
		port, err := getPortFromServerOption(s.Options())
		if err != nil {
			log.Fatalf("unable to get port number from server option, err: %v\n", err)
		}
		localAddr, err := getLocalIP()
		if err != nil {
			log.Fatalf("unable to get local address, err: %v\n", err)
		}

		srvID := fmt.Sprintf("%s-%s", s.Options().Name, s.Options().Id)
		err = d.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:      srvID,
			Name:    s.Options().Name,
			Port:    port,
			Address: localAddr,
		})
		if err != nil {
			log.Fatalf("unable to register service in consul, err: %v\n", err)
		}

		checkID := fmt.Sprintf("service:%s", srvID)
		checkName := fmt.Sprintf("service '%s' check", s.Options().Name)
		err = d.client.Agent().CheckRegister(&api.AgentCheckRegistration{
			ID:                checkID,
			Name:              checkName,
			ServiceID:         srvID,
			AgentServiceCheck: api.AgentServiceCheck{
				Name:   s.Options().Name,
				Status: "passing",
				TTL:    "8640h",
			},
		})
		if err != nil {
			log.Fatalf("unable to register check in consul, err: %v\n", err)
		}

		log.Infof("succeed to registry service and check to consul!! (service id: %s | checker id: %s)", srvID, checkID)
		return
	}
}

// get port number by parsing server.Options.Address
func getPortFromServerOption(opts server.Options) (port int, err error) {
	const portIndex = 3
	portStr := strings.Split(opts.Address, ":")[portIndex]
	port, err = strconv.Atoi(portStr)
	return
}

// get my local ip address to register node in consul
func getLocalIP() (addr string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return
}
