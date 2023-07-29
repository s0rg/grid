package grid

import (
	"image"
	"testing"

	"github.com/s0rg/set"
)

func TestMapBounds(t *testing.T) {
	t.Parallel()

	m := New[struct{}](image.Rect(0, 0, 20, 10))

	if w, h := m.Bounds(); w != 20 || h != 10 {
		t.Fail()
	}
}

func TestMapGetSet(t *testing.T) {
	t.Parallel()

	m := New[int](image.Rect(0, 0, 10, 10))

	p := image.Pt(1, 1)

	if v, ok := m.Get(p); !ok || v != 0 {
		t.Fail()
	}

	if ok := m.Set(p, 1); !ok {
		t.Fail()
	}

	if v, ok := m.Get(p); !ok || v != 1 {
		t.Fail()
	}

	p = image.Pt(11, 11)

	if ok := m.Set(p, 1); ok {
		t.Fail()
	}

	if _, ok := m.Get(p); ok {
		t.Fail()
	}
}

func TestMapMustGet(t *testing.T) {
	t.Parallel()

	m := New[int](image.Rect(0, 0, 10, 10))

	p := image.Pt(1, 1)

	if v, ok := m.Get(p); !ok || v != 0 {
		t.Fail()
	}

	if ok := m.Set(p, 1); !ok {
		t.Fail()
	}

	if m.MustGet(p) != 1 {
		t.Fail()
	}

	var panik bool

	func() {
		defer func() {
			if r := recover(); r != nil {
				panik = true
			}
		}()

		_ = m.MustGet(image.Pt(11, 11))
	}()

	if !panik {
		t.Fail()
	}
}

func TestMapIter(t *testing.T) {
	t.Parallel()

	const W, H = 10, 10

	m := New[int](image.Rect(0, 0, W, H))

	var c int

	m.Iter(func(_ image.Point, _ int) bool {
		c++

		return true
	})

	if c != W*H {
		t.Fail()
	}

	m.Set(image.Pt(0, 0), 1)
	m.Set(image.Pt(1, 0), 2)
	m.Set(image.Pt(2, 0), 3)

	c = 0

	m.Iter(func(p image.Point, v int) bool {
		if p.X > 1 {
			return false
		}

		c += v

		return true
	})

	if c != 3 {
		t.Fail()
	}
}

func TestMapFill(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	m := New[int](image.Rect(0, 0, W, H))

	m.Fill(func() int {
		return 1
	})

	var c int

	m.Iter(func(_ image.Point, v int) bool {
		c += v

		return true
	})

	if c != W*H {
		t.Fail()
	}
}

func neighboursCount(m *Map[struct{}], p image.Point, d []image.Point) (count int) {
	m.Neighbours(p, d, func(_ image.Point, _ struct{}) bool {
		count++

		return true
	})

	return count
}

func TestMapNeighbours(t *testing.T) {
	t.Parallel()

	const (
		W, H          = 5, 5
		min, mid, max = 0, 2, 4
	)

	var cases = []struct {
		Point image.Point
		Count int
	}{
		{Point: image.Pt(min, min), Count: 2},
		{Point: image.Pt(mid, min), Count: 3},
		{Point: image.Pt(max, min), Count: 2},
		{Point: image.Pt(max, mid), Count: 3},
		{Point: image.Pt(max, max), Count: 2},
		{Point: image.Pt(mid, max), Count: 3},
		{Point: image.Pt(min, max), Count: 2},
		{Point: image.Pt(min, mid), Count: 3},
		{Point: image.Pt(mid, mid), Count: 4},
	}

	m := New[struct{}](image.Rect(0, 0, W, H))
	d := Points(DirectionsCardinal...)

	for i, tc := range cases {
		if c := neighboursCount(m, tc.Point, d); c != tc.Count {
			t.Fatalf("case[%d] failed: point: %s want: %d got: %d", i, tc.Point, tc.Count, c)
		}
	}
}

func TestMapNeighboursBreak(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	var (
		m     = New[struct{}](image.Rect(0, 0, W, H))
		start = image.Pt(2, 2)
		north = Points(North)[0]
		seen  image.Point
	)

	m.Neighbours(start, Points(DirectionsCardinal...), func(p image.Point, _ struct{}) bool {
		seen = p

		return false
	})

	if !seen.Eq(start.Add(north)) {
		t.Fail()
	}
}

