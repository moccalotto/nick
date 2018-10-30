package effects

import (
	"github.com/moccalotto/nick/field"
	"math"
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
	nw := int(math.Round(float64(f.Width()) * s.x))
	nh := int(math.Round(float64(f.Height()) * s.y))
	tmp := field.NewField(nw, nh)

	var wg sync.WaitGroup

	for y := 0; y < nh; y++ {
		_y := int(math.Floor(float64(y) / s.y))
		wg.Add(1)
		go func(y, _y int) {
			defer wg.Done()
			for x := 0; x < nw; x++ {
				_x := int(math.Floor(float64(x) / s.x))
				if a, _ := f.On(_x, _y); a {
					_ = tmp.SetOn(x, y, true)
				}
			}
		}(y, _y)
	}

	wg.Wait()

	f.ReplaceCells(nw, nh, tmp.Cells())
}
