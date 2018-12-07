package engine

// Coord is basic x, y coords
type Coord struct {
	X, Y int
}

// NewCoord TBD
func NewCoord(x, y int) *Coord {
	return &Coord{X: x, Y: y}
}

// Shift moves x, y in absolute position
func (c *Coord) Shift(o *Coord) *Coord {
	return &Coord{
		X: c.X + o.X,
		Y: c.Y + o.Y,
	}
}

// Unshift moves x, y in relative position
func (c *Coord) Unshift(o *Coord) *Coord {
	return &Coord{
		X: o.X - c.X,
		Y: o.Y - c.Y,
	}
}
