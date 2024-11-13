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
		return ping(args)
	case "GET":
		return get(args)
	case "SET":
		return set(args)
	case "DEL":
		return del(args)
	case "HGET":
		return hget(args)
	case "HGETALL":
		return hgetAll(args)
	case "HSET":
		return hset(args)
	case "HDEL":
		return hdel(args)
	case "HDELALL":
		return hdelAll(args)
	default:
		return resp.Value{Typ: "error", Str: "Unknown Command"}
	}
}

func ping(args resp.Value) resp.Value {
		if len(args) == 0 {
			return resp.Value{Typ: "string", Str: "PONG"}
		}

		return resp.Value{Typ: "string", Str: args[0].Bulk}
}

func get(args resp.Value) resp.Value {
		if len(args) != 1 {
			return resp.Value{Typ: "error", Str: "Usage: GET [key]"}
		}

		value, found := store.Get(args[0].Bulk)
		if !found {
			return resp.Value{Typ: "null"}
		}

		return resp.Value{Typ: "bulk", Bulk: value}
}

func set(args resp.Value) resp.Value {
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
}

func del(args resp.Value) resp.Value {
		if len(args) != 1 {
			return resp.Value{Typ: "error", Str: "Usage: DEL [key]"}
		}

		err := store.Delete(args[0].Bulk)
		if err != nil {
			return resp.Value{Typ: "error", Str: err.Error()}
		}

		return resp.Value{Typ: "string", Str: "OK"}
}

func hget(args resp.Value) resp.Value {
		if len(args) != 2 {
			return resp.Value{Typ: "error", Str: "Usage: HGET [hash] [key]"}
		}

		hash := args[0].Bulk
		key := args[1].Bulk

		value, found := store.HGet(hash, key)
		if !found {
			return resp.Value{Typ: "null"}
		}

		return resp.Value{Typ: "bulk", Bulk: value}
}

func hgetAll(args resp.Value) resp.Value {
		if len(args) != 1 {
			return resp.Value{Typ: "error", Str: "Usage: HGETALL [hash]"}
		}

		hash := args[0].Bulk

		value, found := store.HGetAll(hash)
		if !found {
			return resp.Value{Typ: "null"}
		}

		result := []resp.Value{}

		for field, val := range value {
			result = append(result, resp.Value{Typ: "string", Str: field})
			result = append(result, resp.Value{Typ: "string", Str: val})
		}

		return resp.Value{Typ: "array", Array: result}
}

func hset(args resp.Value) resp.Value {
		if len(args) != 3 {
			return resp.Value{Typ: "error", Str: "Usage: HSET [hash] [key] [value]"}
		}

		hash := args[0].Bulk
		key := args[1].Bulk
		value := args[2].Bulk

		store.HSet(hash, key, value)
		return resp.Value{Typ: "string", Str: "OK"}
}

func hdel(args resp.Value) resp.Value {
		if len(args) != 2 {
			return resp.Value{Typ: "error", Str: "Usage: HDEL [hash] [key]"}
		}

		hash := args[0].Bulk
		key := args[1].Bulk

		err := store.HDelete(hash, key)
		if err != nil {
			return resp.Value{Typ: "error", Str: err.Error()}
		}

		return resp.Value{Typ: "string", Str: "OK"}
}

func hdelAll(args resp.Value) resp.Value {
		if len(args) != 1 {

			return resp.Value{Typ: "error", Str: "Usage: HDELALL [hash]"}
		}

		hash := args[0].Bulk

		err := store.HDeleteAll(hash)
		if err != nil {
			return resp.Value{Typ: "error", Str: err.Error()}
		}

		return resp.Value{Typ: "string", Str: "OK"}
}
