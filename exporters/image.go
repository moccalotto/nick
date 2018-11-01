package exporters

/*
Export for image.
Considerations:
	Dynamic map names:
		Patterns in filename (such as %rand %Y-%m-%d-%H:%i:s or similar)
		Sequenced filenames (map-%seq that auto-detects previous maps)

	Auto-detect map format from extension

	Define colors via strings á la https://github.com/go-playground/colors
		must support image/color package

	Background images:
		Tiled images
		Offsets for tiled images
		Separate images for areas that are on or off
		Cropping background images to use only a portion of it
		(disintegration/imaging can do cropping)

	Customized grid color and offset
*/

import (
	"bufio"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/moccalotto/nick/machine"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

type GridSettings struct {
	CellWidthPx  float64
	CellHeightPx float64
}

// ImageExporter exports images to files
type ImageExporter struct {
	Machine   *machine.Machine
	FileName  string
	Format    string
	Grid      *GridSettings
	Rect      image.Rectangle
	Algorithm imaging.ResampleFilter
}

// NewImageExporter creates a new ImageExporter
func NewImageExporter(m *machine.Machine) *ImageExporter {
	exporter := ImageExporter{
		Machine:   m,
		FileName:  "map.png",
		Format:    "",
		Algorithm: imaging.Lanczos,
		Rect:      m.Field.Bounds(),
	}

	return &exporter
}

// Calculate the output dimensions of the image
func (this *ImageExporter) makeRect(w, h int) image.Rectangle {
	if h == 0 && w == 0 {
		h = this.Machine.Field.Height()
		w = this.Machine.Field.Width()
	} else if w == 0 {
		ratio := this.Machine.Field.AspectRatio()
		w = int(float64(h) * ratio)
	} else if h == 0 {
		ratio := this.Machine.Field.AspectRatio()
		h = int(float64(w) / ratio)
	}

	return image.Rect(0, 0, w, h)
}

func (this *ImageExporter) detectFormat() (string, error) {
	if this.Format != "" {
		return this.Format, nil
	}

	parts := strings.Split(this.FileName, ".")
	suffix := parts[len(parts)-1]

	switch suffix {
	case "png":
		return "png", nil
	case "jpg":
		return "jpeg", nil
	case "jpeg":
		return "jpeg", nil
	case "gif":
		return "gif", nil
	}

	return "", fmt.Errorf("Could not determine file type from suffix: %s", suffix)
}

func (this *ImageExporter) parseAlgorithmString(algorithm string) (imaging.ResampleFilter, error) {
	switch algorithm {
	case "NearestNeighbor":
		// NearestNeighbor is a nearest-neighbor filter (no anti-aliasing).
		return imaging.NearestNeighbor, nil
	case "Linear":
		// Bilinear interpolation filter, produces reasonably good, smooth output.
		return imaging.Linear, nil
	case "Lanczos":
		// High-quality resampling filter for photographic images yielding sharp results (slow).
		return imaging.Lanczos, nil
	case "CatmullRom":
		// A sharp cubic filter. It's a good filter for both upscaling and downscaling if sharp results are needed.
		return imaging.CatmullRom, nil
	case "MitchellNetravali":
		// A high quality cubic filter that produces smoother results with less ringing artifacts than CatmullRom.
		return imaging.MitchellNetravali, nil
	case "Box":
		// Simple and fast averaging filter appropriate for downscaling.
		// When upscaling it's similar to NearestNeighbor.
		return imaging.Box, nil
	default:
		return imaging.NearestNeighbor, fmt.Errorf("Unknown image scaling algorithm: %s", algorithm)
	}

	panic("Should never be reached!")
}

func (this *ImageExporter) mask() image.Image {
	return imaging.Resize(this.Machine.Field, this.Rect.Max.X, this.Rect.Max.Y, this.Algorithm)
}
func (this *ImageExporter) maskBW() image.Image {
	// create a backup of the existing palette
	orig := make(color.Palette, 255)
	copy(orig, this.Machine.Field.Palette)

	// modify the palette so that free space becomes transparent
	// and the other areas opaque
	this.Machine.Field.Palette[0] = color.Alpha{255}
	for i := 1; i < len(this.Machine.Field.Palette); i++ {
		this.Machine.Field.Palette[i] = color.Alpha{0}
	}

	// generate the image
	img := imaging.Resize(this.Machine.Field, this.Rect.Max.X, this.Rect.Max.Y, this.Algorithm)

	// restore colors
	this.Machine.Field.Palette = orig

	return img
}

func (this *ImageExporter) backgroundImage() (draw.Image, error) {

	// THIS IS A TEMP HACK
	file, err := os.Open("/Users/krh/Desktop/Nick/_backgrounds/paper_by_darkwood67/brown_ice_by_darkwood67.jpg")
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bufio.NewReader(file))

	if err != nil {
		return nil, err
	}

	return imaging.Resize(img, this.Rect.Max.X, this.Rect.Max.Y, this.Algorithm), nil
}

func (this *ImageExporter) grid() (image.Image, error) {

	if this.Grid == nil {
		return nil, fmt.Errorf("No grid settings")
	}

	// the dimensions of the output image.
	img := image.NewRGBA(this.Rect)

	nextX := this.Grid.CellWidthPx
	nextY := this.Grid.CellHeightPx

	for curY := 0; curY < this.Rect.Max.Y; curY++ {
		for curX := 0; curX < this.Rect.Max.X; curX++ {
			if curX == int(nextX) {
				img.Set(curX, curY, color.NRGBA{0x44, 0x44, 0x44, 0x55})
				nextX += this.Grid.CellWidthPx
			}
			if curY == int(nextY) {
				img.Set(curX, curY, color.NRGBA{0x44, 0x44, 0x44, 0x55})
			}
		}
		if curY == int(nextY) {
			nextY += this.Grid.CellHeightPx
		}
		nextX = this.Grid.CellWidthPx
	}

	return img, nil
}

func (this *ImageExporter) GetImage() (image.Image, error) {
	mask := this.mask()
	bg, err := this.backgroundImage()
	if err != nil {
		return nil, err
	}

	// new black image of the given dimensions
	img := image.NewRGBA(this.Rect)

	// draw the background on top of the (black) image through the mask.
	draw.DrawMask(img, this.Rect, bg, this.Rect.Min, mask, this.Rect.Min, draw.Over)

	// draw tiles on the image through a mask that completely blocks drawing on the occupied areas.
	if this.Grid != nil {
		grid, err := this.grid()
		if err != nil {
			return nil, err
		}

		draw.DrawMask(img, this.Rect, grid, this.Rect.Min, this.maskBW(), this.Rect.Min, draw.Over)
	}

	return img, nil
}

// Export the image to a file
func (this *ImageExporter) Export() error {
	file, err := os.Create(this.FileName)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	format, err := this.detectFormat()
	if err != nil {
		return err
	}

	out, err := this.GetImage()

	if err != nil {
		return err
	}

	switch format {
	case "png":
		err = png.Encode(file, out)
	case "gif":
		err = gif.Encode(file, out, &gif.Options{NumColors: 2})
	case "jpeg":
		err = jpeg.Encode(file, out, &jpeg.Options{Quality: 90})
	default:
		log.Fatalf("Unknown file format: %s", this.Format)
	}

	if err != nil {
		return err
	}

	return nil
}
