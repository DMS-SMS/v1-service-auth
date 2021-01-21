// Add package in v.1.1.6
// network package in tool dir is used for utility about network like get available port, get my IP, etc ...
// get_port.go is file to declare various function, not method, about getting port

package network

import (
	"fmt"
	"math/rand"
	"net"
)

func GetRandomPortNotInUsedWithRange(min, max int) (port int) {
	for {
		port = rand.Intn(max - min) + min
		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}
		_ = conn.Close()
		break
	}
	return
}
