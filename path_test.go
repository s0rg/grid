package grid

import (
	"image"
	"testing"
)

func TestPath(t *testing.T) {
	t.Parallel()

	const N = 10

	var (
		p path
		o image.Point
	)

	for i := 0; i < N; i++ {
		o = image.Pt(i, i)

		p = p.Fork(o, float64(i))
	}

	if p.Len() != N {
		t.Fail()
	}

	if l := p.Last(); l.X != N-1 || l.Y != N-1 {
		t.Fail()
	}
}
