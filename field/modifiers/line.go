package modifiers

import (
	"github.com/moccalotto/nick/field"
	"math/rand"
)

type Line struct {
	X0       int
	Y0       int
	X1       int
	Y1       int
	Coverage float64
	Alive    bool
	rng      *rand.Rand
}

func NewLine(startX, startY, endX, endY int, rng *rand.Rand) *Line {
	return &Line{
		X0:       startX,
		Y0:       startY,
		X1:       endX,
		Y1:       endY,
		Coverage: 1.0,
		Alive:    true,
		rng:      rng,
	}
}

// The snow will now add dead cells instead of living cells
func (b *Line) Inverted(dead bool) *Line {
	tmp := *b
	tmp.Alive = !dead

	return &tmp
}

func (l *Line) plot(f *field.Field, x, y int) {
	if f.CoordsInRange(x, y) && l.rng.Float64() < l.Coverage {
		f.SetAlive(x, y, l.Alive)
	}
}

// draw a steep (more than 45 degrees off the x-axis) line between (x0, y1) and (x1, y1)
func (l *Line) plotLineSteep(f *field.Field, x0, y0, x1, y1 int) {
	dx := x1 - x0
	dy := y1 - y0
	x := x0
	xi := 1
	if dx < 0 {
		xi = -1
		dx = -dx
	}

	D := 2*dx - dy

	for y := y0; y <= y1; y++ {
		l.plot(f, x, y)

		if D > 0 {
			x += xi
			D -= 2 * dy
		}

		D += 2 * dx
	}
}

// draw a shallow (less than 45 degrees off the x-axis) line between (x0, y1) and (x1, y1)
func (l *Line) plotLineShallow(f *field.Field, x0, y0, x1, y1 int) {
	dx := x1 - x0
	dy := y1 - y0
	y := y0
	yi := 1
	if dy < 0 {
		yi = -1
		dy = -dy
	}

	D := 2*dy - dx

	for x := x0; x <= x1; x++ {
		l.plot(f, x, y)

		if D > 0 {
			y = y + yi
			D -= 2 * dx
		}
		D += 2 * dy
	}
}

func (l *Line) ApplyToField(f *field.Field) {
	dx := l.X1 - l.X0
	dy := l.Y1 - l.Y0

	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	if dy < dx { // line is steeper than 45 degrees
		if l.X0 > l.X1 {
			// draw shallow line from p1 to p0
			l.plotLineShallow(f, l.X1, l.Y1, l.X0, l.Y0)
		} else {
			// draw shallow line from p0 to p1
			l.plotLineShallow(f, l.X0, l.Y0, l.X1, l.Y1)
		}
	} else { // line is shallower than 45 degrees
		if l.Y0 > l.Y1 {
			// draw steep line from p1 to p0
			l.plotLineSteep(f, l.X1, l.Y1, l.X0, l.Y0)
		} else {
			// draw steep line from p0 to p1
			l.plotLineSteep(f, l.X0, l.Y0, l.X1, l.Y1)
		}
	}
}
