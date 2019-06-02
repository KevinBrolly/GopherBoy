package ppu

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WIDTH  = 160
	HEIGHT = 144
)

const (
	VBLANK_INTERRUPT = 0
	LCDC_INTERRUPT   = 1

	LCDC = 0xFF40 // LCD Control
	STAT = 0xFF41 // LCD Status/Mode
	SCY  = 0xFF42 // Scroll Y
	SCX  = 0xFF43 // Scroll X
	LY   = 0xFF44 // Scanline
	LYC  = 0xFF45 // Scanline Comparison
	DMA  = 0xFF46
	BGP  = 0xFF47 // Background Palette
	OBP0 = 0xFF48 // Object/Sprite Background Palette 0
	OBP1 = 0xFF49 // Object/Sprite Background Palette 1
	WY   = 0xFF4A // Window Y
	WX   = 0xFF4B // Window X
)

const (
	MODE0 = 0x00 // HBlank
	MODE1 = 0x01 // VBlank
	MODE2 = 0x02 // OAM Access
	MODE3 = 0x03 // VRAM Acess
)

type stat struct {
	coincidenceInterruptEnabled bool
	oamInterruptEnabled         bool
	vblankInterruptEnabled      bool
	hblankInterruptEnabled      bool
	coincidenceFlag             bool
	mode                        byte
}

func (s *stat) setStat(value byte) {
	s.coincidenceInterruptEnabled = utils.IsBitSet(value, 6)
	s.oamInterruptEnabled = utils.IsBitSet(value, 5)
	s.vblankInterruptEnabled = utils.IsBitSet(value, 4)
	s.hblankInterruptEnabled = utils.IsBitSet(value, 3)
}

func (s *stat) getStat() byte {
	// Start with 0x80 as bit 7 is unused/always returns 1
	var stat byte = 0x80

	// Bit 6 LYC=LY Coincidence Interrupt
	if s.coincidenceInterruptEnabled {
		stat = utils.SetBit(stat, 6)
	}

	// Bit 5 Mode 2 OAM Interrupt
	if s.oamInterruptEnabled {
		stat = utils.SetBit(stat, 5)
	}

	// Bit 4 Mode 1 V-Blank Interrupt
	if s.vblankInterruptEnabled {
		stat = utils.SetBit(stat, 4)
	}

	// Bit 3 Mode 0 H-Blank Interrupt
	if s.hblankInterruptEnabled {
		stat = utils.SetBit(stat, 3)
	}

	// Bit 2 Coincidence Flag (0:LYC<>LY, 1:LYC=LY)
	if s.coincidenceFlag {
		stat = utils.SetBit(stat, 2)
	}

	// or with the mode to set this on stat
	stat |= s.mode
	return stat
}

type PPU struct {
	mmu *mmu.MMU

	Window *Window

	VRAM [16384]byte
	OAM  [160]byte

	STAT *stat // LCD Status/Mode
	SCY  byte  // Scroll Y
	SCX  byte  // Scroll X
	LY   byte  // Scanline
	LYC  byte
	DMA  byte
	BGP  byte
	OBP0 byte
	OBP1 byte
	WY   byte // Window Y
	WX   byte // Window X

	// LCD Control byte
	LCDC byte
	// LCDC Bit 7 - LCD Display Enable             (0=Off, 1=On)
	lcdEnabled bool
	// LCDC Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	windowMapLocation uint16
	// LCDC Bit 5 - Window Display Enable          (0=Off, 1=On)
	windowEnabled bool
	// LCDC Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
	tileDataLocation uint16
	// LCDC Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
	backgroundMapLocation uint16
	// LCDC Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
	spriteSize byte
	// LCDC Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
	spriteEnabled bool
	// LCDC Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
	backgroundEnabled bool

	Cycles int // Number of cycles since the last LCD Status Mode Change
}

type pixelAttributes struct {
	colorIdentifier byte
	palette         byte
	priority        byte
}

