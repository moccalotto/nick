package main

import (
	tm "github.com/buger/goterm"
	"github.com/moccalotto/nick/field"
	Mod "github.com/moccalotto/nick/field/modifiers"
	"github.com/moccalotto/nick/field/printers"
	"github.com/moccalotto/nick/machine"
	"math/rand"
	"time"
)

func makeField() *field.Field {
	m := machine.MachineFromScript(`
		init 25 x 20		# new canvas
		snow 31%		# add 40% snow
		border 1		# Add a 1-cell border all the way around
		evolve B5678/S345678	# run standard escavator
		loop 3
			scale 1.75		# scale up by 50% in both directions
			evolve B5678/S345678	# run standard escavator
			evolve B5678/S345678	# run standard escavator
		endloop 
		loop 3
			scale 1.8
			evolve B5678/S5678	# run edge smoother
		endloop
		evolve B05678/S05678     # This automaton is good for strengthening edges.
		evolve B012345/S012345   # Running an "inverse" automaton inverses the map
	`)

	m.Execute()

	return m.Field

	f := field.NewField(40, 20)
	// flip := Mod.NewAutomaton("B0/S") 	  // Flip
	// r := Mod.NewAutomaton("B3/S23") 	  // Conway game of life rules
	// r := Mod.NewAutomaton("B5678/S5678")   // Edge Smoother
	// r := Mod.NewAutomaton("B35678/S5678")  // []int{3, 5, 6, 7, 8}, []int{5, 6, 7, 8}) // Diamoeba CA
	// r := Mod.NewAutomaton("B36/S125")      // 2x2 CA
	// r := Mod.NewAutomaton("B2/S")          // Seeds CA
	f.Apply(Mod.NewSnow(0.4))
	f.Apply(Mod.NewBorder())
	// f.Apply(Mod.NewEgress(Mod.NorthWest, 13).WithThickness(0))
	// f.Apply(Mod.NewEgress(Mod.Random, 9).WithThickness(1))
	// f.Apply(Mod.NewEgress(Mod.Random, 5).WithThickness(2))
	// f.Apply(Mod.NewEgress(Mod.Random, 1).WithThickness(3))

	f.Apply(Mod.NewAutomaton("B5678/S345678"))
	f.Apply(Mod.NewAutomaton("B5678/S345678"))

	f.Apply(Mod.NewScaleXY(2.3, 3.3))
	f.Apply(Mod.NewAutomaton("B5678/S5678"))

	f.Apply(Mod.NewScaleXY(1.8, 1.3))
	f.Apply(Mod.NewAutomaton("B5678/S5678"))

	f.Apply(Mod.NewScale(1.7))
	f.Apply(Mod.NewAutomaton("B5678/S5678"))

	// The two effects below will create a sort of "drawing" effect.
	// f.Apply(Mod.NewAutomaton("B05678/S05678"))	// high contrast pencil-like edges, but also inverses
	// f.Apply(Mod.NewAutomaton("B012345/S012345"))	// invert back

	return f
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	tm.Clear()
	tm.MoveCursor(1, 1)
	img := printers.ItermImage(makeField())
	tm.Print(img)
	tm.Flush()
}
