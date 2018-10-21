package field

type Point struct {
	X, Y int
}

type Area []Point

func (p Point) Adjecent() Area {
	return []Point{
		Point{p.X + 1, p.Y}, // east
		Point{p.X - 1, p.Y}, // west
		Point{p.X, p.Y + 1}, // north
		Point{p.X, p.Y - 1}, // south
	}
}
