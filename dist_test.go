package grid

import (
	"image"
	"testing"
)

func TestDistanceEuclidean(t *testing.T) {
	t.Parallel()

	if l := DistanceEuclidean(
		image.Pt(0, 0),
		image.Pt(10, 10),
	); l < 10.0 {
		t.Fail()
	}
}

func TestDistanceManhattan(t *testing.T) {
	t.Parallel()

	if l := DistanceManhattan(
		image.Pt(0, 0),
		image.Pt(10, 10),
	); l < 10.0 {
		t.Fail()
	}
}

func TestDistanceChebyshev(t *testing.T) {
	t.Parallel()

	if l := DistanceChebyshev(
		image.Pt(0, 0),
		image.Pt(10, 10),
	); l < 10.0 {
		t.Fail()
	}
}
