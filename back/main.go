package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/reinerRubin/froppyshima/back/internal/config"
	"github.com/reinerRubin/froppyshima/back/internal/server"
)

func main() {
	if err := runApp(); err != nil {
		log.Fatalf("app is terminated with error: %s", err)
	}
}

func runApp() error {
	log.Println("welcome to `froppyshima' the game")

	config, err := config.New()
	if err != nil {
		return fmt.Errorf("cant init cfg: %s", err)
	}

	server := server.New(config)
	if err := server.Start(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := server.Stop(); err != nil {
		return err
	}

	log.Println("see you space cowboy!")
	return nil
}
