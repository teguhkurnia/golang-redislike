package commands

import (
	"fmt"
	"strconv"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handleExpire(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	key := string(cmd.Args[0])
	seconds, err := strconv.Atoi(string(cmd.Args[1]))
	if err != nil || seconds < 0 {
		return fmt.Appendf(nil, "-ERR invalid expire time in '%s' command\r\n", cmd.Name)
	}

	store.Expire(key, seconds)
	return []byte("+OK\r\n")
}

var ExpireSpec = &CommandSpec{
	Handler:  handleExpire,
	Arity:    2, // Arity is 2 because it expects a key and a TTL
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Sets the expiration time of a key.",
	},
}

func handleTTL(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	key := string(cmd.Args[0])

	ttl := store.TTL(key)
	return fmt.Appendf(nil, ":%d\r\n", ttl)
}

var TTLSpec = &CommandSpec{
	Handler:  handleTTL,
	Arity:    2, // Arity is  1 because it expects a key
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Returns the time (seconds) to live for a key.",
	},
}
