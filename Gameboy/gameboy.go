package Gameboy

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Gameboy struct {
	CPU        *CPU
	GPU        *GPU
	Timer      *Timer
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
	gameboy.Timer = NewTimer(gameboy)
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

	// Check the Timer
	if value := gameboy.Timer.ReadByte(addr); value != 0 {
		return value
	}

	// Check the GPU
	if value := gameboy.GPU.ReadByte(addr); value != 0 {
		return value
	}

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

	// I/O control handling
	case addr == IF:
		return gameboy.CPU.IF
	case addr == IE:
		return gameboy.CPU.IE

	case addr == 0xFF01:
		return gameboy.debug
	}

	return 0
}

func (gameboy *Gameboy) WriteByte(addr uint16, value byte) {

	// Check the timer
	gameboy.Timer.WriteByte(addr, value)

	// Check the GPU
	gameboy.GPU.WriteByte(addr, value)

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

	case addr == 0xFF01:
		gameboy.debug = value
	case addr == 0xFF02 && value == 0x81:
		fmt.Print(string(gameboy.ReadByte(0xFF01)))
	}
}

func (gameboy *Gameboy) requestInterrupt(interrupt byte) {
	gameboy.WriteByte(IF, SetBit(gameboy.CPU.IF, interrupt))
}
