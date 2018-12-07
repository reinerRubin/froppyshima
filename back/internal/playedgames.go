package internal

import (
	"sync"
	"time"

	"github.com/reinerRubin/froppyshima/back/internal/engine"
)

type PlayedGame struct {
	ID        engine.GameID
	StartedAt time.Time
}

type PlayedGameRegister interface {
	TryToPlay(id engine.GameID) bool
	StopPlay(id engine.GameID)
}

type InMemoryGameRegister struct {
	store *sync.Map
}

func NewInMemoryGameRegister() *InMemoryGameRegister {
	return &InMemoryGameRegister{
		store: new(sync.Map),
	}
}

func (m *InMemoryGameRegister) TryToPlay(id engine.GameID) bool {
	_, loaded := m.store.LoadOrStore(id, &PlayedGame{
		ID:        id,
		StartedAt: time.Now(),
	})

	return !loaded
}

func (m *InMemoryGameRegister) StopPlay(id engine.GameID) {
	m.store.Delete(id)
}
