package engine

import "fmt"

type ShipSegmentStatus int

const (
	ShipSegmentStatusKilled ShipSegmentStatus = iota
	ShipSegmentStatusAlive
)

// Ship is a kind of object
type (
	ShipSegment struct {
		Coord  *Coord
		Status ShipSegmentStatus
	}

	Ship struct {
		Lives    int
		Typo     string
		Char     rune
		Segments []*ShipSegment
		Object   *Object
	}

	Ships []*Ship
)

func NewShip(template *ShipTemplate, object *Object) *Ship {
	ship := &Ship{
		Typo:     template.typo,
		Object:   object,
		Char:     template.char,
		Segments: NewShipSegments(object.Template.Layout),
	}
	ship.Lives = len(ship.Segments)
	return ship
}

func NewShipSegments(l Layout) []*ShipSegment {
	// TODO allocate memory properly
	segements := make([]*ShipSegment, 0)

	l.ForEachNotNullYX(func(lCoord *Coord) (stop bool, err error) {
		segements = append(segements, &ShipSegment{
			Status: ShipSegmentStatusAlive,
			Coord:  lCoord,
		})
		return false, nil
	})

	return segements
}

// UnderHit applies a hit to the ship
// TODO: use hit skirt
func (s *Ship) UnderHit(f *Hit) (wounded bool, killed bool) {
	for _, segment := range s.Segments {
		asegmentCoord := s.Object.Coord.Shift(segment.Coord)
		if *asegmentCoord == *f.Object.Coord {
			if segment.Status == ShipSegmentStatusAlive {
				segment.Status = ShipSegmentStatusKilled
				s.Lives--
				return true, s.IsDead()
			}
			return false, false
		}
	}

	return false, false
}

func (s *Ship) Render(c *Coord) string {
	coord := s.Object.Coord.Unshift(c)

	for _, segment := range s.Segments {
		if *segment.Coord == *coord {
			if segment.Status == ShipSegmentStatusAlive {
				return string(s.Char)
			}

			return string("X")
		}
	}

	return ""
}

func (s *Ship) IsDead() bool {
	return s.Lives == 0
}

func (s *Ship) HittedSegmentCoords() []*Coord {
	coords := make([]*Coord, 0, len(s.Segments))

	for _, segment := range s.Segments {
		if segment.Status == ShipSegmentStatusKilled {
			coords = append(coords, s.Object.Coord.Shift(segment.Coord))
		}
	}

	return coords
}

func (s *Ship) InlineInfo() string {
	state := ""
	if s.IsDead() {
		state = "dead"

	} else {
		state = fmt.Sprintf("alive with %d lives", s.Lives)
	}

	return fmt.Sprintf(`"%s" is %s`, s.Typo, state)
}

func (s Ships) IsAllDead() bool {
	for _, ship := range s {
		if !ship.IsDead() {
			return false
		}
	}

	return true
}
