package Gameboy

import (
    "github.com/veandco/go-sdl2/sdl"
    "GoBoy/display"
)

const (
    WIDTH = 160
    HEIGHT = 144
)

// LCD Control Register bit positions
const (
    BACKGROUND_DISPLAY_FLAG = 0
    OBJ_ON_FLAG = 1
    OBJ_BLOCK_COMPOSITION_SELECTION_FLAG = 2
    BG_CODE_AREA_SELECTION_FLAG = 3
    BG_CHARACTER_DATA_SELECTION_FLAG = 4
    WINDOWING_ON_FLAG = 5
    WINDOW_CODE_AREA_SELECTION_FLAG = 6
    LCD_CONTROLLER_OPERATION_STOP_FLAG = 7
)

const (
    MODE0 = 0x00 // HBlank
    MODE1 = 0x01 // VBlank
    MODE2 = 0x02 // OAM Access
    MODE3 = 0x03 // VRAM Acess

    MATCH_FLAG = 2 // LYC=LY Flag
    MODE0_INTERRUPT = 3
    MODE1_INTERRUPT = 4
    MODE2_INTERRUPT = 5
    MATCH_INTERRUPT = 6 // LYC=LY Interrupt
)

type GPU struct {
    gameboy *Gameboy

    window *display.Window

    VRAM [16384]byte
    OAM []byte

    LCDC byte // LCD Control
    STAT byte // LCD Status/Mode
    SCY byte // Scroll Y
    SCX byte // Scroll X
    LY byte // Scanline
    LYC byte
    DMA byte
    BGP byte
    OBP0 byte
    OBP1 byte
    WY byte // Window Y
    WX byte // Window X

    Cycles uint // Number of cycles since the last LCD Status Mode Change
}

func NewGPU(gameboy *Gameboy) *GPU {
    window := display.NewWindow("Gameboy", WIDTH, HEIGHT, gameboy.Quit)

    gpu := &GPU{
        gameboy: gameboy,
        window: window,
        LCDC: 0x91,
        SCY: 0x00,
        SCX: 0x00,
        LYC: 0x00,
        BGP: 0xFC,
        OBP0: 0xFF,
        OBP1: 0xFF,
        WY: 0x00,
        WX: 0x00,
        STAT: 0x85,
    }
    return gpu
}

func (gpu *GPU) GetSTATMode() byte {
    return gpu.STAT & 0x03
}

func (gpu *GPU) SetSTATMode(mode byte) {
    gpu.STAT = (gpu.STAT & 0xFC) | mode
}

