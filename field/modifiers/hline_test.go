package modifiers

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewHLine(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)

	h := NewHLine(20, 20, 5, rng)

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
	if h.Coverage != 1.0 {
		t.Errorf("Expected Coverage is 1.0, actual Coverage: %f", h.Coverage)
	}
}
