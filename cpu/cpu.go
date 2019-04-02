package cpu

import (
	"log"
	"os"

	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
)

// Interrupts
const (
	IF = 0xFF0F
	IE = 0xFFFF

	VBLANK_INTERRUPT              = 0
	LCDC_INTERRUPT                = 1
	TIMER_OVERFLOW_INTERRUPT      = 2
	SERIAL_IO_INTERRUPT           = 3
	JOYPAD_INTERRUPT              = 4
	VBLANK_INTERRUPT_ADDR         = 0x40
	LCDC_INTERRUPT_ADDR           = 0x48
	TIMER_OVERFLOW_INTERRUPT_ADDR = 0x50
	SERIAL_IO_INTERRUPT_ADDR      = 0x58
	JOYPAD_INTERRUPT_ADDR         = 0x60
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
	mmu   *mmu.MMU
	timer *Timer

	Registers          Registers
	SP                 uint16
	PC                 uint16
	CurrentInstruction *Instruction

	IF  byte // Interrupt Flag
	IE  byte // Interrupt Enabled
	IME bool // Interrupt Master Enable

	Halt bool
}

func NewCPU(mmu *mmu.MMU) *CPU {
	cpu := &CPU{
		mmu:   mmu,
		timer: NewTimer(mmu),
	}

	// FF0F - IF - Interrupt Flag
	mmu.MapMemory(cpu, IF)

	// FFFF - IE - Interrupt Enable
	mmu.MapMemory(cpu, IE)

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

	cpu.IE = 0x00
	cpu.IF = 0xE1
}

func (cpu *CPU) GetOpcode() byte {
	return cpu.mmu.ReadByte(cpu.PC)
}

func (cpu *CPU) GetByteOffset(offset uint16) byte {
	return cpu.mmu.ReadByte(cpu.PC + offset)
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

		cycles = instruction.Execute(cpu)

		if initialPC == cpu.PC {
			cpu.PC += cpu.CurrentInstruction.Length
		}
	} else {
		// Halt takes 1 cycle
		cycles = 1
	}

	cpu.timer.Tick(cycles)
	cpu.handleInterrupts()

	return cycles
}

func (cpu *CPU) handleInterrupts() {
	if cpu.IME {
		// if and interrupt is requested (IF) & enabled (IE) then set its bit to 1 in interrupts
		interrupts := cpu.IE & cpu.IF

		// If an interrupt is actionable interrupts will be > 0
		if interrupts != 0 {
			switch {
			case utils.IsBitSet(interrupts, VBLANK_INTERRUPT):
				cpu.handleInterrupt(VBLANK_INTERRUPT, VBLANK_INTERRUPT_ADDR)
			case utils.IsBitSet(interrupts, LCDC_INTERRUPT):
				cpu.handleInterrupt(LCDC_INTERRUPT, LCDC_INTERRUPT_ADDR)
			case utils.IsBitSet(interrupts, TIMER_OVERFLOW_INTERRUPT):
				cpu.handleInterrupt(TIMER_OVERFLOW_INTERRUPT, TIMER_OVERFLOW_INTERRUPT_ADDR)
			case utils.IsBitSet(interrupts, SERIAL_IO_INTERRUPT):
				cpu.handleInterrupt(SERIAL_IO_INTERRUPT, SERIAL_IO_INTERRUPT_ADDR)
			case utils.IsBitSet(interrupts, JOYPAD_INTERRUPT):
				cpu.handleInterrupt(JOYPAD_INTERRUPT, JOYPAD_INTERRUPT_ADDR)
			}
		}
	}
}

func (cpu *CPU) handleInterrupt(interrupt byte, interrupt_addr uint16) {
	cpu.IME = false
	cpu.mmu.WriteByte(IF, utils.ClearBit(cpu.IF, interrupt))
	cpu.pushWordToStack(cpu.PC)
	cpu.PC = interrupt_addr
}

// FLAGS
func (cpu *CPU) SetFlag(flag byte) {
	cpu.Registers.F = utils.SetBit(cpu.Registers.F, flag)
}

func (cpu *CPU) IsFlagSet(flag byte) bool {
	return utils.IsBitSet(cpu.Registers.F, flag)
}

func (cpu *CPU) ResetFlag(flag byte) {
	cpu.Registers.F = utils.ClearBit(cpu.Registers.F, flag)
}

func (cpu *CPU) ReadByte(addr uint16) byte {
	switch {
	// I/O control handling
	case addr == IF:
		return cpu.IF
	case addr == IE:
		return cpu.IE
	}
	return 0
}

func (cpu *CPU) WriteByte(addr uint16, value byte) {
	switch {
	// I/O control handling
	case addr == IF:
		if cpu.Halt {
			if value != cpu.IF {
				cpu.Halt = false
			}
		}
		cpu.IF = value
	case addr == IE:
		cpu.IE = value
	}
}
