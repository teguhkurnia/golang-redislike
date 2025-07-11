package commands

import (
	"fmt"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handlePing(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) == 1 {
		return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(cmd.Args[0]), cmd.Args[0])
	}
	return []byte("+PONG\r\n")
}

var PingSpec = &CommandSpec{
	Handler:  handlePing,
	Arity:    -1,
	Flags:    []string{"readonly", "fast"},
	FirstKey: 0,
	LastKey:  0,
	KeyStep:  0,
	Documentation: map[string]interface{}{
		"summary": "Returns PONG, or the argument if provided.",
	},
}