func NewPPU(mmu *mmu.MMU) *PPU {
	window := NewWindow("Gameboy", WIDTH, HEIGHT)

	ppu := &PPU{
		mmu:    mmu,
		Window: window,
		STAT: &stat{
			coincidenceInterruptEnabled: false,
			oamInterruptEnabled:         false,
			vblankInterruptEnabled:      false,
			hblankInterruptEnabled:      false,
			coincidenceFlag:             true,
			mode:                        MODE1,
		},
		LCDC: 0x91,
		SCY:  0x00,
		SCX:  0x00,
		LYC:  0x00,
		BGP:  0xFC,
		OBP0: 0xFF,
		OBP1: 0xFF,
		WY:   0x00,
		WX:   0x00,
	}

	ppu.setLCDCFields(0x91)

	mmu.MapMemory(ppu, LCDC)
	mmu.MapMemory(ppu, STAT)
	mmu.MapMemory(ppu, SCY)
	mmu.MapMemory(ppu, SCX)
	mmu.MapMemory(ppu, LY)
	mmu.MapMemory(ppu, LYC)
	mmu.MapMemory(ppu, DMA)
	mmu.MapMemory(ppu, BGP)
	mmu.MapMemory(ppu, OBP0)
	mmu.MapMemory(ppu, OBP1)
	mmu.MapMemory(ppu, WY)
	mmu.MapMemory(ppu, WX)

	// VRAM Range
	mmu.MapMemoryRange(ppu, 0x8000, 0x9FFF)

	// OAM RAM
	mmu.MapMemoryRange(ppu, 0xFE00, 0xFE9F)

	return ppu
}

func (ppu *PPU) ReadByte(addr uint16) byte {
	switch {
	case addr == LCDC:
		return ppu.LCDC
	case addr == STAT:
		return ppu.STAT.getStat()
	case addr == SCY:
		return ppu.SCY
	case addr == SCX:
		return ppu.SCX
	case addr == LY:
		return ppu.LY
	case addr == BGP:
		return ppu.BGP
	case addr == OBP0:
		return ppu.OBP0
	case addr == OBP1:
		return ppu.OBP1
	case addr == WY:
		return ppu.WY
	case addr == WX:
		return ppu.WX
	case addr >= 0x8000 && addr <= 0x9FFF:
		return ppu.VRAM[addr&0x1FFF]
	case addr >= 0xFE00 && addr <= 0xFE9F:
		return ppu.OAM[addr&0x9F]
	}
	return 0
}

func (ppu *PPU) WriteByte(addr uint16, value byte) {
	switch {
	case addr == LCDC:
		ppu.LCDC = value
		ppu.setLCDCFields(value)
	case addr == STAT:
		ppu.STAT.setStat(value)
	case addr == SCY:
		ppu.SCY = value
	case addr == SCX:
		ppu.SCX = value
	case addr == LY:
		// If the game writes to scanline it should be unset
		ppu.LY = 0
	case addr == LYC:
		ppu.LYC = value
	case addr == DMA:
		// The value holds the source address of the OAM data divided by 100
		// so we have to multiply it first
		var sourceAddr uint16 = uint16(value) << 8

		for i := 0; i < 160; i++ {
			ppu.OAM[i] = ppu.mmu.ReadByte(sourceAddr + uint16(i))
		}
	case addr == BGP:
		ppu.BGP = value
	case addr == OBP0:
		ppu.OBP0 = value
	case addr == OBP1:
		ppu.OBP1 = value
	case addr == WY:
		ppu.WY = value
	case addr == WX:
		ppu.WX = value
	case addr >= 0x8000 && addr <= 0x9FFF:
		ppu.VRAM[addr&0x1FFF] = value
	}
}

