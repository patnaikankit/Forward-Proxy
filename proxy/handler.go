package proxy

import "net"

func handleConnection(conn net.Conn) {
	defer conn.Close()
}
