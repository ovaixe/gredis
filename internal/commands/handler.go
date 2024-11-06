package commands

import (
	"github.com/ovaixe/gredis/internal/resp"
)

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
}

func ping(args []resp.Value) resp.Value {
	return resp.Value{Typ: "string", Str: "PONG"}
}
