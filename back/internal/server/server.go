package server

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/reinerRubin/froppyshima/back/internal"
	"github.com/reinerRubin/froppyshima/back/internal/client"
)

type Server struct {
	http *http.Server

	stopOnce      sync.Once
	stopChannel   chan struct{}
	stopedChannel chan struct{}

	boltDBProvider *internal.BoltDBProvider

	clientContext *client.ClientContext
}

func New() *Server {
	server := &Server{
		http: &http.Server{
			Addr: ":8080",
		},

		stopChannel:   make(chan struct{}, 0),
		stopedChannel: make(chan struct{}, 0),
	}

	if err := server.Init(); err != nil {
		log.Fatalf("cant init server: %s", err)
	}

	return server
}

func (s *Server) Init() error {
	dbProvider, err := internal.NewBoltDBProvider()
	if err != nil {
		return err
	}
	s.boltDBProvider = dbProvider

	gameRepository, err := internal.NewBoltDBGameRepository(dbProvider)
	if err != nil {
		log.Fatalf("cant init game repository!: %s", err)
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
	<-s.stopedChannel

	return nil
}

func (s *Server) run() error {
	http.HandleFunc("/ws", s.StartGame)
	go s.http.ListenAndServe()

	<-s.stopChannel

	s.stopRoutine()

	close(s.stopedChannel)

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
	CheckOrigin:     func(*http.Request) bool { return true },
}

func (s *Server) StartGame(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client.NewClient(conn, s.clientContext).Start()
}
