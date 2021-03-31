package utils

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"testing"
)

func TestNewColorStateGetLast(t *testing.T) {
	var cs = NewColorState(5, 5, 5)

	var c1 = cs.get(4, 4, 4)

	fmt.Printf("got %v\n", c1)
}

func TestNewColorStateGetIndex(t *testing.T) {
	var cs = NewColorState(5, 5, 5)

	var index = cs.xyzToIndex(0, 0, 0)
	if index != 0 {
		t.Fatalf("Index expected %d, got %d", 0, index)
	}

	index = cs.xyzToIndex(4, 0, 0)
	if index != 4 {
		t.Fatalf("Index expected %d, got %d", 4, index)
	}

	index = cs.xyzToIndex(0, 4, 0)
	if index != 20 {
		t.Fatalf("Index expected %d, got %d", 20, index)
	}

	index = cs.xyzToIndex(4, 4, 0)
	if index != 24 {
		t.Fatalf("Index expected %d, got %d", 24, index)
	}

	index = cs.xyzToIndex(0, 0, 4)
	if index != 100 {
		t.Fatalf("Index expected %d, got %d", 100, index)
	}

	index = cs.xyzToIndex(4, 4, 4)
	if index != 124 {
		t.Fatalf("Index expected %d, got %d", 124, index)
	}
}

func TestColorStateFillZ(t *testing.T) {
	var cs = NewColorState(5, 5, 5)

	for i := 0; i < 5; i++ {
		fillColor := colorful.Color{R: float64(i)}

		cs.fillZ(i, fillColor)
		c := cs.get(0, 0, i)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
		c = cs.get(2, 2, i)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
		c = cs.get(4, 4, i)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
	}
}

func TestColorStateFillY(t *testing.T) {
	var cs = NewColorState(5, 5, 5)

	for i := 0; i < 5; i++ {
		fillColor := colorful.Color{R: float64(i)}

		cs.fillY(i, fillColor)
		c := cs.get(0, i, 0)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
		c = cs.get(4, i, 2)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
		c = cs.get(4, i, 4)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
	}
}

func TestColorStateFillX(t *testing.T) {
	var cs = NewColorState(5, 5, 5)

	for i := 0; i < 5; i++ {
		fillColor := colorful.Color{R: float64(i)}

		cs.fillX(i, fillColor)
		c := cs.get(i, 0, 0)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
		c = cs.get(i, 2, 2)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
		c = cs.get(i, 4, 4)
		if c.R != float64(i) {
			t.Fatalf("Color not correct %d\n", i)
		}
	}
}
