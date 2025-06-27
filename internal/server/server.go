package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/ovaixe/gredis/internal/aof"
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
	listner net.Listener
	storage *storage.Storage
	logger  *log.Logger
	aof     *aof.Aof
}

// NewServer initialized a new server instance
func NewServer(port int) *Server {
	return &Server{
		config:  config{port: port},
		storage: storage.NewStorage(),
		logger:  log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Start creates Listener for incomming connections on the sepecified port
func (s *Server) Start() {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.port))
	if err != nil {
		s.logger.Fatalf("Failed to start the server: %v\n", err)
	}

	aof, err := aof.NewAof("database.aof")
	if err != nil {
		s.logger.Fatalf("Failed to create AOF: %v\n", err)
	}

	defer listner.Close()
	defer aof.Close()

	s.listner = listner
	s.aof = aof

	s.aof.Read(func(value resp.Value) {
		commands.ExecuteCommand(value, s.storage)
	})

	s.serve()

	s.logger.Printf("Server started on port: %v\n", s.config.port)
}

// serve starts accepting incomming connections
func (s *Server) serve() {
	// Accept and handle incomming connections
	for {
		conn, err := s.listner.Accept()
		if err != nil {
			s.logger.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		// Handle each connecton in a new go routine for concurrency
		go s.handleConnection(conn)

	}
}

// HandleConnection processes each client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := resp.NewReader(conn)
	writer := resp.NewWriter(conn)

	for {
		// Read incomming command
		cmd, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				s.logger.Printf("Client disconnected: %v\n", conn.RemoteAddr())
				return
			}

			s.logger.Printf("Error reading command: %v\n", err)
			return
		}

		if cmd.Typ != "array" {
			s.logger.Println("Invalid request, expected array")
			continue
		}

		if len(cmd.Array) == 0 {
			s.logger.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(cmd.Array[0].Bulk)
		if command == "SET" || command == "HSET" || command == "DEL" || command == "HDEL" || command == "HDELALL" {
			err := s.aof.Write(cmd)
			if err != nil {
				s.logger.Printf("Failed to write to the aof: %v\n", err)
			}
		}

		// Execute the command
		result := commands.ExecuteCommand(cmd, s.storage)

		// Send response to client
		writer.Write(result)
	}
}
