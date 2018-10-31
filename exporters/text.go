package exporters

import (
	"fmt"
	"github.com/moccalotto/nick/machine"
	"os"
	"strings"
)

type TextExporter struct {
	Machine  *machine.Machine
	LiveStr  string
	OffStr   string
	FileName string
}

func NewTextExporter(m *machine.Machine) *TextExporter {
	return &TextExporter{
		LiveStr:  "██",
		OffStr:   "  ",
		FileName: "",
	}
}

func (t *TextExporter) String() string {
	f := t.Machine.Field
	var buf strings.Builder
	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			if a, err := f.On(x, y); err != nil {
				panic(err)
			} else if a {
				buf.WriteString(t.LiveStr)
			} else {
				buf.WriteString(t.OffStr)
			}
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

func (t *TextExporter) Export() error {
	output := t.String()

	// Print to screen if filename is empty
	if t.FileName == "" || t.FileName == "-" {
		fmt.Println(output)
		return nil
	}

	file, err := os.Create(t.FileName)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	if _, err := file.WriteString(output); err != nil {
		return err
	}

	return nil
}
