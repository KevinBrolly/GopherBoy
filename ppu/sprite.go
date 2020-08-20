package ppu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type Sprite struct {
	Y          byte
	X          byte
	TileNumber byte
	Attributes byte
}

func (s *Sprite) Priority() byte {
	if utils.IsBitSet(s.Attributes, 7) {
		return 1
	}
	return 0
}

func (s *Sprite) YFlip() bool {
	return utils.IsBitSet(s.Attributes, 6)
}

func (s *Sprite) XFlip() bool {
	return utils.IsBitSet(s.Attributes, 5)
}

func (s *Sprite) DMGPalette() byte {
	return s.Attributes & 0x04
}

func (s *Sprite) VRAMBank() bool {
	return utils.IsBitSet(s.Attributes, 3)
}

func (s *Sprite) GBCPalette() byte {
	return s.Attributes & 0x3
}

func (s *Sprite) Palette() byte {
	return s.DMGPalette()
}
