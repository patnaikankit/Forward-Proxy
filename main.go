package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/patnaikankit/Forward-Proxy/console"
	"github.com/patnaikankit/Forward-Proxy/proxy"
	"github.com/patnaikankit/Forward-Proxy/utils"
)

func main() {
	// Start proxy server in background
	go proxy.StartProxy()

	// Start console UI
	if utils.DebugMode {
		console.StartConsole()
	} else {
		// Wait for termination (useful for debug mode)
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Shutting down proxy...")
	}
}
