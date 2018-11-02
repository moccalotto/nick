package exporters

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
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
func (p *ItermExporter) Export() error {
	img, err := p.ImageExporter.GetImage()

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 60}); err != nil {
		return err
	}

	if p.ClearScreen {
		fmt.Print("\033[2J")
	}

	fmt.Printf(
		"\033]1337;File=inline=1:%s\a\n",
		base64.StdEncoding.EncodeToString(buf.Bytes()),
	)

	return nil
}
