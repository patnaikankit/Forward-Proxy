package console

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patnaikankit/Forward-Proxy.git/utils"
)

func ConsoleRunner() {
	fmt.Println("Proxy Console - Press Ctrl+C to exit!")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Handle Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nExiting Proxy...")
		os.Exit(0)
	}()

	for range ticker.C {
		fmt.Println("---- Recent Requests ----")
		for _, entry := range utils.GetLogs() {
			fmt.Println(entry)
		}
		fmt.Println("---- End of Requests ----")
	}
}
