package Gameboy

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WIDTH  = 160
	HEIGHT = 144
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
	s.coincidenceInterruptEnabled = IsBitSet(value, 6)
	s.oamInterruptEnabled = IsBitSet(value, 5)
	s.vblankInterruptEnabled = IsBitSet(value, 4)
	s.hblankInterruptEnabled = IsBitSet(value, 3)
}

func (s *stat) getStat() byte {
	// Start with 0x80 as bit 7 is unused/always returns 1
	var stat byte = 0x80

	// Bit 6 LYC=LY Coincidence Interrupt
	if s.coincidenceInterruptEnabled {
		stat = SetBit(stat, 6)
	}

	// Bit 5 Mode 2 OAM Interrupt
	if s.oamInterruptEnabled {
		stat = SetBit(stat, 5)
	}

	// Bit 4 Mode 1 V-Blank Interrupt
	if s.vblankInterruptEnabled {
		stat = SetBit(stat, 4)
	}

	// Bit 3 Mode 0 H-Blank Interrupt
	if s.hblankInterruptEnabled {
		stat = SetBit(stat, 3)
	}

	// Bit 2 Coincidence Flag (0:LYC<>LY, 1:LYC=LY)
	if s.coincidenceFlag {
		stat = SetBit(stat, 2)
	}

	// or with the mode to set this on stat
	stat |= s.mode
	return stat
}

type GPU struct {
	gameboy *Gameboy

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

	Cycles uint // Number of cycles since the last LCD Status Mode Change
}

type pixelAttributes struct {
	colorIdentifier byte
	palette         byte
	priority        byte
}

func NewGPU(gameboy *Gameboy) *GPU {
	window := NewWindow("Gameboy", WIDTH, HEIGHT, gameboy.Quit)

	STAT := &stat{
		coincidenceInterruptEnabled: false,
		oamInterruptEnabled:         false,
		vblankInterruptEnabled:      false,
		hblankInterruptEnabled:      false,
		coincidenceFlag:             true,
		mode:                        MODE1,
	}

	gpu := &GPU{
		gameboy: gameboy,
		Window:  window,
		STAT:    STAT,
		LCDC:    0x91,
		SCY:     0x00,
		SCX:     0x00,
		LYC:     0x00,
		BGP:     0xFC,
		OBP0:    0xFF,
		OBP1:    0xFF,
		WY:      0x00,
		WX:      0x00,
	}

	gpu.setLCDCFields(0x91)

	return gpu
}

func (gpu *GPU) ReadByte(addr uint16) byte {
	switch {
	case addr == LCDC:
		return gpu.LCDC
	case addr == STAT:
		return gpu.STAT.getStat()
	case addr == SCY:
		return gpu.SCY
	case addr == SCX:
		return gpu.SCX
	case addr == LY:
		return gpu.LY
	case addr == BGP:
		return gpu.BGP
	case addr == OBP0:
		return gpu.OBP0
	case addr == OBP1:
		return gpu.OBP1
	case addr == WY:
		return gpu.WY
	case addr == WX:
		return gpu.WX
	case addr >= 0x8000 && addr <= 0x9FFF:
		return gpu.VRAM[addr&0x1FFF]
	case addr >= 0xFE00 && addr <= 0xFE9F:
		return gpu.OAM[addr&0x9F]
	}
	return 0
}

func (gpu *GPU) WriteByte(addr uint16, value byte) {
	switch {
	case addr == LCDC:
		gpu.LCDC = value
		gpu.setLCDCFields(value)
	case addr == STAT:
		gpu.STAT.setStat(value)
	case addr == SCY:
		gpu.SCY = value
	case addr == SCX:
		gpu.SCX = value
	case addr == LY:
		// If the game writes to scanline it should be unset
		gpu.LY = 0
	case addr == LYC:
		gpu.LYC = value
	case addr == DMA:
		// The value holds the source address of the OAM data divided by 100
		// so we have to multiply it first
		var sourceAddr uint16 = uint16(value) << 8

		for i := 0; i < 160; i++ {
			gpu.OAM[i] = gpu.gameboy.ReadByte(sourceAddr + uint16(i))
		}
	case addr == BGP:
		gpu.BGP = value
	case addr == OBP0:
		gpu.OBP0 = value
	case addr == OBP1:
		gpu.OBP1 = value
	case addr == WY:
		gpu.WY = value
	case addr == WX:
		gpu.WX = value
	case addr >= 0x8000 && addr <= 0x9FFF:
		gpu.VRAM[addr&0x1FFF] = value
	}
}

