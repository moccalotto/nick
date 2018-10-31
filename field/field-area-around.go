package field

func (f *Field) GetAreaAround(x, y int) (Area, error) {
	w := f.Width()
	h := f.Height()
	queue := make(Area, 1, 100)
	queue[0] = Point{x, y}
	inspected := make([]bool, w*h)

	areaType, err := f.Get(x, y)
	if err != nil {
		return Area{}, err
	}

	result := Area{}

	for len(queue) > 0 {
		_p := queue[0]
		queue = queue[1:]

		// anything on the queue can be appended.
		result = append(result, _p)
		inspected[_p.X+_p.Y*w] = true

		for _, c := range _p.Adjacent() {
			// outside the map?
			if !f.CoordsInRange(c.X, c.Y) {
				continue
			}

			// already inspected?
			if inspected[c.X+c.Y*w] {
				continue
			}

			// does this cell belong to another area?
			if f.s[c.X+c.Y*f.w] != areaType {
				continue
			}

			// Point has not yet been looked at (or marked for inspection)
			queue = append(queue, c) // Add c to the queue
			inspected[c.X+c.Y*w] = true
		}
	}

	return result, nil
}
