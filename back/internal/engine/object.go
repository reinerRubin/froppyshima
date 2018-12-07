package engine

type (
	Object struct {
		Coord    *Coord
		Template *ObjectTemplate
	}
)

// TODO optimaize me
// TODO "skirt" collision
func (o *Object) IsConflicted(other *Object) bool {
	for y := range o.Template.Layout {
		for x := range o.Template.Layout[y] {
			if o.Template.Layout[y][x] == 0 {
				continue
			}

			// check if ships are intersected
			isCollision := func(coord *Coord) bool {
				shiftedCoord := o.Coord.Shift(coord)

				for yo := range other.Template.Layout {
					for xo := range other.Template.Layout[yo] {
						if other.Template.Layout[yo][xo] == 0 {
							continue
						}

						ocoord := other.Coord.Shift(NewCoord(xo, yo))
						if *shiftedCoord == *ocoord {
							return true
						}
					}
				}

				return false
			}

			// TODO check "other object" skirt; it can be bigger;
			for _, skx := range o.Template.Skirt {
				for _, sky := range o.Template.Skirt {
					if isCollision(NewCoord(x+skx, y+sky)) {
						return true
					}
				}
			}
		}
	}

	return false
}
