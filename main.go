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
	go proxy.ProxyInitialization()

	// Start console UI
	if utils.DebugMode {
		console.ConsoleRunner()
	} else {
		// Wait for termination
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Shutting down proxy...")
	}
}
