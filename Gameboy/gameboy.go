package gameboy

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"time"

	"github.com/kevinbrolly/GopherBoy/apu"
	"github.com/kevinbrolly/GopherBoy/control"
	"github.com/kevinbrolly/GopherBoy/cpu"
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/ppu"

	"github.com/veandco/go-sdl2/sdl"
)

// Registers
const (
	DMG_STATUS_REGISTER = 0xFF50 // Signals that the boot ROM has finished
)

type Window interface {
	DrawFrame(frameBuffer *image.RGBA)
}

type Gameboy struct {
	Window     Window
	MMU        *mmu.MMU
	MBC        mmu.Memory
	CPU        *cpu.CPU
	PPU        *ppu.PPU
	APU        *apu.APU
	Controller *control.Controller

	inBootMode        bool
	dmgStatusRegister byte
	WorkingRAM        [8192]byte //0xC000 -> 0xDFFF (8KB Working RAM)
	HRAM              [128]byte  //0xFF80 -> 0xFFFE High RAM (HRAM)

	debug   byte
	running bool
}

func NewGameboy(window Window) (gameboy *Gameboy) {
	mmu := mmu.NewMMU()
	cpu := cpu.NewCPU(mmu)
	ppu := ppu.NewPPU(mmu)
	apu := apu.NewAPU(mmu)
	controller := control.NewController(mmu)

	gameboy = &Gameboy{
		Window:     window,
		MMU:        mmu,
		CPU:        cpu,
		PPU:        ppu,
		APU:        apu,
		Controller: controller,
	}

	// Map memory for outputting result of blargg tests
	mmu.MapMemory(gameboy, 0xFF01)
	mmu.MapMemory(gameboy, 0xFF02)

	// Working RAM
	mmu.MapMemoryRange(gameboy, 0xC000, 0xDFFF)
	// HRAM
	mmu.MapMemoryRange(gameboy, 0xFF80, 0xFFFE)

	return gameboy
}

func (gameboy *Gameboy) LoadCartridgeData(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}
	switch data[0x147] {
	case 0:
		gameboy.MBC = mmu.NewMBC0(gameboy.MMU, data)
	case 1:
		gameboy.MBC = mmu.NewMBC1(gameboy.MMU, data)
	case 2:
		gameboy.MBC = mmu.NewMBC1(gameboy.MMU, data)
	case 3:
		gameboy.MBC = mmu.NewMBC1(gameboy.MMU, data)
	case 4:
		gameboy.MBC = mmu.NewMBC1(gameboy.MMU, data)
	}
}

func (gameboy *Gameboy) Run() {
	frameTime := time.Second / 60

	ticker := time.NewTicker(frameTime)
	start := time.Now()
	frames := 0
	for range ticker.C {
		MAXCYCLES := 69905
		cyclesThisUpdate := 0

		for cyclesThisUpdate < MAXCYCLES {
			cycles := gameboy.CPU.Step()
			gameboy.PPU.Step(cycles)
			gameboy.APU.Tick(cycles)
			cyclesThisUpdate += cycles

			// fmt.Printf("OPCODE: %#x, Desc: %v, LY: %#x, PC: %#x, SP: %#x, IME: %v, IE: %#x, IF: %#x, LCDC: %#x, AF: %#x, BC: %#x, DE: %#x, HL: %#x\n",
			// 	gameboy.CPU.GetOpcode(),
			// 	gameboy.CPU.CurrentInstruction.Description,
			// 	gameboy.PPU.LY,
			// 	gameboy.CPU.PC,
			// 	gameboy.CPU.SP,
			// 	gameboy.CPU.IME,
			// 	gameboy.CPU.IE,
			// 	gameboy.CPU.IF,
			// 	gameboy.PPU.LCDC,
			// 	utils.JoinBytes(gameboy.CPU.Registers.A, gameboy.CPU.Registers.F),
			// 	utils.JoinBytes(gameboy.CPU.Registers.B, gameboy.CPU.Registers.C),
			// 	utils.JoinBytes(gameboy.CPU.Registers.D, gameboy.CPU.Registers.E),
			// 	utils.JoinBytes(gameboy.CPU.Registers.H, gameboy.CPU.Registers.L),
			// )
		}

		// Check for events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				gameboy.Quit()

			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_RIGHT:
						gameboy.Controller.KeyPressed(control.RIGHT)
					case sdl.K_LEFT:
						gameboy.Controller.KeyPressed(control.LEFT)
					case sdl.K_UP:
						gameboy.Controller.KeyPressed(control.UP)
					case sdl.K_DOWN:
						gameboy.Controller.KeyPressed(control.DOWN)
					case sdl.K_s:
						gameboy.Controller.KeyPressed(control.A)
					case sdl.K_a:
						gameboy.Controller.KeyPressed(control.B)
					case sdl.K_SPACE:
						gameboy.Controller.KeyPressed(control.SELECT)
					case sdl.K_RETURN:
						gameboy.Controller.KeyPressed(control.START)
					}
				} else if e.Type == sdl.KEYUP {
					switch e.Keysym.Sym {
					case sdl.K_RIGHT:
						gameboy.Controller.KeyReleased(control.RIGHT)
					case sdl.K_LEFT:
						gameboy.Controller.KeyReleased(control.LEFT)
					case sdl.K_UP:
						gameboy.Controller.KeyReleased(control.UP)
					case sdl.K_DOWN:
						gameboy.Controller.KeyReleased(control.DOWN)
					case sdl.K_s:
						gameboy.Controller.KeyReleased(control.A)
					case sdl.K_a:
						gameboy.Controller.KeyReleased(control.B)
					case sdl.K_SPACE:
						gameboy.Controller.KeyReleased(control.SELECT)
					case sdl.K_RETURN:
						gameboy.Controller.KeyReleased(control.START)

					}
				}
			}
		}

		gameboy.Window.DrawFrame(gameboy.PPU.FrameBuffer)
		frames++
		since := time.Since(start)
		if since > time.Second {
			start = time.Now()
			frames = 0
		}
	}
}

func (gameboy *Gameboy) Quit() {
	gameboy.running = false
}

func (gameboy *Gameboy) ReadByte(addr uint16) byte {
	switch {
	// Working RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		return gameboy.WorkingRAM[addr&0x1FFF]

	// HRAM
	case addr >= 0xFF80 && addr <= 0xFFFE:
		return gameboy.HRAM[addr&0x7F]

	// Registers
	case addr == DMG_STATUS_REGISTER:
		return gameboy.dmgStatusRegister

	case addr == 0xFF01:
		return gameboy.debug
	}

	return 0
}

func (gameboy *Gameboy) WriteByte(addr uint16, value byte) {
	switch {
	//Working RAM
	case addr >= 0xC000 && addr <= 0xDFFF:
		gameboy.WorkingRAM[addr&0x1FFF] = value

	case addr >= 0xFF80 && addr <= 0xFFFE:
		gameboy.HRAM[addr&0x7F] = value

	// Registers
	case addr == DMG_STATUS_REGISTER:
		gameboy.dmgStatusRegister = value

	case addr == 0xFF01:
		gameboy.debug = value
	case addr == 0xFF02 && value == 0x81:
		fmt.Print(string(gameboy.ReadByte(0xFF01)))
	}
}