// setLCDCFields takes a byte written to LCDC
// and extracts the attributes to set fields on the GPU Struct
func (gpu *GPU) setLCDCFields(value byte) {
	// LCDC Bit 7 - LCD Display Enable             (0=Off, 1=On)
	gpu.lcdEnabled = IsBitSet(value, 7)

	// LCDC Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	if IsBitSet(value, 6) {
		gpu.windowMapLocation = 0x9C00
	} else {
		gpu.windowMapLocation = 0x9800
	}

	// LCDC Bit 5 - Window Display Enable          (0=Off, 1=On)
	gpu.windowEnabled = IsBitSet(value, 5)

	// LCDC Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
	if IsBitSet(value, 4) {
		gpu.tileDataLocation = 0x8000
	} else {
		gpu.tileDataLocation = 0x8800
	}

	// LCDC Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
	if IsBitSet(value, 3) {
		gpu.backgroundMapLocation = 0x9C00
	} else {
		gpu.backgroundMapLocation = 0x9800
	}

	// LCDC Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
	if IsBitSet(value, 2) {
		gpu.spriteSize = 16
	} else {
		gpu.spriteSize = 8
	}

	// LCDC Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
	gpu.spriteEnabled = IsBitSet(value, 1)

	// LCDC Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
	gpu.backgroundEnabled = IsBitSet(value, 0)
}

func (gpu *GPU) Step(cycles byte) {
	if gpu.lcdEnabled {
		gpu.Cycles += uint(cycles * 4)

		// STAT indicates the current status of the LCD controller.
		switch gpu.STAT.mode {
		// HBlank
		// After the last HBlank, push the screen data to canvas
		case MODE0:
			if gpu.Cycles >= 204 {
				// Reset the cycle counter
				gpu.Cycles = 0

				// Increase the scanline
				gpu.LY++

				// 143 is the last line, update the screen and enter VBlank
				if gpu.LY == 144 {
					// Request VBLANK interrupt
					gpu.gameboy.requestInterrupt(VBLANK_INTERRUPT)

					// Enter GPU Mode 1/VBlank
					gpu.STAT.mode = MODE1
					if gpu.STAT.vblankInterruptEnabled {
						gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
					}
				} else {
					// Enter GPU Mode 2/OAM Access
					gpu.STAT.mode = MODE2
					if gpu.STAT.oamInterruptEnabled {
						gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
					}
				}
			}

		// VBlank
		// After 10 lines, restart scanline and draw the next frame
		case MODE1:
			if gpu.Cycles >= 4560 {
				// Reset the cycle counter
				gpu.Cycles = 0

				// Increase the scanline
				gpu.LY++

				// If Scanline is 153 we have done 10 lines of VBlank
				if gpu.LY > 153 {
					// Start of next Frame
					// Enter GPU Mode 2/OAM Access
					gpu.STAT.mode = MODE2
					if gpu.STAT.oamInterruptEnabled {
						gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
					}

					// Reset the Scanline
					gpu.LY = 0
				}
			}

		// OAM access mode, scanline active
		case MODE2:
			if gpu.Cycles >= 80 {
				// Reset the cycle counter
				gpu.Cycles = 0
				// Enter GPU Mode 3
				gpu.STAT.mode = MODE3
			}

		// VRAM access mode, scanline active
		// Treat end of mode 3 as end of scanline
		case MODE3:
			if gpu.Cycles >= 172 {
				// Reset the cycle counter
				gpu.Cycles = 0

				// Enter GPU Mode 0/HBlank
				gpu.STAT.mode = MODE0
				if gpu.STAT.hblankInterruptEnabled {
					gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
				}

				// Write a scanline to the framebuffer
				gpu.renderScanline()
			}
		}

		// If LY == LYC then set the STAT match flag and perform
		// the match flag interrupt if it has been requested
		if gpu.LY == gpu.LYC {
			gpu.STAT.coincidenceFlag = true

			if gpu.STAT.coincidenceInterruptEnabled {
				gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
			}
		} else {
			gpu.STAT.coincidenceFlag = false
		}

	} else {
		gpu.Cycles = 456
		gpu.LY = 0
		gpu.STAT.mode = MODE0
	}
}

func (gpu *GPU) generateTileScanline() [160]*pixelAttributes {
	var scanline [160]*pixelAttributes

	windowEnabled := false
	if gpu.windowEnabled {
		if gpu.WY <= gpu.LY {
			windowEnabled = true
		}
	}

	tileMap := gpu.backgroundMapLocation
	pixelY := gpu.LY + gpu.SCY

	if windowEnabled {
		tileMap = gpu.windowMapLocation
		// Add ScrollY to the Scanline to get the current pixel Y position
		pixelY = gpu.LY - gpu.WY
	}

	for pixel := byte(0); pixel < 160; pixel++ {
		// Add pixel being drawn in scanline to scrollX to get the current pixel X position
		pixelX := pixel + gpu.SCX

		// translate the current x pos to window space if necessary
		if windowEnabled {
			if pixel >= gpu.WX {
				pixelX = pixel - gpu.WX
			}
		}

		colorIdentifier := gpu.getTileColorIdentifierForPixel(tileMap, pixelX, pixelY)

		scanline[pixel] = &pixelAttributes{
			colorIdentifier: colorIdentifier,
			palette:         gpu.BGP,
		}
	}

	return scanline
}

