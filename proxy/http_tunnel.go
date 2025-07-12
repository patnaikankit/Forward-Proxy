package proxy

import (
	"io"
	"net"

	"github.com/patnaikankit/Forward-Proxy/utils"
)

func HandleHTTPS(client net.Conn, req *utils.HTTPRequest) {
	if utils.IsBlocked(req.URL) {
		client.Write([]byte("HTTP/1.1 403 Forbidden\r\n\r\n"))
		client.Close()
		return
	}

	server, err := net.Dial("tcp", req.Host+":"+req.Port)
	if err != nil {
		client.Close()
		return
	}

	client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	go io.Copy(server, client)
	io.Copy(client, server)

	server.Close()
	client.Close()
}
