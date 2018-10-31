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
		Machine:  m,
		LiveStr:  "██",
		OffStr:   "  ",
		FileName: "",
	}
}

func (this *TextExporter) String() string {
	f := this.Machine.Field
	var buf strings.Builder
	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			if a, err := f.On(x, y); err != nil {
				panic(err)
			} else if a {
				buf.WriteString(this.LiveStr)
			} else {
				buf.WriteString(this.OffStr)
			}
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

func (this *TextExporter) Export() error {
	output := this.String()

	// Print to screen if filename is empty
	if this.FileName == "" || this.FileName == "-" {
		fmt.Println(output)
		return nil
	}

	file, err := os.Create(this.FileName)
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
