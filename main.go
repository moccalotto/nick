package main

import (
	"flag"
	"fmt"
	"github.com/moccalotto/nick/exporters"
	"github.com/moccalotto/nick/machine"
	"io/ioutil"
	"log"
	"time"
)

func createMachine(filename string) *machine.Machine {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file: '%s'", filename)
	}

	script := string(b)
	m := machine.MachineFromScript(script)

	return m
}

func main() {
	f := flag.String("script", "examples/empty.cave", "Path to script to execute")
	flag.Parse()

	if *f == "" {
		flag.PrintDefaults()
		return
	}

	m := createMachine(*f)

	m.Limits.MaxRuntime, _ = time.ParseDuration("10s")

	if err := m.Execute(); err != nil {
		panic(err)
	}

	runtime := time.Now().Sub(m.StartedAt)

	e := exporters.NewSuggestionExporter(
		m.Vars,
		exporters.NewTextExporter(),
	)

	e.Export(m.Field) // export an image.

	fmt.Printf("Seed: %d\n", m.Seed)
	fmt.Printf("Execution Time: %f seconds\n", runtime.Seconds())
}
