package engine

import "testing"

// Too generic test
func TestMinInfo(t *testing.T) {
	game, err := NewGame()
	if err != nil {
		t.Fatalf("cant init game: %s", err)
	}

	if err := game.PutShips(); err != nil {
		t.Fatalf("cant put ships: %s", err)
	}

	{
		info := game.MinInfo()
		if len(info.FailHits) != 0 {
			t.Error("there are fail hits on empty game")
		}

		if len(info.SuccessHits) != 0 {
			t.Error("there are success hits on empty game")
		}

		if !info.AnyMoreShips {
			t.Error("there are no ships on empty game")
		}
	}

	ship := game.Ships[0]
	result, err := game.DotHit(ship.Object.Coord.Shift(ship.Segments[0].Coord))
	if err != nil {
		t.Errorf("cant hit a ship: %s", err)
	}
	if result != HitResultSuccess {
		t.Errorf("cant hit a ship; result is missed")
	}

	{
		info := game.MinInfo()
		if len(info.FailHits) != 0 {
			t.Error("there are fail hits on empty game")
		}
		if successHits := len(info.SuccessHits) + len(info.FatalHits); successHits != 1 {
			t.Errorf("we have wounded (or killed) only one ship but %d", successHits)
		}
		if !info.AnyMoreShips {
			t.Error("there are no ships after one hit")
		}
	}
}
