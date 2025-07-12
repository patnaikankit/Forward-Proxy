package main

import (
	"github.com/patnaikankit/Forward-Proxy.git/console"
	"github.com/patnaikankit/Forward-Proxy.git/proxy"
)

func main() {
	go proxy.StartProxy()
	console.ConsoleRunner()
}
