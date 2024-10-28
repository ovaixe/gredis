package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ovaixe/gredis/internal/server"
)

const PORT = "6252"

func main() {
	srv := server.NewServer()
	
	// Listen for incomming connections on the specified port
	listner, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	
	defer listner.Close()
	
	fmt.Printf("Server started on port: %v\n", PORT)
	
	// Accept and handle incomming connections
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		
		// Handle each connecton in a new go routine for concurrency
		go srv.HandleConnection(conn)
		
	}
}