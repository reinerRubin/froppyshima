package engine

type ShipTemplate struct {
	typo  string
	char  rune
	ship  *ObjectTemplate
	count int
}

var gameShipTemplates func() []*ShipTemplate = func() []*ShipTemplate {
	return []*ShipTemplate{
		&ShipTemplate{
			typo:  "O shaped",
			char:  'O',
			ship:  OShip,
			count: 1,
		},
		&ShipTemplate{
			typo:  "L shaped",
			char:  'L',
			ship:  LShip,
			count: 1,
		},
		&ShipTemplate{
			typo:  "I shaped",
			char:  'I',
			ship:  IShip,
			count: 1,
		},
		&ShipTemplate{
			typo:  "Dot shaped",
			char:  'G',
			ship:  DotShip,
			count: 2,
		},
	}
}
