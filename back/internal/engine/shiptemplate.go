package engine

var LShip = &ObjectTemplate{
	Layout: Layout{
		LayoutLine{1, 0},
		LayoutLine{1, 0},
		LayoutLine{1, 1},
	},
	Skirt: []int{-1, 0, 1},
}

var IShip = &ObjectTemplate{
	Layout: Layout{
		LayoutLine{1},
		LayoutLine{1},
		LayoutLine{1},
		LayoutLine{1},
	},
	Skirt: []int{-1, 0, 1},
}

var DotShip = &ObjectTemplate{
	Layout: Layout{LayoutLine{1}},
	Skirt:  []int{-1, 0, 1},
}
