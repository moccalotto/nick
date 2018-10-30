package field

func (f *Field) GetAreaAround(x, y int) (Area, error) {
	queue := Area{Point{x, y}}
	inspected := NewField(f.Width(), f.Height())
	areaStatus, err := f.On(x, y)
	if err != nil {
		return Area{}, err
	}

	result := Area{}

	for len(queue) > 0 {
		_p := queue[0]
		queue = queue[1:]

		// anything on the queue can be appended.
		result = append(result, _p)
		_ = inspected.SetOn(_p.X, _p.Y, true)

		for _, c := range _p.Adjecent() {
			// outside the map?
			if !f.CoordsInRange(c.X, c.Y) {
				continue
			}

			// already inspected?
			if a, _ := inspected.On(c.X, c.Y); a {
				continue
			}

			// does this cell belong to another area?
			if a, _ := f.On(c.X, c.Y); a != areaStatus {
				continue
			}

			// Point has not yet been looked at (or marked for inspection)
			queue = append(queue, c)            // Add c to the queue
			_ = inspected.SetOn(c.X, c.Y, true) // Mark c as inspected.
		}
	}

	return result, nil
}
