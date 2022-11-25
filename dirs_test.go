package grid

import (
	"testing"
)

func TestDirsInvert(t *testing.T) {
	t.Parallel()

	for _, d := range DirectionsALL {
		if d != d.Invert().Invert() {
			t.Fail()
		}
	}
}

func TestDirsPoints(t *testing.T) {
	t.Parallel()

	if p := Points(North, East); len(p) != 2 {
		t.Fail()
	}
}
