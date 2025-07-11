package log

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/teguhkurnia/redis-like/internal/protocol/commands"
)

type Log struct {
	logMutex *sync.Mutex
	logFile  string
}

func NewLog(logFile string) *Log {
	// Ensure the log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err := os.Create(logFile)
		if err != nil {
			panic("Error creating log file: " + err.Error())
		}
		file.Close()
	}

	return &Log{
		logFile:  logFile,
		logMutex: &sync.Mutex{},
	}
}

func (Store *Log) StoreWriteCommandToLog(cmd *commands.Command) {
	logEntry, err := cmd.ToLog()
	if err != nil {
		fmt.Printf("Error converting command to log entry: %v\n", err)
		return
	}

	Store.appendToLogFile(logEntry)
}

func (Store *Log) appendToLogFile(logEntry string) {
	Store.logMutex.Lock()
	defer Store.logMutex.Unlock()
	file, err := os.OpenFile(Store.logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logEntry + "\n"); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
		return
	}
}

func (Store *Log) LoadCommandsFromLog() ([]*commands.Command, error) {
	file, err := os.Open(Store.logFile)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	var cmds []*commands.Command
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cmd, err := commands.FromLog(line)
		if err != nil {
			fmt.Printf("Error parsing command from log: %v\n", err)
			continue
		}
		cmds = append(cmds, cmd)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	return cmds, nil
}
