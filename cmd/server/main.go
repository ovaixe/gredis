package main

import (
	"github.com/ovaixe/gredis/internal/server"
)

const PORT = 5262

func main() {
	srv := server.NewServer(PORT)

	srv.Serve()
}
