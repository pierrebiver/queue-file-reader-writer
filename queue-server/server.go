package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"queue-file-reader-writer.com/internal/command"
)

type Server struct {
	q        *Queue
	listener net.Listener
}

func NewServer() *Server {
	return &Server{q: &Queue{}}
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener
	log.Printf("queue server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(connection net.Conn) {
	defer connection.Close()
	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		line := scanner.Text()
		response := s.dispatch(line)
		fmt.Fprintf(connection, "%s\n", response)
	}
}

func (s *Server) dispatch(line string) string {
	for _, cmd := range command.Registry {
		if cmd.Is(line) {
			response, err := cmd.Execute(line, s.q)
			if err != nil {
				return "ERR " + err.Error()
			}
			return response
		}
	}
	return "ERR unknown command"
}
