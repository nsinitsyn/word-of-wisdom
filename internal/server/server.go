package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"word-of-wisdom/internal/config"
	"word-of-wisdom/pkg/challenge"
)

const (
	messageSize int16 = 1024
	complexity  uint8 = 3
)

type Repository interface {
	GetQuote() string
}

type Server struct {
	listener net.Listener
	config   *config.ServerConfig
	repo     Repository
	wg       *sync.WaitGroup
}

func NewServer(config *config.ServerConfig, repo Repository) (*Server, error) {
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("error start server: %w", err)
	}

	return &Server{listener: listener, config: config, repo: repo, wg: &sync.WaitGroup{}}, nil
}

func (s Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			log.Println("listener closed")
			break
		}
		if err != nil {
			log.Printf("error listener accepting: %s\n", err)
			continue
		}
		s.wg.Add(1)
		go func(conn net.Conn) {
			defer func() {
				conn.Close()
				s.wg.Done()
			}()

			s.serve(conn)
		}(conn)
	}
	s.wg.Wait()
}

func (s Server) serve(conn net.Conn) {
	// read request
	_, err := s.readMessage(conn)
	if err != nil {
		log.Printf("error reading: %s\n", err)
		return
	}

	chal := challenge.NewRandomChallenge(s.config.Complexity)

	// send challenge
	conn.Write(append(chal.Input, chal.Target...))

	// read solution
	nonce, err := s.readMessage(conn)
	if err != nil {
		log.Printf("error reading: %s\n", err)
		return
	}

	// verify solution
	if !chal.VerifyNonce(nonce) {
		conn.Write([]byte("incorrect solution"))
		return
	}

	// send quote
	quote := s.repo.GetQuote()
	conn.Write([]byte(quote))
}

func (s Server) Stop() {
	log.Println("server closing...")
	if err := s.listener.Close(); err != nil {
		log.Printf("error server closing: %s\n", err)
	}
}

func (s Server) readMessage(conn net.Conn) ([]byte, error) {
	buf := make([]byte, messageSize)
	bytes, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	// log.Printf("message received: %s", string(buf[:bytes]))
	return buf[:bytes], nil
}
