package exporters

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/buger/goterm"
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
	_ = png.Encode(buf, img)

	if p.ClearScreen {
		goterm.Clear()
		goterm.MoveCursor(1, 1)
	}

	_, _ = goterm.Print(
		fmt.Sprintf(
			"\033]1337;File=inline=1:%s\a",
			base64.StdEncoding.EncodeToString(buf.Bytes()),
		),
	)

	goterm.Flush()
}
