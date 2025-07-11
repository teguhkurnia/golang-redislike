package commands

import (
	"fmt"
	"strconv"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handleLPush(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := string(cmd.Args[0])
	values := cmd.Args[1:]
	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = string(v)
	}

	count := store.LPush(key, strValues)

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var LPushSpec = &CommandSpec{
	Handler:  handleLPush,
	Arity:    -3, // -3 means at least 2 arguments
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Adds one or more elements to the beginning of a list.",
	},
}

func handleRPush(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}
	key := string(cmd.Args[0])
	values := cmd.Args[1:]
	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = string(v)
	}
	count := store.RPush(key, strValues)

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var RPushSpec = &CommandSpec{
	Handler:  handleRPush,
	Arity:    -3, // -3 means at least 2 arguments
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Adds one or more elements to the end of a list.",
	},
}

func handleLRange(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 3 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := string(cmd.Args[0])
	start, err := strconv.Atoi(string(cmd.Args[1]))
	if err != nil {
		return fmt.Appendf(nil, "-ERR invalid start for '%s' command\r\n", cmd.Name)
	}
	end, err := strconv.Atoi(string(cmd.Args[2]))
	if err != nil {
		return fmt.Appendf(nil, "-ERR invalid end for '%s' command\r\n", cmd.Name)
	}

	values := store.LRange(key, start, end)

	if len(values) == 0 {
		return []byte("$-1\r\n")
	}

	response := fmt.Sprintf("*%d\r\n", len(values))
	for _, value := range values {
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
	}

	return []byte(response)
}

var LRangeSpec = &CommandSpec{
	Handler:  handleLRange,
	Arity:    3, // Exactly 3 arguments: key, start, end
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Retrieves a range of elements from a list.",
	},
}

func handleLPop(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 && len(cmd.Args) != 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := string(cmd.Args[0])
	count := 1
	if len(cmd.Args) == 2 {
		var err error
		count, err = strconv.Atoi(string(cmd.Args[1]))
		if err != nil {
			return fmt.Appendf(nil, "-ERR invalid count for '%s' command\r\n", cmd.Name)
		}
	}

	values := store.LPop(key, count)

	if len(values) == 0 {
		return []byte("$-1\r\n")
	}

	response := fmt.Sprintf("*%d\r\n", len(values))
	for _, value := range values {
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
	}

	return []byte(response)
}

var LPopSpec = &CommandSpec{
	Handler:  handleLPop,
	Arity:    -2, // -2 means at least 1 argument and 1 optional count
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Removes and returns the first element of a list.",
	},
}

func handleRPop(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 || len(cmd.Args) > 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := string(cmd.Args[0])
	count := 1
	if len(cmd.Args) == 2 {
		var err error
		count, err = strconv.Atoi(string(cmd.Args[1]))
		if err != nil {
			return fmt.Appendf(nil, "-ERR invalid count for '%s' command\r\n", cmd.Name)
		}
	}

	values := store.RPop(key, count)

	if len(values) == 0 {
		return []byte("$-1\r\n")
	}

	response := fmt.Sprintf("*%d\r\n", len(values))
	for _, value := range values {
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
	}

	return []byte(response)
}

var RPopSpec = &CommandSpec{
	Handler:  handleRPop,
	Arity:    -2, // -2 means at least 1 argument and 1 optional count
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Removes and returns the last element of a list.",
	},
}

func handleLLen(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 1 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := string(cmd.Args[0])
	length, ok := store.LLen(key)
	if !ok {
		return fmt.Appendf(nil, "-ERR WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
	}

	return fmt.Appendf(nil, ":%d\r\n", length)
}

var LLenSpec = &CommandSpec{
	Handler:  handleLLen,
	Arity:    2, // Exactly
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Gets the length of a list.",
	},
}
