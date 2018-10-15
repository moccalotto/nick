package main

import (
	"github.com/moccalotto/nick/exporters"
	"github.com/moccalotto/nick/machine"
	"time"
)

func createMachine() *machine.Machine {
	m := machine.MachineFromScript(`
		suggest export.type      = iterm	# We would like to export an image
		suggest export.format    = png		# in png
		suggest export.width     = 1600		# with a fixed width
		suggest export.algorithm = Lanczos	# using the »box« scaling method

		init 25 x 20		# new canvas
		snow 31%		# add 40% snow
		border 2 @ 85%		# Cover the border with snow at a 85% density
		egress north @ 10 x 2	# create an entrence to the north 10 cells wide and 2 cells thick
		evolve B5678/S345678	# run standard escavator
		loop 3
			scale 1.5
			loop 3
				evolve B5678/S345678	# run standard escavator
			endloop
		endloop 

		loop 7
			scale 1.5
			loop 2
				evolve B5678/S5678	# run edge smoother
			endloop
		endloop

		# evolve B05678/S05678     # This automaton is good for strengthening edges.
		# evolve B012345/S012345   # Running an "inverse" automaton inverses the map
	`)

	return m
}

func main() {
	m := createMachine()
	m.Limits.MaxRuntime, _ = time.ParseDuration("5s")
	err := m.Execute()

	if err != nil {
		panic(err)
	}

	fallback := exporters.NewTextExporter()
	e := exporters.NewSuggestionExporter(m.Vars, fallback)
	e.Export(m.Field) // export an image.
}
