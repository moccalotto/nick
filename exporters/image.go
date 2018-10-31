package exporters

/*
Export for image.
Considerations:
	Dynamic map names:
		Patterns in filename (such as %rand %Y-%m-%d-%H:%i:s or similar)
		Sequenced filenames (map-%seq that auto-detects previous maps)
	Detect map format from extension
	Define colors via strings á la https://github.com/go-playground/colors
		Steal code and rewrite to fit the real go colors.
	Do we draw a grid?
		Do we scale the grid?
		Do we use instructions for that?
	Image behind off-cells?
		Do we use instructions for that?
	Image behind live pixels (possibly drawing a grid on top of that image too)
		Do we use instructions for that?
*/

import (
	"bufio"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/moccalotto/nick/machine"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

// ImageExporter exports images to files
type ImageExporter struct {
	Machine  *machine.Machine
	FileName string
	Format   string

	// If both Width and Height are 0, the image is sized via Scale instead
	Width  int // Scale to new width. If 0, aspect ratio is preserved
	Height int // Scale to new height.If 0, aspect ratio is preserved.

	// Scale the image by a given factor instead of a given pixel count.
	// Useful for scaling x2 or x4 without loss of quality.
	Scale float64

	Algorithm string // algorithm for scaling
}

// NewImageExporter creates a new ImageExporter
func NewImageExporter(m *machine.Machine) *ImageExporter {
	return &ImageExporter{
		Machine:   m,
		FileName:  "map.png",
		Format:    "",
		Width:     0,
		Height:    0,
		Algorithm: "Lanczos",
	}
}

// Calculate the output dimensions of the image
func (this *ImageExporter) targetDimensions() image.Rectangle {
	w := this.Width
	h := this.Height

	if h == 0 && w == 0 {
		h = this.Machine.Field.Height()
		w = this.Machine.Field.Width()
	} else if w == 0 {
		ratio := this.Machine.Field.AspectRatio()
		w = int(float64(this.Height) * ratio)
	} else if h == 0 {
		ratio := this.Machine.Field.AspectRatio()
		h = int(float64(this.Width) / ratio)
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

func (this *ImageExporter) filter() imaging.ResampleFilter {
	switch this.Algorithm {
	case "NearestNeighbor":
		// NearestNeighbor is a nearest-neighbor filter (no anti-aliasing).
		return imaging.NearestNeighbor
	case "Linear":
		// Bilinear interpolation filter, produces reasonably good, smooth output.
		return imaging.Linear
	case "Lanczos":
		// High-quality resampling filter for photographic images yielding sharp results (slow).
		return imaging.Lanczos
	case "CatmullRom":
		// A sharp cubic filter. It's a good filter for both upscaling and downscaling if sharp results are needed.
		return imaging.CatmullRom
	case "MitchellNetravali":
		// A high quality cubic filter that produces smoother results with less ringing artifacts than CatmullRom.
		return imaging.MitchellNetravali
	case "Box":
		// Simple and fast averaging filter appropriate for downscaling.
		// When upscaling it's similar to NearestNeighbor.
		return imaging.Box
	default:
		log.Fatalf("Unknown image scaling algorithm: %s", this.Algorithm)
	}

	panic("Should never be reached!")
}

func (this *ImageExporter) mask() image.Image {
	rect := this.targetDimensions()
	return imaging.Resize(this.Machine.Field, rect.Max.X, rect.Max.Y, this.filter())
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

	rect := this.targetDimensions()

	return imaging.Resize(img, rect.Max.X, rect.Max.Y, this.filter()), nil
}

func (this *ImageExporter) GetImage() (image.Image, error) {
	mask := this.mask()
	bg, err := this.backgroundImage()
	if err != nil {
		return nil, err
	}

	// the dimensions of the output image.
	rect := this.targetDimensions()

	// new black image of the given dimensions
	img := image.NewRGBA(rect)

	// draw the background on top of the (black) image
	// through the mask
	draw.DrawMask(
		img,
		rect,
		bg,
		rect.Min,
		mask,
		rect.Min,
		draw.Over,
	)

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
