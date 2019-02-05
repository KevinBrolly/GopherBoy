package Gameboy

import (
    "os"
    "log"
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
)

type Gameboy struct {
    CPU *CPU
    GPU *GPU

    inBootMode bool
    dmgStatusRegister byte
    bios [256]byte
    ROM []byte
    WorkingRAM [8192]byte //0xC000 -> 0xDFFF (8KB Working RAM)
    ZeroPageRAM [128]byte //0xFF80 - 0xFFFE

    P1 byte
    debug byte
}

func NewGameboy() (gameboy *Gameboy) {
    gameboy = &Gameboy{
        inBootMode: false,
        bios: BIOS,
        ROM: make([]byte, 100000),
        P1: 0xFF,
    }

    gameboy.GPU = NewGPU(gameboy)
    gameboy.CPU = NewCPU(gameboy)
    return
}

func (gameboy *Gameboy) Run() {
    gameboy.GPU.CreateWindow()
    running := true

    for running {
        cycles := gameboy.CPU.Step()
        gameboy.GPU.Step(cycles)

        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch event.(type) {
            case *sdl.QuitEvent:
                println("Quit")
                running = false
                break
            }
        }
    }
}

func (gameboy *Gameboy) LoadROM(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    data := make([]byte, 50000)
    file.Read(data)
    copy(gameboy.ROM, data)
}

func (gameboy *Gameboy) WriteWord(addr uint16, value uint16) {
    hb, lb := SplitBytes(value)
    gameboy.WriteByte(addr, hb)
    gameboy.WriteByte(addr+1, lb)
}

func (gameboy *Gameboy) ReadByte(addr uint16) byte {
    switch {
        //ROM Bank 0
        case addr >= 0x0000 && addr <= 0x3FFF:
            if gameboy.inBootMode && addr < 0x0100 {
                //in bios mode, read from bios
                return gameboy.bios[addr]
            }
            return gameboy.ROM[addr]
        //ROM Bank 1
        case addr >= 0x4000 && addr <= 0x7FFF:
            return gameboy.ROM[addr]
        //Working RAM
        case addr >= 0xC000 && addr <= 0xDFFF:
            return gameboy.WorkingRAM[addr & 0x1FFF]

        case addr >= 0xFF80 && addr <= 0xFFFE:
            return gameboy.ZeroPageRAM[addr & 0x7F]

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

        case addr == P1:
            return gameboy.P1 & 0x0F

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
        case addr == WY:
            return gameboy.GPU.WY
        case addr == WX:
            return gameboy.GPU.WX
        case addr >= 0x8000 && addr <= 0x9FFF:
            return gameboy.GPU.VRAM[addr & 0x1FFF]

        case addr == 0xFF01:
            return gameboy.debug
    }

    return 0
}

func (gameboy *Gameboy) WriteByte(addr uint16, value byte) {
    switch {
        //ROM Bank 0
        case addr >= 0x0000 && addr <= 0x7FFF:
            gameboy.ROM[addr] = value
        //Working RAM
        case addr >= 0xC000 && addr <= 0xDFFF:
            gameboy.WorkingRAM[addr & 0x1FFF] = value

        case addr >= 0xFF80 && addr <= 0xFFFE:
            gameboy.ZeroPageRAM[addr & 0x7F] = value

        // Registers
        case addr == DMG_STATUS_REGISTER:
            gameboy.dmgStatusRegister = value

        // Timer
        case addr == DIV:
            gameboy.CPU.DIV = 0
        case addr == TIMA: // Timer Counter
            gameboy.CPU.TIMA = value
        case addr == TMA: // Timer Modulo
            gameboy.CPU.TMA = value
        case addr == TAC:
            gameboy.CPU.TAC = value

        // I/O control handling
        case addr == IF:
            gameboy.CPU.IF = value
        case addr == IE:
            gameboy.CPU.IE = value

        case addr == P1:
            gameboy.P1 = value & 0x30

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
        case addr == BGP:
            gameboy.GPU.BGP = value
        case addr == WY:
            gameboy.GPU.WY = value
        case addr == WX:
            gameboy.GPU.WX = value
        case addr >= 0x8000 && addr <= 0x9FFF:
            gameboy.GPU.VRAM[addr & 0x1FFF] = value

        case addr == 0xFF01:
            gameboy.debug = value
        case addr == 0xFF02 && value == 0x81:
            fmt.Print(string(gameboy.ReadByte(0xFF01)))
    }
}

func (gameboy *Gameboy) requestInterrupt(interrupt byte) {
    gameboy.WriteByte(IF, SetBit(gameboy.CPU.IF, interrupt))
}
