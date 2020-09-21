package closure

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"net"
	"strconv"
	"strings"
)

func ConsulServiceRegistrar(s server.Server, consul *api.Client) func() error {
	return func() (err error) {
		port, err := getPortFromServerOption(s.Options())
		if err != nil {
			log.Fatalf("unable to get port number from server option, err: %v\n", err)
		}
		localAddr, err := getMyLocalAddr()
		if err != nil {
			log.Fatalf("unable to get local address, err: %v\n", err)
		}

		srvID := fmt.Sprintf("%s-%s", s.Options().Name, s.Options().Id)
		err = consul.Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:      srvID,
			Name:    s.Options().Name,
			Port:    port,
			Address: localAddr.IP.String(),
		})
		if err != nil {
			log.Fatalf("unable to register service in consul, err: %v\n", err)
		}

		checkID := fmt.Sprintf("service:%s", srvID)
		checkName := fmt.Sprintf("service '%s' check", s.Options().Name)
		err = consul.Agent().CheckRegister(&api.AgentCheckRegistration{
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

func ConsulServiceDeregistrar(s server.Server, consul *api.Client) func() error {
	return func() (err error) {
		srvID := fmt.Sprintf("%s-%s", s.Options().Name, s.Options().Id)
		err = consul.Agent().ServiceDeregister(srvID)
		if err != nil {
			log.Fatalf("unable to deregister service in consul, err: %v\n", err)
		}

		checkID := fmt.Sprintf("service:%s", srvID)
		err = consul.Agent().CheckDeregister(checkID)
		if err != nil {
			log.Fatalf("unable to deregister check in consul, err: %v\n", err)
		}

		log.Infof("succeed to deregistry service and check to consul!! (service id: %s | checker id: %s)", srvID, checkID)
		return
	}
}

func getPortFromServerOption(opts server.Options) (port int, err error) {
	const portIndex = 3
	portStr := strings.Split(opts.Address, ":")[portIndex]
	port, err = strconv.Atoi(portStr)
	return
}

func getMyLocalAddr() (addr *net.UDPAddr, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		err = errors.New("unable to assert type to *net.UDPAddr")
		return
	}

	return
}
