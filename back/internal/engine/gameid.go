package engine

import "github.com/google/uuid"

type GameID uuid.UUID

func NewGameID() GameID {
	return GameID(uuid.New())
}
func (id GameID) String() string {
	return uuid.UUID(id).String()
}

// MarshalText implements encoding.TextMarshaler.
func (id GameID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *GameID) UnmarshalText(data []byte) error {
	uuid := new(uuid.UUID)
	if err := uuid.UnmarshalText(data); err != nil {
		return err
	}

	*id = GameID(*uuid)

	return nil
}
