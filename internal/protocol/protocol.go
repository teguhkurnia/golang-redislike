package protocol

import (
	"fmt"
	"strings"

	"github.com/teguhkurnia/redis-like/internal/protocol/commands"
	"github.com/teguhkurnia/redis-like/internal/store"
)

var commandTable = make(map[string]*commands.CommandSpec)

func init() {
	// Connection commands
	commandTable["PING"] = commands.PingSpec
	commandTable["COMMAND"] = CommandHandlerSpec

	// String commands
	commandTable["GET"] = commands.GetSpec
	commandTable["SET"] = commands.SetSpec
	commandTable["DEL"] = commands.DelSpec

	// List commands
	commandTable["LPUSH"] = commands.LPushSpec
	commandTable["RPUSH"] = commands.RPushSpec
	commandTable["LRANGE"] = commands.LRangeSpec
	commandTable["LPOP"] = commands.LPopSpec
	commandTable["RPOP"] = commands.RPopSpec
}

func HandleCommand(cmd *commands.Command, store *store.Store) []byte {
	spec, found := commandTable[cmd.Name]
	if !found {
		return fmt.Appendf(nil, "-ERR unknown command '%s'\r\n", cmd.Name)
	}

	return spec.Handler(cmd, store)
}

// HandleCommand processes the COMMAND command, which introspects the server's command list.
func handleCommand(cmd *commands.Command, store *store.Store) []byte {
	if len(cmd.Args) > 0 {
		subcommand := strings.ToUpper(string(cmd.Args[0]))
		if subcommand == "DOCS" {
			if len(cmd.Args) < 2 {
				return []byte("-ERR wrong number of arguments for 'COMMAND DOCS'\r\n")
			}
			return BuildCommandDocs(cmd.Args[1:])
		}
		return []byte("-ERR Unimplemented subcommand for 'COMMAND'\r\n")
	}
	return BuildCommandInfo()
}

var CommandHandlerSpec = &commands.CommandSpec{
	Handler:  handleCommand,
	Arity:    -1, // Variable arguments
	Flags:    []string{"readonly"},
	FirstKey: 0,
	LastKey:  0,
	KeyStep:  0,
	Documentation: map[string]interface{}{
		"summary": "Introspects the server's command list.",
	},
}

// BuildCommandInfo dynamically creates the response for the COMMAND command.
func BuildCommandInfo() []byte {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("*%d\r\n", len(commandTable)))

	for name, spec := range commandTable {
		b.WriteString("*6\r\n") // 6 elements per command spec
		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(name), strings.ToLower(name)))
		b.WriteString(fmt.Sprintf(":%d\r\n", spec.Arity))
		b.WriteString(fmt.Sprintf("*%d\r\n", len(spec.Flags)))
		for _, flag := range spec.Flags {
			b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(flag), flag))
		}
		b.WriteString(fmt.Sprintf(":%d\r\n", spec.FirstKey))
		b.WriteString(fmt.Sprintf(":%d\r\n", spec.LastKey))
		b.WriteString(fmt.Sprintf(":%d\r\n", spec.KeyStep))
	}
	return []byte(b.String())
}

// BuildCommandDocs dynamically creates the response for COMMAND DOCS.
func BuildCommandDocs(commandNames [][]byte) []byte {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("*%d\r\n", len(commandNames)*2))

	for _, cmdNameBytes := range commandNames {
		cmdName := strings.ToUpper(string(cmdNameBytes))
		spec, ok := commandTable[cmdName]

		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(cmdName), strings.ToLower(cmdName)))

		if !ok || spec.Documentation == nil {
			b.WriteString("$-1\r\n")
			continue
		}
		// ... (Implementation to build RESP map from spec.Documentation) ...
		// This part can be complex, for now we can return a simple map
		doc := spec.Documentation
		b.WriteString(fmt.Sprintf("*%d\r\n", len(doc)*2))
		summary, _ := doc["summary"].(string)
		b.WriteString("$7\r\nsummary\r\n")
		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(summary), summary))
	}
	return []byte(b.String())
}
