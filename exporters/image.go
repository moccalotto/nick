package exporters

/*
Export for image.
Considerations:
	Dynamic map names:
		Patterns in filename (such as %rand %Y-%m-%d-%H:%i:s or similar)
		Sequenced filenames (map-%seq that auto-detects previous maps)

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
	"github.com/moccalotto/nick/effects"
	"github.com/moccalotto/nick/field"
	"github.com/moccalotto/nick/machine"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type GridSettings struct {
	CellWidthPx  float64
	CellHeightPx float64
}

type BackgroundSettings struct {
	FileName string
}

// ImageExporter exports images to files
type ImageExporter struct {
	Machine    *machine.Machine       // The machine to export
	FileName   string                 // Filename of the generated image
	Background *BackgroundSettings    // Background image settings
	Grid       *GridSettings          // Grid settings
	Rect       image.Rectangle        // Dimensions of the output image
	Algorithm  imaging.ResampleFilter // Algorithm used to scale the image
	WallColor  color.Color

	maskBWCache image.Image
	maskCache   image.Image
}

// NewImageExporter creates a new ImageExporter
func NewImageExporter(m *machine.Machine) *ImageExporter {
	exporter := ImageExporter{
		Machine:   m,
		FileName:  "map.png",
		Algorithm: imaging.Lanczos,
		Rect:      m.Cave.Bounds(),
		WallColor: color.NRGBA{255, 0, 0, 144},
	}

	return &exporter
}

// Calculate the output dimensions of the image
func (this *ImageExporter) makeRect(w, h int) image.Rectangle {
	if h == 0 && w == 0 {
		h = this.Machine.Cave.Height()
		w = this.Machine.Cave.Width()
	} else if w == 0 {
		ratio := this.Machine.Cave.AspectRatio()
		w = int(float64(h) * ratio)
	} else if h == 0 {
		ratio := this.Machine.Cave.AspectRatio()
		h = int(float64(w) / ratio)
	}

	return image.Rect(0, 0, w, h)
}

func (this *ImageExporter) detectFormat() (string, error) {
	parts := strings.Split(this.FileName, ".")
	suffix := parts[len(parts)-1]

	switch suffix {
	case "png":
		return "png", nil
	case "jpg":
		return "jpeg", nil
	case "jpeg":
		return "jpeg", nil
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

	if this.maskCache == nil {
		this.maskCache = imaging.Resize(this.Machine.Cave, this.Rect.Max.X, this.Rect.Max.Y, this.Algorithm)
	}

	return this.maskCache
}

func (this *ImageExporter) maskBW() image.Image {

	if this.maskBWCache == nil {

		// backup the palette
		tmp := this.Machine.Cave.Palette

		// use a different palette
		this.Machine.Cave.Palette = field.BinaryPalette()

		// generate the image
		this.maskBWCache = imaging.Resize(this.Machine.Cave, this.Rect.Max.X, this.Rect.Max.Y, this.Algorithm)

		// restore the palette
		this.Machine.Cave.Palette = tmp

	}

	return this.maskBWCache
}

func (this *ImageExporter) backgroundImage() (draw.Image, error) {

	if this.Background == nil {
		return nil, nil
	}

	if this.Background.FileName == "" {
		return nil, nil
	}

	file, err := os.Open(this.Background.FileName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bufio.NewReader(file))

	if err != nil {
		return nil, err
	}

	return imaging.Resize(img, this.Rect.Max.X, this.Rect.Max.Y, this.Algorithm), nil
}

func (this *ImageExporter) outline() (image.Image, error) {

	if this.WallColor == nil {
		return nil, nil
	}

	outliner := effects.NewAutomaton("B678/S1234567")

	f := this.Machine.Cave.MapAsyncToNewField(outliner.NextCellState)
	f.Palette = field.InvertedBinaryPalette()

	return imaging.Resize(f, this.Rect.Max.X, this.Rect.Max.Y, this.Algorithm), nil
}

func (this *ImageExporter) grid() (image.Image, error) {

	if this.Grid == nil {
		return nil, nil
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

func (this *ImageExporter) applyBackground(img draw.Image) error {
	mask := this.mask()

	// if a background image was specified, apply it to the image.
	if bg, err := this.backgroundImage(); err != nil {
		return err
	} else if bg != nil {
		// draw the background on top of the (black) image through the mask.
		draw.DrawMask(img, this.Rect, bg, this.Rect.Min, mask, this.Rect.Min, draw.Over)
	} else {
		draw.DrawMask(img, this.Rect, image.Opaque, this.Rect.Min, mask, this.Rect.Min, draw.Over)
	}

	return nil
}

func (this *ImageExporter) applyGrid(img draw.Image) error {
	if grid, err := this.grid(); err != nil {
		return err
	} else if grid != nil {
		// draw tiles on the image through a mask that completely blocks drawing on the occupied areas.
		draw.DrawMask(img, this.Rect, grid, this.Rect.Min, this.maskBW(), this.Rect.Min, draw.Over)
	}

	return nil
}

func (this *ImageExporter) applyOutline(img draw.Image) error {
	if outline, err := this.outline(); err != nil {
		return err
	} else if outline != nil {
		// draw tiles on the image through a mask that completely blocks drawing on the occupied areas.
		draw.DrawMask(
			img,                              // draw onto img
			this.Rect,                        // on the entire image
			image.NewUniform(this.WallColor), // draw this color
			this.Rect.Min,                    // start drawing at the top left corner
			outline,                          // draw through the mask defined in the outline
			this.Rect.Min,                    // the mask also starts at the top left corner
			draw.Over,                        // draw the outline over img
		)
	}

	return nil
}

func (this *ImageExporter) GetImage() (image.Image, error) {
	// new black image of the given dimensions
	img := image.NewRGBA(this.Rect)

	this.applyBackground(img)
	this.applyGrid(img)
	this.applyOutline(img)

	return img, nil
}

// Export the image to a file
func (this *ImageExporter) Export() error {
	var file *os.File
	var format string
	var err error
	var img image.Image

	// Infer the file format from the given file name
	if format, err = this.detectFormat(); err != nil {
		return err
	}

	// Open the file for writing
	if file, err = os.Create(this.FileName); err != nil {
		return err
	} else {
		defer func() { _ = file.Close() }()
	}

	// Render the image
	if img, err = this.GetImage(); err != nil {
		return err
	}

	// Encode the image
	return encodeImage(format, file, img)
}

// Write binary image data into a file
func encodeImage(format string, file *os.File, img image.Image) error {
	if format == "png" {
		if err := png.Encode(file, img); err != nil {
			return err
		}

		return nil
	}

	if format == "jpeg" {
		if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 80}); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Could not infer format from filename.")
}
