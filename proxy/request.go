package proxy

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strings"

	"github.com/patnaikankit/Forward-Proxy/utils"
)

func ParseRequest(conn net.Conn) (*utils.HTTPRequest, error) {
	reader := bufio.NewReader(conn)
	var buf bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		buf.Write(line)
		if bytes.Equal(line, []byte("\r\n")) {
			break
		}
	}

	raw := buf.Bytes()
	lines := strings.Split(string(raw), "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("invalid request")
	}

	parts := strings.Split(lines[0], " ")
	if len(parts) < 3 {
		return nil, errors.New("malformed request line")
	}

	req := &utils.HTTPRequest{
		Method:  parts[0],
		URL:     parts[1],
		Version: parts[2],
		Raw:     raw,
		IsHTTPS: parts[0] == "CONNECT",
	}

	host := ""
	port := "80"
	if req.IsHTTPS {
		port = "443"
		hostPort := strings.Split(parts[1], ":")
		host = hostPort[0]
		if len(hostPort) > 1 {
			port = hostPort[1]
		}
	} else {
		for _, h := range lines[1:] {
			if strings.HasPrefix(strings.ToLower(h), "host:") {
				host = strings.TrimSpace(strings.Split(h, ":")[1])
				break
			}
		}
	}

	req.Host = host
	req.Port = port
	return req, nil
}
