package commands

import (
	"fmt"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handleSAdd(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command", cmd.Name)
	}

	key := cmd.Args[0]
	members := cmd.Args[1:]
	membersStr := make([]string, len(members))
	for i, member := range members {
		membersStr[i] = string(member)
	}

	count := store.SAdd(string(key), membersStr)

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var SAddSpec = &CommandSpec{
	Handler:  handleSAdd,
	Arity:    -3,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Adds one or more members to a set.",
	},
}

func handleSRem(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) < 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command", cmd.Name)
	}

	key := cmd.Args[0]
	members := cmd.Args[1:]
	membersStr := make([]string, len(members))
	for i, member := range members {
		membersStr[i] = string(member)
	}

	count := store.SRem(string(key), membersStr)

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var SRemSpec = &CommandSpec{
	Handler:  handleSRem,
	Arity:    -3,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Removes one or more members from a set.",
	},
}

func handleSMembers(cmd *Command, store *store.Store) []byte {
	fmt.Printf("Handling SMembers command with args: %v\n", cmd.Args)
	if len(cmd.Args) != 1 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := cmd.Args[0]
	members := store.SMembers(string(key))

	result := fmt.Sprintf("*%d\r\n", len(members))
	for _, member := range members {
		result += fmt.Sprintf("$%d\r\n%s\r\n", len(member), member)
	}

	return []byte(result)
}

var SMembersSpec = &CommandSpec{
	Handler:  handleSMembers,
	Arity:    2,
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Returns all members of the set stored at key.",
	},
}

func handleSIsMember(cmd *Command, store *store.Store) []byte {
	if len(cmd.Args) != 2 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := cmd.Args[0]
	member := cmd.Args[1]

	isMember := store.SIsMember(string(key), string(member))

	if isMember {
		return []byte(":1\r\n")
	}
	return []byte(":0\r\n")
}

var SIsMemberSpec = &CommandSpec{
	Handler:  handleSIsMember,
	Arity:    3,
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Checks if a member is a member of a set.",
	},
}
