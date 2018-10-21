package effects

import (
	"github.com/moccalotto/nick/field"
	"math"
)

type gauss struct {
	x int
	y int
}

type ThielePattern struct {
	Base  complex128
	Fill  bool // should we fill all of the map?
	Alive bool
	a     int // real composant
	b     int // imaginary composant
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func cFloor(n complex128) complex128 {
	return complex(
		math.Floor(real(n)),
		math.Floor(imag(n)),
	)
}

func cGauss(n complex128) gauss {
	return gauss{
		int(math.Round(real(n))),
		int(math.Round(imag(n))),
	}
}

func cRemaind(a, b complex128) gauss {
	div := a / b
	intpart := cFloor(div)
	floatpart := div - intpart

	return cGauss(floatpart * b)
}

func normalizeIntoQuadrant1(a, b int) (int, int) {
	if a < 0 {
		a = -a
	}

	if b < 0 {
		b = -b
	}

	if a == 0 {
		a = b
		b = 0
	}

	return a, b
}

func NewThielePattern(a, b int) *ThielePattern {
	a, b = normalizeIntoQuadrant1(a, b)

	return &ThielePattern{
		complex(float64(a), float64(b)),
		true,
		true,
		a,
		b,
	}
}

// Get the points we need to investigate
// we iterate over the area covered by the base (as a vector)
// and the transposed base (as a vector).
func (t *ThielePattern) pointsToInspect() <-chan complex128 {
	ch := make(chan complex128)
	go func() {
		max := absInt(t.a) + absInt(t.b)
		// naive solution: iterate over a suitably large area,
		// don't mind where the area is located.
		for y := 0; y <= max; y++ {
			for x := 0; x <= max; x++ {
				ch <- complex(float64(x), float64(y))
			}
		}
		close(ch)
	}()

	return ch
}

func (t *ThielePattern) residueNumbers() []gauss {
	numbers := []gauss{}

	for c := range t.pointsToInspect() {
		// square the point, and modulo it into the base area
		p := cRemaind(c*c, t.Base)
		alreadyFound := false

		for _, n := range numbers {
			if n == p {
				alreadyFound = true
				break
			}
		}

		if alreadyFound {
			continue
		}

		numbers = append(numbers, p)
	}

	return numbers
}

func (t *ThielePattern) mapper(f *field.Field, x, y int, c field.Cell) field.Cell {

	numbers := t.residueNumbers()

	c0 := complex(float64(x), float64(y))
	c1 := cRemaind(c0, t.Base)

	for _, c := range numbers {
		if c == c1 {
			return field.LivingCell
			break
		}
	}

	return c
}

func (t *ThielePattern) fillField(f *field.Field) {
	numbers := t.residueNumbers()

	f.MapAsync(func(f *field.Field, x, y int, c field.Cell) field.Cell {
		p := complex(float64(x), float64(y))
		remaind := cRemaind(p, t.Base)

		for _, baseRemainder := range numbers {
			if baseRemainder == remaind {
				return field.LivingCell
				break
			}
		}

		return c
	})
}

func (t *ThielePattern) ApplyToField(f *field.Field) {

	if t.Fill {
		t.fillField(f)
		return
	}

	for _, c := range t.residueNumbers() {
		x := c.x + f.Width()/2
		y := c.y + f.Height()/2
		_ = f.SetAlive(x, y, t.Alive)
	}
}