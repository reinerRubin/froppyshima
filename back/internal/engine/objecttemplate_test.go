package engine

import (
	"testing"
)

func TestRotate(t *testing.T) {
	type testData struct {
		label string

		rotateTimes int
		before      Layout
		after       Layout
	}

	testDataVector := []*testData{
		&testData{
			label:       "zero rotation",
			rotateTimes: 4,
			before: Layout{
				LayoutLine{1, 1, 1, 1},
				LayoutLine{0, 0, 0, 1}},

			after: Layout{
				LayoutLine{1, 1, 1, 1},
				LayoutLine{0, 0, 0, 1}},
		},
		&testData{
			label:       "one rotation",
			rotateTimes: 1,
			before: Layout{
				LayoutLine{1, 1, 1, 1},
				LayoutLine{0, 0, 0, 1}},

			after: Layout{
				LayoutLine{1, 1},
				LayoutLine{1, 0},
				LayoutLine{1, 0},
				LayoutLine{1, 0}},
		},
		&testData{
			label:       "second rotation",
			rotateTimes: 2,
			before: Layout{
				LayoutLine{1, 1, 1, 1},
				LayoutLine{0, 0, 0, 1}},

			after: Layout{
				LayoutLine{1, 0, 0, 0},
				LayoutLine{1, 1, 1, 1}},
		},
		&testData{
			label:       "thrid rotation",
			rotateTimes: 3,
			before: Layout{
				LayoutLine{1, 1, 1, 1},
				LayoutLine{0, 0, 0, 1}},

			after: Layout{
				LayoutLine{0, 1},
				LayoutLine{0, 1},
				LayoutLine{0, 1},
				LayoutLine{1, 1}},
		},
		&testData{
			label:       "thrid rotation",
			rotateTimes: 4,
			before: Layout{
				LayoutLine{1, 1, 1, 1},
				LayoutLine{0, 0, 0, 1}},

			after: Layout{
				LayoutLine{1, 1, 1, 1},
				LayoutLine{0, 0, 0, 1}},
		},
	}

nextTest:
	for _, testData := range testDataVector {
		rotated := testData.before.RotateNTimes(testData.rotateTimes)

		if len(rotated) != len(testData.after) {
			t.Errorf("%s test failed on y len", testData.label)
		}

		for y := range rotated {
			if len(rotated[y]) != len(testData.after[y]) {
				t.Errorf("%s test failed on yx len", testData.label)
			}

			for x := range rotated[y] {
				if rotated[y][x] != testData.after[y][x] {
					t.Errorf("%s test failed; excepted \n%s \n vs actual \n%s\n",
						testData.label, testData.after, rotated)
					continue nextTest
				}
			}
		}
	}
}
