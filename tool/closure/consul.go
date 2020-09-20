package closure

import (
	"errors"
	"github.com/micro/go-micro/v2/server"
	"net"
	"strconv"
	"strings"
)

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
	if ok {
		err = errors.New("unable to assert type to *net.UDPAddr")
		return
	}

	return
}
