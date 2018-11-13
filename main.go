package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/moccalotto/nick/exporters"
	"github.com/moccalotto/nick/machine"
	"github.com/moccalotto/nick/utils"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	m             *machine.Machine
	iterations    int  = 0
	width         int  = 800
	height        int  = 600
	done          bool = false
	scale         float64
	currentImage  *ebiten.Image
	imageExporter *exporters.ImageExporter
)

func init() {
	machine.InstructionHandlers["preview"] = PreviewInstr
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

func PreviewInstr(m *machine.Machine) {
	var err error

	imageExporter, err = newImageExporter(m, width, height)
	m.Assert(err == nil, "Preview failed: %s", err)

	img, err := imageExporter.GetImage()
	m.Assert(err == nil, "Preview failed: %s", err)

	currentImage, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	m.Assert(err == nil, "Preview failed: %s", err)

	if m.HasArg(0) {
		time.Sleep(time.Duration(m.ArgAsFloat(0) * float64(time.Second)))
	}
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

	iterations++

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		log.Println("Done")
		os.Exit(0)
		return nil
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if m == nil {
		return ebitenutil.DebugPrint(screen, "Waiting for machine to execute. "+strconv.Itoa(iterations)+" ...")
	}

	if currentImage == nil {
		return ebitenutil.DebugPrint(screen, "Waiting for preview to become ready. "+strconv.Itoa(m.State.PC)+" ...")
	}

	return screen.DrawImage(currentImage, nil)
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

func main() {
	f := flag.String("script", "examples/empty.cave", "Path to script to execute")
	flag.Parse()

	if *f == "" {
		flag.PrintDefaults()
		return
	}

	m = newMachineFromFile(*f)
	go func() {
		if err := m.Execute(); err != nil {
			panic(err)
		}
		PreviewInstr(m)

		done = true
	}()

	err := ebiten.Run(gameLoop, width, height, scale, "Nick")
	m.Assert(err == nil, "Error in render-loop")
}
