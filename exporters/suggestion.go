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

	if w, ok := e.Machine.Vars["suggestion.export.width"]; ok {
		if ie.Width, err = strconv.Atoi(w); err != nil {
			return nil, err
		}
	}

	if h, ok := e.Machine.Vars["suggestion.export.height"]; ok {
		if ie.Height, err = strconv.Atoi(h); err != nil {
			return nil, err
		}
	}

	if scale, ok := e.Machine.Vars["suggestion.export.scale"]; ok {
		if ie.Scale, err = strconv.ParseFloat(scale, 64); err != nil {
			return nil, err
		}
	}

	if al, ok := e.Machine.Vars["suggestion.export.algorithm"]; ok {
		ie.Algorithm = al
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
