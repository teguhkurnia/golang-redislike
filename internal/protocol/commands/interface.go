package commands

import (
	"fmt"
	"strings"

	"github.com/teguhkurnia/redis-like/internal/store"
)

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

func FromLog(line string) (*Command, error) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid log line: %s", line)
	}
	cmd := &Command{
		Name: parts[0],
		Args: parseArgs(parts[1]),
	}
	return cmd, nil
}

func parseArgs(argsStr string) [][]byte {
	args := strings.Fields(argsStr)
	byteArgs := make([][]byte, len(args))
	for i, arg := range args {
		byteArgs[i] = []byte(arg)
	}
	return byteArgs
}

func (cmd *Command) ToLog() (string, error) {
	// A command must have a name.
	if strings.TrimSpace(cmd.Name) == "" {
		return "", fmt.Errorf("command name cannot be empty")
	}

	// Start with the command name.
	parts := []string{cmd.Name}

	// Append each argument, converting it from []byte to string.
	for _, arg := range cmd.Args {
		parts = append(parts, string(arg))
	}

	// Join all parts with a space to create the final log line.
	return strings.Join(parts, " "), nil
}
