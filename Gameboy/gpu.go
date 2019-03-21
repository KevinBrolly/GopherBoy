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

	MATCH_FLAG      = 2 // LYC=LY Flag
	MODE0_INTERRUPT = 3
	MODE1_INTERRUPT = 4
	MODE2_INTERRUPT = 5
	MATCH_INTERRUPT = 6 // LYC=LY Interrupt
)

type GPU struct {
	gameboy *Gameboy

	Window *Window

	VRAM [16384]byte
	OAM  [160]byte

	STAT byte // LCD Status/Mode
	SCY  byte // Scroll Y
	SCX  byte // Scroll X
	LY   byte // Scanline
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

func NewGPU(gameboy *Gameboy) *GPU {
	window := NewWindow("Gameboy", WIDTH, HEIGHT, gameboy.Quit)

	gpu := &GPU{
		gameboy: gameboy,
		Window:  window,
		LCDC:    0x91,
		SCY:     0x00,
		SCX:     0x00,
		LYC:     0x00,
		BGP:     0xFC,
		OBP0:    0xFF,
		OBP1:    0xFF,
		WY:      0x00,
		WX:      0x00,
		STAT:    0x85,
	}
	return gpu
}

func (gpu *GPU) ReadByte(addr uint16) byte {
	switch {
	case addr == LCDC:
		return gpu.LCDC
	case addr == STAT:
		return gpu.STAT
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
		gpu.STAT = (0xF8 & value) | gpu.GetSTATMode()
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

func (gpu *GPU) GetSTATMode() byte {
	return gpu.STAT & 0x03
}

func (gpu *GPU) SetSTATMode(mode byte) {
	gpu.STAT = (gpu.STAT & 0xFC) | mode
}

func (gpu *GPU) Step(cycles byte) {
	if gpu.lcdEnabled {
		gpu.Cycles += uint(cycles * 4)

		// STAT indicates the current status of the LCD controller.
		switch gpu.GetSTATMode() {
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
					gpu.SetSTATMode(MODE1)
					if IsBitSet(gpu.STAT, MODE1_INTERRUPT) {
						gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
					}
				} else {
					// Enter GPU Mode 2/OAM Access
					gpu.SetSTATMode(MODE2)
					if IsBitSet(gpu.STAT, MODE2_INTERRUPT) {
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
					gpu.SetSTATMode(MODE2)
					if IsBitSet(gpu.STAT, MODE2_INTERRUPT) {
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
				gpu.SetSTATMode(MODE3)
			}

		// VRAM access mode, scanline active
		// Treat end of mode 3 as end of scanline
		case MODE3:
			if gpu.Cycles >= 172 {
				// Reset the cycle counter
				gpu.Cycles = 0

				// Enter GPU Mode 0/HBlank
				gpu.SetSTATMode(MODE0)
				if IsBitSet(gpu.STAT, MODE0_INTERRUPT) {
					gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
				}

				// Write a scanline to the framebuffer
				gpu.renderScanline()
			}
		}

		// If LY == LYC then set the STAT match flag and perform
		// the match flag interrupt if it has been requested
		if gpu.LY == gpu.LYC {
			gpu.STAT = SetBit(gpu.STAT, MATCH_FLAG)

			if IsBitSet(gpu.STAT, MATCH_INTERRUPT) {
				gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
			}
		} else {
			gpu.STAT = ClearBit(gpu.STAT, MATCH_FLAG)
		}

	} else {
		gpu.Cycles = 456
		gpu.LY = 0
		gpu.SetSTATMode(MODE1)
	}
}

func (gpu *GPU) renderTiles() {
	windowEnabled := false
	if gpu.windowEnabled {
		if gpu.WY <= gpu.LY {
			windowEnabled = true
		}
	}

	tileDataAreaSelection := gpu.backgroundMapLocation
	pixelY := gpu.LY + gpu.SCY

	if windowEnabled {
		tileDataAreaSelection = gpu.windowMapLocation
		// Add ScrollY to the Scanline to get the current pixel Y position
		pixelY = gpu.LY - gpu.WY
	}

	// Divide the pixel being drawn Y position by 8 (for 8 pixels in tile)
	// and multiply by 32 (for number of tiles in the background map)
	// to get the row number for the tile in the background map
	tileRow := uint16(pixelY/8) * 32

	for pixel := byte(0); pixel < 160; pixel++ {
		// Add pixel being drawn in scanline to scrollX to get the current pixel X position
		pixelX := pixel + gpu.SCX

		// translate the current x pos to window space if necessary
		if windowEnabled {
			if pixel >= gpu.WX {
				pixelX = pixel - gpu.WX
			}
		}

		// which of the 8 horizontal tiles does this pixel fall within?
		tileCol := uint16(pixelX / 8)

		tileIdentifier := gpu.gameboy.ReadByte(tileDataAreaSelection + tileRow + tileCol)
		tileDataAddress := gpu.getTileDataAddress(tileIdentifier)

		data1, data2 := gpu.getTileData(tileDataAddress, pixelY)

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

		gpu.Window.Framebuffer[int(pixel)+(160*int(gpu.LY))] = gpu.getColorFromBGPalette(colorIdentifier)
	}
}

func (gpu *GPU) renderSprites() {
	for sprite := 0; sprite < 40; sprite++ {
		index := sprite * 4
		yPos := gpu.OAM[index] - 16
		xPos := gpu.OAM[index+1] - 8
		characterCode := gpu.OAM[index+2]
		attributes := gpu.OAM[index+3]

		yFlip := IsBitSet(attributes, 6)
		xFlip := IsBitSet(attributes, 5)

		if (gpu.LY >= yPos) && (gpu.LY < (yPos + gpu.spriteSize)) {
			data1, data2 := gpu.getObjData(characterCode, yPos, yFlip)

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

				var pixel int = int(xPos) + (7 - tilePixel)

				gpu.Window.Framebuffer[pixel+(160*int(gpu.LY))] = gpu.getSpritePalette(colorIdentifier, attributes)
			}
		}
	}
}

func (gpu *GPU) renderScanline() {
	if gpu.backgroundEnabled {
		gpu.renderTiles()
	}

	if gpu.spriteEnabled {
		gpu.renderSprites()
	}
}

func (gpu *GPU) getTileDataAddress(tileIdentifier byte) uint16 {
	// When the BGCharacterDataSelection is 0x8800 the tileIndentifier is
	// a signed byte -127 - 127, the offset corrects for this
	// when looking up the memory location
	offset := uint16(0)

	if gpu.tileDataLocation == 0x8800 {
		offset = 128
	}

	if tileIdentifier > 127 {
		return gpu.tileDataLocation + (uint16(tileIdentifier)-offset)*16 // 16 = tile size in bytes
	} else {
		return gpu.tileDataLocation + (uint16(tileIdentifier)+offset)*16 // 16 = tile size in bytes

	}
}

func (gpu *GPU) getTileData(tileDataAddress uint16, pixelY byte) (byte, byte) {
	// find the correct vertical line we're on of the
	// tile to get the tile data
	// from in memory
	line := pixelY % 8
	line = line * 2 // each vertical line takes up two bytes of memory
	data1 := gpu.gameboy.ReadByte(tileDataAddress + uint16(line))
	data2 := gpu.gameboy.ReadByte(tileDataAddress + uint16(line) + 1)

	return data1, data2
}

func (gpu *GPU) getObjData(characterCode byte, yPos byte, yFlip bool) (byte, byte) {
	line := int(gpu.LY - yPos)

	if yFlip {
		line -= int(gpu.spriteSize)
		line *= -1
	}

	line *= 2 // same as for tiles

	objCharacterDataStorage := 0x8000
	objDataAddress := objCharacterDataStorage + (int(characterCode) * 16) + line // 16 = obj size in bytes

	data1 := gpu.gameboy.ReadByte(uint16(objDataAddress))
	data2 := gpu.gameboy.ReadByte(uint16(objDataAddress + 1))

	return data1, data2
}

func (gpu *GPU) getColorFromBGPalette(colorIdentifier byte) uint32 {
	pixelFormat, _ := sdl.AllocFormat(uint(sdl.PIXELFORMAT_RGBA32))
	var color byte
	var bitmask byte = 0x3
	switch colorIdentifier {
	case 0:
		color = gpu.BGP & bitmask
	case 1:
		color = (gpu.BGP >> 2) & bitmask
	case 2:
		color = (gpu.BGP >> 4) & bitmask
	case 3:
		color = (gpu.BGP >> 6) & bitmask
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

func (gpu *GPU) getSpritePalette(colorIdentifier byte, attributes byte) uint32 {
	pixelFormat, _ := sdl.AllocFormat(uint(sdl.PIXELFORMAT_RGBA32))

	var palette byte
	if IsBitSet(attributes, 4) {
		palette = gpu.OBP1
	} else {
		palette = gpu.OBP0
	}

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
		return sdl.MapRGBA(pixelFormat, 255, 255, 255, 0)
	case 1:
		return sdl.MapRGB(pixelFormat, 192, 192, 192)
	case 2:
		return sdl.MapRGB(pixelFormat, 96, 96, 96)
	case 3:
		return sdl.MapRGB(pixelFormat, 0, 0, 0)
	}
	return 0
}
