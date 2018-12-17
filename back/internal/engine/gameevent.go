package engine

type (
	GameEvent struct {
		GameOver        *GameOver
		ShipIsWounded   *ShipIsWounded
		ShipIsDestroyed *ShipIsDestroyed
	}

	GameOver struct {
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
	if e.GameOver != nil {
		return e.GameOver.Message
	}

	if e.ShipIsWounded != nil {
		return e.ShipIsWounded.Message
	}

	if e.ShipIsDestroyed != nil {
		return e.ShipIsDestroyed.Message
	}

	return "very strange"
}
