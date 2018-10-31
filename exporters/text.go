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
	w := f.Width()

	for i, c := range f.Cells() {
		if i > 0 && i%w == 0 {
			buf.WriteRune('\n')
		}
		if c.On() {
			buf.WriteString(this.LiveStr)
		} else {
			buf.WriteString(this.OffStr)
		}
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
