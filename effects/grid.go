package effects

import (
	"github.com/moccalotto/nick/field"
)

// Grid represents a grid where the "tiles" have a given width and height
type Grid struct {
	TileWidth  int
	TileHeight int
}

func NewGrid(tileWidth, tileHeight int) *Grid {
	return &Grid{
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
	}
}

func (grid *Grid) ApplyToField(f *field.Field) {

	f.MapAsync(func(f *field.Field, x, y int, c field.Cell) field.Cell {
		if c.Alive() {
			return c
		}

		if x == 0 || y == 0 {
			return c
		}

		if (x % grid.TileWidth) == 0 {
			return field.LivingCell
		}
		if (y % grid.TileHeight) == 0 {
			return field.LivingCell
		}

		return c
	})
}