func (gpu *GPU) generateSpriteScanline() [160]*pixelAttributes {
	var scanline [160]*pixelAttributes

	for sprite := 0; sprite < 40; sprite++ {
		index := sprite * 4
		yPos := gpu.OAM[index] - 16
		xPos := gpu.OAM[index+1] - 8
		characterCode := gpu.OAM[index+2]
		attributes := gpu.OAM[index+3]

		yFlip := IsBitSet(attributes, 6)
		xFlip := IsBitSet(attributes, 5)

		if (gpu.LY >= yPos) && (gpu.LY < (yPos + gpu.spriteSize)) {
			line := int(gpu.LY - yPos)

			if yFlip {
				line -= int(gpu.spriteSize)
				line *= -1
			}

			line *= 2 // same as for tiles

			data1, data2 := gpu.getSpriteDataForLine(characterCode, line)

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
				if IsBitSet(data1, byte(pixelBit)) {
					colorIdentifier = SetBit(colorIdentifier, 1)
				}
				if IsBitSet(data2, byte(pixelBit)) {
					colorIdentifier = SetBit(colorIdentifier, 0)
				}

				pixel := int(xPos) + (7 - tilePixel)

				var palette byte
				if IsBitSet(attributes, 4) {
					palette = gpu.OBP1
				} else {
					palette = gpu.OBP0
				}

				prioritySet := IsBitSet(attributes, 7)
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

func (gpu *GPU) renderScanline() {
	if gpu.backgroundEnabled {

		var scanline [2][160]*pixelAttributes
		scanline[0] = gpu.generateTileScanline()

		if gpu.spriteEnabled {
			scanline[1] = gpu.generateSpriteScanline()
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
				pixel = gpu.applyPalette(spritePixel.colorIdentifier, spritePixel.palette)
			} else {
				pixel = gpu.applyPalette(backgroundPixel.colorIdentifier, backgroundPixel.palette)
			}

			gpu.Window.Framebuffer[x+160*int(gpu.LY)] = pixel
		}
	} else {
		for x := 0; x < 160; x++ {
			pixelFormat, _ := sdl.AllocFormat(uint(sdl.PIXELFORMAT_RGBA32))
			gpu.Window.Framebuffer[x+160*int(gpu.LY)] = sdl.MapRGB(pixelFormat, 255, 255, 255)
		}
	}
}

func (gpu *GPU) getTileColorIdentifierForPixel(tileMap uint16, pixelX byte, pixelY byte) byte {
	tileIdentifier := gpu.getTileIdentifierForPixel(tileMap, pixelX, pixelY)
	tileDataAddress := gpu.getTileDataAddress(gpu.tileDataLocation, tileIdentifier)

	// Find the correct vertical line we're on of the
	// tile to get the tile data from memory
	line := pixelY % 8
	line = line * 2 // each vertical line takes up two bytes of memory

	data1 := gpu.gameboy.ReadByte(tileDataAddress + uint16(line))
	data2 := gpu.gameboy.ReadByte(tileDataAddress + uint16(line) + 1)

	// pixel 0 in the tile is it 7 of data 1 and data2.
	// Pixel 1 is bit 6 etc..
	pixelBit := int(pixelX % 8)
	pixelBit -= 7
	pixelBit *= -1

	var colorIdentifier byte
	if IsBitSet(data1, byte(pixelBit)) {
		colorIdentifier = SetBit(colorIdentifier, 1)
	}
	if IsBitSet(data2, byte(pixelBit)) {
		colorIdentifier = SetBit(colorIdentifier, 0)
	}

	return colorIdentifier
}

func (gpu *GPU) getTileIdentifierForPixel(tileMap uint16, pixelX byte, pixelY byte) byte {
	// Divide the pixelY position by 8 (for 8 pixels in tile)
	// and multiply by 32 (for number of tiles in the background map)
	// to get the row number for the tile in the background map
	tileRow := uint16(pixelY/8) * 32

	// Divide the pixelX position by 8 (for 8 tiles in horizontal row)
	// to get the column number for the tile in the background map
	tileCol := uint16(pixelX / 8)

	return gpu.gameboy.ReadByte(tileMap + tileRow + tileCol)
}

func (gpu *GPU) getTileDataAddress(tileMap uint16, tileIdentifier byte) uint16 {
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

func (gpu *GPU) getSpriteDataForLine(characterCode byte, line int) (byte, byte) {
	spriteDataStorage := 0x8000
	spriteDataAddress := spriteDataStorage + (int(characterCode) * 16) + line // 16 = obj size in bytes

	data1 := gpu.gameboy.ReadByte(uint16(spriteDataAddress))
	data2 := gpu.gameboy.ReadByte(uint16(spriteDataAddress + 1))

	return data1, data2
}

func (gpu *GPU) applyPalette(colorIdentifier byte, palette byte) uint32 {
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
		return sdl.MapRGB(pixelFormat, 255, 255, 255)
	case 1:
		return sdl.MapRGB(pixelFormat, 192, 192, 192)
	case 2:
		return sdl.MapRGB(pixelFormat, 96, 96, 96)
	case 3:
		return sdl.MapRGB(pixelFormat, 0, 0, 0)
	}
	return 0
}
