package cpu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type Instruction struct {
	Opcode      byte
	Description string
	Length      uint16
	Execute     func(cpu *CPU) byte
}

//Condition Codes
const (
	CC_NZ = 00
	CC_Z  = 01
	CC_NC = 10
	CC_C  = 11
)

//flags
const (
	Z  = 7 // Zero Flag
	N  = 6 // Subtraction Flag
	H  = 5 // Half-Carry Flag
	CY = 4 // Carry Flag
)

func (i *Instruction) toString() string {
	return i.Description
}

var Instructions map[byte]*Instruction = map[byte]*Instruction{
	0x00: &Instruction{0x00, "NOP", 1, func(cpu *CPU) byte {
		return cpu.NOP()
	}},
	0x01: &Instruction{0x01, "LD BC,d16", 3, func(cpu *CPU) byte {
		return cpu.LD_rr_nn(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x02: &Instruction{0x02, "LD (BC), A", 1, func(cpu *CPU) byte {
		return cpu.LD_BC_A()
	}},
	0x03: &Instruction{0x03, "INC BC", 1, func(cpu *CPU) byte {
		return cpu.INC_rr(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x04: &Instruction{0x04, "INC B", 1, func(cpu *CPU) byte {
		return cpu.INC_r(&cpu.Registers.B)
	}},
	0x05: &Instruction{0x05, "DEC B", 1, func(cpu *CPU) byte {
		return cpu.DEC_r(&cpu.Registers.B)
	}},
	0x06: &Instruction{0x06, "LD B,d8", 2, func(cpu *CPU) byte {
		return cpu.LD_r_n(&cpu.Registers.B)
	}},
	0x07: &Instruction{0x07, "RLCA", 1, func(cpu *CPU) byte {
		return cpu.RLCA()
	}},
	0x08: &Instruction{0x08, "LD (a16),SP", 3, func(cpu *CPU) byte {
		return cpu.LD_nn_SP()
	}},
	0x09: &Instruction{0x09, "ADD HL, BC", 1, func(cpu *CPU) byte {
		return cpu.ADD_HL_rr(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x0A: &Instruction{0x0A, "LD A,(BC)", 1, func(cpu *CPU) byte {
		return cpu.LD_A_BC()
	}},
	0x0B: &Instruction{0x0B, "DEC BC", 1, func(cpu *CPU) byte {
		return cpu.DEC_rr(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x0C: &Instruction{0x0C, "INC C", 1, func(cpu *CPU) byte {
		return cpu.INC_r(&cpu.Registers.C)
	}},
	0x0D: &Instruction{0x0D, "DEC C", 1, func(cpu *CPU) byte {
		return cpu.DEC_r(&cpu.Registers.C)
	}},
	0x0E: &Instruction{0x0E, "LD C,d8", 2, func(cpu *CPU) byte {
		return cpu.LD_r_n(&cpu.Registers.C)
	}},
	0x0F: &Instruction{0x0F, "RRCA", 1, func(cpu *CPU) byte {
		return cpu.RRCA()
	}},
	0x10: &Instruction{0x10, "STOP 0", 2, func(cpu *CPU) byte {
		return 4
	}},
	0x11: &Instruction{0x11, "LD DE,d16", 3, func(cpu *CPU) byte {
		return cpu.LD_rr_nn(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x12: &Instruction{0x12, "LD (DE), A", 1, func(cpu *CPU) byte {
		return cpu.LD_DE_A()
	}},
	0x13: &Instruction{0x13, "INC DE", 1, func(cpu *CPU) byte {
		return cpu.INC_rr(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x14: &Instruction{0x14, "INC D", 1, func(cpu *CPU) byte {
		return cpu.INC_r(&cpu.Registers.D)
	}},
	0x15: &Instruction{0x15, "DEC D", 1, func(cpu *CPU) byte {
		return cpu.DEC_r(&cpu.Registers.D)
	}},
	0x16: &Instruction{0x16, "LD D,d8", 2, func(cpu *CPU) byte {
		return cpu.LD_r_n(&cpu.Registers.D)
	}},
	0x17: &Instruction{0x17, "RLA", 1, func(cpu *CPU) byte {
		return cpu.RLA()
	}},
	0x18: &Instruction{0x18, "JR R8", 2, func(cpu *CPU) byte {
		return cpu.JR_e()
	}},
	0x19: &Instruction{0x19, "ADD HL,DE", 1, func(cpu *CPU) byte {
		return cpu.ADD_HL_rr(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x1A: &Instruction{0x1A, "LD A,(DE)", 1, func(cpu *CPU) byte {
		return cpu.LD_A_DE()
	}},
	0x1B: &Instruction{0x1B, "DEC DE", 1, func(cpu *CPU) byte {
		return cpu.DEC_rr(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x1C: &Instruction{0x1C, "INC E", 1, func(cpu *CPU) byte {
		return cpu.INC_r(&cpu.Registers.E)
	}},
	0x1D: &Instruction{0x1D, "DEC E", 1, func(cpu *CPU) byte {
		return cpu.DEC_r(&cpu.Registers.E)
	}},
	0x1E: &Instruction{0x1E, "LD E,d8", 2, func(cpu *CPU) byte {
		return cpu.LD_r_n(&cpu.Registers.E)
	}},
	0x1F: &Instruction{0x3E, "RRA", 1, func(cpu *CPU) byte {
		return cpu.RRA()
	}},
	0x20: &Instruction{0x20, "JR NZ,r8", 2, func(cpu *CPU) byte {
		return cpu.JR_cc_e(CC_NZ)
	}},
	0x21: &Instruction{0x21, "LD HL,d16", 3, func(cpu *CPU) byte {
		return cpu.LD_rr_nn(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x22: &Instruction{0x22, "LD (HL+),A", 1, func(cpu *CPU) byte {
		return cpu.LD_HLI_A()
	}},
	0x23: &Instruction{0x23, "INC HL", 1, func(cpu *CPU) byte {
		return cpu.INC_rr(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x24: &Instruction{0x24, "INC H", 1, func(cpu *CPU) byte {
		return cpu.INC_r(&cpu.Registers.H)
	}},
	0x25: &Instruction{0x25, "DEC H", 1, func(cpu *CPU) byte {
		return cpu.DEC_r(&cpu.Registers.H)
	}},
	0x26: &Instruction{0x26, "LD H,d8", 2, func(cpu *CPU) byte {
		return cpu.LD_r_n(&cpu.Registers.H)
	}},
	0x27: &Instruction{0x27, "DAA", 1, func(cpu *CPU) byte {
		return cpu.DAA()
	}},
	0x28: &Instruction{0x30, "JR Z,r8", 2, func(cpu *CPU) byte {
		return cpu.JR_cc_e(CC_Z)
	}},
	0x29: &Instruction{0x29, "ADD HL,HL", 1, func(cpu *CPU) byte {
		return cpu.ADD_HL_rr(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x2A: &Instruction{0x2A, "LD A,(HL+)", 1, func(cpu *CPU) byte {
		return cpu.LD_A_HLI()
	}},
	0x2B: &Instruction{0x2B, "DEC HL", 1, func(cpu *CPU) byte {
		return cpu.DEC_rr(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x2C: &Instruction{0x2C, "INC L", 1, func(cpu *CPU) byte {
		return cpu.INC_r(&cpu.Registers.L)
	}},
	0x2D: &Instruction{0x2D, "DEC L", 1, func(cpu *CPU) byte {
		return cpu.DEC_r(&cpu.Registers.L)
	}},
	0x2E: &Instruction{0x2E, "LD L,d8", 2, func(cpu *CPU) byte {
		return cpu.LD_r_n(&cpu.Registers.L)
	}},
	0x2F: &Instruction{0x2F, "CPL", 1, func(cpu *CPU) byte {
		return cpu.CPL()
	}},
	0x30: &Instruction{0x30, "JR NC,r8", 2, func(cpu *CPU) byte {
		return cpu.JR_cc_e(CC_NC)
	}},
	0x31: &Instruction{0x31, "LD SP,d16", 3, func(cpu *CPU) byte {
		return cpu.LD_SP_nn()
	}},
	0x32: &Instruction{0x32, "LD (HL-),A", 1, func(cpu *CPU) byte {
		return cpu.LD_HLD_A()
	}},
	0x33: &Instruction{0x33, "INC SP", 1, func(cpu *CPU) byte {
		return cpu.INC_SP()
	}},
	0x34: &Instruction{0x34, "INC (HL)", 1, func(cpu *CPU) byte {
		return cpu.INC_HL()
	}},
	0x35: &Instruction{0x35, "DEC (HL)", 1, func(cpu *CPU) byte {
		return cpu.DEC_HL()
	}},
	0x36: &Instruction{0x36, "LD (HL),d8", 2, func(cpu *CPU) byte {
		return cpu.LD_HL_n()
	}},
	0x37: &Instruction{0x37, "SCF", 1, func(cpu *CPU) byte {
		return cpu.SCF()
	}},
	0x38: &Instruction{0x30, "JR C,r8", 2, func(cpu *CPU) byte {
		return cpu.JR_cc_e(CC_C)
	}},
	0x39: &Instruction{0x39, "ADD HL,SP", 1, func(cpu *CPU) byte {
		return cpu.ADD_HL_SP()
	}},
	0x3A: &Instruction{0x3A, "LD A,(HL-)", 1, func(cpu *CPU) byte {
		return cpu.LD_A_HLD()
	}},
	0x3B: &Instruction{0x3B, "DEC SP", 1, func(cpu *CPU) byte {
		return cpu.DEC_SP()
	}},
	0x3C: &Instruction{0x3C, "INC A", 1, func(cpu *CPU) byte {
		return cpu.INC_r(&cpu.Registers.A)
	}},
	0x3D: &Instruction{0x3D, "DEC A", 1, func(cpu *CPU) byte {
		return cpu.DEC_r(&cpu.Registers.A)
	}},
	0x3E: &Instruction{0x3E, "LD A,d8", 2, func(cpu *CPU) byte {
		return cpu.LD_r_n(&cpu.Registers.A)
	}},
	0x3F: &Instruction{0x3F, "CCF", 1, func(cpu *CPU) byte {
		return cpu.CCF()
	}},
	0x40: &Instruction{0x40, "LD B,B", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.B)
	}},
	0x41: &Instruction{0x41, "LD B,C", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x42: &Instruction{0x42, "LD B,D", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.D)
	}},
	0x43: &Instruction{0x43, "LD B,E", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.E)
	}},
	0x44: &Instruction{0x44, "LD B,H", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.H)
	}},
	0x45: &Instruction{0x45, "LD B,L", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.L)
	}},
	0x46: &Instruction{0x46, "LD B,(HL)", 1, func(cpu *CPU) byte {
		return cpu.LD_r_HL(&cpu.Registers.B)
	}},
	0x47: &Instruction{0x47, "LD B,A", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.A)
	}},
	0x48: &Instruction{0x48, "LD C,B", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.B)
	}},
	0x49: &Instruction{0x49, "LD C,C", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.C)
	}},
	0x4A: &Instruction{0x4A, "LD C,D", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.D)
	}},
	0x4B: &Instruction{0x4B, "LD C,E", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.E)
	}},
	0x4C: &Instruction{0x4C, "LD C,H", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.H)
	}},
	0x4D: &Instruction{0x4D, "LD C,L", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.L)
	}},
	0x4E: &Instruction{0x4E, "LD C,(HL)", 1, func(cpu *CPU) byte {
		return cpu.LD_r_HL(&cpu.Registers.C)
	}},
	0x4F: &Instruction{0x4F, "LD C,A", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.A)
	}},
	0x50: &Instruction{0x50, "LD D,B", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.B)
	}},
	0x51: &Instruction{0x51, "LD D,C", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.C)
	}},
	0x52: &Instruction{0x52, "LD D,D", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.D)
	}},
	0x53: &Instruction{0x53, "LD D,E", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x54: &Instruction{0x54, "LD D,H", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.H)
	}},
	0x55: &Instruction{0x55, "LD D,L", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.L)
	}},
	0x56: &Instruction{0x56, "LD D,(HL)", 1, func(cpu *CPU) byte {
		return cpu.LD_r_HL(&cpu.Registers.D)
	}},
	0x57: &Instruction{0x57, "LD D,A", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.A)
	}},
	0x58: &Instruction{0x58, "LD E,B", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.B)
	}},
	0x59: &Instruction{0x59, "LD E,C", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.C)
	}},
	0x5A: &Instruction{0x5A, "LD E,D", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.D)
	}},
	0x5B: &Instruction{0x5B, "LD E,E", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.E)
	}},
	0x5C: &Instruction{0x5C, "LD E,H", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.H)
	}},
	0x5D: &Instruction{0x5D, "LD E,L", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.L)
	}},
	0x5E: &Instruction{0x5E, "LD E,(HL)", 1, func(cpu *CPU) byte {
		return cpu.LD_r_HL(&cpu.Registers.E)
	}},
	0x5F: &Instruction{0x5F, "LD E,A", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.A)
	}},
	0x60: &Instruction{0x60, "LD H,B", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.B)
	}},
	0x61: &Instruction{0x61, "LD H,C", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.C)
	}},
	0x62: &Instruction{0x62, "LD H,D", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.D)
	}},
	0x63: &Instruction{0x63, "LD H,E", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.E)
	}},
	0x64: &Instruction{0x64, "LD H,H", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.H)
	}},
	0x65: &Instruction{0x65, "LD H,L", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x66: &Instruction{0x66, "LD H,(HL)", 1, func(cpu *CPU) byte {
		return cpu.LD_r_HL(&cpu.Registers.H)
	}},
	0x67: &Instruction{0x67, "LD H,A", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.A)
	}},
	0x68: &Instruction{0x68, "LD L,B", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.B)
	}},
	0x69: &Instruction{0x69, "LD L,C", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.C)
	}},
	0x6A: &Instruction{0x6A, "LD L,D", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.D)
	}},
	0x6B: &Instruction{0x6B, "LD L,E", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.E)
	}},
	0x6C: &Instruction{0x6C, "LD L,H", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.H)
	}},
	0x6D: &Instruction{0x6D, "LD L,L", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.L)
	}},
	0x6E: &Instruction{0x6E, "LD L,(HL)", 1, func(cpu *CPU) byte {
		return cpu.LD_r_HL(&cpu.Registers.L)
	}},
	0x6F: &Instruction{0x6F, "LD L,A", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.A)
	}},
	0x70: &Instruction{0x70, "LD (HL),B", 1, func(cpu *CPU) byte {
		return cpu.LD_HL_r(&cpu.Registers.B)
	}},
	0x71: &Instruction{0x71, "LD (HL),C", 1, func(cpu *CPU) byte {
		return cpu.LD_HL_r(&cpu.Registers.C)
	}},
	0x72: &Instruction{0x72, "LD (HL),D", 1, func(cpu *CPU) byte {
		return cpu.LD_HL_r(&cpu.Registers.D)
	}},
	0x73: &Instruction{0x73, "LD (HL),E", 1, func(cpu *CPU) byte {
		return cpu.LD_HL_r(&cpu.Registers.E)
	}},
	0x74: &Instruction{0x74, "LD (HL),H", 1, func(cpu *CPU) byte {
		return cpu.LD_HL_r(&cpu.Registers.H)
	}},
	0x75: &Instruction{0x75, "LD (HL),L", 1, func(cpu *CPU) byte {
		return cpu.LD_HL_r(&cpu.Registers.L)
	}},
	0x76: &Instruction{0x76, "HALT", 1, func(cpu *CPU) byte {
		return cpu.HALT()
	}},
	0x77: &Instruction{0x77, "LD (HL), A", 1, func(cpu *CPU) byte {
		return cpu.LD_HL_r(&cpu.Registers.A)
	}},
	0x78: &Instruction{0x78, "LD A,B", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.B)
	}},
	0x79: &Instruction{0x79, "LD A,C", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.C)
	}},
	0x7A: &Instruction{0x7A, "LD A,D", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.D)
	}},
	0x7B: &Instruction{0x7B, "LD A,E", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.E)
	}},
	0x7C: &Instruction{0x7C, "LD A,H", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.H)
	}},
	0x7D: &Instruction{0x7D, "LD A,L", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.L)
	}},
	0x7E: &Instruction{0x7E, "LD A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.LD_r_HL(&cpu.Registers.A)
	}},
	0x7F: &Instruction{0x7F, "LD A,A", 1, func(cpu *CPU) byte {
		return cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.A)
	}},
	0x80: &Instruction{0x80, "ADD A,B", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.Registers.B, 1)
	}},
	0x81: &Instruction{0x81, "ADD A,C", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.Registers.C, 1)
	}},
	0x82: &Instruction{0x82, "ADD A,D", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.Registers.D, 1)
	}},
	0x83: &Instruction{0x83, "ADD A,E", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.Registers.E, 1)
	}},
	0x84: &Instruction{0x84, "ADD A,H", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.Registers.H, 1)
	}},
	0x85: &Instruction{0x85, "ADD A,L", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.Registers.L, 1)
	}},
	0x86: &Instruction{0x86, "ADD A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0x87: &Instruction{0x87, "ADD A,A", 1, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.Registers.A, 1)
	}},
	0x88: &Instruction{0x88, "ADC A,B", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.Registers.B, 1)
	}},
	0x89: &Instruction{0x89, "ADC A,C", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.Registers.C, 1)
	}},
	0x8A: &Instruction{0x8A, "ADC A,D", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.Registers.D, 1)
	}},
	0x8B: &Instruction{0x8B, "ADC A,E", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.Registers.E, 1)
	}},
	0x8C: &Instruction{0x8C, "ADC A,H", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.Registers.H, 1)
	}},
	0x8D: &Instruction{0x8D, "ADC A,L", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.Registers.L, 1)
	}},
	0x8E: &Instruction{0x8E, "ADC A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0x8F: &Instruction{0x8F, "ADC A,A", 1, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.Registers.A, 1)
	}},
	0x90: &Instruction{0x90, "SUB A,B", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.Registers.B, 1)
	}},
	0x91: &Instruction{0x91, "SUB A,C", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.Registers.C, 1)
	}},
	0x92: &Instruction{0x92, "SUB A,D", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.Registers.D, 1)
	}},
	0x93: &Instruction{0x93, "SUB A,E", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.Registers.E, 1)
	}},
	0x94: &Instruction{0x94, "SUB A,H", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.Registers.H, 1)
	}},
	0x95: &Instruction{0x95, "SUB A,L", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.Registers.L, 1)
	}},
	0x96: &Instruction{0x96, "SUB A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0x97: &Instruction{0x97, "SUB A,A", 1, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.Registers.A, 1)
	}},
	0x98: &Instruction{0x98, "SBC A,B", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.Registers.B, 1)
	}},
	0x99: &Instruction{0x99, "SBC A,C", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.Registers.C, 1)
	}},
	0x9A: &Instruction{0x9A, "SBC A,D", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.Registers.D, 1)
	}},
	0x9B: &Instruction{0x9B, "SBC A,E", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.Registers.E, 1)
	}},
	0x9C: &Instruction{0x9C, "SBC A,H", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.Registers.H, 1)
	}},
	0x9D: &Instruction{0x9D, "SBC A,L", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.Registers.L, 1)
	}},
	0x9E: &Instruction{0x9E, "SBC A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0x9F: &Instruction{0x9F, "SBC A,A", 1, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.Registers.A, 1)
	}},
	0xA0: &Instruction{0xA0, "AND A,B", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.Registers.B, 1)
	}},
	0xA1: &Instruction{0xA1, "AND A,C", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.Registers.C, 1)
	}},
	0xA2: &Instruction{0xA2, "AND A,D", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.Registers.D, 1)
	}},
	0xA3: &Instruction{0xA3, "AND A,E", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.Registers.E, 1)
	}},
	0xA4: &Instruction{0xA4, "AND A,H", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.Registers.H, 1)
	}},
	0xA5: &Instruction{0xA5, "AND A,L", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.Registers.L, 1)
	}},
	0xA6: &Instruction{0xA6, "AND A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0xA7: &Instruction{0xA7, "AND A,A", 1, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.Registers.A, 1)
	}},
	0xA8: &Instruction{0xA8, "XOR A,B", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.Registers.B, 1)
	}},
	0xA9: &Instruction{0xA9, "XOR A,C", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.Registers.C, 1)
	}},
	0xAA: &Instruction{0xAA, "XOR A,D", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.Registers.D, 1)
	}},
	0xAB: &Instruction{0xAB, "XOR A,E", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.Registers.E, 1)
	}},
	0xAC: &Instruction{0xAC, "XOR A,H", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.Registers.H, 1)
	}},
	0xAD: &Instruction{0xAD, "XOR A,L", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.Registers.L, 1)
	}},
	0xAE: &Instruction{0xAE, "XOR A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0xAF: &Instruction{0xAF, "XOR A,A", 1, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.Registers.A, 1)
	}},
	0xB0: &Instruction{0xB0, "OR A,B", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.Registers.B, 1)
	}},
	0xB1: &Instruction{0xB1, "OR A,C", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.Registers.C, 1)
	}},
	0xB2: &Instruction{0xB2, "OR A,D", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.Registers.D, 1)
	}},
	0xB3: &Instruction{0xB3, "OR A,E", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.Registers.E, 1)
	}},
	0xB4: &Instruction{0xB4, "OR A,H", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.Registers.H, 1)
	}},
	0xB5: &Instruction{0xB5, "OR A,L", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.Registers.L, 1)
	}},
	0xB6: &Instruction{0xB6, "OR A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0xB7: &Instruction{0xB7, "OR A,A", 1, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.Registers.A, 1)
	}},
	0xB8: &Instruction{0xB8, "CP A,B", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.Registers.B, 1)
	}},
	0xB9: &Instruction{0xB9, "CP A,C", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.Registers.C, 1)
	}},
	0xBA: &Instruction{0xBA, "CP A,D", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.Registers.D, 1)
	}},
	0xBB: &Instruction{0xBB, "CP A,E", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.Registers.E, 1)
	}},
	0xBC: &Instruction{0xBC, "CP A,H", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.Registers.H, 1)
	}},
	0xBD: &Instruction{0xBD, "CP A,L", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.Registers.L, 1)
	}},
	0xBE: &Instruction{0xBE, "CP A,(HL)", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)), 2)
	}},
	0xBF: &Instruction{0xBF, "CP A,A", 1, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.Registers.A, 1)
	}},
	0xC0: &Instruction{0xC0, "RET NZ", 1, func(cpu *CPU) byte {
		return cpu.RET_cc(CC_NZ)
	}},
	0xC1: &Instruction{0xC1, "POP BC", 1, func(cpu *CPU) byte {
		return cpu.POP_qq(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0xC2: &Instruction{0xC2, "JP NZ, a16", 3, func(cpu *CPU) byte {
		return cpu.JP_cc_nn(CC_NZ)
	}},
	0xC3: &Instruction{0xC3, "JP d16", 3, func(cpu *CPU) byte {
		return cpu.JP_nn()
	}},
	0xC4: &Instruction{0xC4, "CALL NZ, a16", 3, func(cpu *CPU) byte {
		return cpu.CALL_cc(CC_NZ)
	}},
	0xC5: &Instruction{0xC5, "PUSH BC", 1, func(cpu *CPU) byte {
		return cpu.PUSH_qq(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0xC6: &Instruction{0xC6, "ADD A,D8", 2, func(cpu *CPU) byte {
		return cpu.ADD_s(cpu.GetByteOffset(1), 2)
	}},
	0xC7: &Instruction{0xC7, "RST 00H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x00)
	}},
	0xC8: &Instruction{0xC8, "RET Z", 1, func(cpu *CPU) byte {
		return cpu.RET_cc(CC_Z)
	}},
	0xC9: &Instruction{0xC9, "RET", 1, func(cpu *CPU) byte {
		return cpu.RET()
	}},
	0xCA: &Instruction{0xCA, "JP Z, a16", 3, func(cpu *CPU) byte {
		return cpu.JP_cc_nn(CC_Z)
	}},
	// 0xCB handled in main CPU loop
	0xCC: &Instruction{0xCC, "CALL Z, a16", 3, func(cpu *CPU) byte {
		return cpu.CALL_cc(CC_Z)
	}},
	0xCD: &Instruction{0xCD, "CALL addr", 3, func(cpu *CPU) byte {
		return cpu.CALL()
	}},
	0xCE: &Instruction{0xCE, "ADC A,D8", 2, func(cpu *CPU) byte {
		return cpu.ADC_A_s(cpu.GetByteOffset(1), 2)
	}},
	0xCF: &Instruction{0xCF, "RST 08H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x08)
	}},
	0xD0: &Instruction{0xD0, "RET NC", 1, func(cpu *CPU) byte {
		return cpu.RET_cc(CC_NC)
	}},
	0xD1: &Instruction{0xD1, "POP DE", 1, func(cpu *CPU) byte {
		return cpu.POP_qq(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0xD2: &Instruction{0xD2, "JP NC,a16", 3, func(cpu *CPU) byte {
		return cpu.JP_cc_nn(CC_NC)
	}},
	0xD4: &Instruction{0xD4, "CALL NC,a16", 3, func(cpu *CPU) byte {
		return cpu.CALL_cc(CC_NC)
	}},
	0xD5: &Instruction{0xD5, "PUSH DE", 1, func(cpu *CPU) byte {
		return cpu.PUSH_qq(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0xD6: &Instruction{0xD6, "SUB A,D8", 2, func(cpu *CPU) byte {
		return cpu.SUB_s(cpu.GetByteOffset(1), 2)
	}},
	0xD7: &Instruction{0xD7, "RST 10H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x10)
	}},
	0xD8: &Instruction{0xD8, "RET C", 1, func(cpu *CPU) byte {
		return cpu.RET_cc(CC_C)
	}},
	0xD9: &Instruction{0xD9, "RETI", 1, func(cpu *CPU) byte {
		return cpu.RETI()
	}},
	0xDA: &Instruction{0xDA, "JP C,a16", 3, func(cpu *CPU) byte {
		return cpu.JP_cc_nn(CC_C)
	}},
	0xDC: &Instruction{0xDC, "CALL C,a16", 3, func(cpu *CPU) byte {
		return cpu.CALL_cc(CC_C)
	}},
	0xDE: &Instruction{0xDE, "SBC A,D8", 2, func(cpu *CPU) byte {
		return cpu.SBC_s(cpu.GetByteOffset(1), 2)
	}},
	0xDF: &Instruction{0xDF, "RST 18H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x18)
	}},
	0xE0: &Instruction{0xE0, "LDH (a8),A", 2, func(cpu *CPU) byte {
		return cpu.LDH_n_A()
	}},
	0xE1: &Instruction{0xE1, "POP HL", 1, func(cpu *CPU) byte {
		return cpu.POP_qq(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0xE2: &Instruction{0xE2, "LD (C),A", 1, func(cpu *CPU) byte {
		return cpu.LD_C_A()
	}},
	0xE5: &Instruction{0xC5, "PUSH HL", 1, func(cpu *CPU) byte {
		return cpu.PUSH_qq(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0xE6: &Instruction{0xE6, "AND A,D8", 2, func(cpu *CPU) byte {
		return cpu.AND_s(cpu.GetByteOffset(1), 2)
	}},
	0xE7: &Instruction{0xE7, "RST 20H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x20)
	}},
	0xE8: &Instruction{0xE8, "ADD SP,r8", 2, func(cpu *CPU) byte {
		return cpu.ADD_SP_e()
	}},
	0xE9: &Instruction{0xE9, "JP (HL)", 1, func(cpu *CPU) byte {
		return cpu.JP_HL()
	}},
	0xEA: &Instruction{0xEA, "LD (a16),A", 3, func(cpu *CPU) byte {
		return cpu.LD_nn_A()
	}},
	0xEE: &Instruction{0xEE, "XOR A,D8", 2, func(cpu *CPU) byte {
		return cpu.XOR_s(cpu.GetByteOffset(1), 2)
	}},
	0xEF: &Instruction{0xEF, "RST 28H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x28)
	}},
	0xF0: &Instruction{0xF0, "LDH A,(a8)", 2, func(cpu *CPU) byte {
		return cpu.LDH_A_n()
	}},
	0xF1: &Instruction{0xF1, "POP AF", 1, func(cpu *CPU) byte {
		return cpu.POP_AF()
	}},
	0xF2: &Instruction{0xF2, "LD A,(C)", 1, func(cpu *CPU) byte {
		return cpu.LD_A_C()
	}},
	0xF3: &Instruction{0xF3, "DI", 1, func(cpu *CPU) byte {
		return cpu.DI()
	}},
	0xF5: &Instruction{0xF5, "PUSH AF", 1, func(cpu *CPU) byte {
		return cpu.PUSH_qq(&cpu.Registers.A, &cpu.Registers.F)
	}},
	0xF6: &Instruction{0xF6, "OR A,D8", 2, func(cpu *CPU) byte {
		return cpu.OR_s(cpu.GetByteOffset(1), 2)
	}},
	0xF7: &Instruction{0xF7, "RST 30H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x30)
	}},
	0xF8: &Instruction{0xF8, "LD HL,SP+r8", 2, func(cpu *CPU) byte {
		return cpu.LD_HL_SP_e()
	}},
	0xF9: &Instruction{0xF9, "LD SP,HL", 1, func(cpu *CPU) byte {
		return cpu.LD_SP_HL()
	}},
	0xFA: &Instruction{0xFA, "LD A,(a16)", 3, func(cpu *CPU) byte {
		return cpu.LD_A_nn()
	}},
	0xFB: &Instruction{0xFB, "EI", 1, func(cpu *CPU) byte {
		return cpu.EI()
	}},
	0xFE: &Instruction{0xFE, "CP D8", 2, func(cpu *CPU) byte {
		return cpu.CP_s(cpu.GetByteOffset(1), 2)
	}},
	0xFF: &Instruction{0xFF, "RST 38H", 1, func(cpu *CPU) byte {
		return cpu.RST(0x38)
	}},
}

var CBInstructions map[byte]*Instruction = map[byte]*Instruction{
	0x00: &Instruction{0x00, "RLC B", 2, func(cpu *CPU) byte {
		return cpu.RLC_r(&cpu.Registers.B)
	}},
	0x01: &Instruction{0x01, "RLC C", 2, func(cpu *CPU) byte {
		return cpu.RLC_r(&cpu.Registers.C)
	}},
	0x02: &Instruction{0x02, "RLC D", 2, func(cpu *CPU) byte {
		return cpu.RLC_r(&cpu.Registers.D)
	}},
	0x03: &Instruction{0x03, "RLC E", 2, func(cpu *CPU) byte {
		return cpu.RLC_r(&cpu.Registers.E)
	}},
	0x04: &Instruction{0x04, "RLC H", 2, func(cpu *CPU) byte {
		return cpu.RLC_r(&cpu.Registers.H)
	}},
	0x05: &Instruction{0x05, "RLC L", 2, func(cpu *CPU) byte {
		return cpu.RLC_r(&cpu.Registers.L)
	}},
	0x06: &Instruction{0x06, "RLC (HL)", 2, func(cpu *CPU) byte {
		return cpu.RLC_HL()
	}},
	0x07: &Instruction{0x07, "RLC A", 2, func(cpu *CPU) byte {
		return cpu.RLC_r(&cpu.Registers.A)
	}},
	0x08: &Instruction{0x08, "RRC B", 2, func(cpu *CPU) byte {
		return cpu.RRC_r(&cpu.Registers.B)
	}},
	0x09: &Instruction{0x09, "RRC C", 2, func(cpu *CPU) byte {
		return cpu.RRC_r(&cpu.Registers.C)
	}},
	0x0A: &Instruction{0x0A, "RRC D", 2, func(cpu *CPU) byte {
		return cpu.RRC_r(&cpu.Registers.D)
	}},
	0x0B: &Instruction{0x0B, "RRC E", 2, func(cpu *CPU) byte {
		return cpu.RRC_r(&cpu.Registers.E)
	}},
	0x0C: &Instruction{0x0C, "RRC H", 2, func(cpu *CPU) byte {
		return cpu.RRC_r(&cpu.Registers.H)
	}},
	0x0D: &Instruction{0x0D, "RRC L", 2, func(cpu *CPU) byte {
		return cpu.RRC_r(&cpu.Registers.L)
	}},
	0x0E: &Instruction{0x0E, "RRC (HL)", 2, func(cpu *CPU) byte {
		return cpu.RRC_HL()
	}},
	0x0F: &Instruction{0x0F, "RRC A", 2, func(cpu *CPU) byte {
		return cpu.RRC_r(&cpu.Registers.A)
	}},
	0x10: &Instruction{0x10, "RL B", 2, func(cpu *CPU) byte {
		return cpu.RL_r(&cpu.Registers.B)
	}},
	0x11: &Instruction{0x11, "RL C", 2, func(cpu *CPU) byte {
		return cpu.RL_r(&cpu.Registers.C)
	}},
	0x12: &Instruction{0x12, "RL D", 2, func(cpu *CPU) byte {
		return cpu.RL_r(&cpu.Registers.D)
	}},
	0x13: &Instruction{0x13, "RL E", 2, func(cpu *CPU) byte {
		return cpu.RL_r(&cpu.Registers.E)
	}},
	0x14: &Instruction{0x14, "RL H", 2, func(cpu *CPU) byte {
		return cpu.RL_r(&cpu.Registers.H)
	}},
	0x15: &Instruction{0x15, "RL L", 2, func(cpu *CPU) byte {
		return cpu.RL_r(&cpu.Registers.L)
	}},
	0x16: &Instruction{0x16, "RL (HL)", 2, func(cpu *CPU) byte {
		return cpu.RL_HL()
	}},
	0x17: &Instruction{0x17, "RL A", 2, func(cpu *CPU) byte {
		return cpu.RL_r(&cpu.Registers.A)
	}},
	0x18: &Instruction{0x18, "RR B", 2, func(cpu *CPU) byte {
		return cpu.RR_r(&cpu.Registers.B)
	}},
	0x19: &Instruction{0x19, "RR C", 2, func(cpu *CPU) byte {
		return cpu.RR_r(&cpu.Registers.C)
	}},
	0x1A: &Instruction{0x1A, "RR D", 2, func(cpu *CPU) byte {
		return cpu.RR_r(&cpu.Registers.D)
	}},
	0x1B: &Instruction{0x1B, "RR E", 2, func(cpu *CPU) byte {
		return cpu.RR_r(&cpu.Registers.E)
	}},
	0x1C: &Instruction{0x1C, "RR H", 2, func(cpu *CPU) byte {
		return cpu.RR_r(&cpu.Registers.H)
	}},
	0x1D: &Instruction{0x1D, "RR L", 2, func(cpu *CPU) byte {
		return cpu.RR_r(&cpu.Registers.L)
	}},
	0x1E: &Instruction{0x1E, "RR (HL)", 2, func(cpu *CPU) byte {
		return cpu.RR_HL()
	}},
	0x1F: &Instruction{0x3E, "RR A", 2, func(cpu *CPU) byte {
		return cpu.RR_r(&cpu.Registers.A)
	}},
	0x20: &Instruction{0x20, "SLA B", 2, func(cpu *CPU) byte {
		return cpu.SLA_r(&cpu.Registers.B)
	}},
	0x21: &Instruction{0x21, "SLA C", 2, func(cpu *CPU) byte {
		return cpu.SLA_r(&cpu.Registers.C)
	}},
	0x22: &Instruction{0x22, "SLA D", 2, func(cpu *CPU) byte {
		return cpu.SLA_r(&cpu.Registers.D)
	}},
	0x23: &Instruction{0x23, "SLA E", 2, func(cpu *CPU) byte {
		return cpu.SLA_r(&cpu.Registers.E)
	}},
	0x24: &Instruction{0x24, "SLA H", 2, func(cpu *CPU) byte {
		return cpu.SLA_r(&cpu.Registers.H)
	}},
	0x25: &Instruction{0x25, "SLA L", 2, func(cpu *CPU) byte {
		return cpu.SLA_r(&cpu.Registers.L)
	}},
	0x26: &Instruction{0x26, "SLA (HL)", 2, func(cpu *CPU) byte {
		return cpu.SLA_HL()
	}},
	0x27: &Instruction{0x27, "SLA A", 2, func(cpu *CPU) byte {
		return cpu.SLA_r(&cpu.Registers.A)
	}},
	0x28: &Instruction{0x30, "SRA B", 2, func(cpu *CPU) byte {
		return cpu.SRA_r(&cpu.Registers.B)
	}},
	0x29: &Instruction{0x29, "SRA C", 2, func(cpu *CPU) byte {
		return cpu.SRA_r(&cpu.Registers.C)
	}},
	0x2A: &Instruction{0x2A, "SRA D", 2, func(cpu *CPU) byte {
		return cpu.SRA_r(&cpu.Registers.D)
	}},
	0x2B: &Instruction{0x2B, "SRA E", 2, func(cpu *CPU) byte {
		return cpu.SRA_r(&cpu.Registers.E)
	}},
	0x2C: &Instruction{0x2C, "SRA H", 2, func(cpu *CPU) byte {
		return cpu.SRA_r(&cpu.Registers.H)
	}},
	0x2D: &Instruction{0x2D, "SRA L", 2, func(cpu *CPU) byte {
		return cpu.SRA_r(&cpu.Registers.L)
	}},
	0x2E: &Instruction{0x2E, "SRA (HL)", 2, func(cpu *CPU) byte {
		return cpu.SRA_HL()
	}},
	0x2F: &Instruction{0x2F, "SRA A", 2, func(cpu *CPU) byte {
		return cpu.SRA_r(&cpu.Registers.A)
	}},
	0x30: &Instruction{0x30, "SWAP B", 2, func(cpu *CPU) byte {
		return cpu.SWAP_r(&cpu.Registers.B)
	}},
	0x31: &Instruction{0x31, "SWAP C", 2, func(cpu *CPU) byte {
		return cpu.SWAP_r(&cpu.Registers.C)
	}},
	0x32: &Instruction{0x32, "SWAP D", 2, func(cpu *CPU) byte {
		return cpu.SWAP_r(&cpu.Registers.D)
	}},
	0x33: &Instruction{0x33, "SWAP E", 2, func(cpu *CPU) byte {
		return cpu.SWAP_r(&cpu.Registers.E)
	}},
	0x34: &Instruction{0x34, "SWAP H", 2, func(cpu *CPU) byte {
		return cpu.SWAP_r(&cpu.Registers.H)
	}},
	0x35: &Instruction{0x35, "SWAP L", 2, func(cpu *CPU) byte {
		return cpu.SWAP_r(&cpu.Registers.L)
	}},
	0x36: &Instruction{0x36, "SWAP (HL)", 2, func(cpu *CPU) byte {
		return cpu.SWAP_HL()
	}},
	0x37: &Instruction{0x37, "SWAP A", 2, func(cpu *CPU) byte {
		return cpu.SWAP_r(&cpu.Registers.A)
	}},
	0x38: &Instruction{0x30, "SRL B", 2, func(cpu *CPU) byte {
		return cpu.SRL_r(&cpu.Registers.B)
	}},
	0x39: &Instruction{0x39, "SRL C", 2, func(cpu *CPU) byte {
		return cpu.SRL_r(&cpu.Registers.C)
	}},
	0x3A: &Instruction{0x3A, "SRL D", 2, func(cpu *CPU) byte {
		return cpu.SRL_r(&cpu.Registers.D)
	}},
	0x3B: &Instruction{0x3B, "SRL E", 2, func(cpu *CPU) byte {
		return cpu.SRL_r(&cpu.Registers.E)
	}},
	0x3C: &Instruction{0x3C, "SRL H", 2, func(cpu *CPU) byte {
		return cpu.SRL_r(&cpu.Registers.H)
	}},
	0x3D: &Instruction{0x3D, "SRL L", 2, func(cpu *CPU) byte {
		return cpu.SRL_r(&cpu.Registers.L)
	}},
	0x3E: &Instruction{0x3E, "SRL (HL)", 2, func(cpu *CPU) byte {
		return cpu.SRL_HL()
	}},
	0x3F: &Instruction{0x3F, "SRL A", 2, func(cpu *CPU) byte {
		return cpu.SRL_r(&cpu.Registers.A)
	}},
	0x40: &Instruction{0x40, "BIT 0,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(0, &cpu.Registers.B)
	}},
	0x41: &Instruction{0x41, "BIT 0,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(0, &cpu.Registers.C)
	}},
	0x42: &Instruction{0x42, "BIT 0,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(0, &cpu.Registers.D)
	}},
	0x43: &Instruction{0x43, "BIT 0,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(0, &cpu.Registers.E)
	}},
	0x44: &Instruction{0x44, "BIT 0,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(0, &cpu.Registers.H)
	}},
	0x45: &Instruction{0x45, "BIT 0,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(0, &cpu.Registers.L)
	}},
	0x46: &Instruction{0x46, "BIT 0,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(0)
	}},
	0x47: &Instruction{0x47, "BIT 0,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(0, &cpu.Registers.A)
	}},
	0x48: &Instruction{0x48, "BIT 1,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(1, &cpu.Registers.B)
	}},
	0x49: &Instruction{0x49, "BIT 1,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(1, &cpu.Registers.C)
	}},
	0x4A: &Instruction{0x4A, "BIT 1,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(1, &cpu.Registers.D)
	}},
	0x4B: &Instruction{0x4B, "BIT 1,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(1, &cpu.Registers.E)
	}},
	0x4C: &Instruction{0x4C, "BIT 1,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(1, &cpu.Registers.H)
	}},
	0x4D: &Instruction{0x4D, "BIT 1,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(1, &cpu.Registers.L)
	}},
	0x4E: &Instruction{0x4E, "BIT 1,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(1)
	}},
	0x4F: &Instruction{0x4F, "BIT 1,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(1, &cpu.Registers.A)
	}},
	0x50: &Instruction{0x50, "BIT 2,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(2, &cpu.Registers.B)
	}},
	0x51: &Instruction{0x51, "BIT 2,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(2, &cpu.Registers.C)
	}},
	0x52: &Instruction{0x52, "BIT 2,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(2, &cpu.Registers.D)
	}},
	0x53: &Instruction{0x53, "BIT 2,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(2, &cpu.Registers.E)
	}},
	0x54: &Instruction{0x54, "BIT 2,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(2, &cpu.Registers.H)
	}},
	0x55: &Instruction{0x55, "BIT 2,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(2, &cpu.Registers.L)
	}},
	0x56: &Instruction{0x56, "BIT 2,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(2)
	}},
	0x57: &Instruction{0x57, "BIT 2,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(2, &cpu.Registers.A)
	}},
	0x58: &Instruction{0x58, "BIT 3,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(3, &cpu.Registers.B)
	}},
	0x59: &Instruction{0x59, "BIT 3,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(3, &cpu.Registers.C)
	}},
	0x5A: &Instruction{0x5A, "BIT 3,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(3, &cpu.Registers.D)
	}},
	0x5B: &Instruction{0x5B, "BIT 3,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(3, &cpu.Registers.E)
	}},
	0x5C: &Instruction{0x5C, "BIT 3,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(3, &cpu.Registers.H)
	}},
	0x5D: &Instruction{0x5D, "BIT 3,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(3, &cpu.Registers.L)
	}},
	0x5E: &Instruction{0x5E, "BIT 3,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(3)
	}},
	0x5F: &Instruction{0x5F, "BIT 3,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(3, &cpu.Registers.A)
	}},
	0x60: &Instruction{0x60, "BIT 4,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(4, &cpu.Registers.B)
	}},
	0x61: &Instruction{0x61, "BIT 4,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(4, &cpu.Registers.C)
	}},
	0x62: &Instruction{0x62, "BIT 4,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(4, &cpu.Registers.D)
	}},
	0x63: &Instruction{0x63, "BIT 4,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(4, &cpu.Registers.E)
	}},
	0x64: &Instruction{0x64, "BIT 4,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(4, &cpu.Registers.H)
	}},
	0x65: &Instruction{0x65, "BIT 4,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(4, &cpu.Registers.L)
	}},
	0x66: &Instruction{0x66, "BIT 4,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(4)
	}},
	0x67: &Instruction{0x67, "BIT 4,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(4, &cpu.Registers.A)
	}},
	0x68: &Instruction{0x68, "BIT 5,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(5, &cpu.Registers.B)
	}},
	0x69: &Instruction{0x69, "BIT 5,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(5, &cpu.Registers.C)
	}},
	0x6A: &Instruction{0x6A, "BIT 45,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(5, &cpu.Registers.D)
	}},
	0x6B: &Instruction{0x6B, "BIT 5,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(5, &cpu.Registers.E)
	}},
	0x6C: &Instruction{0x6C, "BIT 5,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(5, &cpu.Registers.H)
	}},
	0x6D: &Instruction{0x6D, "BIT 5,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(5, &cpu.Registers.L)
	}},
	0x6E: &Instruction{0x6E, "BIT 5,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(5)
	}},
	0x6F: &Instruction{0x6F, "BIT 5,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(5, &cpu.Registers.A)
	}},
	0x70: &Instruction{0x70, "BIT 6,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(6, &cpu.Registers.B)
	}},
	0x71: &Instruction{0x71, "BIT 6,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(6, &cpu.Registers.C)
	}},
	0x72: &Instruction{0x72, "BIT 6,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(6, &cpu.Registers.D)
	}},
	0x73: &Instruction{0x73, "BIT 6,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(6, &cpu.Registers.E)
	}},
	0x74: &Instruction{0x74, "BIT 6,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(6, &cpu.Registers.H)
	}},
	0x75: &Instruction{0x75, "BIT 6,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(6, &cpu.Registers.L)
	}},
	0x76: &Instruction{0x76, "BIT 6,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(6)
	}},
	0x77: &Instruction{0x77, "BIT 6,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(6, &cpu.Registers.A)
	}},
	0x78: &Instruction{0x78, "BIT 7,B", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(7, &cpu.Registers.B)
	}},
	0x79: &Instruction{0x79, "BIT 7,C", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(7, &cpu.Registers.C)
	}},
	0x7A: &Instruction{0x7A, "BIT 7,D", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(7, &cpu.Registers.D)
	}},
	0x7B: &Instruction{0x7B, "BIT 7,E", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(7, &cpu.Registers.E)
	}},
	0x7C: &Instruction{0x7C, "BIT 7,H", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(7, &cpu.Registers.H)
	}},
	0x7D: &Instruction{0x7D, "BIT 7,L", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(7, &cpu.Registers.L)
	}},
	0x7E: &Instruction{0x7E, "BIT 7,(HL)", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_HL(7)
	}},
	0x7F: &Instruction{0x7F, "BIT 7,A", 2, func(cpu *CPU) byte {
		return cpu.BIT_b_r(7, &cpu.Registers.A)
	}},
	0x80: &Instruction{0x80, "RES 0,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(0, &cpu.Registers.B)
	}},
	0x81: &Instruction{0x81, "RES 0,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(0, &cpu.Registers.C)
	}},
	0x82: &Instruction{0x82, "RES 0,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(0, &cpu.Registers.D)
	}},
	0x83: &Instruction{0x83, "RES 0,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(0, &cpu.Registers.E)
	}},
	0x84: &Instruction{0x84, "RES 0,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(0, &cpu.Registers.H)
	}},
	0x85: &Instruction{0x85, "RES 0,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(0, &cpu.Registers.L)
	}},
	0x86: &Instruction{0x86, "RES 0,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(0)
	}},
	0x87: &Instruction{0x87, "RES 0,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(0, &cpu.Registers.A)
	}},
	0x88: &Instruction{0x88, "RES 1,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(1, &cpu.Registers.B)
	}},
	0x89: &Instruction{0x89, "RES 1,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(1, &cpu.Registers.C)
	}},
	0x8A: &Instruction{0x8A, "RES 1,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(1, &cpu.Registers.D)
	}},
	0x8B: &Instruction{0x8B, "RES 1,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(1, &cpu.Registers.E)
	}},
	0x8C: &Instruction{0x8C, "RES 1,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(1, &cpu.Registers.H)
	}},
	0x8D: &Instruction{0x8D, "RES 1,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(1, &cpu.Registers.L)
	}},
	0x8E: &Instruction{0x8E, "RES 1,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(1)
	}},
	0x8F: &Instruction{0x8F, "RES 1,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(1, &cpu.Registers.A)
	}},
	0x90: &Instruction{0x90, "RES 2,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(2, &cpu.Registers.B)
	}},
	0x91: &Instruction{0x91, "RES 2,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(2, &cpu.Registers.C)
	}},
	0x92: &Instruction{0x92, "RES 2,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(2, &cpu.Registers.D)
	}},
	0x93: &Instruction{0x93, "RES 2,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(2, &cpu.Registers.E)
	}},
	0x94: &Instruction{0x94, "RES 2,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(2, &cpu.Registers.H)
	}},
	0x95: &Instruction{0x95, "RES 2,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(2, &cpu.Registers.L)
	}},
	0x96: &Instruction{0x96, "RES 2,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(2)
	}},
	0x97: &Instruction{0x97, "RES 2,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(2, &cpu.Registers.A)
	}},
	0x98: &Instruction{0x98, "RES 3,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(3, &cpu.Registers.B)
	}},
	0x99: &Instruction{0x99, "RES 3,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(3, &cpu.Registers.C)
	}},
	0x9A: &Instruction{0x9A, "RES 3,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(3, &cpu.Registers.D)
	}},
	0x9B: &Instruction{0x9B, "RES 3,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(3, &cpu.Registers.E)
	}},
	0x9C: &Instruction{0x9C, "RES 3,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(3, &cpu.Registers.H)
	}},
	0x9D: &Instruction{0x9D, "RES 3,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(3, &cpu.Registers.L)
	}},
	0x9E: &Instruction{0x9E, "RES 3,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(3)
	}},
	0x9F: &Instruction{0x9F, "RES 3,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(3, &cpu.Registers.A)
	}},
	0xA0: &Instruction{0xA0, "RES 4,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(4, &cpu.Registers.B)
	}},
	0xA1: &Instruction{0xA1, "RES 4,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(4, &cpu.Registers.C)
	}},
	0xA2: &Instruction{0xA2, "RES 4,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(4, &cpu.Registers.D)
	}},
	0xA3: &Instruction{0xA3, "RES 4,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(4, &cpu.Registers.E)
	}},
	0xA4: &Instruction{0xA4, "RES 4,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(4, &cpu.Registers.H)
	}},
	0xA5: &Instruction{0xA5, "RES 4,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(4, &cpu.Registers.L)
	}},
	0xA6: &Instruction{0xA6, "RES 4,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(4)
	}},
	0xA7: &Instruction{0xA7, "RES 4,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(4, &cpu.Registers.A)
	}},
	0xA8: &Instruction{0xA8, "RES 5,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(5, &cpu.Registers.B)
	}},
	0xA9: &Instruction{0xA9, "RES 5,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(5, &cpu.Registers.C)
	}},
	0xAA: &Instruction{0xAA, "RES 5,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(5, &cpu.Registers.D)
	}},
	0xAB: &Instruction{0xAB, "RES 5,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(5, &cpu.Registers.E)
	}},
	0xAC: &Instruction{0xAC, "RES 5,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(5, &cpu.Registers.H)
	}},
	0xAD: &Instruction{0xAD, "RES 5,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(5, &cpu.Registers.L)
	}},
	0xAE: &Instruction{0xAE, "RES 5,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(5)
	}},
	0xAF: &Instruction{0xAF, "RES 5,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(5, &cpu.Registers.A)
	}},
	0xB0: &Instruction{0xB0, "RES 6,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(6, &cpu.Registers.B)
	}},
	0xB1: &Instruction{0xB1, "RES 6,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(6, &cpu.Registers.C)
	}},
	0xB2: &Instruction{0xB2, "RES 6,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(6, &cpu.Registers.D)
	}},
	0xB3: &Instruction{0xB3, "RES 6,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(6, &cpu.Registers.E)
	}},
	0xB4: &Instruction{0xB4, "RES 6,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(6, &cpu.Registers.H)
	}},
	0xB5: &Instruction{0xB5, "RES 6,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(6, &cpu.Registers.L)
	}},
	0xB6: &Instruction{0xB6, "RES 6,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(6)
	}},
	0xB7: &Instruction{0xB7, "RES 6,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(6, &cpu.Registers.A)
	}},
	0xB8: &Instruction{0xB8, "RES 7,B", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(7, &cpu.Registers.B)
	}},
	0xB9: &Instruction{0xB9, "RES 7,C", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(7, &cpu.Registers.C)
	}},
	0xBA: &Instruction{0xBA, "RES 7,D", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(7, &cpu.Registers.D)
	}},
	0xBB: &Instruction{0xBB, "RES 7,E", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(7, &cpu.Registers.E)
	}},
	0xBC: &Instruction{0xBC, "RES 7,H", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(7, &cpu.Registers.H)
	}},
	0xBD: &Instruction{0xBD, "RES 7,L", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(7, &cpu.Registers.L)
	}},
	0xBE: &Instruction{0xBE, "RES 7,(HL)", 2, func(cpu *CPU) byte {
		return cpu.RES_b_HL(7)
	}},
	0xBF: &Instruction{0xBF, "RES 7,A", 2, func(cpu *CPU) byte {
		return cpu.RES_b_r(7, &cpu.Registers.A)
	}},
	0xC0: &Instruction{0xC0, "SET 0,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(0, &cpu.Registers.B)
	}},
	0xC1: &Instruction{0xC1, "SET 0,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(0, &cpu.Registers.C)
	}},
	0xC2: &Instruction{0xC2, "SET 0,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(0, &cpu.Registers.D)
	}},
	0xC3: &Instruction{0xC3, "SET 0,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(0, &cpu.Registers.E)
	}},
	0xC4: &Instruction{0xC4, "SET 0,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(0, &cpu.Registers.H)
	}},
	0xC5: &Instruction{0xC5, "SET 0,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(0, &cpu.Registers.L)
	}},
	0xC6: &Instruction{0xC6, "SET 0,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(0)
	}},
	0xC7: &Instruction{0xC7, "SET 0,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(0, &cpu.Registers.A)
	}},
	0xC8: &Instruction{0xC8, "SET 1,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(1, &cpu.Registers.B)
	}},
	0xC9: &Instruction{0xC9, "SET 1,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(1, &cpu.Registers.C)
	}},
	0xCA: &Instruction{0xCA, "SET 1,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(1, &cpu.Registers.D)
	}},
	0xCB: &Instruction{0xCB, "SET 1,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(1, &cpu.Registers.E)
	}},
	0xCC: &Instruction{0xCC, "SET 1,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(1, &cpu.Registers.H)
	}},
	0xCD: &Instruction{0xCD, "SET 1,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(1, &cpu.Registers.L)
	}},
	0xCE: &Instruction{0xCE, "SET 1,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(1)
	}},
	0xCF: &Instruction{0xCF, "SET 1,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(1, &cpu.Registers.A)
	}},
	0xD0: &Instruction{0xD0, "SET 2,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(2, &cpu.Registers.B)
	}},
	0xD1: &Instruction{0xD1, "SET 2,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(2, &cpu.Registers.C)
	}},
	0xD2: &Instruction{0xD2, "SET 2,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(2, &cpu.Registers.D)
	}},
	0xD3: &Instruction{0xD3, "SET 2,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(2, &cpu.Registers.E)
	}},
	0xD4: &Instruction{0xD4, "SET 2,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(2, &cpu.Registers.H)
	}},
	0xD5: &Instruction{0xD5, "SET 2,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(2, &cpu.Registers.L)
	}},
	0xD6: &Instruction{0xD6, "SET 2,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(2)
	}},
	0xD7: &Instruction{0xD7, "SET 2,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(2, &cpu.Registers.A)
	}},
	0xD8: &Instruction{0xD8, "SET 3,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(3, &cpu.Registers.B)
	}},
	0xD9: &Instruction{0xD9, "SET 3,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(3, &cpu.Registers.C)
	}},
	0xDA: &Instruction{0xDA, "SET 3,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(3, &cpu.Registers.D)
	}},
	0xDB: &Instruction{0xDB, "SET 3,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(3, &cpu.Registers.E)
	}},
	0xDC: &Instruction{0xDC, "SET 3,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(3, &cpu.Registers.H)
	}},
	0xDD: &Instruction{0xDD, "SET 3,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(3, &cpu.Registers.L)
	}},
	0xDE: &Instruction{0xDE, "SET 3,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(3)
	}},
	0xDF: &Instruction{0xDF, "SET 3,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(3, &cpu.Registers.A)
	}},
	0xE0: &Instruction{0xE0, "SET 4,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(4, &cpu.Registers.B)
	}},
	0xE1: &Instruction{0xE1, "SET 4,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(4, &cpu.Registers.C)
	}},
	0xE2: &Instruction{0xE2, "SET 4,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(4, &cpu.Registers.D)
	}},
	0xE3: &Instruction{0xE3, "SET 4,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(4, &cpu.Registers.E)
	}},
	0xE4: &Instruction{0xE4, "SET 4,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(4, &cpu.Registers.H)
	}},
	0xE5: &Instruction{0xE5, "SET 4,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(4, &cpu.Registers.L)
	}},
	0xE6: &Instruction{0xE6, "SET 4,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(4)
	}},
	0xE7: &Instruction{0xE7, "SET 4,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(4, &cpu.Registers.A)
	}},
	0xE8: &Instruction{0xE8, "SET 5,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(5, &cpu.Registers.B)
	}},
	0xE9: &Instruction{0xE9, "SET 5,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(5, &cpu.Registers.C)
	}},
	0xEA: &Instruction{0xEA, "SET 5,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(5, &cpu.Registers.D)
	}},
	0xEB: &Instruction{0xEB, "SET 5,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(5, &cpu.Registers.E)
	}},
	0xEC: &Instruction{0xEC, "SET 5,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(5, &cpu.Registers.H)
	}},
	0xED: &Instruction{0xED, "SET 5,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(5, &cpu.Registers.L)
	}},
	0xEE: &Instruction{0xEE, "SET 5,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(5)
	}},
	0xEF: &Instruction{0xEF, "SET 5,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(5, &cpu.Registers.A)
	}},
	0xF0: &Instruction{0xF0, "SET 6,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(6, &cpu.Registers.B)
	}},
	0xF1: &Instruction{0xF1, "SET 6,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(6, &cpu.Registers.C)
	}},
	0xF2: &Instruction{0xF2, "SET 6,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(6, &cpu.Registers.D)
	}},
	0xF3: &Instruction{0xF3, "SET 6,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(6, &cpu.Registers.E)
	}},
	0xF4: &Instruction{0xF4, "SET 6,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(6, &cpu.Registers.H)
	}},
	0xF5: &Instruction{0xF5, "SET 6,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(6, &cpu.Registers.L)
	}},
	0xF6: &Instruction{0xF6, "SET 6,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(6)
	}},
	0xF7: &Instruction{0xF7, "SET 6,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(6, &cpu.Registers.A)
	}},
	0xF8: &Instruction{0xF8, "SET 7,B", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(7, &cpu.Registers.B)
	}},
	0xF9: &Instruction{0xF9, "SET 7,C", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(7, &cpu.Registers.C)
	}},
	0xFA: &Instruction{0xFA, "SET 7,D", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(7, &cpu.Registers.D)
	}},
	0xFB: &Instruction{0xFB, "SET 7,E", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(7, &cpu.Registers.E)
	}},
	0xFC: &Instruction{0xFC, "SET 7,H", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(7, &cpu.Registers.H)
	}},
	0xFD: &Instruction{0xFD, "SET 7,L", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(7, &cpu.Registers.L)
	}},
	0xFE: &Instruction{0xFE, "SET 7,(HL)", 2, func(cpu *CPU) byte {
		return cpu.SET_b_HL(7)
	}},
	0xFF: &Instruction{0xFF, "SET 7,A", 2, func(cpu *CPU) byte {
		return cpu.SET_b_r(7, &cpu.Registers.A)
	}},
}

