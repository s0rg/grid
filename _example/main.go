package main

import (
	"fmt"
	"image"
	"math/rand"
	"strings"
	"time"

	"github.com/s0rg/grid"
)

const (
	mapW = 60
	mapH = 20
)

type cell struct {
	Wall bool
	Rune rune
}

func showMap(m *grid.Map[*cell]) {
	w, h := m.Bounds()

	cells := make([][]rune, h)

	for i := 0; i < h; i++ {
		cells[i] = make([]rune, w)
	}

	m.Iter(func(p image.Point, c *cell) (next bool) {
		r := c.Rune

		switch {
		case r != 0:
		case c.Wall:
			r = '#'
		default:
			r = ' '
		}

		cells[p.Y][p.X] = r

		return true
	})

	var sb strings.Builder

	for y := 0; y < h; y++ {
		sb.WriteString(string(cells[y]))
		sb.WriteByte('\n')
	}

	fmt.Print(sb.String())
}

func main() {
	g := grid.New[*cell](image.Rect(0, 0, mapW, mapH))

	// pre-fill grid items
	g.Fill(func() *cell {
		return &cell{}
	})

	// set borders
	g.Iter(func(p image.Point, c *cell) (next bool) {
		switch {
		case p.X == 0 || p.X == mapW-1:
			fallthrough
		case p.Y == 0 || p.Y == mapH-1:
			c.Wall = true
		}

		return true
	})

	// add random walls
	rand.Seed(time.Now().UnixNano())

	g.Iter(func(p image.Point, c *cell) (next bool) {
		if rand.Float64() < 0.13 {
			c.Wall = true
		}

		return true
	})

	// cast DDA ray
	var (
		steps    int
		src, dst image.Point
	)

	src = image.Pt(3, 3)

	g.CastRay(src, 10, 50.0, func(_ image.Point, _ float64, c *cell) (next bool) {
		if c.Wall {
			c.Rune = '%'
		} else {
			c.Rune = '.'
		}

		steps++

		// this is a ray-casting example, so dont stop on walls
		return true
	})

	g.MustGet(src).Rune = '@'

	showMap(g)

	fmt.Printf("ray took: %d steps\n\n", steps)

	// clear ray marks + find non walls in corners
	src.X, src.Y = 0, 0
	dst.X, dst.Y = 0, 0

	g.Iter(func(p image.Point, c *cell) (next bool) {
		c.Rune = 0

		switch {
		case c.Wall:
		case src.Eq(image.ZP) && p.X < 5 && p.Y < 5:
			src = p
		case dst.Eq(image.ZP) && p.X > 55 && p.Y > 15:
			dst = p
		}

		return true
	})

	// A-Star pathfinding

	fmt.Printf("building path from %s to %s\n", src, dst)

	points, ok := g.Path(
		src,
		dst,
		grid.Points(grid.DirectionsCardinal...),
		grid.DistanceManhattan,
		func(_ image.Point, dist float64, c *cell) (cost float64, walkable bool) {
			return dist, !c.Wall
		},
	)
	if ok {
		for _, p := range points {
			g.MustGet(p).Rune = '.'
		}
	}

	g.MustGet(src).Rune = '@'
	g.MustGet(dst).Rune = 'X'

	showMap(g)

	if ok {
		fmt.Println("path length: ", len(points))
	} else {
		fmt.Println("no path found")
	}

	// Ray-based line-of-sight

	src.X, src.Y = 0, 0

	// clear ray marks + find non wall in center
	g.Iter(func(p image.Point, c *cell) (next bool) {
		c.Rune = 0

		switch {
		case c.Wall:
		case src.Eq(image.ZP) && p.X > 25 && p.X < 35 && p.Y > 7 && p.Y < 13:
			src = p
		}

		return true
	})

	fmt.Printf("\nline of sight from %s\n", src)

	g.LineOfSight(src, 10.0, func(_ image.Point, _ float64, c *cell) (next bool) {
		if c.Wall {
			return false
		}

		c.Rune = '.'

		return true
	})

	g.MustGet(src).Rune = '@'

	showMap(g)

	// Shadow casting

	// clear ray marks + find non wall in center
	g.Iter(func(p image.Point, c *cell) (next bool) {
		c.Rune = 0

		switch {
		case c.Wall:
		case src.Eq(image.ZP) && p.X > 25 && p.X < 35 && p.Y > 7 && p.Y < 13:
			src = p
		}

		return true
	})

	fmt.Printf("\nshadow cast from %s\n", src)

	g.CastShadow(src, 10.0, func(_ image.Point, _ float64, c *cell) (next bool) {
		if c.Wall {
			return false
		}

		c.Rune = '.'

		return true
	})

	g.MustGet(src).Rune = '@'

	showMap(g)

	targets := []image.Point{
		image.Pt(1, 1),
		image.Pt(5, 6),
	}

	dm := g.DijkstraMap(targets, func(_ image.Point, c *cell) (next bool) {
		return !c.Wall
	})

	to, ok := dm.GetTarget(
		src,
		grid.Points(grid.DirectionsCardinal...),
	)

	if ok {
		fmt.Printf("best step from %s is to %s\n", src, to)
	}
}
