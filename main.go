package main

import (
	"github.com/moccalotto/nick/exporters"
	"github.com/moccalotto/nick/field"
	"github.com/moccalotto/nick/machine"
	"math/rand"
	"time"
)

func makeField() *field.Field {
	m := machine.MachineFromScript(`
		suggest export.type      = image	# We would like to export an image
		suggest export.format    = png		# in png
		suggest export.width     = 1600		# with a fixed width
		suggest export.algorithm = box		# using the »box« scaling method

		init 25 x 20		# new canvas
		snow 31%		# add 40% snow
		border 1		# Add a 1-cell border all the way around
		evolve B5678/S345678	# run standard escavator
		loop 3
			scale 2
			loop 3
				evolve B5678/S345678	# run standard escavator
			endloop
		endloop 

		loop 4
			scale 2
			loop 3
				evolve B5678/S5678	# run edge smoother
			endloop
		endloop

		# evolve B05678/S05678     # This automaton is good for strengthening edges.
		# evolve B012345/S012345   # Running an "inverse" automaton inverses the map
	`)

	m.Execute()

	return m.Field
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	f := makeField()

	ie := exporters.NewImageExporter()
	// ie.Width = 1000
	ie.Algorithm = "Box"
	ie.Export(f)                             // export an image.
	exporters.NewItermExporter(ie).Export(f) // export image to iterm.
}
