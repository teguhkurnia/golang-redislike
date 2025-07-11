package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/teguhkurnia/redis-like/internal/protocol/commands"
)

type RESPValue struct {
	Type  byte
	Str   string
	Bulk  []byte
	Array []RESPValue
}

func ParseNextValue(reader *bufio.Reader) (RESPValue, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return RESPValue{}, err
	}
	if len(line) < 2 {
		return RESPValue{}, fmt.Errorf("invalid RESP line")
	}
	line = line[:len(line)-2] // Trim trailing \r\n

	valueType := line[0]
	valueContent := line[1:]

	switch valueType {
	case '+': // Simple String
		str := string(valueContent)
		// Handle quoted strings
		if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
			// Remove quotes and handle escape sequences
			unquoted := str[1 : len(str)-1]
			unquoted = strings.ReplaceAll(unquoted, `\"`, `"`)
			unquoted = strings.ReplaceAll(unquoted, `\\`, `\`)
			return RESPValue{Type: '+', Str: unquoted}, nil
		}
		return RESPValue{Type: '+', Str: str}, nil
	case '$': // Bulk String
		length, err := strconv.Atoi(string(valueContent))
		if err != nil {
			return RESPValue{}, fmt.Errorf("invalid bulk string length: %s", valueContent)
		}
		if length == -1 {
			return RESPValue{Type: '$'}, nil // Null Bulk String
		}
		buf := make([]byte, length+2)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return RESPValue{}, err
		}
		return RESPValue{Type: '$', Bulk: buf[:length]}, nil
	case '*': // Array
		count, err := strconv.Atoi(string(valueContent))
		if err != nil {
			return RESPValue{}, fmt.Errorf("invalid array count: %s", valueContent)
		}
		array := make([]RESPValue, count)
		for i := 0; i < count; i++ {
			val, err := ParseNextValue(reader)
			if err != nil {
				return RESPValue{}, err
			}
			array[i] = val
		}
		return RESPValue{Type: '*', Array: array}, nil
	}
	return RESPValue{}, fmt.Errorf("unsupported RESP type: %c", valueType)
}

// ToCommand converts a RESPValue into our Command struct
func (v RESPValue) ToCommand() (*commands.Command, error) {
	if v.Type != '*' {
		return nil, fmt.Errorf("invalid command: not a RESP array")
	}
	if len(v.Array) == 0 {
		return nil, fmt.Errorf("invalid command: empty array")
	}
	
	// Get command name from first element
	var cmdName string
	switch v.Array[0].Type {
	case '$':
		cmdName = string(v.Array[0].Bulk)
	case '+':
		cmdName = v.Array[0].Str
	default:
		return nil, fmt.Errorf("invalid command name type: %c", v.Array[0].Type)
	}
	
	cmd := &commands.Command{
		Name: strings.ToUpper(cmdName),
	}
	
	// Process arguments
	for i := 1; i < len(v.Array); i++ {
		var arg []byte
		switch v.Array[i].Type {
		case '$':
			arg = v.Array[i].Bulk
		case '+':
			arg = []byte(v.Array[i].Str)
		default:
			return nil, fmt.Errorf("invalid argument type: %c", v.Array[i].Type)
		}
		cmd.Args = append(cmd.Args, arg)
	}
	return cmd, nil
}
