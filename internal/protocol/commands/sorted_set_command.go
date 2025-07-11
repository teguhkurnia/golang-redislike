package commands

import (
	"fmt"

	"strconv"

	"github.com/teguhkurnia/redis-like/internal/store"
)

func handleZAdd(cmd *Command, s *store.Store) []byte {
	if len(cmd.Args) < 3 || (len(cmd.Args)-1)%2 != 0 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := cmd.Args[0]
	members := cmd.Args[1:]
	membersStruct := make([]store.SortedSet, 0, len(members)/2)
	for i := 0; i < len(members); i += 2 {
		score, err := strconv.ParseFloat(string(members[i]), 64)
		if err != nil {
			return fmt.Appendf(nil, "-ERR score is not a valid number: %s\r\n", members[i])
		}
		member := string(members[i+1])
		membersStruct = append(membersStruct, store.SortedSet{
			Score:  score,
			Member: member,
		})
	}

	count := s.ZAdd(string(key), membersStruct)

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var ZAddSpec = &CommandSpec{
	Handler:  handleZAdd,
	Arity:    -3,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Adds one or more members to a sorted set.",
	},
}

func handleZRange(cmd *Command, s *store.Store) []byte {
	if len(cmd.Args) < 3 || len(cmd.Args) > 4 {
		return fmt.Appendf(nil, "-ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := cmd.Args[0]
	start, err := strconv.Atoi(string(cmd.Args[1]))
	if err != nil {
		return fmt.Appendf(nil, "-ERR start index is not a valid integer: %s\r\n", cmd.Args[1])
	}
	end := -1
	if len(cmd.Args) > 2 {
		end, err = strconv.Atoi(string(cmd.Args[2]))
		if err != nil {
			return fmt.Appendf(nil, "-ERR end index is not a valid integer: %s\r\n", cmd.Args[2])
		}
	}

	withScores := false
	if len(cmd.Args) == 4 && string(cmd.Args[3]) == "WITHSCORES" {
		withScores = true
	}

	members := s.ZRange(string(key), start, end)
	results := ""
	if withScores {
		results += fmt.Sprintf("*%d\r\n", len(members)*2)
	} else {
		results += fmt.Sprintf("*%d\r\n", len(members))
	}

	for _, member := range members {
		if withScores {
			scoreStr := fmt.Sprintf("%.17g", member.Score)
			results += fmt.Sprintf("$%d\r\n%s\r\n", len(member.Member), member.Member)
			results += fmt.Sprintf("$%d\r\n%s\r\n", len(scoreStr), scoreStr)
		} else {
			results += fmt.Sprintf("$%d\r\n%s\r\n", len(member.Member), member.Member)
		}
	}

	return []byte(results)
}

var ZRangeSpec = &CommandSpec{
	Handler:  handleZRange,
	Arity:    -4,
	Flags:    []string{"readonly"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Returns a range of members in a sorted set.",
	},
}

func handleZRem(cmd *Command, s *store.Store) []byte {
	if len(cmd.Args) < 2 {
		return fmt.Appendf(nil, "ERR wrong number of arguments for '%s' command\r\n", cmd.Name)
	}

	key := cmd.Args[0]
	members := cmd.Args[1:]
	membersStr := make([]string, len(members))
	for i, member := range members {
		membersStr[i] = string(member)
	}

	count := s.ZRem(string(key), membersStr)

	return fmt.Appendf(nil, ":%d\r\n", count)
}

var ZRemSpec = &CommandSpec{
	Handler:  handleZRem,
	Arity:    -3,
	Flags:    []string{"write"},
	FirstKey: 1,
	LastKey:  1,
	KeyStep:  1,
	Documentation: map[string]any{
		"summary": "Removes one or more members from a sorted set.",
	},
}