// 8-Bit Transfer/Input-Output Instructions

// LD r,r | 1 | ---- | r=r
func (cpu *CPU) LD_r_r(register_1 *byte, register_2 *byte) (cycles byte) {
	*register_1 = *register_2
	return 1
}

// LD r,n | 2 | ---- | r=n
func (cpu *CPU) LD_r_n(register *byte) (cycles byte) {
	*register = cpu.GetByteOffset(1)
	return 2
}

// LD r,(HL) | 2 | ---- | r=(HL)
func (cpu *CPU) LD_r_HL(register *byte) (cycles byte) {
	*register = cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L))
	return 2
}

// LD (HL),r | 2 | ---- | (HL)=r
func (cpu *CPU) LD_HL_r(register *byte) (cycles byte) {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L), *register)
	return 2
}

// LD (HL),n | 3 | ---- | (HL)=n
func (cpu *CPU) LD_HL_n() (cycles byte) {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L), cpu.GetByteOffset(1))
	return 3
}

// LD A,(BC) | 2 | ---- | A=(BC)
func (cpu *CPU) LD_A_BC() (cycles byte) {
	cpu.Registers.A = cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.B, cpu.Registers.C))
	return 2
}

// LD A,(DE) | 2 | ---- | A=(DE)
func (cpu *CPU) LD_A_DE() (cycles byte) {
	cpu.Registers.A = cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.D, cpu.Registers.E))
	return 2
}

