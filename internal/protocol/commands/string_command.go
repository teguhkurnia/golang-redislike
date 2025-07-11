package commands

import (
	"fmt"
	"strconv"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handleGet(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	value, found := store.Get(string(cmd.Args[0]))
	if !found {
		return fmt.Appendf(nil, "$-1\r\n")
	}

	// parse to integer
	number, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(value), value)
	}

	// return integer response
	return fmt.Appendf(nil, ":%d\r\n", number)
}

var GetSpec = &CommandSpec{
	Handler:  handleGet,
	Arity:    2,
	Flags:    []string{"read"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]interface{}{
		"summary": "Gets the value of a key.",
	},
}

func handleSet(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	store.Set(string(cmd.Args[0]), string(cmd.Args[1]))
	return []byte("+OK\r\n")
}

var SetSpec = &CommandSpec{
	Handler:  handleSet,
	Arity:    3,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Sets the value of a key.",
	},
}

func handleDel(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 1 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	count := 0
	for _, key := range cmd.Args {
		deleted := store.Del(string(key))
		if deleted {
			count++
		}
	}

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var DelSpec = &CommandSpec{
	Handler:  handleDel,
	Arity:    -1,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  -1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Deletes a key.",
	},
}

func handleIncr(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	value, ok := store.Incr(string(cmd.Args[0]))
	if !ok {
		return fmt.Appendf(nil, "-ERR value is not an integer\r\n")
	}
	return fmt.Appendf(nil, ":%d\r\n", value)
}

var IncrSpec = &CommandSpec{
	Handler:  handleIncr,
	Arity:    2,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]interface{}{
		"summary": "Increments the integer value of a key by one.",
	},
}

func handleDecr(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	value, ok := store.Decr(string(cmd.Args[0]))
	if !ok {
		return fmt.Appendf(nil, "-ERR value is not an integer\r\n")
	}

	return fmt.Appendf(nil, ":%d\r\n", value)
}

var DecrSpec = &CommandSpec{
	Handler:  handleDecr,
	Arity:    2,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]interface{}{
		"summary": "Decrements the integer value of a key by one.",
	},
}
