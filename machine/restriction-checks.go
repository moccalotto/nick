package machine

import (
	"fmt"
	"time"
)

// checkRestrictions returns an error if m.Limits is not adhered to
func (m *Machine) checkRestrictions() error {
	if err := m.timedOut(); err != nil {
		return err
	}

	if err := m.tooManyCells(); err != nil {
		return err
	}

	if err := m.tooWide(); err != nil {
		return err
	}

	if err := m.tooTall(); err != nil {
		return err
	}

	return nil
}

// timeOut returns an error if a machine's execution has taken too long.
// the time Limits is given by Machine.MaxRunTime
func (m *Machine) timedOut() error {
	// no max time specified, i.e. we can go on forever.
	if m.MaxRuntime == 0 {
		return nil
	}

	runtime := time.Now().Sub(m.StartedAt)

	if runtime > m.MaxRuntime {
		return fmt.Errorf("Timed out after %f seconds", runtime.Seconds())
	}

	return nil
}

// tooManyCells returns an error if the field has too many cells
// the Limits on cell count is given in Machine.MaxCells
func (m *Machine) tooManyCells() error {
	// we don't have a field yet, so it can't be too large
	if m.Field == nil {
		return nil
	}

	w := m.Field.Width()
	h := m.Field.Height()
	max := m.MaxCells

	// there is no maximum number of cells
	if max <= 0 {
		return nil
	}

	if w*h > max {
		return fmt.Errorf(
			"Grid is too large. Max number of cells allowed is %d, but the current size is (%dx%d) %d",
			max,
			w,
			h,
			w*h,
		)
	}

	return nil
}

func (m *Machine) tooWide() error {
	// we don't have a field yet, so it can't be too large
	if m.Field == nil {
		return nil
	}

	w := m.Field.Width()
	max := m.MaxWidth

	// there is no maximum width
	if max <= 0 {
		return nil
	}

	if w > max {
		return fmt.Errorf(
			"Grid is too large. Max width is %d, but the current width is %d",
			max,
			w,
		)
	}

	return nil
}

func (m *Machine) tooTall() error {
	// we don't have a field yet, so it can't be too large
	if m.Field == nil {
		return nil
	}

	h := m.Field.Height()
	max := m.MaxHeight

	// there is no maximum height
	if max <= 0 {
		return nil
	}

	if h > max {
		return fmt.Errorf(
			"Grid is too tall. Max height is %d, but the current height is %d",
			max,
			h,
		)
	}

	return nil
}
