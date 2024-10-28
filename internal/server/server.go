package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/ovaixe/gredis/internal/commands"
	"github.com/ovaixe/gredis/internal/storage"
)

// Server struct represents the Redis server
type Server struct {
	storage *storage.Storage
}

// NewServer initialized a new server instance
func NewServer () *Server {
	return &Server{
		storage: storage.NewStorage(),
	}
}

// HandleConnection processes each client connection
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	
	for {
		// Read incomming command
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(conn, "Error reading command: %V\n", err)
		}
		
		// Parse and execute the command
		cmd = strings.TrimSpace(cmd)
		response, err := commands.ExecuteCommand(cmd, s.storage)
		if err != nil {
			fmt.Fprintf(conn, "[Error]: %v", err)
		}
		
		// Send response to client
		fmt.Fprintln(conn, response)
	}
}