func TestMapPath(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	var (
		src    = image.Pt(1, 1)
		dst    = image.Pt(3, 2)
		dirs   = Points(DirectionsCardinal...)
		walls  = make(set.Set[image.Point])
		coster = func(p image.Point, d float64, _ struct{}) (cost float64, walkable bool) {
			return d, !walls.Has(p)
		}
	)

	m := New[struct{}](image.Rect(0, 0, W, H))

	m.Iter(func(p image.Point, _ struct{}) (next bool) {
		if p.X == 0 || p.X == W-1 || p.Y == 0 || p.Y == H-1 {
			walls.Add(p)
		}

		return true
	})

	// step 1: clean map
	p, ok := m.Path(
		src,
		dst,
		dirs,
		DistanceManhattan,
		coster,
	)
	if !ok {
		t.Fail()
	}

	if len(p) != 4 {
		t.Fail()
	}

	// step 2: first wall
	walls.Add(image.Pt(2, 1))

	p, ok = m.Path(
		src,
		dst,
		dirs,
		DistanceManhattan,
		coster,
	)
	if !ok {
		t.Fail()
	}

	if len(p) != 4 {
		t.Fail()
	}

	// step 3: second wall
	walls.Add(image.Pt(2, 2))

	p, ok = m.Path(
		src,
		dst,
		dirs,
		DistanceManhattan,
		coster,
	)
	if !ok {
		t.Fail()
	}

	if len(p) != 6 {
		t.Fail()
	}

	// step 4: the last diagonal wall
	walls.Add(image.Pt(1, 3))

	_, ok = m.Path(
		src,
		dst,
		dirs,
		DistanceManhattan,
		coster,
	)
	if ok {
		t.Fail()
	}

	// OOB cases
	_, ok = m.Path(
		image.Pt(10, 10),
		dst,
		dirs,
		DistanceManhattan,
		coster,
	)
	if ok {
		t.Fail()
	}

	_, ok = m.Path(
		src,
		image.Pt(10, 10),
		dirs,
		DistanceManhattan,
		coster,
	)
	if ok {
		t.Fail()
	}

	_, ok = m.Path(
		src,
		image.Pt(4, 4),
		dirs,
		DistanceManhattan,
		coster,
	)
	if ok {
		t.Fail()
	}
}

func TestMapLOS(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	var (
		src    = image.Pt(1, 1)
		walls  = make(set.Set[image.Point])
		seen   = make(set.Set[image.Point])
		caster = func(p image.Point, _ float64, _ struct{}) (walkable bool) {
			if walls.Has(p) {
				return false
			}

			seen.Add(p)

			return true
		}
	)

	m := New[struct{}](image.Rect(0, 0, W, H))

	m.Iter(func(p image.Point, _ struct{}) (next bool) {
		if p.X == 0 || p.X == W-1 || p.Y == 0 || p.Y == H-1 {
			walls.Add(p)
		}

		return true
	})

	m.LineOfSight(src, 6.0, caster)

	if len(seen) != 8 {
		t.Fail()
	}

	seen = make(set.Set[image.Point])

	m.LineOfSight(image.Pt(10, 10), 6.0, caster)

	if len(seen) != 0 {
		t.Fail()
	}
}

func TestMapRayOOB(t *testing.T) {
	t.Parallel()

	const (
		W, H  = 5, 5
		angle = 30.0
		dist  = 10.0
	)

	var (
		seen   = make(set.Set[image.Point])
		caster = func(p image.Point, _ float64, _ struct{}) (walkable bool) {
			seen.Add(p)

			return true
		}
		cases = []struct {
			Point image.Point
			Seen  int
		}{
			{Point: image.Pt(6, 6), Seen: 0},
			{Point: image.Pt(1, 1), Seen: 5},
		}
	)

	m := New[struct{}](image.Rect(0, 0, W, H))

	for i, tc := range cases {
		m.CastRay(tc.Point, angle, dist, caster)

		if tc.Seen != len(seen) {
			t.Log(seen)
			t.Fatalf("case[%d] failed want: %d got: %d", i, tc.Seen, len(seen))
		}
	}
}

