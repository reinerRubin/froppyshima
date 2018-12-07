package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/reinerRubin/froppyshima/back/internal/server"
)

func main() {
	fmt.Println("welcome to `froppyshima' the game")

	server := server.New()
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.Stop()

	fmt.Println("ya later!")
}
