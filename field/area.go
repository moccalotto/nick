package field

type Area []Point

func (this Area) FindClosestPoints(other Area) (Point, Point, uint64) {
	bestDistSq := uint64(0xFFFFFFFFFFFF)
	var bestP1 Point
	var bestP2 Point
	for _, p1 := range this {
		for _, p2 := range other {
			distSq := uint64(p1.DistSq(p2))

			if distSq == 0 {
				return p1, p2, distSq
			}

			if distSq < bestDistSq {
				bestDistSq = distSq
				bestP1 = p1
				bestP2 = p2
			}
		}
	}

	return bestP1, bestP2, bestDistSq
}