// setLCDCFields takes a byte written to LCDC
// and extracts the attributes to set fields on the GPU Struct
func (ppu *PPU) setLCDCFields(value byte) {
	// LCDC Bit 7 - LCD Display Enable             (0=Off, 1=On)
	ppu.lcdEnabled = utils.IsBitSet(value, 7)

	// LCDC Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	if utils.IsBitSet(value, 6) {
		ppu.windowMapLocation = 0x9C00
	} else {
		ppu.windowMapLocation = 0x9800
	}

	// LCDC Bit 5 - Window Display Enable          (0=Off, 1=On)
	ppu.windowEnabled = utils.IsBitSet(value, 5)

	// LCDC Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
	if utils.IsBitSet(value, 4) {
		ppu.tileDataLocation = 0x8000
	} else {
		ppu.tileDataLocation = 0x8800
	}

	// LCDC Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
	if utils.IsBitSet(value, 3) {
		ppu.backgroundMapLocation = 0x9C00
	} else {
		ppu.backgroundMapLocation = 0x9800
	}

	// LCDC Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
	if utils.IsBitSet(value, 2) {
		ppu.spriteSize = 16
	} else {
		ppu.spriteSize = 8
	}

	// LCDC Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
	ppu.spriteEnabled = utils.IsBitSet(value, 1)

	// LCDC Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
	ppu.backgroundEnabled = utils.IsBitSet(value, 0)
}

func (ppu *PPU) Step(cycles int) {
	if ppu.lcdEnabled {
		ppu.Cycles += cycles

		// STAT indicates the current status of the LCD controller.
		switch ppu.STAT.mode {
		// HBlank
		// After the last HBlank, push the screen data to canvas
		case MODE0:
			if ppu.Cycles >= 204 {
				// Reset the cycle counter
				ppu.Cycles = 0

				// Increase the scanline
				ppu.LY++

				// 143 is the last line, update the screen and enter VBlank
				if ppu.LY == 144 {
					// Request VBLANK interrupt
					ppu.mmu.RequestInterrupt(VBLANK_INTERRUPT)

					// Enter GPU Mode 1/VBlank
					ppu.STAT.mode = MODE1
					if ppu.STAT.vblankInterruptEnabled {
						ppu.mmu.RequestInterrupt(LCDC_INTERRUPT)
					}
				} else {
					// Enter GPU Mode 2/OAM Access
					ppu.STAT.mode = MODE2
					if ppu.STAT.oamInterruptEnabled {
						ppu.mmu.RequestInterrupt(LCDC_INTERRUPT)
					}
				}
			}

		// VBlank
		// After 10 lines, restart scanline and draw the next frame
		case MODE1:
			if ppu.Cycles >= 456 {
				// Reset the cycle counter
				ppu.Cycles = 0

				// Increase the scanline
				ppu.LY++

				// If Scanline is 153 we have done 10 lines of VBlank
				if ppu.LY > 153 {
					// Start of next Frame
					// Enter GPU Mode 2/OAM Access
					ppu.STAT.mode = MODE2
					if ppu.STAT.oamInterruptEnabled {
						ppu.mmu.RequestInterrupt(LCDC_INTERRUPT)
					}

					// Reset the Scanline
					ppu.LY = 0
				}
			}

		// OAM access mode, scanline active
		case MODE2:
			if ppu.Cycles >= 80 {
				// Reset the cycle counter
				ppu.Cycles = 0
				// Enter GPU Mode 3
				ppu.STAT.mode = MODE3
			}

		// VRAM access mode, scanline active
		// Treat end of mode 3 as end of scanline
		case MODE3:
			if ppu.Cycles >= 172 {
				// Reset the cycle counter
				ppu.Cycles = 0

				// Enter GPU Mode 0/HBlank
				ppu.STAT.mode = MODE0
				if ppu.STAT.hblankInterruptEnabled {
					ppu.mmu.RequestInterrupt(LCDC_INTERRUPT)
				}

				// Write a scanline to the framebuffer
				ppu.renderScanline()
			}
		}

		// If LY == LYC then set the STAT match flag and perform
		// the match flag interrupt if it has been requested
		if ppu.LY == ppu.LYC {
			ppu.STAT.coincidenceFlag = true

			if ppu.STAT.coincidenceInterruptEnabled {
				ppu.mmu.RequestInterrupt(LCDC_INTERRUPT)
			}
		} else {
			ppu.STAT.coincidenceFlag = false
		}

	} else {
		ppu.Cycles = 456
		ppu.LY = 0
		ppu.STAT.mode = MODE0
	}
}

