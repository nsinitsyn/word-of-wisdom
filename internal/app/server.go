package app

import (
	"log"
	"time"
	"word-of-wisdom/internal/config"
	"word-of-wisdom/internal/repository"
	"word-of-wisdom/internal/server"
)

type ServerApp struct {
	server  *server.Server
	stopped chan struct{}
}

func NewServerApp() ServerApp {
	config := config.GetServerConfig()

	repo := repository.NewFileRepository()
	server, err := server.NewServer(config, &repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %s\n", config.Addr)

	app := ServerApp{server: server, stopped: make(chan struct{})}

	go func() {
		server.Run()
		close(app.stopped)
	}()

	return app
}

func (app ServerApp) GracefulStop(timeout time.Duration) {
	app.server.Stop()
	select {
	case <-app.stopped:
	case <-time.After(timeout):
		break
	}

	log.Println("application stopped")
}
