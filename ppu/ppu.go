package ppu

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
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

	// Color-related Addresses
	VRAMBank = 0xFF4F // VRAM Bank
	BGPI     = 0xFF68 // Background Palette Index
	BGPD     = 0xFF69 // Background Palette Data
	OBPI     = 0xFF6A // Sprite Palette Index
	OBPD     = 0xFF6B // Sprite Palette Data
	HDMA1    = 0xFF51 // New DMA Source, High
	HDMA3    = 0xFF53 // New DMA Destination, High
	HDMA2    = 0xFF52 // New DMA Source, Low
	HDMA4    = 0xFF54 // New DMA Destination, Low
	HDMA5    = 0xFF55 // New DMA Length/Mode/Start
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

	FrameBuffer *image.RGBA

	VRAM [16384]byte

	OAM            []*Sprite
	VisibleSprites []*Sprite

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

	// Color-related Addresses
	VRAMBank byte // VRAM Bank

	// FF68 - BCPS/BGPI - CGB Mode Only - Background Palette Index
	backgroundPaletteIndex         byte
	backgroundPaletteAutoIncrement bool

	// FF69 - BCPD/BGPD - CGB Mode Only - Background Palette Data
	backgroundPaletteData [0x40]byte

	// FF6A - OCPS/OBPI - CGB Mode Only - Sprite Palette Index
	spritePaletteIndex         byte
	spritePaletteAutoIncrement bool

	// FF6B - OCPD/OBPD - CGB Mode Only - Sprite Palette Data
	spritePaletteData [0x40]byte

	HDMA1 byte // New DMA Source, High
	HDMA3 byte // New DMA Destination, High
	HDMA2 byte // New DMA Source, Low
	HDMA4 byte // New DMA Destination, Low
	HDMA5 byte // New DMA Length/Mode/Start

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

	rectImage := image.NewRGBA(image.Rect(0, 0, 160, 144))
	draw.Draw(rectImage, rectImage.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	oam := make([]*Sprite, 40)

	for i := range oam {
		oam[i] = &Sprite{}
	}

	ppu := &PPU{
		mmu:            mmu,
		OAM:            oam,
		VisibleSprites: make([]*Sprite, 10),
		FrameBuffer:    rectImage,
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
	case addr == VRAMBank:
		return ppu.VRAMBank
	case addr == BGPI:
		var value byte
		value = ppu.backgroundPaletteIndex
		if ppu.backgroundPaletteAutoIncrement {
			utils.SetBit(value, 7)
		}
		return value
	case addr == BGPD:
		return ppu.backgroundPaletteData[ppu.backgroundPaletteIndex]
	case addr == OBPI:
		var value byte
		value = ppu.spritePaletteIndex
		if ppu.spritePaletteAutoIncrement {
			utils.SetBit(value, 7)
		}
		return value
	case addr == OBPD:
		return ppu.spritePaletteData[ppu.spritePaletteIndex]
	case addr >= 0x8000 && addr <= 0x9FFF:
		return ppu.VRAM[addr&0x1FFF]
	case addr >= 0xFE00 && addr <= 0xFE9F:
		oamAddr := addr & 0x9F
		sprite := ppu.OAM[oamAddr/4] // 4 bits per sprite
		spriteBit := oamAddr % 4

		switch spriteBit {
		case 0:
			return sprite.Y
		case 1:
			return sprite.X
		case 2:
			return sprite.TileNumber
		case 3:
			return sprite.Attributes
		}
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
		for _, sprite := range ppu.OAM {
			sprite.Y = ppu.mmu.ReadByte(sourceAddr)
			sprite.X = ppu.mmu.ReadByte(sourceAddr + 1)
			sprite.TileNumber = ppu.mmu.ReadByte(sourceAddr + 2)
			sprite.Attributes = ppu.mmu.ReadByte(sourceAddr + 3)
			sourceAddr = sourceAddr + 4
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
	case addr == VRAMBank:
		ppu.VRAMBank = value
	case addr == BGPI:
		ppu.backgroundPaletteIndex = value & 0x1F
		if utils.IsBitSet(value, 7) {
			ppu.backgroundPaletteAutoIncrement = true
		}
	case addr == BGPD:
		ppu.backgroundPaletteData[ppu.backgroundPaletteIndex] = value
	case addr >= 0x8000 && addr <= 0x9FFF:
		ppu.VRAM[addr&0x1FFF] = value
	case addr >= 0xFE00 && addr <= 0xFE9F:
		oamAddr := addr & 0x9F
		sprite := ppu.OAM[oamAddr/4] // 4 bits per sprite
		spriteBit := oamAddr % 4

		switch spriteBit {
		case 0:
			sprite.Y = value
		case 1:
			sprite.X = value
		case 2:
			sprite.TileNumber = value
		case 3:
			sprite.Attributes = value
		}
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
				// Do OAMSearch
				ppu.OAMSearch()
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

func (ppu *PPU) OAMSearch() {
	visibleSprites := make([]*Sprite, 0)

	for _, sprite := range ppu.OAM {
		if len(visibleSprites) == 10 {
			break
		}

		if sprite.X != 0 && ((ppu.LY + 16) >= sprite.Y) && ((ppu.LY + 16) < (sprite.Y + ppu.spriteSize)) {
			visibleSprites = append(visibleSprites, sprite)
		}
	}

	ppu.VisibleSprites = visibleSprites
}

func (ppu *PPU) renderScanline() {
	if ppu.backgroundEnabled {

		fifo := &Fifo{}

		fetcher := &Fetcher{
			tileMapAddress:  ppu.backgroundMapLocation,
			tileDataAddress: ppu.tileDataLocation,
			SCY:             ppu.SCY,
			LY:              ppu.LY,
			memory:          ppu.mmu,
		}

		x := 0
		pushedDots := 0
		for pushedDots <= 159 {
			if fifo.Length() <= 8 {
				fifo.PushDots(fetcher.NextTileLine())
			}

			xSprites := make([]*Sprite, 0)
			for _, sprite := range ppu.VisibleSprites {
				if sprite != nil && (int(sprite.X-8)+int(ppu.SCX)) == x {
					xSprites = append(xSprites, sprite)
				}
			}
			if len(xSprites) != 0 {
				sprite := xSprites[len(xSprites)-1]
				dots := fetcher.fetchSpriteLine(sprite)
				fifo.OverlaySprite(dots, sprite.Palette(), sprite.Priority(), sprite.X)
			}

			if ppu.windowEnabled {
				if ppu.WY <= ppu.LY {
					if ppu.WX == byte(x) {
						fifo.Clear()
						fetcher.tileMapAddress = ppu.windowMapLocation
						fifo.PushDots(fetcher.NextTileLine())
					}
				}
			}

			dot := fifo.PopDot()

			palette := ppu.BGP
			if dot.Type == SPRITE {
				if dot.Palette == 0 {
					palette = ppu.OBP0
				} else {
					palette = ppu.OBP1
				}
			}

			if int(ppu.SCX) <= x {
				ppu.FrameBuffer.SetRGBA(pushedDots, int(ppu.LY), dot.ToRGBA(palette))
				pushedDots++

				// Background wraps around the screen (i.e. when part of it goes off the screen, it appears on the opposite side.)
				// So if pushed dots and background x start position == 256 then we reset the fetcher tileMapAddress
				// as the background wraps round
				if int(ppu.SCX)+pushedDots == 256 {
					fetcher.tileMapAddress = ppu.backgroundMapLocation
					fifo.Clear()
				}
			}

			x++
		}
	}
}