// LD A,(C) | 2 | ---- | A=(0xFF00+C)
func (cpu *CPU) LD_A_C() (cycles byte) {
	cpu.Registers.A = cpu.mmu.ReadByte(uint16(0xFF00 + uint16(cpu.Registers.C)))
	return 2
}

// LD (C),A | 2 | ---- | (0xFF00+C)=A
func (cpu *CPU) LD_C_A() (cycles byte) {
	cpu.mmu.WriteByte(uint16(0xFF00+uint16(cpu.Registers.C)), cpu.Registers.A)
	return 2
}

// LDH A,(n) | 3 | ---- | A=(n)
func (cpu *CPU) LDH_A_n() (cycles byte) {
	cpu.Registers.A = cpu.mmu.ReadByte(uint16(0xFF00 + uint16(cpu.GetByteOffset(1))))
	return 3
}

// LDH (n),A | 3 | ---- | (n)=A
func (cpu *CPU) LDH_n_A() (cycles byte) {
	cpu.mmu.WriteByte(uint16(0xFF00+uint16(cpu.GetByteOffset(1))), cpu.Registers.A)
	return 3
}

// LD A,(nn) | 4 | ---- | A=(nn)
func (cpu *CPU) LD_A_nn() (cycles byte) {
	cpu.Registers.A = cpu.mmu.ReadByte(utils.JoinBytes(cpu.GetByteOffset(2), cpu.GetByteOffset(1)))
	return 4
}

