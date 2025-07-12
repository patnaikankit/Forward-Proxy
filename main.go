package main

import (
	"github.com/patnaikankit/Forward-Proxy/console"
	"github.com/patnaikankit/Forward-Proxy/proxy"
)

func main() {
	go proxy.StartProxy()
	console.ConsoleRunner()
}
