[![PkgGoDev](https://pkg.go.dev/badge/github.com/s0rg/grid)](https://pkg.go.dev/github.com/s0rg/grid)
[![License](https://img.shields.io/github/license/s0rg/grid)](https://github.com/s0rg/grid/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/s0rg/grid)](go.mod)
[![Tag](https://img.shields.io/github/v/tag/s0rg/grid?sort=semver)](https://github.com/s0rg/grid/tags)

[![CI](https://github.com/s0rg/grid/workflows/ci/badge.svg)](https://github.com/s0rg/grid/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/s0rg/grid)](https://goreportcard.com/report/github.com/s0rg/grid)
[![Maintainability](https://api.codeclimate.com/v1/badges/8478f67a6b72d9e67cab/maintainability)](https://codeclimate.com/github/s0rg/grid/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/8478f67a6b72d9e67cab/test_coverage)](https://codeclimate.com/github/s0rg/grid/test_coverage)
![Issues](https://img.shields.io/github/issues/s0rg/grid)

# grid

Generic 2D grid

# features

- [DDA RayCasting](https://lodev.org/cgtutor/raycasting.html)
- [A-Star pathfinding](https://en.wikipedia.org/wiki/A*_search_algorithm)
- [Ray-based line of sight](https://en.wikipedia.org/wiki/Line_of_sight_(video_games))
- [Recursive ShadowCasting](http://www.roguebasin.com/index.php/Shadow_casting)
- [Dijkstra maps](http://www.roguebasin.com/index.php/Dijkstra_Maps_Visualized)
- [Bresenham's lines](https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm)
- 100% test cover

# usage

```go
import (
    "image"

    "github.com/s0rg/grid"
)

const mapW, mapH = 100, 100

func valueExample() {
    // working with value-types is straightforward
    g := grid.New[int](image.Rect(0, 0, mapW, mapH))

    // now grid is filled with nil-value for your type
    // you still can re-fill it with some other values:
    g.Fill(func() int {
        return 1
    })
}

func pointerExample() {
    // working with pointer-types is same, but you now you must to pre-fill them
    type mycell struct {}

    g := grid.New[*mycell](image.Rect(0, 0, mapW, mapH))

    // now grid is filled with nil's, so you need pre-fill it with some values,
    // otherwise you will access those nil's with Get / MustGet methods.
    g.Fill(func() *mycell {
        return &mycell{}
    })
}

func usageExample() {
    type mycell struct {
        wall bool
    }

    g := grid.New[*mycell](image.Rect(0, 0, mapW, mapH))

    g.Fill(func() *mycell {
        return &mycell{}
    })

    pt := image.Pt(10, 10)

    // set new value
    g.Set(pt, &mycell{wall: true})

    // update existing value
    if v, ok := g.Get(pt); ok {
        v.wall = false
    }

    // shorthand, for above, will panic on out-of-bounds access
    g.MustGet(pt).wall = true

    // iterate items
    g.Iter(func(p image.Point, c *mycell) (next bool) {
        if c.wall {
            // wall found
        }

        return true
    })
}
```

# example

[Here](https://github.com/s0rg/grid/blob/master/_example/main.go) is a full example.

You can run it with `go run _example/main.go` to see results.

# benchmarks

run:

```bash
make bench
```

results:

```
goos: linux
goarch: amd64
pkg: github.com/s0rg/grid
cpu: AMD Ryzen 5 5500U with Radeon Graphics
BenchmarkGrid/Set-12         	1000000000	        0.8108 ns/op	      0 B/op	      0 allocs/op
BenchmarkGrid/Get-12         	641611768	        1.764 ns/op	      0 B/op	      0 allocs/op
BenchmarkGrid/Neighbours-12  	52243890	       23.41 ns/op	      0 B/op	      0 allocs/op
BenchmarkGrid/LineBresenham-12         	4416172	      269.0 ns/op	      0 B/op	      0 allocs/op
BenchmarkGrid/CastRay-12               	3829839	      321.1 ns/op	      0 B/op	      0 allocs/op
BenchmarkGrid/CastShadow-12            	  32648	    36950 ns/op	      0 B/op	      0 allocs/op
BenchmarkGrid/LineOfSight-12           	   9897	   114576 ns/op	      0 B/op	      0 allocs/op
BenchmarkGrid/DijkstraMap-12           	   1029	  1190195 ns/op	  20656 B/op	      3 allocs/op
BenchmarkGrid/Path-12                  	    372	  3225325 ns/op	 997588 B/op	  13643 allocs/op
PASS
ok  	github.com/s0rg/grid	12.098s
```
