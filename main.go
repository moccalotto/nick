package main

import (
	"fmt"
	"github.com/moccalotto/nick/exporters"
	"github.com/moccalotto/nick/machine"
	"time"
)

func createMachine() *machine.Machine {
	m := machine.MachineFromScript(`
		suggest export.type      = iterm	# We would like to export an image
		suggest export.format    = png		# in png
		# suggest export.width     = 1600		# with a fixed width
		suggest export.algorithm = Lanczos	# using the »box« scaling method

		# Small initial sizes often yield simple caves
		# Large initial sizes yields multi-room or complex caves with stalactites/columns, etc.
		init 25 x 20		# New canvas. Small initial canvas yields simple caves
		snow 31%		# Add random dots.
		border 1 80%    	# Cover the border with snow at a 80 percent density

		loop 3
			set-rand-int $width = 6 to 10	# determine the width of the egress
			set-rand-int $depth  = 1 to 3	# determine the depth of the egress
			egress random $width x $depth	# create an opening

			evolve B5678/S345678	# run standard escavator
			evolve B5678/S345678	# run standard escavator

			scale 1.75 x 1.5
		endloop

		loop 7
			scale 1.5
			evolve B5678/S5678	# run edge smoother
			evolve B5678/S5678	# run edge smoother
		endloop

		# evolve B2345678/S0 		# Produce an outline. Living edge on dead background
		# border 1 100% (dead)		# The outliner above will remove any egress. We revive it by killing the border.

		# evolve B05678/S05678		# This automaton reduces the map to just the edges (but also reverses it).
		# evolve B012345/S012345	# Running an "inverse" automaton re-inverses the map

	`)

	return m
}

func main() {
	m := createMachine()
	m.Limits.MaxRuntime, _ = time.ParseDuration("5s")
	start := time.Now()
	err := m.Execute()
	elapsed := time.Now().Sub(start).Seconds()

	if err != nil {
		panic(err)
	}

	fallback := exporters.NewTextExporter()
	e := exporters.NewSuggestionExporter(m.Vars, fallback)
	e.Export(m.Field) // export an image.

	fmt.Printf("Seed: %d\n", m.Seed)
	fmt.Printf("Time elapsed: %f\n", elapsed)
}
