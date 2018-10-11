package main

import (
	tm "github.com/buger/goterm"
	"github.com/moccaloto/nick/field"
	Mod "github.com/moccaloto/nick/field/modifiers"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	f := field.NewField(150, 37)
	// flip := Mod.NewAutomaton("B0/S") 	  // Flip
	gameOfLife := Mod.NewAutomaton("B3/S23") 	  // Conway game of life rules
	// r := Mod.NewAutomaton("B5678/S5678")   // Edge Smoother
	// r := Mod.NewAutomaton("B35678/S5678")  // []int{3, 5, 6, 7, 8}, []int{5, 6, 7, 8}) // Diamoeba CA
	// r := Mod.NewAutomaton("B36/S125")      // 2x2 CA
	// r := Mod.NewAutomaton("B2/S")          // Seeds CA
	f.Apply(Mod.NewSnow(0.1))

	// glider
	f.Set(3, 1, true)
	f.Set(1, 2, true)
	f.Set(3, 2, true)
	f.Set(2, 3, true)
	f.Set(3, 3, true)

	for {
		f.Apply(gameOfLife)
		tm.Clear()
		tm.MoveCursor(1, 1)
		tm.Println(f.ToText())
		tm.Println(gameOfLife.String())
		tm.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}