func (gpu *GPU) Step(cycles byte) {
    if gpu.isLCDCEnabled() {
        gpu.Cycles += uint(cycles*4)

        // STAT indicates the current status of the LCD controller.
        switch gpu.GetSTATMode() {
            // HBlank
            // After the last HBlank, push the screen data to canvas
            case MODE0:
                if gpu.Cycles >= 204 {
                    // Reset the cycle counter
                    gpu.Cycles = 0

                    // Increase the scanline
                    gpu.LY += 1

                    // 143 is the last line, update the screen and enter VBlank
                    if gpu.LY == 144 {
                        // Request VBLANK interrupt
                        gpu.gameboy.requestInterrupt(VBLANK_INTERRUPT)
                        gpu.window.Update()

                        // Enter GPU Mode 1/VBlank
                        gpu.SetSTATMode(MODE1)
                        if IsBitSet(gpu.STAT, MODE1_INTERRUPT){
                            gpu.gameboy.requestInterrupt(LCDC_INTERRUPT)
                        }
                    } else {
                        // Enter GPU Mode 2/OAM Access
                        gpu.SetSTATMode(MODE2)
                        if IsBitSet(gpu.STAT, MODE2_INTERRUPT){
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
                    gpu.LY += 1

                    // If Scanline is 153 we have done 10 lines of VBlank
                    if gpu.LY > 153 {
                        // Start of next Frame
                        // Enter GPU Mode 2/OAM Access
                        gpu.SetSTATMode(MODE2)
                        if IsBitSet(gpu.STAT, MODE2_INTERRUPT){
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
                    if IsBitSet(gpu.STAT, MODE0_INTERRUPT){
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

func (gpu *GPU) isLCDCEnabled() bool {
    return IsBitSet(gpu.LCDC, LCD_CONTROLLER_OPERATION_STOP_FLAG)
}


func (gpu *GPU) renderTiles() {
    BGCodeAreaSelection := gpu.getBGCodeAreaSelection()

    // We need the address of the tile in the background map

    // Add ScrollY to the Scanline to get the current pixel Y position
    pixelY := gpu.SCY + gpu.LY

    // Divide the pixel being drawn Y position by 8 (for 8 pixels in tile)
    // and multiply by 32 (for number of tiles in the background map)
    // to get the row number for the tile in the background map
    tileRow := uint16(pixelY/8) * 32

    for pixel := byte(0); pixel < 160; pixel++ {
        // Add pixel being drawn in scanline to scrollX to get the current pixel X position
        pixelX := pixel + gpu.SCX

        // which of the 8 horizontal tiles does this pixel fall within?
        tileCol := uint16(pixelX / 8)

        tileIdentifier := gpu.gameboy.ReadByte(BGCodeAreaSelection + tileRow + tileCol)
        tileDataAddress := gpu.getTileDataAddress(tileIdentifier)

        data1, data2 := gpu.getTileData(tileDataAddress, pixelY)

        // pixel 0 in the tile is it 7 of data 1 and data2.
        // Pixel 1 is bit 6 etc..
        pixelBit := int(pixelX % 8)
        pixelBit -= 7
        pixelBit *= -1

        var colorIdentifier byte
        if IsBitSet(data1, byte(pixelBit)) {
            colorIdentifier = 0x02
        }
        if IsBitSet(data2, byte(pixelBit)) {
            SetBit(colorIdentifier, 0)
        }
        gpu.window.Framebuffer[int(pixel)+(160*int(gpu.LY))] = gpu.getColorFromBGPalette(colorIdentifier)
    }
}

func (gpu *GPU) renderSprites() {
    for sprite := 0; sprite < 40; sprite++ {
        // sprite occupies 4 bytes in the sprite attributes table
        index := sprite*4
        yPos := gpu.gameboy.ReadByte(uint16(0xFE00 + index)) - 16
        xPos := gpu.gameboy.ReadByte(uint16(0xFE00 + index + 1)) - 8
        characterCode := gpu.gameboy.ReadByte(uint16(0xFE00 + index + 2))
        attributes := gpu.gameboy.ReadByte(uint16(0xFE00 + index + 3))

        //yFlip := IsBitSet(attributes, 6)
        xFlip := IsBitSet(attributes, 5)

        if ((gpu.LY >= yPos) && (gpu.LY < (yPos+8))) {

            data1, data2 := gpu.getObjData(characterCode, yPos)

            // its easier to read in from right to left as pixel 0 is
            // bit 7 in the colour data, pixel 1 is bit 6 etc...
            for tilePixel := 7; tilePixel >= 0; tilePixel-- {
                pixelBit := tilePixel
                // read the sprite in backwards for the x axis
                if (xFlip) {
                    pixelBit -= 7
                    pixelBit *= -1
                }

                var colorIdentifier byte
                if IsBitSet(data1, byte(pixelBit)) {
                    colorIdentifier = 0x02
                }
                if IsBitSet(data2, byte(pixelBit)) {
                    SetBit(colorIdentifier, 0)
                }

                gpu.window.Framebuffer[int(xPos)+int(yPos)] = gpu.getSpritePalette(colorIdentifier, attributes)
            }
        }
    }
}

func (gpu *GPU) renderScanline() {
    if (IsBitSet(gpu.LCDC, BACKGROUND_DISPLAY_FLAG)) {
        gpu.renderTiles()
    }

    if (IsBitSet(gpu.LCDC, OBJ_ON_FLAG)) {
        gpu.renderSprites()
    }
}

func (gpu *GPU) getBGCodeAreaSelection() uint16 {
    var BGCodeAreaSelection uint16
    if(IsBitSet(gpu.LCDC, BG_CODE_AREA_SELECTION_FLAG)) {
        BGCodeAreaSelection = 0x9C00
    } else {
        BGCodeAreaSelection = 0x9800
    }

    return BGCodeAreaSelection
}

func (gpu *GPU) getBGCharacterDataSelection() (uint16, uint16) {
    var BGCharacterDataSelection uint16
    var offset uint16
    if(IsBitSet(gpu.LCDC, BG_CHARACTER_DATA_SELECTION_FLAG)) {
        BGCharacterDataSelection = 0x8000
        offset = 0
    } else {
        BGCharacterDataSelection = 0x8800
        offset = 128
    }

    return BGCharacterDataSelection, offset
}

func (gpu *GPU) getWindowCodeAreaSelection() uint16 {
    var WindowCodeAreaSelection uint16
    if(IsBitSet(gpu.LCDC, WINDOW_CODE_AREA_SELECTION_FLAG)) {
        WindowCodeAreaSelection = 0x9C00
    } else {
        WindowCodeAreaSelection = 0x9800
    }

    return WindowCodeAreaSelection
}

func (gpu *GPU) getTileDataAddress(tileIdentifier byte) uint16 {
    // When the BGCharacterDataSelection is 0x8800 the tileIndentifier is
    // a signed byte -127 - 127, the offset corrects for this
    // when looking up the memory location
    BGCharacterDataSelection, offset := gpu.getBGCharacterDataSelection()
    return BGCharacterDataSelection + ((uint16(tileIdentifier) + offset) * 16) // 16 = tile size in bytes
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

func (gpu *GPU) getObjData(characterCode byte, yPos byte) (byte, byte) {
    line := gpu.LY - yPos
    line *= 2; // same as for tiles

    var objCharacterDataStorage uint16 = 0x8000
    objDataAddress := objCharacterDataStorage + (uint16(characterCode) * 16) + uint16(line) // 16 = obj size in bytes

    data1 := gpu.gameboy.ReadByte(objDataAddress)
    data2 := gpu.gameboy.ReadByte(objDataAddress + 1)

    return data1, data2
}

func (gpu *GPU) getColorFromBGPalette(colorIdentifier byte) uint32 {
    pixelFormat, _ := sdl.AllocFormat(uint(sdl.PIXELFORMAT_RGBA32))
    var color byte
    var bitmask byte = 0x3
    switch (colorIdentifier) {
        case 0:
            color = gpu.BGP & bitmask
        case 1:
            color = (gpu.BGP >> 2) & bitmask
        case 2:
            color = (gpu.BGP >> 4) & bitmask
        case 3:
            color = (gpu.BGP >> 6) & bitmask
    }

    switch (color) {
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
    switch (colorIdentifier) {
        case 0:
            color = palette & bitmask
        case 1:
            color = (palette >> 2) & bitmask
        case 2:
            color = (palette >> 4) & bitmask
        case 3:
            color = (palette >> 6) & bitmask
    }

    switch (color) {
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
