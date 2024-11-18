package main

import (
	"github.com/ovaixe/gredis/internal/server"
)

const PORT = 6379

func main() {
	srv := server.NewServer(PORT)

	srv.Start()
}
