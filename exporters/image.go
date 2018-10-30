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
	"github.com/moccalotto/nick/field"
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

// ImageExporter exports images to files
type ImageExporter struct {
	FileName string
	Format   string

	// If both Width and Height are 0, the image is sized via Scale instead
	Width  int // Scale to new width. If 0, aspect ratio is preserved
	Height int // Scale to new height.If 0, aspect ratio is preserved.

	// Scale the image by a given factor instead of a given pixel count.
	// Useful for scaling x2 or x4 without loss of quality.
	Scale float64

	Algorithm string // algorithm for scaling

	OffColor color.Color
	OnColor  color.Color
}

// NewImageExporter creates a new ImageExporter
func NewImageExporter() *ImageExporter {
	return &ImageExporter{
		FileName:  "map.png",
		Format:    "",
		Width:     0,
		Height:    0,
		Algorithm: "Lanczos",
		OnColor:   color.Alpha{0x99},
		OffColor:  color.Alpha{0xff},
	}
}

// Calculate the output dimensions of the image
func (this *ImageExporter) targetDimensions(f *field.Field) image.Rectangle {
	if this.Width == 0 && this.Height == 0 {
		return image.Rect(0, 0, f.Width(), f.Height())
	}

	return image.Rect(0, 0, this.Width, this.Height)
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

// GetMask returns a raw NRGBA image (for use in other exporters, etc.)
func (this *ImageExporter) GetMask(f *field.Field) image.Image {
	fw := f.Width()
	fh := f.Height()

	// create an image the size of the field, it will be scaled later
	img := image.NewAlpha(image.Rect(0, 0, fw, fh))

	f.Walk(func(x, y int, c field.Cell) {
		if c.On() {
			img.Set(x, y, this.OnColor)
		} else {
			img.Set(x, y, this.OffColor)
		}
	})

	rect := this.targetDimensions(f)

	return imaging.Resize(img, rect.Max.X, rect.Max.Y, this.filter())
}

func (this *ImageExporter) LoadBackgroundImage(r image.Rectangle) image.Image {

	// THIS IS A TEMP HACK
	f, err := os.Open("/Users/krh/Desktop/Nick/_backgrounds/paper_by_darkwood67/brown_ice_by_darkwood67.jpg")
	if err != nil {
		log.Fatal(err)
	}

	if img, _, err := image.Decode(bufio.NewReader(f)); err == nil {
		return imaging.Resize(img, r.Max.X, r.Max.Y, this.filter())
	}

	panic(err)
}

func (this *ImageExporter) GetImage(f *field.Field) image.Image {
	mask := this.GetMask(f)
	rect := this.targetDimensions(f)
	bg := this.LoadBackgroundImage(rect)

	img := image.NewRGBA(rect)

	draw.DrawMask(
		img,
		rect,
		bg,
		rect.Min,
		mask,
		rect.Min,
		draw.Over,
	)

	return img
}

// Export the image to a file
func (this *ImageExporter) Export(f *field.Field) {
	file, err := os.Create(this.FileName)
	if err != nil {
		log.Fatalf("Could not open file '%s': %s", this.FileName, err)
	}

	defer func() {
		_ = file.Close()
	}()

	format, err := this.detectFormat()
	if err != nil {
		log.Fatal("Could not auto-detect image format")
	}

	out := this.GetImage(f)

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
		log.Fatalf("Could not save image: %s", err)
	}
}
