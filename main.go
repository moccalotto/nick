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
	m := createMachine("example.cave")
	m.Limits.MaxRuntime, _ = time.ParseDuration("5s")
	start := time.Now()
	err := m.Execute()

	if err != nil {
		panic(err)
	}

	elapsed := time.Now().Sub(start).Seconds()

	fallback := exporters.NewTextExporter()
	e := exporters.NewSuggestionExporter(m.Vars, fallback)
	e.Export(m.Field) // export an image.

	fmt.Printf("Seed: %d\n", m.Seed)
	fmt.Printf("Time elapsed: %f\n", elapsed)
}
