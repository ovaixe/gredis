package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/ovaixe/gredis/internal/resp"
	"github.com/ovaixe/gredis/internal/storage"
)

// ExecuteCommand function parses and executes a given command on a key-value store.
// It supports commands such as SET, GET, and DEL to store, retrieve, and delete values associated with keys.
func ExecuteCommand(cmd resp.Value, store *storage.Storage) resp.Value {
	command := strings.ToUpper(cmd.Array[0].Bulk)
	args := cmd.Array[1:]

	switch command {
	case "PING":
		if len(args) == 0 {
			return resp.Value{Typ: "string", Str: "PONG"}
		}

		return resp.Value{Typ: "string", Str: args[0].Bulk}
	case "SET":
		if len(args) < 2 || len(args) > 3 {
			return resp.Value{Typ: "error", Str: "Usage: SET [key] [value] [TTL]"}
		}

		key := args[0].Bulk
		value := args[1].Bulk
		var ttl time.Duration
		if len(args) == 3 {
			ttlSeconds, err := strconv.Atoi(args[2].Bulk)
			if err != nil {
				return resp.Value{Typ: "error", Str: "Invalid TTL value"}
			}

			ttl = time.Duration(ttlSeconds) * time.Second
		}

		store.Set(key, value, ttl)
		return resp.Value{Typ: "string", Str: "OK"}
	case "GET":
		if len(args) != 1 {
			return resp.Value{Typ: "error", Str: "Usage: GET [key]"}
		}

		value, found := store.Get(args[0].Bulk)
		if !found {
			return resp.Value{Typ: "null"}
		}

		return resp.Value{Typ: "bulk", Bulk: value}
	case "DEL":
		if len(args) != 1 {
			return resp.Value{Typ: "error", Str: "Usage: DEL [key]"}
		}

		err := store.Delete(args[0].Bulk)
		if err != nil {
			return resp.Value{Typ: "error", Str: err.Error()}
		}

		return resp.Value{Typ: "string", Str: "OK"}
	default:
		return resp.Value{Typ: "error", Str: "Unknown Command"}
	}
}
