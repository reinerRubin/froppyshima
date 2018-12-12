package engine

import (
	"testing"
)

func TestNewShipSegments(t *testing.T) {
	type testData struct {
		label    string
		layout   Layout
		segments []*ShipSegment
	}

	testDataVector := []*testData{
		&testData{
			label: "L test",
			layout: Layout{
				LayoutLine{1, 1},
				LayoutLine{0, 1},
			},
			segments: []*ShipSegment{
				&ShipSegment{
					Coord:  &Coord{0, 0},
					Status: ShipSegmentStatusAlive,
				},
				&ShipSegment{
					Coord:  &Coord{1, 0},
					Status: ShipSegmentStatusAlive,
				},
				&ShipSegment{
					Coord:  &Coord{1, 1},
					Status: ShipSegmentStatusAlive,
				},
			},
		},
		&testData{
			label: "dot test",
			layout: Layout{
				LayoutLine{1},
			},
			segments: []*ShipSegment{
				&ShipSegment{
					Coord:  &Coord{0, 0},
					Status: ShipSegmentStatusAlive,
				},
			},
		},
	}

	for _, testData := range testDataVector {
		actualSegments := newShipSegments(testData.layout)
		if len(actualSegments) != len(testData.segments) {
			t.Errorf("%s test failed on len check", testData.label)
		}

		for _, actualSegment := range actualSegments {
			if actualSegment.Status != ShipSegmentStatusAlive {
				t.Errorf("%s test failed; segment is dead", testData.label)
			}
		}

		for _, actualSegment := range actualSegments {
			found := false
			for _, segment := range testData.segments {
				if *segment.Coord == *actualSegment.Coord {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s test failed: segments does not match", testData.label)
			}
		}
	}
}

func TestUnderHit(t *testing.T) {
	type testData struct {
		label string
		ship  *Ship
		hits  []*Hit

		// last hit result
		wounded bool
		killed  bool
	}

	LShip := &ObjectTemplate{
		Layout: Layout{
			LayoutLine{1, 0},
			LayoutLine{1, 0},
			LayoutLine{1, 1},
		},
		Skirt: []int{-1, 0, 1},
	}

	newDotHit := func(x, y int) *Hit {
		DotHit := &ObjectTemplate{
			Layout: Layout{LayoutLine{1}},
			Skirt:  []int{0},
		}

		return &Hit{
			Object: &Object{
				Coord:    NewCoord(x, y),
				Template: DotHit.Copy(),
			},
		}
	}

	testDataVector := []*testData{
		&testData{
			label: "L test",
			ship: &Ship{
				Typo: "l test ship",
				Object: &Object{
					Coord:    NewCoord(0, 0),
					Template: LShip.Copy(),
				},
				Char:     'L',
				Segments: newShipSegments(LShip.Layout),
				Lives:    4,
			},
			hits:    []*Hit{newDotHit(0, 0)},
			wounded: true,
			killed:  false,
		},
		&testData{
			label: "L test killed",
			ship: &Ship{
				Typo: "l test ship",
				Object: &Object{
					Coord:    NewCoord(0, 0),
					Template: LShip.Copy(),
				},
				Char:     'L',
				Segments: newShipSegments(LShip.Layout),
				Lives:    4,
			},
			hits: []*Hit{
				newDotHit(0, 0),
				newDotHit(0, 1), newDotHit(30, 4), newDotHit(12, 15),
				newDotHit(0, 2), newDotHit(1, 2),
			},
			wounded: true,
			killed:  true,
		},
	}

	for _, testData := range testDataVector {
		var wounded, killed bool
		for _, testHit := range testData.hits {
			wounded, killed = testData.ship.UnderHit(testHit)
		}

		if testData.wounded != wounded {
			t.Errorf("%s failed at wounded test (%t)", testData.label, wounded)
		}
		if testData.killed != killed {
			t.Errorf("%s failed at killed test (%t)", testData.label, killed)
		}
	}
}
