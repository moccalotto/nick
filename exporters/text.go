package exporters

import (
	"fmt"
	"github.com/moccalotto/nick/field"
	"log"
	"os"
	"strings"
)

type TextExporter struct {
	LiveStr  string
	DeadStr  string
	FileName string
}

func NewTextExporter() *TextExporter {
	return &TextExporter{
		LiveStr:  "██",
		DeadStr:  "  ",
		FileName: "",
	}
}

func (t *TextExporter) String(f *field.Field) string {
	var buf strings.Builder
	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			if f.Alive(x, y) {
				buf.WriteString(t.LiveStr)
			} else {
				buf.WriteString(t.DeadStr)
			}
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

func (t *TextExporter) Export(f *field.Field) {
	// Maybe just print to screen
	if t.FileName == "" {
		fmt.Println(t.String(f))
		return
	}

	file, err := os.Create(t.FileName)

	if err != nil {
		log.Fatalf("Could not open file '%s': %s", t.FileName, err)
	}

	defer file.Close()

	if _, err := file.WriteString(t.String(f)); err != nil {
		log.Fatalf("Could not write to file '%s': %s", t.FileName, err)
	}
}
