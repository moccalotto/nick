package exporters

import (
	"bytes"
	"encoding/base64"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/moccalotto/nick/field"
	"image/png"
)

// ItermExporter exports an image directly to an iTerm screen (not a file)
type ItermExporter struct {
	ClearScreen bool
	*ImageExporter
}

// NewItermExporter creates a new ItermExporter
func NewItermExporter(i *ImageExporter) *ItermExporter {
	return &ItermExporter{
		ClearScreen:   true,
		ImageExporter: i,
	}
}

// Export a field directly to iTerm screen
func (p *ItermExporter) Export(f *field.Field) {
	img := p.ImageExporter.GetImage(f)
	buf := new(bytes.Buffer)
	png.Encode(buf, img)

	if p.ClearScreen {
		tm.Clear()
		tm.MoveCursor(1, 1)
	}

	tm.Print(
		fmt.Sprintf(
			"\033]1337;File=inline=1:%s\a",
			base64.StdEncoding.EncodeToString(buf.Bytes()),
		),
	)

	tm.Flush()
}
