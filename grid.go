package grid

import (
	"image"
	"math"

	"github.com/s0rg/array2d"
	"github.com/s0rg/set"
	"github.com/s0rg/vec2d"
	"github.com/zyedidia/generic/heap"
)

const one = 1.0

// Iter is an iteration callback.
type Iter[T any] func(image.Point, T) bool

// Cast is a ray-casting callback.
type Cast[T any] func(image.Point, float64, T) bool

// Cost is a path-finding callback.
type Cost[T any] func(image.Point, float64, T) (float64, bool)

// Distance is a distance-measurement function.
type Distance func(a, b image.Point) float64

// Map represents generic 2D grid map.
type Map[T any] struct {
	cells array2d.Array[T]
	rc    image.Rectangle
}

// New return empty [Map] with given bounding rectangle.
func New[T any](rc image.Rectangle) (rv *Map[T]) {
	return &Map[T]{
		rc:    rc,
		cells: array2d.New[T](rc.Dx(), rc.Dy()),
	}
}

// Bounds returns grid width and height.
func (m *Map[T]) Bounds() (w, h int) {
	return m.cells.Bounds()
}

// Rectangle returns grid bounding rectangle.
func (m *Map[T]) Rectangle() image.Rectangle {
	return m.rc
}

// Get returns value (if any) at given point.
func (m *Map[T]) Get(p image.Point) (c T, ok bool) {
	return m.cells.Get(p.X, p.Y)
}

// MustGet returns value at given point, it will panic on out-of-bound access.
func (m *Map[T]) MustGet(p image.Point) (c T) {
	var ok bool

	if c, ok = m.Get(p); ok {
		return c
	}

	panic("grid: out-of-bounds access")
}

// Set sets value at given point.
func (m *Map[T]) Set(p image.Point, v T) (ok bool) {
	return m.cells.Set(p.X, p.Y, v)
}

// Iter iterates over map cells.
func (m *Map[T]) Iter(it Iter[T]) {
	m.cells.Iter(func(x, y int, v T) (next bool) {
		return it(image.Pt(x, y), v)
	})
}

// Fill fills map with given constructor.
func (m *Map[T]) Fill(filler func() T) {
	m.cells.Fill(filler)
}

// Neighbours iterates grid cell neighbours in given directions and order.
func (m *Map[T]) Neighbours(
	src image.Point,
	dirs []image.Point,
	iter Iter[T],
) {
	var (
		cur image.Point
		val T
		ok  bool
	)

	for _, d := range dirs {
		cur = src.Add(d)

		if val, ok = m.cells.Get(cur.X, cur.Y); !ok {
			continue
		}

		if !iter(cur, val) {
			break
		}
	}
}

// Path performs A-Star path finding in map.
func (m *Map[T]) Path(
	src, dst image.Point,
	dirs []image.Point,
	dist Distance,
	cost Cost[T],
) (rv []image.Point, ok bool) {
	if !src.In(m.rc) {
		return rv, false
	}

	var val T

	if val, ok = m.cells.Get(dst.X, dst.Y); !ok {
		return rv, false
	}

	tdist := dist(dst, src)

	if _, ok = cost(dst, tdist, val); !ok {
		return rv, false
	}

	var (
		road   path
		last   image.Point
		closed = make(set.Set[image.Point])
	)

	queue := heap.New[path](func(a, b path) bool {
		return a.Cost < b.Cost
	})

	queue.Push(road.Fork(src, tdist))

	for queue.Size() > 0 {
		road, _ = queue.Pop()
		last = road.Last()

		if !closed.TryAdd(last) {
			continue
		}

		if last.Eq(dst) {
			return road.Points, true
		}

		m.Neighbours(last, dirs, func(p image.Point, t T) (ok bool) {
			var ncost float64

			if ncost, ok = cost(p, dist(dst, p), t); ok {
				queue.Push(road.Fork(p, ncost))
			}

			return true
		})
	}

	return nil, false
}

// LineOfSight iterates visible cells within given distance.
func (m *Map[T]) LineOfSight(
	src image.Point,
	distMax float64,
	cast Cast[T],
) {
	if !src.In(m.rc) {
		return
	}

	const maxDegrees = 360.0

	for t := float64(0); t < maxDegrees; t++ {
		m.CastRay(src, t, distMax, cast)
	}
}

