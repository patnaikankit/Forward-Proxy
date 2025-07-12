package proxy

import (
	"net"

	"github.com/patnaikankit/Forward-Proxy/utils"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseRequest(conn)
	if err != nil {
		return
	}

	utils.LogRequest(req)

	if req.IsHTTPS {
		HandleHTTPS(conn, req)
	} else {
		HandleHTTP(conn, req)
	}
}