// LD (nn),A | 4 | ---- | (nn)=A
func (cpu *CPU) LD_nn_A() (cycles byte) {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.GetByteOffset(2), cpu.GetByteOffset(1)), cpu.Registers.A)
	return 4
}

// LD A,(HLI) | 2 | ---- | A=(HL) HL=HL+1
func (cpu *CPU) LD_A_HLI() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.Registers.A = cpu.mmu.ReadByte(HL)
	HL += 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
	return 2
}

// LD A,(HLD) | 2 | ---- | A=(HL) HL=HL-1
func (cpu *CPU) LD_A_HLD() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.Registers.A = cpu.mmu.ReadByte(HL)
	HL -= 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
	return 2
}

// LD (BC),A | 2 | ---- | (BC)=A
func (cpu *CPU) LD_BC_A() (cycles byte) {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.B, cpu.Registers.C), cpu.Registers.A)
	return 2
}

// LD (DE),A | 2 | ---- | (DE)=A
func (cpu *CPU) LD_DE_A() (cycles byte) {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.D, cpu.Registers.E), cpu.Registers.A)
	return 2
}

// LD (HLI),A | 2 | ---- | (HL)=A HL=HL+1
func (cpu *CPU) LD_HLI_A() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.mmu.WriteByte(HL, cpu.Registers.A)
	HL += 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
	return 2
}