// CastRay performs DDA ray cast from point at map with given angle (in degrees), limited by given max distance.
func (m *Map[T]) CastRay(
	src image.Point,
	angle, distMax float64,
	cast Cast[T],
) {
	if !src.In(m.rc) {
		return
	}

	var (
		start      = vec2d.New(float64(src.X), float64(src.Y))
		s, c       = math.Sincos(radians(angle))
		dest       = start.Add(vec2d.New(c, s))
		rdir       = dest.Sub(start).Norm()
		step, rlen vec2d.V[float64]
	)

	if rdir.X < 0 {
		rlen.X = start.X - start.X
		step.X = -one
	} else {
		rlen.X = (start.X + one) - start.X
		step.X = one
	}

	if rdir.Y < 0 {
		rlen.Y = start.Y - start.Y
		step.Y = -one
	} else {
		rlen.Y = (start.Y + one) - start.Y
		step.Y = one
	}

	var (
		unit = vec2d.New(one, one).Div(rdir).Abs()
		mpt  image.Point
		dist float64
		val  T
		ok   bool
	)

	rlen = rlen.Mul(unit)

	for {
		if rlen.X < rlen.Y {
			start.X += step.X
			rlen.X += unit.X
		} else {
			start.Y += step.Y
			rlen.Y += unit.Y
		}

		if dist = math.Max(rlen.X, rlen.Y); dist > distMax {
			break
		}

		mpt = image.Pt(int(start.X), int(start.Y))

		if val, ok = m.cells.Get(mpt.X, mpt.Y); !ok {
			break
		}

		if !cast(mpt, dist, val) {
			break
		}
	}
}

// CastShadow performs recursive shadow-casting.
func (m *Map[T]) CastShadow(
	src image.Point,
	distMax float64,
	cast Cast[T],
) {
	const (
		octetMin = 1
		octetMax = 8
	)

	val, ok := m.cells.Get(src.X, src.Y)
	if !ok {
		return
	}

	cast(src, 0, val)

	for oct := octetMin; oct <= octetMax; oct++ {
		m.emitShadow(src, oct, one, distMax, 0.0, one, cast)
	}
}

// DijkstraMap calculates 'Dijkstra' map for given points.
func (m *Map[T]) DijkstraMap(
	targets []image.Point,
	iter Iter[T],
) (rv *DijkstraMap) {
	rv = &DijkstraMap{
		ranks: array2d.New[uint16](m.cells.Bounds()),
	}

	rv.update(targets, func(p image.Point) (ok bool) {
		val, _ := m.cells.Get(p.X, p.Y)

		return iter(p, val)
	})

	return rv
}

// Line by Bresenham's algorithm.
func (m *Map[T]) LineBresenham(
	src, dst image.Point,
	iter Iter[T],
) {
	if !src.In(m.rc) {
		return
	}

	const two = 2

	var (
		sx, sy = 1, 1
		dx, dy = abs(dst.X - src.X), -abs(dst.Y - src.Y)
		e1     = dx + dy
		e2     int
		val    T
		ok     bool
	)

	if src.X > dst.X {
		sx = -1
	}

	if src.Y > dst.Y {
		sy = -1
	}

	cur := src

	for {
		if val, ok = m.cells.Get(cur.X, cur.Y); !ok {
			break
		}

		if !iter(cur, val) {
			break
		}

		if cur.Eq(dst) {
			break
		}

		e2 = e1 * two

		if e2 >= dy {
			cur.X += sx
			e1 += dy
		}

		if e2 <= dx {
			cur.Y += sy
			e1 += dx
		}
	}
}

func (m *Map[T]) emitShadow(
	src image.Point,
	oct int,
	dist, distMax, slopeLow, slopeHigh float64,
	cast Cast[T],
) {
	if dist > distMax {
		return
	}

	const half = 0.5

	var (
		pt      image.Point
		low     = math.Floor(slopeLow*dist + half)
		high    = math.Ceil(slopeHigh*dist + half)
		val     T
		pdist   float64
		gap, ok bool
	)

	for h := low; h < high; h++ {
		pt = octantPoint(src, oct, int(dist), int(h))

		if val, ok = m.cells.Get(pt.X, pt.Y); !ok {
			continue
		}

		if pdist = dist + DistanceEuclidean(src, pt); pdist > distMax {
			continue
		}

		switch {
		case cast(pt, pdist, val):
			gap = true
		case gap:
			m.emitShadow(src, oct, dist+1, distMax, slopeLow, (h-half)/dist, cast)

			slopeLow = (h + half) / dist
			gap = false
		}
	}

	m.emitShadow(src, oct, dist+1, distMax, slopeLow, slopeHigh, cast)
}

func radians(v float64) (d float64) {
	const rad2deg = 180.0 / math.Pi

	return v / rad2deg
}

func octantPoint(p image.Point, oct, d, h int) (rv image.Point) {
	if oct&0x1 > 0 {
		d = -d
	}

	if oct&0x2 > 0 {
		h = -h
	}

	rv.X, rv.Y = d, h

	if oct&0x4 > 0 {
		rv.X, rv.Y = h, d
	}

	return p.Add(rv)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}
