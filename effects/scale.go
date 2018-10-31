package effects

import (
	"github.com/moccalotto/nick/field"
	"sync"
)

type Scale struct {
	x, y float64
}

func NewScale(f float64) *Scale {
	return &Scale{f, f}
}

func NewScaleXY(x, y float64) *Scale {
	return &Scale{x, y}
}

func NewScaleTo(startW, newW, startH, newH int) *Scale {
	x := float64(newW) / float64(startW)
	y := float64(newH) / float64(startH)

	return &Scale{x, y}
}

func (s *Scale) ApplyToField(f *field.Field) {
	newW := int(float64(f.Width()) * s.x)
	newH := int(float64(f.Height()) * s.y)
	offsets := make([]int, newW)
	tmp := make([]field.Cell, newH*newW)

	// pre-calculate a map between new x-values and old x-values
	for x := range offsets {
		offsets[x] = int(float64(x) / s.x)
	}

	var wg sync.WaitGroup

	for newY := 0; newY < newH; newY++ {
		oldY := int(float64(newY) / s.y)
		wg.Add(1)
		go func(newY, oldY int) {
			defer wg.Done()
			for newX, oldX := range offsets {
				tmp[newX+newY*newW], _ = f.Get(oldX, oldY)
			}
		}(newY, oldY)
	}

	wg.Wait()

	f.ReplaceCells(newW, newH, tmp)
}
