package proxy

import (
	"io"
	"net"

	"github.com/patnaikankit/Forward-Proxy/utils"
)

func HandleHTTP(conn net.Conn, req *utils.HTTPRequest) {
	defer conn.Close()

	if utils.IsBlocked(req.URL) {
		conn.Write([]byte("HTTP/1.1 403 Forbidden\r\n\r\n"))
		return
	}

	if cached := utils.GetCache(req.URL); cached != nil {
		conn.Write(cached)
		return
	}

	serverConn, err := net.Dial("tcp", req.Host+":"+req.Port)
	if err != nil {
		return
	}

	defer serverConn.Close()

	serverConn.Write(req.Raw)

	response, err := io.ReadAll(serverConn)
	if err != nil {
		return
	}

	utils.SetCache(req.URL, response)
	conn.Write(response)
}
