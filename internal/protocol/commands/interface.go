package commands

import "github.com/teguhkurnia/redis-like/internal/store"

type Command struct {
	Name string
	Args [][]byte
}

type CommandSpec struct {
	Handler       func(cmd *Command, store *store.Store) []byte
	Documentation map[string]any
	Arity         int
	Flags         []string
	FirstKey      int
	LastKey       int
	KeyStep       int
}
