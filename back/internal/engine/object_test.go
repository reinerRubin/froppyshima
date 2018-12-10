package engine

import "testing"

func TestIsConflicted(t *testing.T) {
	type testData struct {
		sample   string
		o1, o2   *Object
		conflict bool
	}

	testDataVector := []*testData{
		&testData{
			sample: "dot",
			o1: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1},
					},
					Skirt: []int{-1, 0, 1},
				},
			},
			o2: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1},
					},
					Skirt: []int{-1, 0, 1},
				},
			},
			conflict: true,
		},
		&testData{
			sample: "Yin & yang",
			o1: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1, 1, 1, 1},
						LayoutLine{0, 0, 0, 1},
						LayoutLine{0, 0, 0, 1},
						LayoutLine{0, 0, 0, 0},
						LayoutLine{0, 0, 0, 0}},
					Skirt: []int{-1, 0, 1},
				},
			},
			o2: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{0, 0, 0, 0},
						LayoutLine{0, 0, 0, 0},
						LayoutLine{1, 0, 0, 0},
						LayoutLine{1, 0, 0, 0},
						LayoutLine{1, 1, 1, 1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			conflict: false,
		},
		&testData{
			sample: "relative inland",
			o1: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1, 1, 1, 1, 1},
						LayoutLine{1, 0, 0, 0, 1},
						LayoutLine{1, 0, 0, 0, 1},
						LayoutLine{1, 0, 0, 0, 1},
						LayoutLine{1, 0, 0, 0, 1},
						LayoutLine{1, 1, 1, 1, 1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			o2: &Object{
				Coord: NewCoord(2, 2),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			conflict: false,
		},
		&testData{
			sample: "Ship collision ('naval')",
			o1: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			o2: &Object{
				Coord: NewCoord(1, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			conflict: true,
		},
		&testData{
			sample: "resolved ship collision ('naval')",
			o1: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			o2: &Object{
				Coord: NewCoord(2, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1},
						LayoutLine{1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			conflict: false,
		},
		&testData{
			sample: "sparse templates",
			o1: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1, 0, 0},
						LayoutLine{1, 0, 0},
						LayoutLine{1, 0, 0},
						LayoutLine{1, 0, 0}},
					Skirt: []int{-1, 0, 1},
				},
			},
			o2: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{0, 0, 1},
						LayoutLine{0, 0, 1},
						LayoutLine{0, 0, 1},
						LayoutLine{0, 0, 1}},
					Skirt: []int{-1, 0, 1},
				},
			},
			conflict: false,
		},
		&testData{
			sample: "sparse templates 2",
			o1: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{1, 0, 0},
						LayoutLine{1, 0, 0},
						LayoutLine{1, 0, 0},
						LayoutLine{1, 0, 0}},
					Skirt: []int{-1, 0, 1},
				},
			},
			o2: &Object{
				Coord: NewCoord(0, 0),
				Template: &ObjectTemplate{
					Layout: Layout{
						LayoutLine{0, 0, 0},
						LayoutLine{0, 0, 1},
						LayoutLine{0, 0, 0},
						LayoutLine{0, 0, 0}},
					Skirt: []int{-1, 0, 1},
				},
			},
			conflict: false,
		},
	}

	for _, testData := range testDataVector {
		actual := testData.o1.IsConflicted(testData.o2)
		if actual != testData.conflict {
			t.Errorf("%s sample is broken must be %t but %t",
				testData.sample, testData.conflict, actual)
		}
	}
}
