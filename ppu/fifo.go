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

func (fifo *Fifo) OverlaySprite(dots []Dot, palette, bgPriority, spritePriority byte) {
	for i, targetDot := range fifo.dots[0:8] {
		if targetDot.Type == SPRITE {
			if spritePriority < targetDot.Priority && dots[i].ColorIdentifier != 0 {
				targetDot.ColorIdentifier = dots[i].ColorIdentifier
				targetDot.Palette = palette
				targetDot.Priority = spritePriority
				targetDot.Type = SPRITE
			}
		} else {
			// If there is a sprite at this position in the scanline
			// and the sprite priority is 0 or the background dots
			// colorIdentifier is 0, then the sprite is rendered on top
			// of the background, otherwise the background is rendered.
			if dots[i].ColorIdentifier != 0 && (bgPriority == 0 || targetDot.ColorIdentifier == 0) {
				targetDot.ColorIdentifier = dots[i].ColorIdentifier
				targetDot.Palette = palette
				targetDot.Priority = spritePriority
				targetDot.Type = SPRITE
			}
		}
	}
}
