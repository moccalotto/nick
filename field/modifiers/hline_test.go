package modifiers

import (
	//	"github.com/moccaloto/nick/field"
	"testing"
)

func TestNewHLine(t *testing.T) {
	h := NewHLine(20, 20, 5)

	if h.StartX != 20 {
		t.Errorf("StartX != 20")
	}
	if h.StartY != 20 {
		t.Errorf("StartY != 20")
	}
	if h.Length != 5 {
		t.Errorf("Expected Length is 5, actual Length: %d", h.Length)
	}
	if h.Thickness != 1 {
		t.Errorf("Expected Thickness is 1, actual thickness: %d", h.Thickness)
	}
	if h.Probability != 1.0 {
		t.Errorf("Expected Probability is 1.0, actual Probability: %f", h.Probability)
	}
}

func TestWithThickness(t *testing.T) {
	h0 := NewHLine(20, 20, 5)
	h := h0.WithThickness(4)

	if h0.Thickness != 1 {
		t.Errorf("Mutation: %+v, %+v", h0, h)
	}

	if h.StartX != 20 {
		t.Errorf("StartX != 20")
	}
	if h.StartY != 20 {
		t.Errorf("StartY != 20")
	}
	if h.Length != 5 {
		t.Errorf("Expected Length is 5, actual Length: %d", h.Length)
	}
	if h.Thickness != 4 {
		t.Errorf("Expected Thickness is 4, actual thickness: %d", h.Thickness)
	}
	if h.Probability != 1.0 {
		t.Errorf("Expected Probability is 1.0, actual Probability: %f", h.Probability)
	}
}