// LD (HLD),A | 2 | ---- | (HL)=A HL=HL-1
func (cpu *CPU) LD_HLD_A() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.mmu.WriteByte(HL, cpu.Registers.A)
	HL -= 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
	return 2
}

// 16-Bit Transfer Instructions

// LD rr,nn | 3 | ---- | rr=nn
func (cpu *CPU) LD_rr_nn(r1 *byte, r2 *byte) (cycles byte) {
	*r1 = cpu.GetByteOffset(2)
	*r2 = cpu.GetByteOffset(1)
	return 3
}

// LD SP,nn | 3 | ---- | SP=nn
func (cpu *CPU) LD_SP_nn() (cycles byte) {
	cpu.SP = utils.JoinBytes(cpu.GetByteOffset(2), cpu.GetByteOffset(1))
	return 3
}

// LD SP,HL | 2 | ---- | SP=HL
func (cpu *CPU) LD_SP_HL() (cycles byte) {
	cpu.SP = utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	return 2
}

// PUSH qq | 4 | ---- | (SP-1)=qqH (SP-2)=qqL SP=SP-2
func (cpu *CPU) PUSH_qq(r1 *byte, r2 *byte) (cycles byte) {
	cpu.pushWordToStack(utils.JoinBytes(*r1, *r2))
	return 4
}

