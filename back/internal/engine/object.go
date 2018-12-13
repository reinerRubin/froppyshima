package engine

type Object struct {
	Coord    *Coord
	Template *ObjectTemplate
}

// TODO optimize this, "skirt" collision
// we can sum two matrices and check collisions for one iterate
func (o *Object) IsConflicted(other *Object) (collision bool) {
	_ = o.Template.Layout.ForEachNotNullYX(func(lCoord *Coord) (stop bool, err error) {
		isCollision := func(coord *Coord) (collision bool) {
			shiftedCoord := o.Coord.Shift(coord)

			_ = other.Template.Layout.ForEachNotNullYX(func(olCoord *Coord) (bool, error) {
				oCoord := other.Coord.Shift(olCoord)
				if *shiftedCoord == *oCoord {
					collision = true
					return true, nil
				}

				return false, nil
			})

			return
		}

		// TODO check "other object" skirt; it can be bigger
		for _, skx := range o.Template.Skirt {
			for _, sky := range o.Template.Skirt {
				if isCollision(NewCoord(skx, sky).Shift(lCoord)) {
					collision = true
					return true, nil
				}
			}
		}

		return false, nil
	})

	return
}