func TestMapShadow(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	var (
		src    = image.Pt(1, 1)
		walls  = make(set.Set[image.Point])
		seen   = make(set.Set[image.Point])
		caster = func(p image.Point, _ float64, _ struct{}) (walkable bool) {
			if walls.Has(p) {
				return false
			}

			seen.Add(p)

			return true
		}
	)

	m := New[struct{}](image.Rect(0, 0, W, H))

	m.Iter(func(p image.Point, _ struct{}) (next bool) {
		if p.X == 0 || p.X == W-1 || p.Y == 0 || p.Y == H-1 {
			walls.Add(p)
		}

		return true
	})

	walls.Add(image.Pt(2, 1))
	walls.Add(image.Pt(2, 2))
	walls.Add(image.Pt(2, 3))

	m.CastShadow(src, 3.0, caster)

	if seen.Has(image.Pt(3, 1)) {
		t.Fail()
	}

	m.CastShadow(image.Pt(6, 1), 3.0, caster)

	if seen.Has(image.Pt(3, 1)) {
		t.Fail()
	}
}

func TestMapDijkstra(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	var (
		src     = image.Pt(2, 1)
		walls   = make(set.Set[image.Point])
		targets = []image.Point{image.Pt(2, 4)}
		dirs    = Points(DirectionsCardinal...)
	)

	m := New[struct{}](image.Rect(0, 0, W, H))

	walls.Add(image.Pt(1, 1))
	walls.Add(image.Pt(3, 1))

	d := m.DijkstraMap(targets, func(p image.Point, _ struct{}) (ok bool) {
		return !walls.Has(p)
	})

	dst, ok := d.GetTarget(src, dirs)
	if !ok {
		t.Fail()
	}

	if dst.Y != 2 {
		t.Fail()
	}

	dst, ok = d.GetTarget(dst, dirs)
	if !ok {
		t.Fail()
	}

	if dst.Y != 3 {
		t.Fail()
	}

	if _, ok = d.GetTarget(image.Pt(2, 6), dirs); ok {
		t.Fail()
	}
}

func TestMapDijkstraEmpty(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	var (
		src  = image.Pt(2, 1)
		dirs = Points(DirectionsCardinal...)
	)

	m := New[struct{}](image.Rect(0, 0, W, H))

	d := m.DijkstraMap([]image.Point{}, func(p image.Point, _ struct{}) (ok bool) {
		return false
	})

	if p, ok := d.GetTarget(src, dirs); ok {
		t.Log(p)
		t.Fail()
	}
}

func TestLineBresenham(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	var cases = []struct {
		Path []image.Point
		Src  image.Point
		Dst  image.Point
	}{
		{
			Src: image.Pt(1, 1),
			Dst: image.Pt(6, 6),
			Path: []image.Point{
				image.Pt(1, 1),
				image.Pt(2, 2),
				image.Pt(3, 3),
				image.Pt(4, 4),
			},
		},
		{
			Src: image.Pt(1, 4),
			Dst: image.Pt(1, 1),
			Path: []image.Point{
				image.Pt(1, 1),
				image.Pt(1, 2),
				image.Pt(1, 3),
				image.Pt(1, 4),
			},
		},
		{
			Src: image.Pt(4, 1),
			Dst: image.Pt(1, 1),
			Path: []image.Point{
				image.Pt(1, 1),
				image.Pt(2, 1),
				image.Pt(3, 1),
				image.Pt(4, 1),
			},
		},
		{
			Src: image.Pt(-1, -1),
		},
	}

	m := New[struct{}](image.Rect(0, 0, W, H))

	for i, c := range cases {
		seen := make(set.Set[image.Point])

		m.LineBresenham(c.Src, c.Dst, func(p image.Point, _ struct{}) (ok bool) {
			seen.Add(p)

			return true
		})

		if len(seen) != len(c.Path) {
			t.Fatalf("case %d failed", i)
		}

		if len(c.Path) == 0 {
			continue
		}

		for _, p := range c.Path {
			if !seen.Has(p) {
				t.Fatalf("not seen: %v", p)
			}
		}
	}
}

func TestLineBresenhamBreak(t *testing.T) {
	t.Parallel()

	const W, H = 5, 5

	m := New[struct{}](image.Rect(0, 0, W, H))
	x := image.Pt(2, 2)

	m.LineBresenham(image.Pt(1, 1), image.Pt(3, 3), func(p image.Point, _ struct{}) (ok bool) {
		if p.Eq(x) {
			t.Fail()
		}

		return false
	})
}