// POP qq | 3 | ---- | qqL=(SP) qqH=(SP+1) SP=SP+2
func (cpu *CPU) POP_qq(r1 *byte, r2 *byte) (cycles byte) {
	*r1, *r2 = utils.SplitBytes(cpu.popWordFromStack())
	return 3
}

// POP AF | 3 | **** | qqL=(SP) qqH=(SP+1) SP=SP+2
func (cpu *CPU) POP_AF() (cycles byte) {
	cpu.Registers.A, cpu.Registers.F = utils.SplitBytes(cpu.popWordFromStack())
	// Since F only holds 4 bits/flags, we mask to ensure only bits 4-7 are set
	cpu.Registers.F &= 0xF0
	return 3
}

// LDHL SP,e | 3 | **00 | HL=SP+e
func (cpu *CPU) LD_HL_SP_e() (cycles byte) {
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(cpu.addSignedByte(cpu.SP, int8(cpu.GetByteOffset(1))))
	return 3
}

// LD (nn),SP | 5 | ---- | (nn)=SPL (nn+1)==SPH
func (cpu *CPU) LD_nn_SP() (cycles byte) {
	hb, lb := utils.SplitBytes(cpu.SP)
	nn := utils.JoinBytes(cpu.GetByteOffset(2), cpu.GetByteOffset(1))
	cpu.mmu.WriteByte(nn, lb)
	cpu.mmu.WriteByte(nn+1, hb)
	return 5
}

