package engine

type (
	// Field describes a field geometry
	Field struct {
		Maxx, Maxy int
	}
)

// IsPossiblePosition checks if object overlap borders
func (f *Field) IsPossiblePosition(object *Object) bool {
	if object.Coord.X < 0 || object.Coord.Y < 0 {
		return false
	}

	for y := range object.Template.Layout {
		for x := range object.Template.Layout[y] {
			if object.Template.Layout[y][x] == 0 {
				continue
			}

			coord := object.Coord.Shift(NewCoord(x, y))
			if coord.X >= f.Maxx {
				return false
			}
			if coord.Y >= f.Maxy {
				return false
			}
		}
	}

	return true
}

func (f *Field) forEachSquare(fn func(coord *Coord) error) error {
	for y0 := 0; y0 < f.Maxy; y0++ {
		for x0 := 0; x0 < f.Maxx; x0++ {
			if err := fn(NewCoord(x0, y0)); err != nil {
				return err
			}
		}

	}

	return nil
}
