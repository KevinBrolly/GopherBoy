package ppu

type Fifo struct {
	dots []*Dot
}

func (fifo *Fifo) PushDots(dots []*Dot) {
	for i := 0; i <= 7; i++ {
		fifo.dots = append(fifo.dots, dots[i])
	}
}

func (fifo *Fifo) PopDot() (dot *Dot) {
	dot, fifo.dots = fifo.dots[0], fifo.dots[1:]
	return
}

func (fifo *Fifo) Clear() {
	fifo.dots = nil
}

func (fifo *Fifo) Length() int {
	return len(fifo.dots)
}
