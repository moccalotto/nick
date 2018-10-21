package field

type Point struct {
	X, Y int
}

func (p Point) Adjecent() Area {
	return []Point{
		Point{p.X + 1, p.Y}, // east
		Point{p.X - 1, p.Y}, // west
		Point{p.X, p.Y + 1}, // north
		Point{p.X, p.Y - 1}, // south
	}
}

func (p Point) WithinRadius(r float64) Area {
	r2 := r * r

	result := Area{}

	// iterate over a quarter of the area.
	for i := 0; float64(i) <= r; i++ {
		for j := 0; j <= i; j++ {
			d2 := i*i + j*j

			// point outside radius
			if float64(d2) > r2 {
				continue
			}

			result = append(
				result,
				Point{p.X, p.Y},
				Point{p.X + i, p.Y + j},
				Point{p.X + i, p.Y - j},
				Point{p.X - i, p.Y + j},
				Point{p.X - i, p.Y - j},
				Point{p.X + j, p.Y + i},
				Point{p.X + j, p.Y - i},
				Point{p.X - j, p.Y + i},
				Point{p.X - j, p.Y - i},
			)
		}
	}

	return result
}
