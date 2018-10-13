package machine

import (
	"github.com/moccalotto/nick/field"
	"testing"
)

func TestEvolve(t *testing.T) {
	script := `evolve B3/S23`

	m := MachineFromScript(script)
	m.Field = field.NewField(5, 5)

	// draw a horizontal bar.
	// the horizontal bar will become vertical after running B2/S23
	m.Field.Set(1, 2, true)
	m.Field.Set(2, 2, true)
	m.Field.Set(3, 2, true)

	Evolve(m)

	success := m.Field.Alive(2, 1) &&
		m.Field.Alive(2, 2) &&
		m.Field.Alive(2, 3) &&
		m.Field.Dead(1, 2) &&
		m.Field.Dead(3, 2)

	if !success {
		t.Errorf("Expected the horizontal bar to become vertical: %+v", m.Field)
	}

	// run the same instruction again, reverting back to the original horizontal bar
	Evolve(m)

	success = m.Field.Alive(2, 2) &&
		m.Field.Alive(1, 2) &&
		m.Field.Alive(3, 2) &&
		m.Field.Dead(2, 3) &&
		m.Field.Dead(2, 1)

	if !success {
		t.Errorf("Expected the vertical bar to become horizontal: %+v", m.Field)
	}

}