// 8-Bit Arithmetic and Logical Operation Instructions

// ADD s | 1,2 | **0* | A=A+s
func (cpu *CPU) ADD_s(s byte, cycles byte) byte {
	cpu.Registers.A = cpu.addBytes(cpu.Registers.A, s)
	return cycles
}

// ADC A,s | 1,2 | **0* | A=A+s+CY
func (cpu *CPU) ADC_A_s(s byte, cycles byte) byte {
	cpu.Registers.A = cpu.addBytesWithCarry(cpu.Registers.A, s)
	return cycles
}

// SUB s | 1,2 | **1* | A=A-s
func (cpu *CPU) SUB_s(s byte, cycles byte) byte {
	cpu.Registers.A = cpu.subBytes(cpu.Registers.A, s)
	return cycles
}

// SBC A,s | 1,2 | **1* | A=A-s-CY
func (cpu *CPU) SBC_s(s byte, cycles byte) byte {
	cpu.Registers.A = cpu.subBytesWithCarry(cpu.Registers.A, s)
	return cycles
}

// AND s | 1,2 | 010* | A=A&s
func (cpu *CPU) AND_s(s byte, cycles byte) byte {
	cpu.Registers.A = cpu.Registers.A & s

	cpu.ResetFlag(CY)
	cpu.SetFlag(H)
	cpu.ResetFlag(N)

	if cpu.Registers.A == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return cycles
}

// OR s | 1,2 | 000* | A=A|s
func (cpu *CPU) OR_s(s byte, cycles byte) byte {
	cpu.Registers.A = cpu.Registers.A | s

	cpu.ResetFlag(CY)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if cpu.Registers.A == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return cycles
}

// XOR s | 1,2 | 000* | A=A^s
func (cpu *CPU) XOR_s(s byte, cycles byte) byte {
	cpu.Registers.A = cpu.Registers.A ^ s

	cpu.ResetFlag(CY)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if cpu.Registers.A == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return cycles
}

//CP s | 1,2 | **1* | A-s
func (cpu *CPU) CP_s(s byte, cycles byte) byte {
	cpu.subBytes(cpu.Registers.A, s)
	return cycles
}

// INC r | 1 | -*0* | r=r+1
func (cpu *CPU) INC_r(r *byte) (cycles byte) {
	*r = cpu.incByte(*r)
	return 1
}

// INC (HL) | 3 | -*0* | (HL)=(HL)+1
func (cpu *CPU) INC_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	cpu.mmu.WriteByte(HL, cpu.incByte(value))
	return 3
}

// DEC r | 1 | -*1* | r=r-1
func (cpu *CPU) DEC_r(r *byte) (cycles byte) {
	*r = cpu.decByte(*r)
	return 1
}

// DEC (HL) | 3 | -*1* | (HL)=(HL)-1
func (cpu *CPU) DEC_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	cpu.mmu.WriteByte(HL, cpu.decByte(value))
	return 3
}

// 16-Bit Arithmetic and Logical Operation Instructions

// ADD HL,rr | 2 | **0- | HL=HL+rr
func (cpu *CPU) ADD_HL_rr(r1 *byte, r2 *byte) (cycles byte) {
	ss := utils.JoinBytes(*r1, *r2)

	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(cpu.addHL_rr(ss))

	return 2
}

// ADD HL,SP | 2 | **0- | HL=HL+SP
func (cpu *CPU) ADD_HL_SP() (cycles byte) {
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(cpu.addHL_rr(cpu.SP))

	return 2
}

// ADD SP,e | 4 | **00 | SP=SP+e
func (cpu *CPU) ADD_SP_e() (cycles byte) {
	cpu.SP = cpu.addSignedByte(cpu.SP, int8(cpu.GetByteOffset(1)))

	return 4
}

// INC rr | 2 | ---- | rr=rr+1
func (cpu *CPU) INC_rr(r1 *byte, r2 *byte) (cycles byte) {
	var ss uint16 = utils.JoinBytes(*r1, *r2)
	ss += 0x01
	*r1, *r2 = utils.SplitBytes(ss)
	return 2
}

// DEC rr | 2 | ---- | rr=rr-1
func (cpu *CPU) DEC_rr(r1 *byte, r2 *byte) (cycles byte) {
	var rr uint16 = utils.JoinBytes(*r1, *r2)
	rr -= 0x01
	*r1, *r2 = utils.SplitBytes(rr)
	return 2
}

// INC SP | 2 | ---- | SP=SP+1
func (cpu *CPU) INC_SP() (cycles byte) {
	cpu.SP += 1
	return 2
}

// DEC SP | 2 | ---- | SP=SP-1
func (cpu *CPU) DEC_SP() (cycles byte) {
	cpu.SP -= 1
	return 2
}

// Rotate Shift Instructions

// RLCA | 1 | *000 | A<<1 A0=A7 CY=A7
func (cpu *CPU) RLCA() (cycles byte) {
	cpu.Registers.A = cpu.rotateLeft(cpu.Registers.A)
	// The rotateLeft function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
	return 1
}

// RLA | 1 | *000 | A<<1 A0=CY CY=A7
func (cpu *CPU) RLA() (cycles byte) {
	cpu.Registers.A = cpu.rotateLeftThroughCarry(cpu.Registers.A)
	// The rotateLeftThroughCarry function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
	return 1
}

// RRCA | 1 | *000 | A>>1 A7=A0 CY=A0
func (cpu *CPU) RRCA() (cycles byte) {
	cpu.Registers.A = cpu.rotateRight(cpu.Registers.A)
	// The rotateRight function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
	return 1
}

// RRA | 1 | *000 | A>>1 A7=CY CY=A0
func (cpu *CPU) RRA() (cycles byte) {
	cpu.Registers.A = cpu.rotateRightThroughCarry(cpu.Registers.A)
	// The rotateRightThroughCarry function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
	return 1
}

// RLC r | 2 | *00* | r<<1 r0=r7 CY=r7
func (cpu *CPU) RLC_r(r *byte) (cycles byte) {
	*r = cpu.rotateLeft(*r)
	return 2
}

