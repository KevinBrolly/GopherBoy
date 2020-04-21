package ppu

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
)

type Fetcher struct {
	tileMapAddress  uint16
	tileDataAddress uint16
	SCY             byte
	LY              byte
	memory          mmu.Memory
}

func (f *Fetcher) NextTileLine() []*Dot {
	tileLine := f.fetchTileLine(f.getTileIdentifier())
	f.tileMapAddress++
	return tileLine
}

func (f *Fetcher) getTileIdentifier() byte {
	// Divide the Y position by 8 (for 8 pixels in tile)
	// and multiply by 32 (for number of tiles in the background map)
	// to get the row number for the tile in the background map
	tileRow := uint16((f.SCY+f.LY)/8) * 32
	tileAddress := f.tileMapAddress + tileRow
	return f.memory.ReadByte(tileAddress)
}

func (f *Fetcher) getTileDataAddress(tileIdentifier byte) uint16 {
	// When the tileDataAddress used is 0x8800 the tileIndentifier is
	// a signed byte -127 - 127, the offset corrects for this
	// when looking up the memory location
	offset := uint16(0)
	if f.tileDataAddress == 0x8800 {
		offset = 128
	}

	if tileIdentifier > 127 {
		return f.tileDataAddress + ((uint16(tileIdentifier) - offset) * 16) // 16 = tile size in bytes
	} else {
		return f.tileDataAddress + ((uint16(tileIdentifier) + offset) * 16) // 16 = tile size in bytes
	}
}

func (f *Fetcher) getTileData(dataAddress uint16, line byte) (data1, data2 byte) {
	data1 = f.memory.ReadByte(dataAddress + uint16(line))
	data2 = f.memory.ReadByte(dataAddress + uint16(line) + 1)
	return data1, data2
}

func (f *Fetcher) fetchTileLine(tileIdentifier byte) []*Dot {
	tileDataAddress := f.getTileDataAddress(tileIdentifier)
	// Find the correct vertical line we're on of the
	// tile to get the tile data from memory
	verticalLine := (f.LY + f.SCY) % 8
	verticalLine = verticalLine * 2 // each vertical line takes up two bytes of memory

	data1, data2 := f.getTileData(tileDataAddress, verticalLine)

	line := make([]*Dot, 8)
	dataBit := 7
	for i := byte(0); i <= 7; i++ {
		var colorIdentifier byte
		if utils.IsBitSet(data1, byte(dataBit)) {
			colorIdentifier = utils.SetBit(colorIdentifier, 0)
		}
		if utils.IsBitSet(data2, byte(dataBit)) {
			colorIdentifier = utils.SetBit(colorIdentifier, 1)
		}

		line[i] = &Dot{
			ColorIdentifier: colorIdentifier,
			Type:            BG,
		}
		dataBit--
	}

	return line
}

func (f *Fetcher) fetchSpriteLine(sprite *Sprite) []Dot {
	verticalLine := f.LY - (sprite.Y - 16)
	verticalLine = verticalLine * 2 // each vertical line takes up two bytes of memory
	spriteDataAddress := 0x8000 + uint16(sprite.TileNumber)*16
	data1, data2 := f.getTileData(spriteDataAddress, verticalLine)

	line := make([]Dot, 8)
	dataBit := 7
	for i := 0; i <= 7; i++ {
		var colorIdentifier byte
		if utils.IsBitSet(data1, byte(dataBit)) {
			colorIdentifier = utils.SetBit(colorIdentifier, 0)
		}
		if utils.IsBitSet(data2, byte(dataBit)) {
			colorIdentifier = utils.SetBit(colorIdentifier, 1)
		}

		line[i] = Dot{
			ColorIdentifier: colorIdentifier,
		}

		dataBit--
	}

	// reverse the sprite if XFlip is set
	if sprite.XFlip() {
		// https://github.com/golang/go/wiki/SliceTricks#reversing
		for i := len(line)/2-1; i >= 0; i-- {
			opp := len(line)-1-i
			line[i], line[opp] = line[opp], line[i]
		}
	}

	return line
}
