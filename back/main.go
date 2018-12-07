package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/reinerRubin/froppyshima/back/internal/server"
)

func main() {
	log.Println("welcome to `froppyshima' the game")

	server := server.New()
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	server.Stop()

	log.Println("see you space cowboy!")
}
