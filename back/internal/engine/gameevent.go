package engine

type (
	GameEvent struct {
		GameOwer        *GameOwer
		ShipIsWounded   *ShipIsWounded
		ShipIsDestroyed *ShipIsDestroyed
	}

	GameOwer struct {
		Message string
	}

	ShipIsWounded struct {
		Message string
	}

	ShipIsDestroyed struct {
		Message string
		Ship    *Ship
	}
)

func (e *GameEvent) String() string {
	if e.GameOwer != nil {
		return e.GameOwer.Message
	}

	if e.ShipIsWounded != nil {
		return e.ShipIsWounded.Message
	}

	if e.ShipIsDestroyed != nil {
		return e.ShipIsDestroyed.Message
	}

	return "very strange"
}
