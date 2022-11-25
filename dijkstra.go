package grid

import (
	"image"
	"math"

	"github.com/s0rg/array2d"
)

const (
	minRank = 0
	maxRank = math.MaxUint16 - 2
)

type DijkstraMap struct {
	ranks array2d.Array[uint16]
}

func (dm *DijkstraMap) GetTarget(
	src image.Point,
	dirs []image.Point,
) (rv image.Point, ok bool) {
	var r1, r2 uint16

	if r1, ok = dm.ranks.Get(src.X, src.Y); !ok {
		return
	}

	rv, r2 = dm.lowest(src, dirs)

	return rv, r2 < r1
}

func (dm *DijkstraMap) update(
	targets []image.Point,
	canpass func(image.Point) bool,
) {
	dm.ranks.Fill(func() uint16 {
		return maxRank
	})

	for _, pt := range targets {
		dm.ranks.Set(pt.X, pt.Y, minRank)
	}

	var (
		mW, mH       = dm.ranks.Bounds()
		dirs         = Points(DirectionsALL...)
		changed      = true
		pt           image.Point
		srank, lrank uint16
	)

	for changed {
		changed = false

		for x := 0; x < mW; x++ {
			for y := 0; y < mH; y++ {
				for _, pt = range []image.Point{
					image.Pt(x, y),
					image.Pt((mW-1)-x, (mH-1)-y),
				} {
					if !canpass(pt) {
						continue
					}

					srank, _ = dm.ranks.Get(pt.X, pt.Y)

					if _, lrank = dm.lowest(pt, dirs); srank > lrank+1 {
						dm.ranks.Set(pt.X, pt.Y, lrank+1)

						changed = true
					}
				}
			}
		}
	}
}

func (dm *DijkstraMap) lowest(
	src image.Point,
	dirs []image.Point,
) (rv image.Point, rank uint16) {
	rank = maxRank

	var (
		p  image.Point
		r  uint16
		ok bool
	)

	for _, d := range dirs {
		p = src.Add(d)

		if r, ok = dm.ranks.Get(p.X, p.Y); !ok {
			continue
		}

		if r < rank {
			rv, rank = p, r
		}
	}

	return rv, rank
}
