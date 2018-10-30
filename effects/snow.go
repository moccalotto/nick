package effects

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

// Randomly bring cells to life.
// Each cell has »Coverage« chance to be born.
// NOTE: cells do not die via this method, they are only brought to Coverage.
type Snow struct {
	Coverage float64
	Cell     field.Cell
	rng      *rand.Rand
}

func NewSnow(p float64, rng *rand.Rand) *Snow {
	return &Snow{p, field.LivingCell, rng}
}

// Turn on cells randomly throughout the field.
func (s *Snow) ApplyToField(f *field.Field) {
	// We cannot generate snow asynchronously because it would intorduce
	// actual randomness due to race conditions.
	f.Map(func(f *field.Field, x, y int, c field.Cell) field.Cell {
		if s.rng.Float64() < s.Coverage {
			return s.Cell
		}
		return c
	})
}
