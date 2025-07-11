package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/teguhkurnia/redis-like/internal/log"
	"github.com/teguhkurnia/redis-like/internal/protocol"
	"github.com/teguhkurnia/redis-like/internal/protocol/parser"
	"github.com/teguhkurnia/redis-like/internal/store"
)

type Message struct {
	From    net.Addr
	Payload []byte
}

type Server struct {
	ListenAddr string
	Store      *store.Store
	Log        *log.Log
	ln         net.Listener
	clients    map[string]net.Conn

	quitChan chan struct{}
	msgChan  chan *Message
}

func NewServer(listenAddr string, store *store.Store) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Store:      store,
		clients:    make(map[string]net.Conn),
		quitChan:   make(chan struct{}),
		msgChan:    make(chan *Message, 100),
		Log:        log.NewLog("server.log"),
	}
}

func (s *Server) Start() {
	// Initialize the store and log
	cmds, err := s.Log.LoadCommandsFromLog()
	if err != nil {
		panic(fmt.Sprintf("Failed to load commands from log: %v", err))
	}

	for _, cmd := range cmds {
		protocol.HandleCommand(cmd, s.Store, s.Log, true)
	}

	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		panic(err)
	}
	s.ln = ln
	go s.acceptLoop()

	fmt.Printf("🚀 Server started on %s\n", s.ListenAddr)
	<-s.quitChan
	s.ln.Close()
	close(s.msgChan)
	close(s.quitChan)
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			continue
		}
		fmt.Printf("💬 New connection from %s\n", conn.RemoteAddr().String())
		s.clients[conn.RemoteAddr().String()] = conn
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		value, err := parser.ParseNextValue(reader)
		if err != nil {
			continue
		}
		cmd, err := value.ToCommand()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("❌ Connection closed by %s\n", conn.RemoteAddr().String())
				conn.Close()
				delete(s.clients, conn.RemoteAddr().String())
				return
			}

			continue
		}

		response := protocol.HandleCommand(cmd, s.Store, s.Log, false)
		if response != nil {
			_, err := conn.Write(response)
			if err != nil {
				conn.Close()
				delete(s.clients, conn.RemoteAddr().String())
				return
			}
		}
	}
}
