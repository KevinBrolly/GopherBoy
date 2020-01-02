package ppu

import (
	"testing"
)

func TestPushDots(t *testing.T) {
	dots := make([]*Dot, 8)

	for i := 0; i <= 7; i++ {
		dots[i] = &Dot{
			ColorIdentifier: 0x1,
		}
	}
	fifo := &Fifo{}

	fifo.PushDots(dots)

	for i, dot := range fifo.dots {
		if dot != dots[i] {
			t.Errorf("PushDots not pushing dots")
		}
	}
}

func TestLength(t *testing.T) {
	dots := make([]*Dot, 8)

	fifo := &Fifo{}

	length := fifo.Length()

	if length != 0 {
		t.Errorf("TestLength failed, expected 0 got %v", length)
	}

	fifo.PushDots(dots)

	length = fifo.Length()

	if length != len(dots) {
		t.Errorf("TestLength failed, expected %v got %v", len(dots), length)
	}
}

func TestClear(t *testing.T) {
	dots := make([]*Dot, 8)

	fifo := &Fifo{}
	fifo.PushDots(dots)

	fifo.Clear()

	if fifo.dots != nil {
		t.Errorf("TestClear failed, dots not nil")
	}
}

func TestPopDot(t *testing.T) {
	dots := make([]*Dot, 8)

	dots[0] = &Dot{
		ColorIdentifier: 0x1,
	}

	fifo := &Fifo{}
	fifo.PushDots(dots)

	dot := fifo.PopDot()

	if fifo.Length() != 7 {
		t.Errorf("TestPopDot failed, expected length %v got %v", 7, fifo.Length())
	}

	if dot.ColorIdentifier != dots[0].ColorIdentifier {
		t.Errorf("TestPopDot failed, expected ColorIdentifier %v got %v", dots[0].ColorIdentifier, dot.ColorIdentifier)
	}

	fifo.PopDot()

	if fifo.Length() != 6 {
		t.Errorf("TestPopDot failed, expected length %v got %v", 6, fifo.Length())
	}
}
