package engine

type (
	Hit struct {
		Object *Object
	}

	HitResult int
)

const (
	HitResultMissed HitResult = iota
	HitResultSuccess
)

func (f *Hit) Render(c *Coord) string {
	coord := f.Object.Coord.Unshift(c)
	for y := range f.Object.Template.Layout {
		for x := range f.Object.Template.Layout[y] {
			if *NewCoord(x, y) == *coord {
				return "*"
			}
		}
	}

	return ""
}
