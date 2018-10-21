package effects

import (
	"github.com/moccalotto/nick/field"
	"math"
)

type ThielePattern struct {
	Base complex128
	Cell field.Cell
	a    int // real composant
	b    int // imaginary composant
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

func cGauss(n complex128) field.Point {
	return field.Point{
		int(math.Round(real(n))),
		int(math.Round(imag(n))),
	}
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
		field.LivingCell,
		a,
		b,
	}
}

func (t *ThielePattern) residueNumbers() field.Area {
	result := field.Area{}
	max := absInt(t.a) + absInt(t.b)

	for y := 0; y <= max; y++ {
		for x := 0; x <= max; x++ {
			c := complex(float64(x), float64(y))
			// square the point, and modulo it into the base area
			p := t.remaind(c * c)
			alreadyFound := false

			for _, n := range result {
				if n == p {
					alreadyFound = true
					break
				}
			}

			if alreadyFound {
				continue
			}

			result = append(result, p)
		}
	}

	return result
}

func (t *ThielePattern) mapper(f *field.Field, x, y int, c field.Cell) field.Cell {

	numbers := t.residueNumbers()

	c0 := complex(float64(x), float64(y))
	c1 := t.remaind(c0)

	for _, c := range numbers {
		if c == c1 {
			return field.LivingCell
			break
		}
	}

	return c
}

func (t *ThielePattern) remaind(c complex128) field.Point {
	div := c / t.Base
	intpart := cFloor(div)
	floatpart := div - intpart
	return cGauss(floatpart * t.Base)
}

func (t *ThielePattern) ApplyToField(f *field.Field) {
	numbers := t.residueNumbers()

	f.MapAsync(func(f *field.Field, x, y int, c field.Cell) field.Cell {
		p := complex(float64(x), float64(y))
		remaind := t.remaind(p)

		for _, baseRemainder := range numbers {
			if baseRemainder == remaind {
				return t.Cell
				break
			}
		}

		return c
	})
}
