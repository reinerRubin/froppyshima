package engine

import "math/rand"

type (
	LayoutLine     []int
	Layout         []LayoutLine
	ObjectTemplate struct {
		Layout Layout
		Skirt  []int
	}
)

func (s ObjectTemplate) Copy() *ObjectTemplate {
	return &ObjectTemplate{
		Layout: s.Layout.Copy(),
		Skirt:  s.Skirt,
	}
}

// Rotate returns the rotated layout
// TODO use matrix multiplication
func (s Layout) Rotate() (rotated Layout) {
	rotated = make(Layout, len(s[0]))
	for y := range rotated {
		rotated[y] = make(LayoutLine, len(s))
	}

	for y := range s {
		for x := range s[y] {
			rotated[x][len(s)-y-1] = s[y][x]
		}
	}

	return
}

// TODO we can archive a rotation via the matrix multiplication
func (s Layout) RotateNTimes(times int) (rotated Layout) {
	times = times % 4

	if times == 0 {
		return s.Copy()
	}

	rotated = s.Rotate()
	rotated = rotated.RotateNTimes(times - 1)

	return
}

func (s Layout) Copy() (copied Layout) {
	copied = make(Layout, len(s))

	for i, line := range s {
		copied[i] = make(LayoutLine, len(line))
		copy(copied[i], line)
	}

	return
}

func (s Layout) ForEachNotNullYX(fn func(coord *Coord) (stop bool, err error)) error {
	for y := range s {
		for x := range s[y] {
			if s[y][x] == 0 {
				continue
			}

			stop, err := fn(NewCoord(x, y))
			if err != nil {
				return err
			}

			if stop {
				return nil
			}
		}
	}

	return nil
}

func (st Layout) String() (rendered string) {
	for lineNumber, y := range st {
		for _, x := range y {
			if x > 0 {
				rendered += "1"
			} else {
				rendered += "0"
			}
		}

		if lineNumber != len(st)-1 {
			rendered += "\n"
		}
	}

	return
}

func SuffleRotateVariants() (times []int) {
	rotateVariants := []int{0, 1, 2, 3}
	rand.Shuffle(len(rotateVariants), func(i, j int) {
		rotateVariants[i], rotateVariants[j] = rotateVariants[j], rotateVariants[i]
	})

	return rotateVariants
}
