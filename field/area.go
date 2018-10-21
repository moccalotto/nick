package field

type Area []Point

func (a Area) FindClosestPoints(a2 Area) (Point, Point) {
	bestDistSq := uint64(0xFFFFFFFFFFFF)
	var bestP1 Point
	var bestP2 Point
	for _, p1 := range a {
		for _, p2 := range a2 {
			dx := p1.X - p2.X
			dy := p1.Y - p2.Y
			distSq := uint64(dx*dx + dy*dy)

			if distSq < bestDistSq {
				bestDistSq = distSq
				bestP1 = p1
				bestP2 = p2
			}
		}
	}

	return bestP1, bestP2
}