func (ppu *PPU) generateTileScanline() [160]*pixelAttributes {
	var scanline [160]*pixelAttributes

	windowEnabled := false
	if ppu.windowEnabled {
		if ppu.WY <= ppu.LY {
			windowEnabled = true
		}
	}

	tileMap := ppu.backgroundMapLocation
	pixelY := ppu.LY + ppu.SCY

	if windowEnabled {
		tileMap = ppu.windowMapLocation
		// Add ScrollY to the Scanline to get the current pixel Y position
		pixelY = ppu.LY - ppu.WY
	}

	for pixel := byte(0); pixel < 160; pixel++ {
		// Add pixel being drawn in scanline to scrollX to get the current pixel X position
		pixelX := pixel + ppu.SCX

		// translate the current x pos to window space if necessary
		if windowEnabled {
			if pixel >= ppu.WX {
				pixelX = pixel - (ppu.WX - 7)
			}
		}

		colorIdentifier := ppu.getTileColorIdentifierForPixel(tileMap, pixelX, pixelY)

		scanline[pixel] = &pixelAttributes{
			colorIdentifier: colorIdentifier,
			palette:         ppu.BGP,
		}
	}

	return scanline
}

func (ppu *PPU) generateSpriteScanline() [160]*pixelAttributes {
	var scanline [160]*pixelAttributes

	for sprite := 0; sprite < 40; sprite++ {
		index := sprite * 4
		yPos := ppu.OAM[index] - 16
		xPos := ppu.OAM[index+1] - 8
		characterCode := ppu.OAM[index+2]
		attributes := ppu.OAM[index+3]

		yFlip := utils.IsBitSet(attributes, 6)
		xFlip := utils.IsBitSet(attributes, 5)

		if (ppu.LY >= yPos) && (ppu.LY < (yPos + ppu.spriteSize)) {
			line := int(ppu.LY - yPos)

			if yFlip {
				line -= int(ppu.spriteSize)
				line *= -1
			}

			line *= 2 // same as for tiles

			data1, data2 := ppu.getSpriteDataForLine(characterCode, line)

			// its easier to read in from right to left as pixel 0 is
			// bit 7 in the colour data, pixel 1 is bit 6 etc...
			for tilePixel := 7; tilePixel >= 0; tilePixel-- {
				pixelBit := tilePixel
				// read the sprite in backwards for the x axis
				if xFlip {
					pixelBit -= 7
					pixelBit *= -1
				}

				var colorIdentifier byte
				if utils.IsBitSet(data1, byte(pixelBit)) {
					colorIdentifier = utils.SetBit(colorIdentifier, 0)
				}
				if utils.IsBitSet(data2, byte(pixelBit)) {
					colorIdentifier = utils.SetBit(colorIdentifier, 1)
				}

				pixel := int(xPos) + (7 - tilePixel)

				var palette byte
				if utils.IsBitSet(attributes, 4) {
					palette = ppu.OBP1
				} else {
					palette = ppu.OBP0
				}

				prioritySet := utils.IsBitSet(attributes, 7)
				var priority byte
				if prioritySet {
					priority = 1
				}

				if (pixel >= 0) && (pixel <= 159) && colorIdentifier != 0 {
					scanline[pixel] = &pixelAttributes{
						colorIdentifier: colorIdentifier,
						palette:         palette,
						priority:        priority,
					}
				}
			}
		}
	}
	return scanline
}

func (ppu *PPU) renderScanline() {
	if ppu.backgroundEnabled {

		var scanline [2][160]*pixelAttributes
		scanline[0] = ppu.generateTileScanline()

		if ppu.spriteEnabled {
			scanline[1] = ppu.generateSpriteScanline()
		}

		for x := 0; x < 160; x++ {
			var pixel uint32
			backgroundPixel := scanline[0][x]
			spritePixel := scanline[1][x]

			// If there is a sprite at this position in the scanline
			// and the sprite priority is 0 or the background pixels
			// colorIdentifier is 0, then the sprite is rendered on top
			// of the background, otherwise the background is rendered.
			if spritePixel != nil && (spritePixel.priority == 0 || backgroundPixel.colorIdentifier == 0) {
				pixel = ppu.applyPalette(spritePixel.colorIdentifier, spritePixel.palette)
			} else {
				pixel = ppu.applyPalette(backgroundPixel.colorIdentifier, backgroundPixel.palette)
			}

			ppu.Window.Framebuffer[x+160*int(ppu.LY)] = pixel
		}
	}
}

