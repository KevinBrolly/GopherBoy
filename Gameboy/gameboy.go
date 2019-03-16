package Gameboy

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Gameboy struct {
	CPU        *CPU
	GPU        *GPU
	Cartridge  *Cartridge
	Controller *Controller

	inBootMode        bool
	dmgStatusRegister byte
	bios              [256]byte
	WorkingRAM        [8192]byte //0xC000 -> 0xDFFF (8KB Working RAM)
	ZeroPageRAM       [128]byte  //0xFF80 - 0xFFFE

	P1      byte
	debug   byte
	running bool
}

func NewGameboy() (gameboy *Gameboy) {
	gameboy = &Gameboy{
		inBootMode: false,
		bios:       BIOS,
		P1:         0xFF,
	}

	gameboy.GPU = NewGPU(gameboy)
	gameboy.CPU = NewCPU(gameboy)
	gameboy.Controller = NewController(gameboy)
	gameboy.Cartridge = NewCartridge(gameboy)
	return
}

func (gameboy *Gameboy) Run() {
	timePerFrame := time.Second / 60

	ticker := time.NewTicker(timePerFrame)

	for range ticker.C {
		MAXCYCLES := 17476
		cyclesThisUpdate := 0

		for cyclesThisUpdate < MAXCYCLES {
			cycles := gameboy.CPU.Step()
			gameboy.GPU.Step(cycles)
			cyclesThisUpdate += int(cycles)
		}

		// Check for events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				gameboy.GPU.Window.Quit()
				gameboy.Quit()

			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_RIGHT:
						gameboy.Controller.KeyPressed(RIGHT)
					case sdl.K_LEFT:
						gameboy.Controller.KeyPressed(LEFT)
					case sdl.K_UP:
						gameboy.Controller.KeyPressed(UP)
					case sdl.K_DOWN:
						gameboy.Controller.KeyPressed(DOWN)
					case sdl.K_a:
						gameboy.Controller.KeyPressed(A)
					case sdl.K_b:
						gameboy.Controller.KeyPressed(B)
					case sdl.K_RETURN:
						gameboy.Controller.KeyPressed(SELECT)
					case sdl.K_SPACE:
						gameboy.Controller.KeyPressed(START)

					}
				} else if e.Type == sdl.KEYUP {
					switch e.Keysym.Sym {
					case sdl.K_RIGHT:
						gameboy.Controller.KeyReleased(RIGHT)
					case sdl.K_LEFT:
						gameboy.Controller.KeyReleased(LEFT)
					case sdl.K_UP:
						gameboy.Controller.KeyReleased(UP)
					case sdl.K_DOWN:
						gameboy.Controller.KeyReleased(DOWN)
					case sdl.K_a:
						gameboy.Controller.KeyReleased(A)
					case sdl.K_b:
						gameboy.Controller.KeyReleased(B)
					case sdl.K_RETURN:
						gameboy.Controller.KeyReleased(SELECT)
					case sdl.K_SPACE:
						gameboy.Controller.KeyReleased(START)

					}
				}
			}
		}

		gameboy.GPU.Window.Update()
	}
}

func (gameboy *Gameboy) Quit() {
	gameboy.running = false
}

func (gameboy *Gameboy) WriteWord(addr uint16, value uint16) {
	hb, lb := SplitBytes(value)
	gameboy.WriteByte(addr, hb)
	gameboy.WriteByte(addr+1, lb)
}

func (gameboy *Gameboy) ReadByte(addr uint16) byte {

	// Check the cartridge
	if value := gameboy.Cartridge.ReadByte(addr); value != 0 {
		return value
	}

	// Check the controls
	if value := gameboy.Controller.ReadByte(addr); value != 0 {
		return value
	}

	switch {
	//Working RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		return gameboy.WorkingRAM[addr&0x1FFF]

	case addr >= 0xFF80 && addr <= 0xFFFE:
		return gameboy.ZeroPageRAM[addr&0x7F]

	// Registers
	case addr == DMG_STATUS_REGISTER:
		return gameboy.dmgStatusRegister

	// Timer
	case addr == DIV:
		return gameboy.CPU.DIV
	case addr == TIMA:
		return gameboy.CPU.TIMA
	case addr == TMA:
		return gameboy.CPU.TMA
	case addr == TAC:
		return gameboy.CPU.TAC

	// I/O control handling
	case addr == IF:
		return gameboy.CPU.IF
	case addr == IE:
		return gameboy.CPU.IE

	case addr == LCDC:
		return gameboy.GPU.LCDC
	case addr == STAT:
		return gameboy.GPU.STAT
	case addr == SCY:
		return gameboy.GPU.SCY
	case addr == SCX:
		return gameboy.GPU.SCX
	case addr == LY:
		return gameboy.GPU.LY
	case addr == BGP:
		return gameboy.GPU.BGP
	case addr == OBP0:
		return gameboy.GPU.OBP0
	case addr == OBP1:
		return gameboy.GPU.OBP1
	case addr == WY:
		return gameboy.GPU.WY
	case addr == WX:
		return gameboy.GPU.WX
	case addr >= 0x8000 && addr <= 0x9FFF:
		return gameboy.GPU.VRAM[addr&0x1FFF]
	case addr >= 0xFE00 && addr <= 0xFE9F:
		return gameboy.GPU.OAM[addr&0x9F]

	case addr == 0xFF01:
		return gameboy.debug
	}

	return 0
}

