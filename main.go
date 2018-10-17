package main

import (
	"fmt"
	"github.com/moccalotto/nick/exporters"
	"github.com/moccalotto/nick/machine"
	"io/ioutil"
	"time"
)

func createMachine(filename string) *machine.Machine {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	script := string(b)
	m := machine.MachineFromScript(script)

	return m
}

func main() {
	var m *machine.Machine
	m = createMachine("example.cave")
	m.Limits.MaxRuntime, _ = time.ParseDuration("5s")

	if err := m.Execute(); err != nil {
		panic(err)
	}

	e := exporters.NewSuggestionExporter(
		m.Vars,
		exporters.NewTextExporter(),
	)

	e.Export(m.Field) // export an image.

	fmt.Printf("Seed: %d\n", m.Seed)
}
