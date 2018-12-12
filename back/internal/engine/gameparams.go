package engine

type shipTemplate struct {
	typo  string
	char  rune
	ship  *ObjectTemplate
	count int
}

// hardcoded game params
func getGameShipTemplates() []*shipTemplate {
	return []*shipTemplate{
		&shipTemplate{
			typo:  "O shaped",
			char:  'O',
			ship:  OShip,
			count: 1,
		},
		&shipTemplate{
			typo:  "L shaped",
			char:  'L',
			ship:  LShip,
			count: 1,
		},
		&shipTemplate{
			typo:  "I shaped",
			char:  'I',
			ship:  IShip,
			count: 1,
		},
		&shipTemplate{
			typo:  "Dot shaped",
			char:  'G',
			ship:  DotShip,
			count: 2,
		},
	}
}
