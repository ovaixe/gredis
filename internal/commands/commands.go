package commands

import (
	"errors"
	"strings"
	
	"github.com/ovaixe/gredis/internal/storage"
)

func ExecuteCommand(cmd string, store *storage.Storage) (string, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", errors.New("Invalid Command")
	}
	
	command := parts[0]
	switch command {
		case "SET":
			if len(parts) != 3 {
				return "", errors.New("Usage: SET [key] [value]")
			}
			store.Set(parts[1], parts[2])
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