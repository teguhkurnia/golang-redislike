package commands

import (
	"fmt"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handleHSet(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 3 || (len(cmd.Args)-1)%2 != 0 {
		return []byte("-ERR wrong number of arguments for 'HSET' command\r\n")
	}

	key := string(cmd.Args[0])
	count := 0

	for i := 1; i < len(cmd.Args); i += 2 {
		field := string(cmd.Args[i])
		value := string(cmd.Args[i+1])
		count += store.HSet(key, field, value)
	}

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var HSetSpec = &CommandSpec{
	Handler:  handleHSet,
	Arity:    -3, // Arity is -3 because it expects at least one key and pairs of field-value
	Flags:    []string{"write", "deny-oom"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Sets the value of a hash field.",
	},
}

func handleHGet(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 2 {
		return []byte("-ERR wrong number of arguments for 'HGET' command\r\n")
	}

	key := string(cmd.Args[0])
	field := string(cmd.Args[1])
	value, exists := store.HGet(key, field)
	if !exists {
		return []byte("$-1\r\n")
	}

	return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(value), value)
}

var HGetSpec = &CommandSpec{
	Handler:  handleHGet,
	Arity:    3, // Arity is 3 because it expects one key and one field
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Gets the value of a hash field.",
	},
}

func handleHGetAll(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 {
		return []byte("-ERR wrong number of arguments for 'HGETALL' command\r\n")
	}

	key := string(cmd.Args[0])
	hash, exists := store.HGetAll(key)
	if !exists {
		return []byte("*\r\n")
	}

	result := fmt.Sprintf("*%d\r\n", len(hash)*2)
	for field, value := range hash {
		result += fmt.Sprintf("$%d\r\n%s\r\n$%d\r\n%s\r\n", len(field), field, len(value), value)
	}
	return []byte(result)
}

var HGetAllSpec = &CommandSpec{
	Handler:  handleHGetAll,
	Arity:    2, // Arity is 2 because it expects one key
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Gets all fields and values in a hash.",
	},
}

func handleHDel(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 1 {
		return []byte("-ERR wrong number of arguments for 'HDEL' command\r\n")
	}

	deleted := 0
	key := string(cmd.Args[0])
	fields := cmd.Args[1:]
	for _, field := range fields {
		deleted += store.HDel(key, string(field))
	}

	return fmt.Appendf(nil, ":%d\r\n", deleted)
}

var HDelSpec = &CommandSpec{
	Handler:  handleHDel,
	Arity:    -2, // Arity is -2 because it expects at least one key and one field
	Flags:    []string{"write", "deny-oom"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Deletes a field from a hash.",
	},
}
