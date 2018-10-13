package exporters

/*
Export for image.
Considerations:
	Dynamic map names:
		Patterns in filename (such as %rand %Y-%m-%d-%H:%i:s or similar)
		Sequenced filenames (map-%seq that auto-detects previous maps)
	Detect map format from extension
	Do we draw a grid?
		Do we scale the grid?
		Do we use instructions for that?
	Image behind dead pixels
		Do we use instructions for that?
	Image behind live pixels (possibly drawing a grid on top of that image too)
		Do we use instructions for that?
*/

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/moccalotto/nick/field"
	"image"
	"image/color"
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

	DeadColor color.NRGBA
	LiveColor color.NRGBA
}

// NewImageExporter creates a new ImageExporter
func NewImageExporter() *ImageExporter {
	return &ImageExporter{
		FileName:  "map.png",
		Format:    "",
		Width:     0,
		Height:    0,
		Scale:     1,
		Algorithm: "Lanczos",
		DeadColor: color.NRGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 0,
		}, // transparent
		LiveColor: color.NRGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		}, // white opaque
	}
}

// Calculate the output dimensions of the image
func (p *ImageExporter) dimensions(f *field.Field) (int, int) {
	if p.Width == 0 && p.Height == 0 {
		return int(float64(f.Width()) * p.Scale), int(float64(f.Height()) * p.Scale)
	}

	return p.Width, p.Height
}

func (p *ImageExporter) detectFormat() (string, error) {
	if p.Format != "" {
		return p.Format, nil
	}

	parts := strings.Split(p.FileName, ".")
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

func (p *ImageExporter) filter() imaging.ResampleFilter {
	switch p.Algorithm {
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
		log.Fatalf("Unknown image scaling algorithm: %s", p.Algorithm)
	}

	panic("Should never be reached!")
}

// GetImage returns a raw NRGBA image (for use in other exporters, etc.)
func (p *ImageExporter) GetImage(f *field.Field) *image.NRGBA {
	fw := f.Width()
	fh := f.Height()

	// create an image the size of the field, it will be scaled later
	img := image.NewNRGBA(image.Rect(0, 0, fw, fh))

	for y := 0; y < fh; y++ {
		for x := 0; x < fw; x++ {
			if f.Alive(x, y) {
				img.Set(x, y, p.LiveColor)
			} else {
				img.Set(x, y, p.DeadColor)
			}
		}
	}

	imgW, imgH := p.dimensions(f)

	return imaging.Resize(img, imgW, imgH, p.filter())
}

// Export the image to a file
func (p *ImageExporter) Export(f *field.Field) {
	img := p.GetImage(f)

	file, err := os.Create(p.FileName)
	if err != nil {
		log.Fatalf("Could not open file '%s': %s", p.FileName, err)
	}

	defer file.Close()

	format, err := p.detectFormat()
	if err != nil {
		log.Fatal("Could not auto-detect image format")
	}

	switch format {
	case "png":
		err = png.Encode(file, img)
	case "gif":
		err = gif.Encode(file, img, &gif.Options{NumColors: 2})
	case "jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	default:
		log.Fatalf("Unknown file format: %s", p.Format)
	}

	if err != nil {
		log.Fatalf("Could not save image: %s", err)
	}
}
