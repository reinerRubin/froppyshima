package engine

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewGameID(t *testing.T) {
	gameID := NewGameID()
	gameIDStr := gameID.String()
	restoredGameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		t.Fatalf("cant restore gameID: %s", gameIDStr)
	}
	if GameID(restoredGameID) != gameID {
		t.Fatalf("restored gameID does not match the old one: %s vs %s", gameID, restoredGameID)
	}
}
