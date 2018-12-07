package engine

type GameMinInfo struct {
	Maxx int
	Maxy int

	FailHits    []*Coord
	FatalHits   []*Coord
	SuccessHits []*Coord

	AnyMoreShips bool
}
