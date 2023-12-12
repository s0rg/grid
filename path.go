package grid

import "image"

type path struct {
	Parent *path
	Point  image.Point
	length int
	Cost   float64
}

func (p *path) Len() (rv int) {
	return p.length
}

func (p *path) Last() (v image.Point) {
	return p.Point
}

func (p *path) Points() (rv []image.Point) {
	rv = make([]image.Point, p.length)

	for i := p.length - 1; i > 0; i-- {
		rv[i], p = p.Point, p.Parent
	}

	return rv
}

func (p *path) Fork(pt image.Point, dist float64) (rv *path) {
	rv = &path{
		Point:  pt,
		Cost:   dist,
		length: 1,
	}

	if p == nil {
		return
	}

	rv.Parent = p
	rv.Cost += p.Cost
	rv.length += p.length

	return rv
}
