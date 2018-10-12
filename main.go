package main

import (
	tm "github.com/buger/goterm"
	"github.com/moccaloto/nick/field"
	Mod "github.com/moccaloto/nick/field/modifiers"
	"github.com/moccaloto/nick/field/printers"
	"github.com/moccaloto/nick/machine"
	"math/rand"
	"time"
)

func makeField() *field.Field {
	m := machine.MachineFromScript(`
		init 50 x 40		# new canvas
		snow 40%		# add 40% snow
		evolve B5678/S345678	# run standard escavator
		scale 1.3		# scale up by 30%
		snow 0.1 (dead)		# add 10% holes
		evolve B5678/S345678	# run standard escavator
		scale 150% x 270%	# scale up
		evolve B5678/S5678	# run edge smoother
		scale 2.7 x 1.5		# scale up
		evolve B5678/S5678	# run edge smoother
		scale 2.3		# scape up
		loop 2
			log Loop: @loop
			evolve B5678/S5678	# run edge smoother
			exit			# exit - so the loop only runs once
		endloop
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
