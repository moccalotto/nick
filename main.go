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

	m.MaxRuntime, _ = time.ParseDuration("10s")

	if err := m.Execute(); err != nil {
		panic(err)
	}

	runtime := time.Now().Sub(m.StartedAt)

	exportStart := time.Now()
	err := exporters.NewSuggestionExporter(m, exporters.NewTextExporter(m)).Export()
	exportTime := time.Now().Sub(exportStart)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Printf("% -30s = %d\n", "Seed", m.Seed)
	fmt.Printf("% -30s = %f seconds\n", "Execution Time", runtime.Seconds())
	fmt.Printf("% -30s = %f seconds\n", "Export Time", exportTime.Seconds())

	for n, v := range m.Vars {
		n = "$" + n
		fmt.Printf("% -30s = %s\n", n, v)
	}
}
