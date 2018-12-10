package engine

import "log"

type Object struct {
	Coord    *Coord
	Template *ObjectTemplate
}

// TODO optimaize; "skirt" collision; remove complexity
func (o *Object) IsConflicted(other *Object) (collision bool) {
	err := o.Template.Layout.ForEachNotNullYX(func(lCoord *Coord) (stop bool, err error) {
		isCollision := func(coord *Coord) (collision bool) {
			shiftedCoord := o.Coord.Shift(coord)

			other.Template.Layout.ForEachNotNullYX(func(olCoord *Coord) (bool, error) {
				oCoord := other.Coord.Shift(olCoord)
				if *shiftedCoord == *oCoord {
					collision = true
					return true, nil
				}

				return false, nil
			})

			return
		}

		// TODO check "other object" skirt; it can be bigger;
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
	if err != nil {
		log.Fatal("unreachable collision")
	}

	return
}
