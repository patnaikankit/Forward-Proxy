package proxy

import (
	"net"

	"github.com/patnaikankit/Forward-Proxy/utils"
)

func HandleHTTP(conn net.Conn, req *utils.HTTPRequest) {
	defer conn.Close()
}
