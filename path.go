package grid

import "image"

type path struct {
	Points []image.Point
	Cost   float64
}

func (p *path) Len() (rv int) {
	return len(p.Points)
}

func (p *path) Last() (v image.Point) {
	return p.Points[p.Len()-1]
}

func (p *path) Fork(pt image.Point, dist float64) (rv path) {
	l := p.Len()
	add := make([]image.Point, l+1)
	copy(add, p.Points)

	add[l] = pt

	return path{
		Points: add,
		Cost:   p.Cost + dist,
	}
}
