package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"word-of-wisdom/internal/app"
)

const GracefulShutdownTimeoutSec = 10

func main() {
	app := app.NewServerApp()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.GracefulStop(GracefulShutdownTimeoutSec * time.Second)
}
