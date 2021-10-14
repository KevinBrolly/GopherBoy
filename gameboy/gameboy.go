package gameboy

import (
	"fmt"
	"image"
	"time"

	"github.com/kevinbrolly/GopherBoy/apu"
	"github.com/kevinbrolly/GopherBoy/cartridge"
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
	CPU        *cpu.CPU
	PPU        *ppu.PPU
	APU        *apu.APU
	Controller *control.Controller
	Cartridge  *cartridge.Cartridge

	inBootMode        bool
	dmgStatusRegister byte
	WorkingRAM        [8192]byte //0xC000 -> 0xDFFF (8KB Working RAM)
	HRAM              [128]byte  //0xFF80 -> 0xFFFE High RAM (HRAM)

	debug   byte
	running bool

	cycleChannel chan int
}

func NewGameboy(window Window) (gameboy *Gameboy) {
	cycleChannel := make(chan int)

	mmu := mmu.NewMMU(cycleChannel)
	cpu := cpu.NewCPU(mmu, cycleChannel)
	ppu := ppu.NewPPU(mmu, cycleChannel)
	apu := apu.NewAPU(mmu, cycleChannel)
	controller := control.NewController(mmu)

	gameboy = &Gameboy{
		Window:       window,
		MMU:          mmu,
		CPU:          cpu,
		PPU:          ppu,
		APU:          apu,
		Controller:   controller,
		cycleChannel: cycleChannel,
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

func (gameboy *Gameboy) LoadCartridge(filename string) {
	gameboy.Cartridge = cartridge.NewCartridge(filename, gameboy.MMU)
}

func (gameboy *Gameboy) Run() {

	frameTime := time.Second / 60

	ticker := time.NewTicker(frameTime)
	start := time.Now()
	frames := 0

	// 4.194304MHz == 4194304hz / 60 = ~69905 cycles in one second
	cyclesThisUpdate := 0
	MAXCYCLES := 69905

	for {
		select {
		case cycle := <-gameboy.cycleChannel:
			fmt.Printf("%v/n", cyclesThisUpdate)
			cyclesThisUpdate += cycle
			if cyclesThisUpdate == MAXCYCLES {
				gameboy.CPU.IsRunning = false
			}
		case <-ticker.C:
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
						case sdl.K_z:
							gameboy.Controller.KeyPressed(control.DEBUG)
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
						case sdl.K_z:
							gameboy.Controller.KeyPressed(control.DEBUG)
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

			cyclesThisUpdate = 0
			gameboy.CPU.IsRunning = true
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
