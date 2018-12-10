package engine

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type (

	// Game is main "engine" object
	Game struct {
		Ships  Ships
		Hits   []*Hit
		Events []*GameEvent `json:"-"`

		Field *Field
	}

	// Render is interface for all objects that want to be drawn must implement this interface
	Render interface {
		Render(c *Coord) string
	}
)

// NewGame is the game constructor
func NewGame() (*Game, error) {
	game := &Game{
		Ships: make([]*Ship, 0),
		Hits:  make([]*Hit, 0),
		Field: &Field{
			Maxx: 10,
			Maxy: 10,
		},
	}

	if err := game.Init(); err != nil {
		return nil, err
	}

	return game, nil
}

// Init inits the game
func (g *Game) Init() error {
	g.Events = make([]*GameEvent, 0)
	return nil
}

// PutShips puts ships to random positions
func (g *Game) PutShips() error {
	templates := gameShipTemplates()

	for len(templates) > 0 {
		templateNumber := rand.Intn(len(templates))
		template := templates[templateNumber]

		object, err := g.putToRandomPosition(template.ship)
		if err != nil {
			return err
		}

		g.Ships = append(g.Ships, NewShip(template, object))

		template.count--
		if template.count == 0 {
			templates = append(templates[:templateNumber], templates[templateNumber+1:]...)
		}
	}

	return nil
}

// oh boy; It's kind of Knapsack problem. This is a naive and terrible implementation
// because a full solution seems too complicated for a test application;
// so we only try to rotate
// =
// We get all the variants for the object and choose one instead of trying to put it randomly
// It does not good for sparse fields but work well with small one;
// We need a smart switch between these strategies
func (g *Game) putToRandomPosition(template *ObjectTemplate) (*Object, error) {
	newTemplate := template.Copy()

tryNextRotation:
	for _, rotationTimes := range SuffleRotateVariants() {
		rotatedLayout := template.Layout.RotateNTimes(rotationTimes)
		newTemplate.Layout = rotatedLayout

		var toPut []*Object

		err := g.Field.forEachSquare(func(coord *Coord) error {
			objectCandidate := &Object{
				Coord:    coord,
				Template: newTemplate,
			}

			if !g.Field.IsPossiblePosition(objectCandidate) {
				return nil
			}

			if len(g.Ships) == 0 {
				toPut = append(toPut, objectCandidate)
			}

			for _, object := range g.ConflictedObjects() {
				if objectCandidate.IsConflicted(object) {
					return nil
				}
			}
			toPut = append(toPut, objectCandidate)

			return nil
		})
		if err != nil {
			return nil, err
		}
		if len(toPut) == 0 {
			continue tryNextRotation
		}

		return toPut[rand.Intn(len(toPut))], nil
	}

	return nil, fmt.Errorf("cant put ship: %s", template.Layout)
}

// Render shows game state; Current implementation is very primitive;
// TODO ask objects to draw them self and merge results
func (g *Game) Render() (render string) {
	for yf := 0; yf < g.Field.Maxy; yf++ {
		for xf := 0; xf < g.Field.Maxx; xf++ {
			cell := "." // emptyCell

			for _, toRender := range g.ToRender() {
				newCell := toRender.Render(NewCoord(xf, yf))
				if newCell != "" {
					cell = newCell
				}
			}

			render += cell
		}

		if yf != g.Field.Maxy-1 {
			render += "\n"
		}
	}

	return
}

// ConflictedObjects resurns all in game objects
func (g *Game) ConflictedObjects() []*Object {
	shipObjects := g.ShipObjects()

	objects := make([]*Object, 0, len(shipObjects))
	objects = append(objects, shipObjects...)

	return objects
}

// ShipObjects returns all ship objects
func (g *Game) ShipObjects() []*Object {
	objects := make([]*Object, 0, len(g.Ships))
	for _, ship := range g.Ships {
		objects = append(objects, ship.Object)
	}

	return objects
}

// ToRender resurns all in game objects
func (g *Game) ToRender() []Render {
	objects := make([]Render, 0, len(g.Ships))

	for _, ship := range g.Ships {
		objects = append(objects, ship)
	}

	for _, hit := range g.Hits {
		objects = append(objects, hit)
	}

	return objects
}

// DotHit emulates dot hit to coord
func (g *Game) DotHit(c *Coord) (HitResult, error) {
	hit := &Hit{
		Object: &Object{
			Coord:    c,
			Template: DotHit.Copy(),
		},
	}

	for _, oldHit := range g.Hits {
		if *oldHit.Object.Coord == *hit.Object.Coord {
			return HitResultMissed, nil
		}
	}

	result := HitResultMissed
	for _, ship := range g.Ships {
		if hit.Object.IsConflicted(ship.Object) {
			if wounded, killed := ship.UnderHit(hit); wounded {
				result = HitResultSuccess

				if !killed {
					g.WoundedShip(ship)
					continue
				}
				g.KilledShip(ship)
				if g.Ships.IsAllDead() {
					g.GameOwer()
				}

			}
		}
	}

	if result == HitResultMissed {
		g.Hits = append(g.Hits, hit)
	}

	return result, nil
}

func (g *Game) GameOwer() {
	g.AddEvent(&GameEvent{
		GameOwer: &GameOwer{
			Message: "this is the end my only friend @ congratulation!",
		},
	})
}

func (g *Game) WoundedShip(s *Ship) {
	g.AddEvent(&GameEvent{
		ShipIsWounded: &ShipIsWounded{
			Message: s.InlineInfo(),
		},
	})
}

func (g *Game) KilledShip(s *Ship) {
	g.AddEvent(&GameEvent{
		ShipIsDestroyed: &ShipIsDestroyed{
			Message: s.InlineInfo(),
			Ship:    s,
		},
	})
}

func (g *Game) AddEvent(e *GameEvent) {
	g.Events = append(g.Events, e)
}

func (g *Game) PullEvents() []*GameEvent {
	events := g.Events
	g.Events = make([]*GameEvent, 0)
	return events
}

func (g *Game) MinInfo() *GameMinInfo {
	info := &GameMinInfo{
		Maxx:         g.Field.Maxx,
		Maxy:         g.Field.Maxy,
		FailHits:     make([]*Coord, 0, len(g.Hits)),
		SuccessHits:  make([]*Coord, 0),
		AnyMoreShips: !g.Ships.IsAllDead(),
	}

	for _, hit := range g.Hits {
		info.FailHits = append(info.FailHits, hit.Object.Coord)
	}

	for _, ship := range g.Ships {
		if !ship.IsDead() {
			info.SuccessHits = append(info.SuccessHits, ship.HittedSegmentCoords()...)
		} else {
			info.FatalHits = append(info.FatalHits, ship.HittedSegmentCoords()...)
		}
	}

	return info
}
