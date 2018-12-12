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

// RotateNTimes rotates  the layout counterclockwise
func (s Layout) RotateNTimes(times int) (rotated Layout) {
	times = times % 4
	if times == 0 {
		return s.Copy()
	}

	cfc := coefficients(times)
	ylen := abs(len(s[0])*cfc.sin + len(s)*cfc.cos)
	xlen := abs(len(s)*cfc.sin + len(s[0])*cfc.cos)

	rotated = make(Layout, ylen)
	for y := range rotated {
		rotated[y] = make(LayoutLine, xlen)
	}

	for y := range s {
		for x := range s[y] {
			ynew := x*cfc.sin + y*cfc.cos
			xnew := x*cfc.cos - y*cfc.sin

			// shift coords back to the positive I quarter
			if times == 1 || times == 2 {
				ynew += ylen - 1
			}
			if times == 2 || times == 3 {
				xnew += xlen - 1
			}

			rotated[ynew][xnew] = s[y][x]
		}
	}

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
