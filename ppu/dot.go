package ppu

import "image/color"

type DotType int

const (
	BG     DotType = 0
	SPRITE DotType = 1
)

type Dot struct {
	ColorIdentifier byte
	Palette         byte
	Priority        byte
	Type            DotType
}

func (d *Dot) ToRGBA(palette byte) color.RGBA {
	var paletteNum byte
	var bitmask byte = 0x3
	switch d.ColorIdentifier {
	case 0:
		paletteNum = palette & bitmask
	case 1:
		paletteNum = (palette >> 2) & bitmask
	case 2:
		paletteNum = (palette >> 4) & bitmask
	case 3:
		paletteNum = (palette >> 6) & bitmask
	}

	switch paletteNum {
	case 0:
		return color.RGBA{224, 248, 208, 0}
	case 1:
		return color.RGBA{136, 192, 112, 0}
	case 2:
		return color.RGBA{52, 104, 86, 0}
	case 3:
		return color.RGBA{8, 24, 32, 0}
	}

	return color.RGBA{}
}
