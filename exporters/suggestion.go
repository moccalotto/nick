package exporters

import (
	"fmt"
	"github.com/moccalotto/nick/machine"
	"strconv"
)

type SuggestionExporter struct {
	Machine  *machine.Machine
	Fallback Exporter
}

func NewSuggestionExporter(m *machine.Machine, fallback Exporter) *SuggestionExporter {
	return &SuggestionExporter{m, fallback}
}

func (e *SuggestionExporter) image() (*ImageExporter, error) {
	ie := NewImageExporter(e.Machine)

	var err error = nil

	if fn, ok := e.Machine.Vars["suggestion.export.file"]; ok {
		ie.FileName = fn
	}

	if fmt, ok := e.Machine.Vars["suggestion.export.format"]; ok {
		ie.Format = fmt
	}

	var width, height int

	if w, ok := e.Machine.Vars["suggestion.export.width"]; ok {
		if width, err = strconv.Atoi(w); err != nil {
			return nil, err
		}
	}

	if h, ok := e.Machine.Vars["suggestion.export.height"]; ok {
		if height, err = strconv.Atoi(h); err != nil {
			return nil, err
		}
	}

	ie.Rect = ie.makeRect(width, height)

	if str, ok := e.Machine.Vars["suggestion.export.algorithm"]; ok {
		if ie.Algorithm, err = ie.parseAlgorithmString(str); err != nil {
			return nil, err
		}
	}

	var tileWidth, tileHeight float64

	if str, ok := e.Machine.Vars["suggestion.grid.cols"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileWidth = float64(ie.Rect.Max.X) / num
		} else {
			return nil, err
		}
	}
	if str, ok := e.Machine.Vars["suggestion.grid.rows"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileHeight = float64(ie.Rect.Max.Y) / num
		} else {
			return nil, err
		}
	}
	if str, ok := e.Machine.Vars["suggestion.grid.width"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileWidth = num
		} else {
			return nil, err
		}
	}
	if str, ok := e.Machine.Vars["suggestion.grid.height"]; ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			tileHeight = num
		} else {
			return nil, err
		}
	}

	if tileWidth > 0 && tileHeight > 0 {
		ie.Grid = &GridSettings{tileWidth, tileHeight}
	}

	return ie, nil
}

func (e *SuggestionExporter) iterm() (*ItermExporter, error) {
	if i, e := e.image(); e == nil {
		return NewItermExporter(i), nil
	} else {
		return nil, e
	}
}

func (e *SuggestionExporter) text() (*TextExporter, error) {
	te := NewTextExporter(e.Machine)

	if fn, ok := e.Machine.Vars["suggestion.export.file"]; ok {
		te.FileName = fn
	}
	if l, ok := e.Machine.Vars["suggestion.export.on"]; ok {
		te.LiveStr = l
	}
	if d, ok := e.Machine.Vars["suggestion.export.off"]; ok {
		te.OffStr = d
	}

	return te, nil
}

func (e *SuggestionExporter) Export() error {
	ex, ok := e.Machine.Vars["suggestion.export.type"]

	if !ok {
		return e.Fallback.Export()
	}

	var exporter Exporter = nil
	var err error = nil
	switch ex {
	case "image":
		exporter, err = e.image()
	case "iterm":
		exporter, err = e.iterm()
	case "text":
		exporter, err = e.text()
	default:
		return fmt.Errorf("Unknown exporter: %s", ex)
	}

	if err != nil {
		return err
	}

	return exporter.Export()
}
