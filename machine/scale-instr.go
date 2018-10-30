package machine

import (
	"github.com/moccalotto/nick/effects"
)

func init() {
	InstructionHandlers["scale"] = Scale
	InstructionHandlers["scale-to"] = ScaleTo
	InstructionHandlers["scale-to-width"] = ScaleToWidth
	InstructionHandlers["scale-to-height"] = ScaleToHeight
}

func Scale(m *Machine) {
	m.Assert(m.Field != nil, "Cannot scale a non-initialized field!")
	x := m.ArgAsFloat(0)
	y := x
	if m.ArgCount() > 1 {
		y = m.ArgAsFloat(2)
		m.Assert(m.ArgAsString(1) == "x", "Second arg to scale must be an 'x', but '%s' was given", m.ArgAsString(1))

	}
	effects.NewScaleXY(x, y).ApplyToField(m.Field)
}

func ScaleTo(m *Machine) {
	m.Assert(m.Field != nil, "Cannot scale a non-initialized field!")

	newW := m.ArgAsInt(0)
	m.Assert(m.ArgAsString(1) == "x", "Second arg to scale must be an 'x', but '%s' was given", m.ArgAsString(1))
	newH := m.ArgAsInt(2)

	effects.NewScaleTo(
		m.Field.Width(),
		newW,
		m.Field.Height(),
		newH,
	).ApplyToField(m.Field)
}

func ScaleToWidth(m *Machine) {
	m.Assert(m.Field != nil, "Cannot scale a non-initialized field!")

	aspect := m.Field.AspectRatio()

	newW := m.ArgAsInt(0)

	newH := int(aspect / float64(newW))

	effects.NewScaleTo(
		m.Field.Width(),
		newW,
		m.Field.Height(),
		newH,
	).ApplyToField(m.Field)
}

func ScaleToHeight(m *Machine) {
	m.Assert(m.Field != nil, "Cannot scale a non-initialized field!")

	aspect := m.Field.AspectRatio()

	newH := m.ArgAsInt(0)

	newW := int(aspect * float64(newH))

	effects.NewScaleTo(
		m.Field.Width(),
		newW,
		m.Field.Height(),
		newH,
	).ApplyToField(m.Field)
}
