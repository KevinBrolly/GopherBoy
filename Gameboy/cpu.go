package Gameboy

import (
    "os"
    "log"
)

type Registers struct {
    A byte
    B byte
    C byte
    D byte
    E byte
    H byte
    L byte
    F byte
}

type CPU struct {
    gameboy *Gameboy
    Registers Registers
    SP uint16
    PC uint16
    CurrentInstruction *Instruction

    IF byte  // Interrupt Flag
    IE byte  // Interrupt Enabled
    IME bool // Interrupt Master Enable

    DIV byte  // Divider
    TIMA byte // Timer Counter
    TMA byte  // Timer Modulo
    TAC byte  // Timer Controller

    DividerCounter int
    TimerCounter int

    Halt bool
}

func NewCPU(gameboy *Gameboy) *CPU {
    cpu := &CPU{gameboy: gameboy}
    cpu.Reset()
    return cpu
}

func (cpu *CPU) Reset() {
    cpu.PC = 0x100

    cpu.Registers.A = 0x01
    cpu.Registers.B = 0x00
    cpu.Registers.C = 0x13
    cpu.Registers.D = 0x00
    cpu.Registers.E = 0xD8
    cpu.Registers.F = 0xB0
    cpu.Registers.H = 0x01
    cpu.Registers.L = 0x4D

    cpu.SP = 0xFFFE

    cpu.TIMA = 0x00
    cpu.TMA = 0x00
    cpu.TAC = 0x05

    cpu.IE = 0x00
    cpu.IF = 0xE1

    cpu.TimerCounter = 1024
}

func (cpu *CPU) GetOpcode() byte {
    return cpu.gameboy.ReadByte(cpu.PC)
}

func (cpu *CPU) GetByteOffset(offset uint16) byte {
    return cpu.gameboy.ReadByte(cpu.PC+offset)
}

func (cpu *CPU) getInstruction(opcode byte) *Instruction {
    var instruction *Instruction
    var ok bool

    if opcode == 0xCB {
        opcode = cpu.GetByteOffset(1)

        instruction, ok = CBInstructions[opcode]
    } else {
        instruction, ok = Instructions[opcode]
    }

    if !ok {
        log.Panicf("No instruction found for opcode: %x\n", opcode)
    }

    return instruction
}

func UnimplementedInstruction(cpu *CPU) {
    //fmt.Printf("Error: Unimplemented instruction - %#x\n", cpu.GetOpcode())
    os.Exit(1)
}

func (cpu *CPU) Step() (cycles byte) {
    if !cpu.Halt {
        initialPC := cpu.PC

        var opcode byte = cpu.GetOpcode()

        instruction := cpu.getInstruction(opcode)
        cpu.CurrentInstruction = instruction

        //fmt.Printf("OPCODE: %#x, Desc: %v, LY: %#x, PC: %#x, SP: %#x, IME: %v, IE: %#x, IF: %#x, LCDC: %#x, STAT: %#x, AF: %#x, BC: %#x, DE: %#x, HL: %#x\n", cpu.GetOpcode(), cpu.CurrentInstruction.Description, cpu.gameboy.GPU.LY, cpu.PC, cpu.SP, cpu.IME, cpu.IE, cpu.IF, cpu.gameboy.GPU.LCDC, cpu.gameboy.GPU.STAT, JoinBytes(cpu.Registers.A, cpu.Registers.F), JoinBytes(cpu.Registers.B, cpu.Registers.C), JoinBytes(cpu.Registers.D, cpu.Registers.E), JoinBytes(cpu.Registers.H, cpu.Registers.L))

        cycles = instruction.Execute(cpu)

        if initialPC == cpu.PC {
            cpu.PC += cpu.CurrentInstruction.Length
        }
    } else {
        // Halt takes 1 cycle
        cycles = 1
    }

    cpu.updateTimer(cycles)
    cpu.handleInterrupts()

    return cycles
}

func (cpu *CPU) updateTimer(cycles byte) {
    cycles = cycles*4

    cpu.updateDividerRegister(cycles)

    if IsBitSet(cpu.TAC, TIMER_STOP){

        cpu.TimerCounter -= int(cycles)

        if cpu.TimerCounter <= 0 {

            frequency := cpu.getClockFrequency()

            switch (frequency) {
                case 0:
                    cpu.TimerCounter = 1024  // frequency 4096
                case 1:
                    cpu.TimerCounter = 16    // frequency 262144
                case 2:
                    cpu.TimerCounter = 64    // frequency 65536
                case 3:
                    cpu.TimerCounter = 256   // frequency 16382
            }

            // If timer is about to overflow
            if cpu.TIMA == 255 {
                cpu.TIMA = cpu.TMA
                cpu.gameboy.requestInterrupt(TIMER_OVERFLOW_INTERRUPT)
            } else {
                cpu.TIMA = cpu.TIMA + 1
            }
        }
    }
}

func (cpu *CPU) updateDividerRegister(cycles byte) {
    cpu.DividerCounter += int(cycles)

    if cpu.DividerCounter >= 255 {
        cpu.DividerCounter = 0
        cpu.DIV += 1
    }
}

func (cpu *CPU) getClockFrequency() byte {
    return cpu.TAC & 0x03
}

func (cpu *CPU) handleInterrupts() {
    if cpu.IME {
        // if and interrupt is requested (IF) & enabled (IE) then set its bit to 1 in interrupts
        interrupts := cpu.IE & cpu.IF
        // If an interrupt is actionable interrupts will be > 0
        if interrupts != 0 {
            switch {
            case IsBitSet(interrupts, VBLANK_INTERRUPT):
                cpu.handleInterrupt(VBLANK_INTERRUPT, VBLANK_INTERRUPT_ADDR)
            case IsBitSet(interrupts, LCDC_INTERRUPT):
                cpu.handleInterrupt(LCDC_INTERRUPT, LCDC_INTERRUPT_ADDR)
            case IsBitSet(interrupts, TIMER_OVERFLOW_INTERRUPT):
                cpu.handleInterrupt(TIMER_OVERFLOW_INTERRUPT, TIMER_OVERFLOW_INTERRUPT_ADDR)
            case IsBitSet(interrupts, SERIAL_IO_INTERRUPT):
                cpu.handleInterrupt(SERIAL_IO_INTERRUPT, SERIAL_IO_INTERRUPT_ADDR)
            case IsBitSet(interrupts, JOYPAD_INTERRUPT):
                cpu.handleInterrupt(JOYPAD_INTERRUPT, JOYPAD_INTERRUPT_ADDR)
            }
        }
    }
}

func (cpu *CPU) handleInterrupt(interrupt byte, interrupt_addr uint16) {
    cpu.IME = false
    cpu.gameboy.WriteByte(IF, ClearBit(cpu.IF, interrupt))
    cpu.pushWordToStack(cpu.PC)
    cpu.PC = interrupt_addr
}

// FLAGS
func (cpu *CPU) SetFlag(flag byte) {
    cpu.Registers.F = SetBit(cpu.Registers.F, flag)
}

func (cpu *CPU) IsFlagSet(flag byte) bool {
    return IsBitSet(cpu.Registers.F, flag)
}

func (cpu *CPU) ResetFlag(flag byte) {
    cpu.Registers.F = ClearBit(cpu.Registers.F, flag)
}

