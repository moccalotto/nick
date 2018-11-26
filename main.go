package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/moccalotto/nick/exporters"
	"github.com/moccalotto/nick/machine"
	"github.com/moccalotto/nick/utils"
)

var (
	m             *machine.Machine
	caveFileName  string
	ticks         int = 0
	width         int = 800
	height        int = 600
	scale         float64
	currentImage  *ebiten.Image
	imageExporter *exporters.ImageExporter
)

func init() {
	machine.InstructionHandlers["window"] = SetWindowSizeInstr

	scale = 1.0 / ebiten.DeviceScaleFactor() // Disable HiDPI compensation
}

func SetWindowSizeInstr(m *machine.Machine) {
	m.Assert(m.HasArg(2), "Correct usage: 'preview-size WIDTH x HEIGHT'")
	m.Assert(m.ArgAsString(1) == "x", "Correct usage: 'preview-size WIDTH x HEIGHT'")

	width = m.ArgAsInt(0)
	height = m.ArgAsInt(2)

	ebiten.SetScreenSize(width, height)
}

func render(m *machine.Machine) {
	if m.Cave == nil {
		return
	}

	var err error

	imageExporter, err = newImageExporter(m, width, height)
	m.Assert(err == nil, "Preview failed: %s", err)

	img, err := imageExporter.GetImage()
	m.Assert(err == nil, "Preview failed: %s", err)

	currentImage, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	m.Assert(err == nil, "Preview failed: %s", err)

	return
}

func newMachineFromFile(filename string) *machine.Machine {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file: '%s'", filename)
	}

	script := string(b)
	m := machine.MachineFromScript(script)

	return m
}

func gameLoop(screen *ebiten.Image) error {
	ticks++

	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		log.Println("Done")
		os.Exit(0)
		return nil
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if currentImage == nil {
		return ebitenutil.DebugPrint(screen, "Creating Map")
	}

	err := screen.DrawImage(currentImage, nil)

	if ebiten.IsKeyPressed(ebiten.KeyN) {
		currentImage = nil
		drawAnImage()
	}

	return err
}

func newImageExporter(m *machine.Machine, width, height int) (*exporters.ImageExporter, error) {
	var err error
	ie := exporters.NewImageExporter(m)

	var tileWidth, tileHeight float64

	ie.Background = &exporters.BackgroundSettings{}
	ie.Grid = &exporters.GridSettings{}

	ie.Rect = ie.MakeRect(width, height)

	if str, ok := m.Vars[".export.algorithm"]; ok {
		if ie.Algorithm, err = ie.ParseAlgorithmString(str); err != nil {
			return nil, err
		}
	}

	if str, ok := m.Vars[".grid.cols"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileWidth = float64(ie.Rect.Max.X) / num
		} else {
			return nil, err
		}
	}
	if str, ok := m.Vars[".grid.rows"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileHeight = float64(ie.Rect.Max.Y) / num
		} else {
			return nil, err
		}
	}
	if str, ok := m.Vars[".grid.width"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileWidth = num
		} else {
			return nil, err
		}
	}
	if str, ok := m.Vars[".grid.height"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileHeight = num
		} else {
			return nil, err
		}
	}

	if str, ok := m.Vars[".grid.color"]; ok {
		if col, err := utils.ParseColorString(str); err != nil {
			return nil, err
		} else {
			ie.Grid.Color = col
		}
	}

	if tileWidth > 0 && tileHeight > 0 {
		ie.Grid.CellWidthPx = tileWidth
		ie.Grid.CellHeightPx = tileHeight
	}

	if str, ok := m.Vars[".background.file"]; ok {
		ie.Background.FileName = str
	}

	if str, ok := m.Vars[".background.color"]; ok {
		if col, err := utils.ParseColorString(str); err == nil {
			ie.Background.Color = col
		} else {
			return nil, err
		}
	}

	if str, ok := m.Vars[".wall.color"]; ok {
		if col, err := utils.ParseColorString(str); err != nil {
			return nil, err
		} else {
			ie.WallColor = col
		}
	}

	return ie, nil
}

func drawAnImage() {
	m = newMachineFromFile(caveFileName)
	tick := func(m *machine.Machine, i *machine.Instruction) error {
		if i == nil {
			render(m)
			return nil
		}

		return nil
	}
	go func() {
		if err := m.Execute(tick); err != nil {
			panic(err)
		}
	}()

}

func main() {
	f := flag.String("script", "examples/empty.cave", "Path to script to execute")
	flag.Parse()

	if *f == "" {
		flag.PrintDefaults()
		return
	}

	caveFileName = *f

	drawAnImage()
	err := ebiten.Run(gameLoop, width, height, scale, "Nick")
	m.Assert(err == nil, "Error in render-loop")
}