// RLC (HL) | 4 | *00* | (HL)<<1 (HL)0=(HL)7 CY=(HL)7
func (cpu *CPU) RLC_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.rotateLeft(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// RL r | 2 | *00* | r<<1 r0=CY CY=r7
func (cpu *CPU) RL_r(r *byte) (cycles byte) {
	*r = cpu.rotateLeftThroughCarry(*r)
	return 2
}

// RL (HL) | 4 | *00* | (HL)<<1 (HL)0=CY CY=(HL)7
func (cpu *CPU) RL_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.rotateLeftThroughCarry(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// RRC r | 2 | *00* | r>>1 r7=r0 CY=r0
func (cpu *CPU) RRC_r(r *byte) (cycles byte) {
	*r = cpu.rotateRight(*r)
	return 2
}

// RRC (HL) | 4 | *00* | (HL)>>1 (HL)7=(HL)0 CY=(HL)0
func (cpu *CPU) RRC_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.rotateRight(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// RR r | 2 | *00* | r>>1 r7=CY CY=r0
func (cpu *CPU) RR_r(r *byte) (cycles byte) {
	*r = cpu.rotateRightThroughCarry(*r)
	return 2
}

// RR (HL) | 4 | *00* | (HL)>>1 (HL)7=CY CY=(HL)0
func (cpu *CPU) RR_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.rotateRightThroughCarry(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// SLA r | 2 | *00* | r<<1 CY=r7
func (cpu *CPU) SLA_r(r *byte) (cycles byte) {
	*r = cpu.shiftLeftArithmetic(*r)
	return 2
}

// SLA (HL) | 4 | *00* | (HL)<<1 CY=(HL)7
func (cpu *CPU) SLA_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.shiftLeftArithmetic(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// SRA r | 2 | *00* | r>>1 r7=r7 CY=r0
func (cpu *CPU) SRA_r(r *byte) (cycles byte) {
	*r = cpu.shiftRightArithmetic(*r)
	return 2
}

// SRA (HL) | 4 | *00* | (HL)>>1 (HL)7=(HL)7 CY=(HL)0
func (cpu *CPU) SRA_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.shiftRightArithmetic(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// SRL r | 2 | *00* | r>>1 CY=r0
func (cpu *CPU) SRL_r(r *byte) (cycles byte) {
	*r = cpu.shiftRightLogical(*r)
	return 2
}

// SRL (HL) | 4 | *00* | (HL)>>1 CY=(HL)0
func (cpu *CPU) SRL_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.shiftRightLogical(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// SWAP r | 2 | 000* | r=r[4:7]&r[0:3]
func (cpu *CPU) SWAP_r(r *byte) (cycles byte) {
	*r = cpu.swapNibbles(*r)
	return 2
}

// SWAP (HL) | 4 | 000* | (HL)=(HL)[4:7]&(HL)[0:3]
func (cpu *CPU) SWAP_HL() (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.swapNibbles(value)
	cpu.mmu.WriteByte(HL, value)
	return 4
}

// Bit Operations

// BIT b,r | 2 | -10* | Z=~rb
func (cpu *CPU) BIT_b_r(bit byte, r *byte) (cycles byte) {
	cpu.SetFlag(H)
	cpu.ResetFlag(N)

	if utils.IsBitSet(*r, bit) {
		cpu.ResetFlag(Z)
	} else {
		cpu.SetFlag(Z)
	}
	return 2
}

// BIT b,(HL) | 3 | -10* | Z=^(HL)b
func (cpu *CPU) BIT_b_HL(bit byte) (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	var value byte = cpu.mmu.ReadByte(HL)

	cpu.SetFlag(H)
	cpu.ResetFlag(N)

	if utils.IsBitSet(value, bit) {
		cpu.ResetFlag(Z)
	} else {
		cpu.SetFlag(Z)
	}

	return 3
}

// SET b,r | 2 | ---- | rb=1
func (cpu *CPU) SET_b_r(bit byte, r *byte) (cycles byte) {
	*r = utils.SetBit(*r, bit)
	return 2
}

// SET b,(HL) | 4 | ---- | (HL)b=1
func (cpu *CPU) SET_b_HL(bit byte) (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)

	var value byte = cpu.mmu.ReadByte(HL)

	cpu.mmu.WriteByte(HL, utils.SetBit(value, bit))
	return 4
}

// RES b,r | 2 | ---- | rb=0
func (cpu *CPU) RES_b_r(bit byte, r *byte) (cycles byte) {
	*r = utils.ClearBit(*r, bit)
	return 2
}

// RES b,(HL) | 4 | ---- | (HL)b=0
func (cpu *CPU) RES_b_HL(bit byte) (cycles byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)

	var value byte = cpu.mmu.ReadByte(HL)

	cpu.mmu.WriteByte(HL, utils.ClearBit(value, bit))

	return 4
}

// Jump Instructions

// JP nn | 4 | ---- | PC=nn
func (cpu *CPU) JP_nn() (cycles byte) {
	cpu.PC = utils.JoinBytes(cpu.GetByteOffset(2), cpu.GetByteOffset(1))
	return 4
}

// JP cc,nn | 4,3 | ---- | if cc true, PC=nn
func (cpu *CPU) JP_cc_nn(conditionCode int) (cycles byte) {
	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		cpu.PC = utils.JoinBytes(cpu.GetByteOffset(2), cpu.GetByteOffset(1))

		return 4
	} else {
		return 3
	}
}

// JR e | 3 | ---- | PC=PC+e
func (cpu *CPU) JR_e() (cycles byte) {
	e := cpu.GetByteOffset(1)

	cpu.PC += cpu.CurrentInstruction.Length

	// e is signed, if it is more than 127 then it is negative
	if e > 127 {
		cpu.PC -= uint16(-e)
	} else {
		cpu.PC += uint16(e)
	}

	return 3
}

// JR cc,e | 3/2 | ---- | if cc true, PC=PC+e
func (cpu *CPU) JR_cc_e(conditionCode int) (cycles byte) {
	e := cpu.GetByteOffset(1)

	cpu.PC += cpu.CurrentInstruction.Length

	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		if e > 127 {
			cpu.PC -= uint16(-e)
		} else {
			cpu.PC += uint16(e)
		}

		return 3
	} else {
		return 2
	}
}

// JP HL | 1 | ---- | PC=HL
func (cpu *CPU) JP_HL() (cycles byte) {
	cpu.PC = utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	return 1
}

// Call/Return Instructions

// CALL nn | 6 | ---- | (SP-1)=PCH (SP-2)=PCL PC=nn SP=SP-2
func (cpu *CPU) CALL() (cycles byte) {
	nextInstruction := cpu.PC + 3
	cpu.pushWordToStack(nextInstruction)
	cpu.PC = utils.JoinBytes(cpu.GetByteOffset(2), cpu.GetByteOffset(1))
	return 6
}

// CALL cc,nn | 6/3 | ---- | (SP-1)=PCH (SP-2)=PCL PC=nn SP=SP-2
func (cpu *CPU) CALL_cc(conditionCode int) (cycles byte) {
	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		return cpu.CALL()
	} else {
		return 3
	}
}

// RET | 4 | ---- | PCL=(SP) PCH=(SP+1) SP=SP+2
func (cpu *CPU) RET() (cycles byte) {
	cpu.PC = cpu.popWordFromStack()
	return 4
}

// RETI | 4 | ---- | PCL=(SP) PCH=(SP+1) SP=SP+2 IME=true
func (cpu *CPU) RETI() (cycles byte) {
	cpu.RET()
	cpu.IME = true
	return 4
}

// RET cc | 5/2 | ---- | if cc true, PCL=(SP) PCH=(SP+1) SP=SP+2
func (cpu *CPU) RET_cc(conditionCode int) (cycles byte) {
	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		return cpu.RET() + 1
	} else {
		return 2
	}
}

// RST t | 4 | ---- | (SP-1)=PCH (SP-2)=PCL SP=SP-2 PCH=0 PCL=t
func (cpu *CPU) RST(t byte) (cycles byte) {
	cpu.pushWordToStack(cpu.PC + 1)
	cpu.PC = uint16(t)
	return 4
}

// General-Purpose Arithmetic Operations and CPU Control Instructions

// DAA | 1 | z-0x | Decimal Adjust acc
func (cpu *CPU) DAA() (cycles byte) {
	a := int(cpu.Registers.A)

	if cpu.IsFlagSet(N) == false {
		if cpu.IsFlagSet(H) || a&0x0F > 9 {
			a += 0x06
		}

		if cpu.IsFlagSet(CY) || a > 0x9F {
			a += 0x60
		}
	} else {
		if cpu.IsFlagSet(H) {
			a = (a - 6) & 0xFF
		}

		if cpu.IsFlagSet(CY) {
			a -= 0x60
		}
	}

	cpu.ResetFlag(H)

	if a&0x100 == 0x100 {
		cpu.SetFlag(CY)
	}

	a &= 0xFF

	if a == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	cpu.Registers.A = byte(a)

	return 1
}

// CPL | 1 | -11- | A=^A
func (cpu *CPU) CPL() (cycles byte) {
	cpu.Registers.A = ^cpu.Registers.A
	cpu.SetFlag(N)
	cpu.SetFlag(H)
	return 1
}

// CCF | 1 | ---- | CY=~CY
func (cpu *CPU) CCF() (cycles byte) {
	if cpu.IsFlagSet(CY) {
		cpu.ResetFlag(CY)
	} else {
		cpu.SetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	return 1
}

// SCF | 1 | ---- | CY=1
func (cpu *CPU) SCF() (cycles byte) {
	cpu.SetFlag(CY)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	return 1
}

// NOP | 1 | ---- | No operation
func (cpu *CPU) NOP() (cycles byte) {
	cpu.PC += cpu.CurrentInstruction.Length
	return 1
}

// NOP | 1 | ---- | Halt until interrupt occurs
func (cpu *CPU) HALT() (cycles byte) {
	cpu.Halt = true
	return 1
}

// DI | 1 | ---- | Disable interrupts, IME=0
func (cpu *CPU) DI() (cycles byte) {
	cpu.IME = false
	return 1
}

// EI | 1 | ---- | Enable interrupts, IME=1
func (cpu *CPU) EI() (cycles byte) {
	cpu.IME = true
	return 1
}

// UTILITIES

func (cpu *CPU) pushWordToStack(word uint16) {
	hb, lb := utils.SplitBytes(word)

	cpu.mmu.WriteByte(cpu.SP-1, hb)
	cpu.mmu.WriteByte(cpu.SP-2, lb)

	cpu.SP -= 2
}

func (cpu *CPU) popWordFromStack() uint16 {
	lb := cpu.mmu.ReadByte(cpu.SP)
	hb := cpu.mmu.ReadByte(cpu.SP + 1)
	cpu.SP += 2

	return utils.JoinBytes(hb, lb)
}

func (cpu *CPU) addBytes(a byte, b byte) byte {
	calculation := a + b

	// Reset subtraction flag
	cpu.ResetFlag(N)

	//If calulation is zero
	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	// Set to 1 when an operation results in carrying from or borrowing to bit 3.
	if ((uint16(a) & 0xF) + (uint16(b) & 0xF)) > 0xF {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	// Set to 1 when an operation results in carrying from or borrowing to bit 7.
	if uint16(a)+uint16(b) > 0xFF {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	return calculation
}

func (cpu *CPU) addBytesWithCarry(a byte, b byte) byte {
	var carry byte = 0
	if cpu.IsFlagSet(CY) {
		carry = 1
	}

	calculation := a + b + carry

	// Reset subtraction flag
	cpu.ResetFlag(N)

	// Set to 1 when an operation results in carrying from or borrowing to bit 3.
	if ((uint16(a) & 0xF) + (uint16(b) & 0xF) + uint16(carry)) > 0xF {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	// Set to 1 when an operation results in carrying from or borrowing to bit 7.
	if (uint16(a) + uint16(b) + uint16(carry)) > 0xFF {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	//If calulation is zero
	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) addHL_rr(rr uint16) uint16 {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)

	calculation := HL + rr

	// Set if there is a carry from bit 15, otherwise reset
	if uint32(HL)+uint32(rr) > 0xFFFF {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	// reset subtraction flag
	cpu.ResetFlag(N)

	// Set if there is a carry from bit 11, otherwise reset
	if ((uint32(HL) & 0xFFF) + (uint32(rr) & 0xFFF)) > 0xFFF {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	return calculation
}

func (cpu *CPU) addSignedByte(a uint16, e int8) uint16 {
	calculation := int32(a) + int32(e)

	flagCheck := uint16(cpu.SP ^ uint16(e) ^ ((cpu.SP + uint16(e)) & 0xffff))

	if (flagCheck & 0x100) == 0x100 {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	if (flagCheck & 0x10) == 0x10 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	cpu.ResetFlag(N)
	cpu.ResetFlag(Z)

	return uint16(calculation)
}

func (cpu *CPU) subBytes(a byte, b byte) byte {
	calculation := a - b

	// Set subtraction flag
	cpu.SetFlag(N)

	// Set to 1 when an operation results in carrying from or borrowing to bit 3.
	if (int(a) & 0xF) < (int(b) & 0xF) {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	// Set to 1 when an operation results in carrying from or borrowing to bit 7.
	if (int(a) & 0xFF) < (int(b) & 0xFF) {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	//If calulation is zero
	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) subBytesWithCarry(a byte, b byte) byte {
	var carry byte = 0x00
	if cpu.IsFlagSet(CY) {
		carry = 0x01
	}

	calculation := a - b - carry

	// Set subtraction flag
	cpu.SetFlag(N)

	// Set to 1 when an operation results in carrying from or borrowing to bit 3.
	if ((int(a) & 0xF) - (int(b) & 0xF) - int(carry)) < 0 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	// Set to 1 when an operation results in carrying from or borrowing to bit 7.
	if (int(a) - int(b) - int(carry)) < 0 {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	//If calulation is zero
	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) incByte(value byte) byte {
	value += 0x01

	// Set to 1 when an operation results in carrying from or borrowing to bit 3.
	if (value & 0x0F) == 0x00 {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	// Reset subtraction flag
	cpu.ResetFlag(N)

	//If calulation is zero
	if value == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return value
}

func (cpu *CPU) decByte(value byte) byte {
	value -= 0x01

	// Set to 1 when an operation results in carrying from or borrowing to bit 3.
	if (value & 0x0F) == 0x0F {
		cpu.SetFlag(H)
	} else {
		cpu.ResetFlag(H)
	}

	// Set subtraction flag
	cpu.SetFlag(N)

	//If calulation is zero
	if value == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return value
}

func (cpu *CPU) rotateLeft(value byte) byte {
	var calculation byte = value << 1

	if utils.IsBitSet(value, 7) {
		cpu.SetFlag(CY)
		calculation = utils.SetBit(calculation, 0)
	} else {
		cpu.ResetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) rotateLeftThroughCarry(value byte) byte {
	var calculation byte = value << 1

	if cpu.IsFlagSet(CY) {
		calculation = utils.SetBit(calculation, 0)
	}

	if utils.IsBitSet(value, 7) {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) rotateRight(value byte) byte {
	var calculation byte = value >> 1

	if utils.IsBitSet(value, 0) {
		cpu.SetFlag(CY)
		calculation = utils.SetBit(calculation, 7)
	} else {
		cpu.ResetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) rotateRightThroughCarry(value byte) byte {
	var calculation byte = value >> 1

	if cpu.IsFlagSet(CY) {
		calculation = utils.SetBit(calculation, 7)
	}

	if utils.IsBitSet(value, 0) {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) shiftLeftArithmetic(value byte) byte {
	var calculation byte = value << 1

	if utils.IsBitSet(value, 7) {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) shiftRightLogical(value byte) byte {
	var calculation byte = value >> 1

	if utils.IsBitSet(value, 0) {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) shiftRightArithmetic(value byte) byte {
	var calculation byte = value >> 1

	if utils.IsBitSet(value, 7) {
		calculation = utils.SetBit(calculation, 7)
	} else {
		calculation = utils.ClearBit(calculation, 7)
	}

	if utils.IsBitSet(value, 0) {
		cpu.SetFlag(CY)
	} else {
		cpu.ResetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}

func (cpu *CPU) swapNibbles(value byte) byte {
	var calculation byte = utils.SwapNibbles(value)

	cpu.ResetFlag(CY)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if calculation == 0 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}

	return calculation
}
