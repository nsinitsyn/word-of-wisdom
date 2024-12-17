package app

import (
	"context"
	"log"
	"net"
	"word-of-wisdom/internal/config"
	"word-of-wisdom/pkg/challenge"
)

const messageSize int16 = 1024

func RunClientApp() {
	config := config.GetClientConfig()

	var dialer net.Dialer
	conn, err := dialer.DialContext(context.Background(), "tcp", config.ServerAddr)
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection: %s", err)
		}
	}()

	// send request
	Write(conn, []byte("get quote"))

	// get challenge
	resp := Read(conn)

	// solve challenge
	chal := challenge.NewChallenge(resp)
	solution := chal.Solve()
	if solution == nil {
		log.Fatal("cannot found the solution")
	}

	// send solution
	Write(conn, solution)

	// get quoute
	quote := Read(conn)
	log.Println(string(quote))
}

func Read(conn net.Conn) []byte {
	buf := make([]byte, messageSize)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("failed to read response: %s", err)
	}
	return buf[:n]
}

func Write(conn net.Conn, msg []byte) {
	_, err := conn.Write(msg)
	if err != nil {
		log.Fatalf("failed to send request: %s", err)
	}
}
