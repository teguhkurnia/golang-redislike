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
	cmd := &commands.Command{
		Name: strings.ToUpper(string(v.Array[0].Bulk)),
	}
	for i := 1; i < len(v.Array); i++ {
		cmd.Args = append(cmd.Args, v.Array[i].Bulk)
	}
	return cmd, nil
}
