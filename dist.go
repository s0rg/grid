package grid

import (
	"image"
	"math"

	"github.com/s0rg/vec2d"
)

// DistanceEuclidean calculates Euclidean (the shortest) distance between two points.
func DistanceEuclidean(a, b image.Point) (rv float64) {
	var (
		va = vec2d.New(float64(a.X), float64(a.Y))
		vb = vec2d.New(float64(b.X), float64(b.Y))
	)

	return va.Sub(vb).Len()
}

// DistanceManhattan calculates Manhattan distance between two points.
func DistanceManhattan(a, b image.Point) (rv float64) {
	var (
		dx = float64(a.X - b.X)
		dy = float64(a.Y - b.Y)
	)

	return math.Abs(dx) + math.Abs(dy)
}

// DistanceChebyshev calculates Chebyshev distance between two points.
func DistanceChebyshev(a, b image.Point) (rv float64) {
	var (
		dx = float64(a.X - b.X)
		dy = float64(a.Y - b.Y)
	)

	return math.Max(math.Abs(dx), math.Abs(dy))
}
