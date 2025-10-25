package grid

import "image"

type dir uint8

const (
	North dir = iota
	East
	South
	West
	NorthWest
	NorthEast
	SouthEast
	SouthWest
)

const (
	Up        dir = North
	Right         = East
	Down          = South
	Left          = West
	UpLeft        = NorthWest
	UpRight       = NorthEast
	DownRight     = SouthEast
	DownLeft      = SouthWest
)

var coords = []image.Point{
	{X: 0, Y: -1},  // N
	{X: 1, Y: 0},   // E
	{X: 0, Y: 1},   // S
	{X: -1, Y: 0},  // W
	{X: -1, Y: -1}, // NW
	{X: 1, Y: -1},  // NE
	{X: 1, Y: 1},   // SE
	{X: -1, Y: 1},  // SW
}

var (
	DirectionsCardinal = []dir{
		North,
		East,
		South,
		West,
	}

	DirectionsDiagonal = []dir{
		NorthWest,
		NorthEast,
		SouthEast,
		SouthWest,
	}

	DirectionsALL = []dir{
		NorthWest,
		North,
		NorthEast,
		East,
		SouthEast,
		South,
		SouthWest,
		West,
	}
)

// Points returns displacements for given directions, in requested order.
func Points(dirs ...dir) (rv []image.Point) {
	rv = make([]image.Point, len(dirs))

	for i, d := range dirs {
		rv[i] = coords[d]
	}

	return rv
}

// Invert returns opposite direction.
func (d dir) Invert() (rv dir) {
	switch d {
	case North:
		return South
	case East:
		return West
	case West:
		return East
	case NorthEast:
		return SouthWest
	case SouthWest:
		return NorthEast
	case NorthWest:
		return SouthEast
	case SouthEast:
		return NorthWest
	case South:
	}

	return North
}
