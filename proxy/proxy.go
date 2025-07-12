package proxy

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func StartProxy() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("PROXY_HOST")
	port := os.Getenv("PROXY_PORT")

	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
	log.Printf("Proxy running on %s:%s\n", host, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
