package grid

import "image"

type path struct {
	Points []image.Point
	Cost   float64
}

func (p path) Len() (rv int) {
	return len(p.Points)
}

func (p path) Last() (v image.Point) {
	return p.Points[p.Len()-1]
}

func (p path) Fork(pt image.Point, dist float64) (rv path) {
	l := p.Len()
	points := make([]image.Point, l, l+1)
	copy(points, p.Points)

	rv.Points = append(points, pt)
	rv.Cost = p.Cost + dist

	return rv
}
