package commands

import (
	"fmt"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handleExists(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 1 {
		return []byte("-ERR wrong number of arguments for 'EXISTS' command\r\n")
	}

	existsCount := 0
	for _, arg := range cmd.Args {
		key := string(arg)
		if store.Exists(key) {
			existsCount++
		}
	}

	return fmt.Appendf(nil, ":%d\r\n", existsCount)
}

var ExistsSpec = &CommandSpec{
	Handler:  handleExists,
	Arity:    -2, // Arity is -2 because it expects at least one key and can take multiple keys
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  -1, // Last key is variable
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Checks if a key exists.",
	},
}
