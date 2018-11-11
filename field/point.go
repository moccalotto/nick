package field

type Point struct {
	X, Y int
}

func (p1 Point) DistSq(p2 Point) int {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return dx*dx + dy*dy
}

// OrthogonalAdjacent returns the 4 points to the north, east, south and west of p
func (p Point) OrthogonalAdjacent() Area {
	return []Point{
		Point{p.X + 1, p.Y}, // east
		Point{p.X - 1, p.Y}, // west
		Point{p.X, p.Y + 1}, // north
		Point{p.X, p.Y - 1}, // south
	}
}

func (p Point) WithinRadius(r float64) Area {
	r2 := r * r

	result := Area{Point{p.X, p.Y}}

	// Iterate over a an eighth of the area.:
	// All the points generated by [i, j] lie inside a 45 right-angled triangle,
	// with it's pointy end positioned at P.
	// If we disregard the points in that triangle, that are outside the
	// radius around P, we have found an eighth of the points in the circle.
	// We can now infer the remaining points due to symmetry.
	for i := 0; float64(i) <= r; i++ {
		for j := 0; j <= i; j++ {
			// d2 is the squared distance from p.
			d2 := float64(i*i + j*j)

			// point outside radius
			if d2 > r2 {
				continue
			}

			result = append(
				result,
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
