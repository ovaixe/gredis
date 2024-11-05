package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/ovaixe/gredis/internal/commands"
	"github.com/ovaixe/gredis/internal/storage"
)

// Define a config struct to hold all the configuration settings for server
type config struct {
	port int
}

// Server struct represents the Redis server
type Server struct {
	config config
	storage *storage.Storage
	logger *log.Logger
}

// NewServer initialized a new server instance
func NewServer (port int) *Server {
	return &Server{
		config: config{port: port},
		storage: storage.NewStorage(),
		logger: log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Serve creates Listener for incomming connections on the specified port
func (s *Server) Serve() {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.port))
		if err != nil {
			s.logger.Fatalf("Failed to start server: %v\n", err)
		}
		
		defer listner.Close()
		
		s.logger.Printf("Server started on port: %v\n", s.config.port)
		
		// Accept and handle incomming connections
		for {
			conn, err := listner.Accept()
			if err != nil {
				s.logger.Printf("Failed to accept connection: %v\n", err)
				continue
			}
			
			// Handle each connecton in a new go routine for concurrency
			go s.HandleConnection(conn)
			
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
			fmt.Fprintf(conn, "Error reading command: %v\n", err)
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