func (ppu *PPU) getTileColorIdentifierForPixel(tileMap uint16, pixelX byte, pixelY byte) byte {
	tileIdentifier := ppu.getTileIdentifierForPixel(tileMap, pixelX, pixelY)
	tileDataAddress := ppu.getTileDataAddress(ppu.tileDataLocation, tileIdentifier)

	// Find the correct vertical line we're on of the
	// tile to get the tile data from memory
	line := pixelY % 8
	line = line * 2 // each vertical line takes up two bytes of memory

	data1 := ppu.mmu.ReadByte(tileDataAddress + uint16(line))
	data2 := ppu.mmu.ReadByte(tileDataAddress + uint16(line) + 1)

	// pixel 0 in the tile is it 7 of data 1 and data2.
	// Pixel 1 is bit 6 etc..
	pixelBit := int(pixelX % 8)
	pixelBit -= 7
	pixelBit *= -1

	var colorIdentifier byte
	if utils.IsBitSet(data1, byte(pixelBit)) {
		colorIdentifier = utils.SetBit(colorIdentifier, 0)
	}
	if utils.IsBitSet(data2, byte(pixelBit)) {
		colorIdentifier = utils.SetBit(colorIdentifier, 1)
	}

	return colorIdentifier
}

func (ppu *PPU) getTileIdentifierForPixel(tileMap uint16, pixelX byte, pixelY byte) byte {
	// Divide the pixelY position by 8 (for 8 pixels in tile)
	// and multiply by 32 (for number of tiles in the background map)
	// to get the row number for the tile in the background map
	tileRow := uint16(pixelY/8) * 32

	// Divide the pixelX position by 8 (for 8 tiles in horizontal row)
	// to get the column number for the tile in the background map
	tileCol := uint16(pixelX / 8)

	return ppu.mmu.ReadByte(tileMap + tileRow + tileCol)
}

func (ppu *PPU) getTileDataAddress(tileMap uint16, tileIdentifier byte) uint16 {
	// When the tileMap used is 0x8800 the tileIndentifier is
	// a signed byte -127 - 127, the offset corrects for this
	// when looking up the memory location
	offset := uint16(0)
	if tileMap == 0x8800 {
		offset = 128
	}

	if tileIdentifier > 127 {
		return tileMap + (uint16(tileIdentifier)-offset)*16 // 16 = tile size in bytes
	} else {
		return tileMap + (uint16(tileIdentifier)+offset)*16 // 16 = tile size in bytes

	}
}

func (ppu *PPU) getSpriteDataForLine(characterCode byte, line int) (byte, byte) {
	spriteDataStorage := 0x8000
	spriteDataAddress := spriteDataStorage + (int(characterCode) * 16) + line // 16 = obj size in bytes

	data1 := ppu.mmu.ReadByte(uint16(spriteDataAddress))
	data2 := ppu.mmu.ReadByte(uint16(spriteDataAddress + 1))

	return data1, data2
}

func (ppu *PPU) applyPalette(colorIdentifier byte, palette byte) uint32 {
	pixelFormat, _ := sdl.AllocFormat(uint(sdl.PIXELFORMAT_RGBA32))
	var color byte
	var bitmask byte = 0x3
	switch colorIdentifier {
	case 0:
		color = palette & bitmask
	case 1:
		color = (palette >> 2) & bitmask
	case 2:
		color = (palette >> 4) & bitmask
	case 3:
		color = (palette >> 6) & bitmask
	}

	switch color {
	case 0:
		return sdl.MapRGB(pixelFormat, 224, 248, 208)
	case 1:
		return sdl.MapRGB(pixelFormat, 136, 192, 112)
	case 2:
		return sdl.MapRGB(pixelFormat, 52, 104, 86)
	case 3:
		return sdl.MapRGB(pixelFormat, 8, 24, 32)
	}
	return 0
}
