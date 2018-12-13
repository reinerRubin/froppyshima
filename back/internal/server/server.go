package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/reinerRubin/froppyshima/back/internal"
	"github.com/reinerRubin/froppyshima/back/internal/client"
	"github.com/reinerRubin/froppyshima/back/internal/config"
)

type Server struct {
	http *http.Server

	stopOnce       sync.Once
	stopChannel    chan struct{}
	stoppedChannel chan struct{}

	boltDBProvider *internal.BoltDBProvider

	clientContext *client.ClientContext
}

func New(config *config.Config) (*Server, error) {
	server := &Server{
		stopChannel:    make(chan struct{}),
		stoppedChannel: make(chan struct{}),
	}

	if err := server.Init(config); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) Init(config *config.Config) error {
	s.http = &http.Server{
		Addr: s.ServerAddr(config.Server.Port),
	}

	dbProvider, err := internal.NewBoltDBProvider(config.BoltDB)
	if err != nil {
		return err
	}
	s.boltDBProvider = dbProvider

	gameRepository, err := internal.NewBoltDBGameRepository(dbProvider)
	if err != nil {
		return err
	}

	s.clientContext = &client.ClientContext{
		GameRepository:     gameRepository,
		PlayedGameRegister: internal.NewInMemoryGameRegister(),
	}

	return nil
}

func (s *Server) Start() error {
	go s.run()

	return nil
}

func (s *Server) Stop() error {
	s.stopOnce.Do(func() {
		close(s.stopChannel)
	})
	<-s.stoppedChannel

	return nil
}

func (s *Server) run() error {
	http.HandleFunc("/ws", s.StartGame)
	go log.Fatalf("cant listen and serve: %s", s.http.ListenAndServe())

	<-s.stopChannel

	s.stopRoutine()

	close(s.stoppedChannel)

	return nil
}

func (s *Server) stopRoutine() {
	err := s.http.Shutdown(context.Background())
	if err != nil {
		log.Printf("error on shutdown: %s", err)
	}

	err = s.boltDBProvider.DB.Close()
	if err != nil {
		log.Printf("cant close db: %s", err)
	}
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// not worth for a test application
	CheckOrigin: func(*http.Request) bool { return true },
}

func (s *Server) StartGame(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client.NewClient(conn, s.clientContext).Start()
}

func (s *Server) ServerAddr(port int) string {
	return fmt.Sprintf(":%d", port)
}
