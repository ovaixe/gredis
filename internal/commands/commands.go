package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ovaixe/gredis/internal/storage"
)

// ExecuteCommand function parses and executes a given command on a key-value store.
// It supports commands such as SET, GET, and DEL to store, retrieve, and delete values associated with keys.
func ExecuteCommand(cmd string, store *storage.Storage) (string, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", errors.New("Invalid Command")
	}
	
	command := parts[0]
	switch command {
		case "SET":
			if len(parts) < 3 || len(parts) > 4 {
				return "", errors.New("Usage: SET [key] [value] [TTL]")
			}
			
			key := parts[1]
			value := parts[2]
			var ttl time.Duration
			if len(parts) == 4 {
				ttlSeconds, err := strconv.Atoi(parts[3])
				if err != nil {
					return "", errors.New("Invalid TTL value")
				}
				
				ttl = time.Duration(ttlSeconds) * time.Second
			}
			store.Set(key, value, ttl)
			return "OK", nil
		case "GET":
			if len(parts) != 2 {
				return "", errors.New("Usage: GET [key]")
			}
			value, found := store.Get(parts[1])
			if !found {
				return "(nil)", nil
			}
			return value, nil
		case "DEL":
			if len(parts) != 2 {
				return "", errors.New("Usage: DEL [key]")
			}
			err := store.Delete(parts[1])
			if err != nil {
				return "", err
			}
			return "OK", nil
		default:
			return "", errors.New("Unknown Command")
	}
}