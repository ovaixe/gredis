package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/ovaixe/gredis/internal/commands"
	"github.com/ovaixe/gredis/internal/resp"
	"github.com/ovaixe/gredis/internal/storage"
)

// Define a config struct to hold all the configuration settings for server
type config struct {
	port int
}

// Server struct represents the Redis server
type Server struct {
	config  config
	storage *storage.Storage
	logger  *log.Logger
}

// NewServer initialized a new server instance
func NewServer(port int) *Server {
	return &Server{
		config:  config{port: port},
		storage: storage.NewStorage(),
		logger:  log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile),
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
	reader := resp.NewResp(conn)

	for {
		// Read incomming command
		value, err := reader.Read()
		if err != nil {
			fmt.Printf("Error reading command: %v\n", err)
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// Parse the command
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		// Execute the command
		handler, ok := commands.Handlers[command]
		if !ok {
			writer.Write(resp.Value{Typ: "string", Str: "Invalid command"})
			continue
		}

		result := handler(args)

		// Send response to client
		writer.Write(result)
	}
}
