package main

import (
	tm "github.com/buger/goterm"
	"github.com/moccaloto/nick/field"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	f := field.NewField(150, 37)
	r := field.NewRules("B0/S") // Flip
	// r := field.NewRules("B3/S23") // Conway game of life rules
	// r := field.NewRules("B5678/S5678")   // Edge Smoother
	// r := field.NewRules("B35678/S5678")  // []int{3, 5, 6, 7, 8}, []int{5, 6, 7, 8}) // Diamoeba CA
	// r := field.NewRules("B36/S125")      // 2x2 CA
	// r := field.NewRules("B2/S")          // Seeds CA
	f.Seed(0.5)

	// glider
	f.Set(3, 1, true)
	f.Set(1, 2, true)
	f.Set(3, 2, true)
	f.Set(2, 3, true)
	f.Set(3, 3, true)

	for {
		f = f.Evolve(r)
		tm.Clear()
		tm.MoveCursor(1, 1)
		tm.Println(f.ToText())
		tm.Println(r.String())
		tm.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}