func (gameboy *Gameboy) WriteByte(addr uint16, value byte) {

	// Check the cartridge
	gameboy.Cartridge.WriteByte(addr, value)

	// Check the controls
	gameboy.Controller.WriteByte(addr, value)

	switch {
	//Working RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		gameboy.WorkingRAM[addr&0x1FFF] = value

	case addr >= 0xFF80 && addr <= 0xFFFE:
		gameboy.ZeroPageRAM[addr&0x7F] = value

	// Registers
	case addr == DMG_STATUS_REGISTER:
		gameboy.dmgStatusRegister = value

	// Timer
	case addr == DIV: // Divider
		gameboy.CPU.DIV = 0
	case addr == TIMA: // Timer Counter
		gameboy.CPU.TIMA = value
	case addr == TMA: // Timer Modulo
		gameboy.CPU.TMA = value
	case addr == TAC:
		currentfreq := gameboy.CPU.getClockFrequency()
		gameboy.CPU.TAC = value

		newfreq := gameboy.CPU.getClockFrequency()

		if currentfreq != newfreq {
			switch newfreq {
			case 0:
				gameboy.CPU.TimerCounter = 1024 // frequency 4096
			case 1:
				gameboy.CPU.TimerCounter = 16 // frequency 262144
			case 2:
				gameboy.CPU.TimerCounter = 64 // frequency 65536
			case 3:
				gameboy.CPU.TimerCounter = 256 // frequency 16382
			}
		}

	// I/O control handling
	case addr == IF:
		if gameboy.CPU.Halt {
			if value != gameboy.CPU.IF {
				gameboy.CPU.Halt = false
			}
		}
		gameboy.CPU.IF = value
	case addr == IE:
		gameboy.CPU.IE = value

	case addr == LCDC:
		gameboy.GPU.LCDC = value
	case addr == STAT:
		gameboy.GPU.STAT = (0xF8 & value) | gameboy.GPU.GetSTATMode()
	case addr == SCY:
		gameboy.GPU.SCY = value
	case addr == SCX:
		gameboy.GPU.SCX = value
	case addr == LY:
		// If the game writes to scanline it should be unset
		gameboy.GPU.LY = 0
	case addr == LYC:
		gameboy.GPU.LYC = value
	case addr == DMA:
		// The value holds the source address of the OAM data divided by 100
		// so we have to multiply it first
		var sourceAddr uint16 = uint16(value) << 8

		for i := 0; i < 160; i++ {
			gameboy.GPU.OAM[i] = gameboy.ReadByte(sourceAddr + uint16(i))
		}
	case addr == BGP:
		gameboy.GPU.BGP = value
	case addr == OBP0:
		gameboy.GPU.OBP0 = value
	case addr == OBP1:
		gameboy.GPU.OBP1 = value
	case addr == WY:
		gameboy.GPU.WY = value
	case addr == WX:
		gameboy.GPU.WX = value
	case addr >= 0x8000 && addr <= 0x9FFF:
		gameboy.GPU.VRAM[addr&0x1FFF] = value

	case addr == 0xFF01:
		gameboy.debug = value
	case addr == 0xFF02 && value == 0x81:
		fmt.Print(string(gameboy.ReadByte(0xFF01)))
	}
}

func (gameboy *Gameboy) requestInterrupt(interrupt byte) {
	gameboy.WriteByte(IF, SetBit(gameboy.CPU.IF, interrupt))
}
