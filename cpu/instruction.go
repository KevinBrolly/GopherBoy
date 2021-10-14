package cpu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type Instruction struct {
	Opcode      byte
	Description string
	Length      uint16
	Execute     func(cpu *CPU)
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
	0x00: {0x00, "NOP", 1, func(cpu *CPU) {
		cpu.NOP()
	}},
	0x01: {0x01, "LD BC,d16", 3, func(cpu *CPU) {
		cpu.LD_rr_nn(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x02: {0x02, "LD (BC), A", 1, func(cpu *CPU) {
		cpu.LD_BC_A()
	}},
	0x03: {0x03, "INC BC", 2, func(cpu *CPU) {
		cpu.INC_rr(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x04: {0x04, "INC B", 1, func(cpu *CPU) {
		cpu.INC_r(&cpu.Registers.B)
	}},
	0x05: {0x05, "DEC B", 1, func(cpu *CPU) {
		cpu.DEC_r(&cpu.Registers.B)
	}},
	0x06: {0x06, "LD B,d8", 2, func(cpu *CPU) {
		cpu.LD_r_n(&cpu.Registers.B)
	}},
	0x07: {0x07, "RLCA", 1, func(cpu *CPU) {
		cpu.RLCA()
	}},
	0x08: {0x08, "LD (a16),SP", 3, func(cpu *CPU) {
		cpu.LD_nn_SP()
	}},
	0x09: {0x09, "ADD HL, BC", 1, func(cpu *CPU) {
		cpu.ADD_HL_rr(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x0A: {0x0A, "LD A,(BC)", 1, func(cpu *CPU) {
		cpu.LD_A_BC()
	}},
	0x0B: {0x0B, "DEC BC", 1, func(cpu *CPU) {
		cpu.DEC_rr(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x0C: {0x0C, "INC C", 1, func(cpu *CPU) {
		cpu.INC_r(&cpu.Registers.C)
	}},
	0x0D: {0x0D, "DEC C", 1, func(cpu *CPU) {
		cpu.DEC_r(&cpu.Registers.C)
	}},
	0x0E: {0x0E, "LD C,d8", 2, func(cpu *CPU) {
		cpu.LD_r_n(&cpu.Registers.C)
	}},
	0x0F: {0x0F, "RRCA", 1, func(cpu *CPU) {
		cpu.RRCA()
	}},
	0x10: {0x10, "STOP 0", 2, func(cpu *CPU) {
		cpu.cycleChannel <- 1
	}},
	0x11: {0x11, "LD DE,d16", 3, func(cpu *CPU) {
		cpu.LD_rr_nn(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x12: {0x12, "LD (DE), A", 1, func(cpu *CPU) {
		cpu.LD_DE_A()
	}},
	0x13: {0x13, "INC DE", 1, func(cpu *CPU) {
		cpu.INC_rr(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x14: {0x14, "INC D", 1, func(cpu *CPU) {
		cpu.INC_r(&cpu.Registers.D)
	}},
	0x15: {0x15, "DEC D", 1, func(cpu *CPU) {
		cpu.DEC_r(&cpu.Registers.D)
	}},
	0x16: {0x16, "LD D,d8", 2, func(cpu *CPU) {
		cpu.LD_r_n(&cpu.Registers.D)
	}},
	0x17: {0x17, "RLA", 1, func(cpu *CPU) {
		cpu.RLA()
	}},
	0x18: {0x18, "JR R8", 2, func(cpu *CPU) {
		cpu.JR_e()
	}},
	0x19: {0x19, "ADD HL,DE", 1, func(cpu *CPU) {
		cpu.ADD_HL_rr(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x1A: {0x1A, "LD A,(DE)", 1, func(cpu *CPU) {
		cpu.LD_A_DE()
	}},
	0x1B: {0x1B, "DEC DE", 1, func(cpu *CPU) {
		cpu.DEC_rr(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x1C: {0x1C, "INC E", 1, func(cpu *CPU) {
		cpu.INC_r(&cpu.Registers.E)
	}},
	0x1D: {0x1D, "DEC E", 1, func(cpu *CPU) {
		cpu.DEC_r(&cpu.Registers.E)
	}},
	0x1E: {0x1E, "LD E,d8", 2, func(cpu *CPU) {
		cpu.LD_r_n(&cpu.Registers.E)
	}},
	0x1F: {0x3E, "RRA", 1, func(cpu *CPU) {
		cpu.RRA()
	}},
	0x20: {0x20, "JR NZ,r8", 2, func(cpu *CPU) {
		cpu.JR_cc_e(CC_NZ)
	}},
	0x21: {0x21, "LD HL,d16", 3, func(cpu *CPU) {
		cpu.LD_rr_nn(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x22: {0x22, "LD (HL+),A", 1, func(cpu *CPU) {
		cpu.LD_HLI_A()
	}},
	0x23: {0x23, "INC HL", 1, func(cpu *CPU) {
		cpu.INC_rr(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x24: {0x24, "INC H", 1, func(cpu *CPU) {
		cpu.INC_r(&cpu.Registers.H)
	}},
	0x25: {0x25, "DEC H", 1, func(cpu *CPU) {
		cpu.DEC_r(&cpu.Registers.H)
	}},
	0x26: {0x26, "LD H,d8", 2, func(cpu *CPU) {
		cpu.LD_r_n(&cpu.Registers.H)
	}},
	0x27: {0x27, "DAA", 1, func(cpu *CPU) {
		cpu.DAA()
	}},
	0x28: {0x30, "JR Z,r8", 2, func(cpu *CPU) {
		cpu.JR_cc_e(CC_Z)
	}},
	0x29: {0x29, "ADD HL,HL", 1, func(cpu *CPU) {
		cpu.ADD_HL_rr(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x2A: {0x2A, "LD A,(HL+)", 1, func(cpu *CPU) {
		cpu.LD_A_HLI()
	}},
	0x2B: {0x2B, "DEC HL", 1, func(cpu *CPU) {
		cpu.DEC_rr(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x2C: {0x2C, "INC L", 1, func(cpu *CPU) {
		cpu.INC_r(&cpu.Registers.L)
	}},
	0x2D: {0x2D, "DEC L", 1, func(cpu *CPU) {
		cpu.DEC_r(&cpu.Registers.L)
	}},
	0x2E: {0x2E, "LD L,d8", 2, func(cpu *CPU) {
		cpu.LD_r_n(&cpu.Registers.L)
	}},
	0x2F: {0x2F, "CPL", 1, func(cpu *CPU) {
		cpu.CPL()
	}},
	0x30: {0x30, "JR NC,r8", 2, func(cpu *CPU) {
		cpu.JR_cc_e(CC_NC)
	}},
	0x31: {0x31, "LD SP,d16", 3, func(cpu *CPU) {
		cpu.LD_SP_nn()
	}},
	0x32: {0x32, "LD (HL-),A", 1, func(cpu *CPU) {
		cpu.LD_HLD_A()
	}},
	0x33: {0x33, "INC SP", 1, func(cpu *CPU) {
		cpu.INC_SP()
	}},
	0x34: {0x34, "INC (HL)", 1, func(cpu *CPU) {
		cpu.INC_HL()
	}},
	0x35: {0x35, "DEC (HL)", 1, func(cpu *CPU) {
		cpu.DEC_HL()
	}},
	0x36: {0x36, "LD (HL),d8", 2, func(cpu *CPU) {
		cpu.LD_HL_n()
	}},
	0x37: {0x37, "SCF", 1, func(cpu *CPU) {
		cpu.SCF()
	}},
	0x38: {0x30, "JR C,r8", 2, func(cpu *CPU) {
		cpu.JR_cc_e(CC_C)
	}},
	0x39: {0x39, "ADD HL,SP", 1, func(cpu *CPU) {
		cpu.ADD_HL_SP()
	}},
	0x3A: {0x3A, "LD A,(HL-)", 1, func(cpu *CPU) {
		cpu.LD_A_HLD()
	}},
	0x3B: {0x3B, "DEC SP", 1, func(cpu *CPU) {
		cpu.DEC_SP()
	}},
	0x3C: {0x3C, "INC A", 1, func(cpu *CPU) {
		cpu.INC_r(&cpu.Registers.A)
	}},
	0x3D: {0x3D, "DEC A", 1, func(cpu *CPU) {
		cpu.DEC_r(&cpu.Registers.A)
	}},
	0x3E: {0x3E, "LD A,d8", 2, func(cpu *CPU) {
		cpu.LD_r_n(&cpu.Registers.A)
	}},
	0x3F: {0x3F, "CCF", 1, func(cpu *CPU) {
		cpu.CCF()
	}},
	0x40: {0x40, "LD B,B", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.B)
	}},
	0x41: {0x41, "LD B,C", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0x42: {0x42, "LD B,D", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.D)
	}},
	0x43: {0x43, "LD B,E", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.E)
	}},
	0x44: {0x44, "LD B,H", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.H)
	}},
	0x45: {0x45, "LD B,L", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.L)
	}},
	0x46: {0x46, "LD B,(HL)", 1, func(cpu *CPU) {
		cpu.LD_r_HL(&cpu.Registers.B)
	}},
	0x47: {0x47, "LD B,A", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.B, &cpu.Registers.A)
	}},
	0x48: {0x48, "LD C,B", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.B)
	}},
	0x49: {0x49, "LD C,C", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.C)
	}},
	0x4A: {0x4A, "LD C,D", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.D)
	}},
	0x4B: {0x4B, "LD C,E", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.E)
	}},
	0x4C: {0x4C, "LD C,H", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.H)
	}},
	0x4D: {0x4D, "LD C,L", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.L)
	}},
	0x4E: {0x4E, "LD C,(HL)", 1, func(cpu *CPU) {
		cpu.LD_r_HL(&cpu.Registers.C)
	}},
	0x4F: {0x4F, "LD C,A", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.C, &cpu.Registers.A)
	}},
	0x50: {0x50, "LD D,B", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.B)
	}},
	0x51: {0x51, "LD D,C", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.C)
	}},
	0x52: {0x52, "LD D,D", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.D)
	}},
	0x53: {0x53, "LD D,E", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0x54: {0x54, "LD D,H", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.H)
	}},
	0x55: {0x55, "LD D,L", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.L)
	}},
	0x56: {0x56, "LD D,(HL)", 1, func(cpu *CPU) {
		cpu.LD_r_HL(&cpu.Registers.D)
	}},
	0x57: {0x57, "LD D,A", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.D, &cpu.Registers.A)
	}},
	0x58: {0x58, "LD E,B", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.B)
	}},
	0x59: {0x59, "LD E,C", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.C)
	}},
	0x5A: {0x5A, "LD E,D", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.D)
	}},
	0x5B: {0x5B, "LD E,E", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.E)
	}},
	0x5C: {0x5C, "LD E,H", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.H)
	}},
	0x5D: {0x5D, "LD E,L", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.L)
	}},
	0x5E: {0x5E, "LD E,(HL)", 1, func(cpu *CPU) {
		cpu.LD_r_HL(&cpu.Registers.E)
	}},
	0x5F: {0x5F, "LD E,A", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.E, &cpu.Registers.A)
	}},
	0x60: {0x60, "LD H,B", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.B)
	}},
	0x61: {0x61, "LD H,C", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.C)
	}},
	0x62: {0x62, "LD H,D", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.D)
	}},
	0x63: {0x63, "LD H,E", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.E)
	}},
	0x64: {0x64, "LD H,H", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.H)
	}},
	0x65: {0x65, "LD H,L", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0x66: {0x66, "LD H,(HL)", 1, func(cpu *CPU) {
		cpu.LD_r_HL(&cpu.Registers.H)
	}},
	0x67: {0x67, "LD H,A", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.H, &cpu.Registers.A)
	}},
	0x68: {0x68, "LD L,B", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.B)
	}},
	0x69: {0x69, "LD L,C", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.C)
	}},
	0x6A: {0x6A, "LD L,D", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.D)
	}},
	0x6B: {0x6B, "LD L,E", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.E)
	}},
	0x6C: {0x6C, "LD L,H", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.H)
	}},
	0x6D: {0x6D, "LD L,L", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.L)
	}},
	0x6E: {0x6E, "LD L,(HL)", 1, func(cpu *CPU) {
		cpu.LD_r_HL(&cpu.Registers.L)
	}},
	0x6F: {0x6F, "LD L,A", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.L, &cpu.Registers.A)
	}},
	0x70: {0x70, "LD (HL),B", 1, func(cpu *CPU) {
		cpu.LD_HL_r(&cpu.Registers.B)
	}},
	0x71: {0x71, "LD (HL),C", 1, func(cpu *CPU) {
		cpu.LD_HL_r(&cpu.Registers.C)
	}},
	0x72: {0x72, "LD (HL),D", 1, func(cpu *CPU) {
		cpu.LD_HL_r(&cpu.Registers.D)
	}},
	0x73: {0x73, "LD (HL),E", 1, func(cpu *CPU) {
		cpu.LD_HL_r(&cpu.Registers.E)
	}},
	0x74: {0x74, "LD (HL),H", 1, func(cpu *CPU) {
		cpu.LD_HL_r(&cpu.Registers.H)
	}},
	0x75: {0x75, "LD (HL),L", 1, func(cpu *CPU) {
		cpu.LD_HL_r(&cpu.Registers.L)
	}},
	0x76: {0x76, "HALT", 1, func(cpu *CPU) {
		cpu.HALT()
	}},
	0x77: {0x77, "LD (HL), A", 1, func(cpu *CPU) {
		cpu.LD_HL_r(&cpu.Registers.A)
	}},
	0x78: {0x78, "LD A,B", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.B)
	}},
	0x79: {0x79, "LD A,C", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.C)
	}},
	0x7A: {0x7A, "LD A,D", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.D)
	}},
	0x7B: {0x7B, "LD A,E", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.E)
	}},
	0x7C: {0x7C, "LD A,H", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.H)
	}},
	0x7D: {0x7D, "LD A,L", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.L)
	}},
	0x7E: {0x7E, "LD A,(HL)", 1, func(cpu *CPU) {
		cpu.LD_r_HL(&cpu.Registers.A)
	}},
	0x7F: {0x7F, "LD A,A", 1, func(cpu *CPU) {
		cpu.LD_r_r(&cpu.Registers.A, &cpu.Registers.A)
	}},
	0x80: {0x80, "ADD A,B", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.Registers.B)
	}},
	0x81: {0x81, "ADD A,C", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.Registers.C)
	}},
	0x82: {0x82, "ADD A,D", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.Registers.D)
	}},
	0x83: {0x83, "ADD A,E", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.Registers.E)
	}},
	0x84: {0x84, "ADD A,H", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.Registers.H)
	}},
	0x85: {0x85, "ADD A,L", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.Registers.L)
	}},
	0x86: {0x86, "ADD A,(HL)", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0x87: {0x87, "ADD A,A", 1, func(cpu *CPU) {
		cpu.ADD_s(cpu.Registers.A)
	}},
	0x88: {0x88, "ADC A,B", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.Registers.B)
	}},
	0x89: {0x89, "ADC A,C", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.Registers.C)
	}},
	0x8A: {0x8A, "ADC A,D", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.Registers.D)
	}},
	0x8B: {0x8B, "ADC A,E", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.Registers.E)
	}},
	0x8C: {0x8C, "ADC A,H", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.Registers.H)
	}},
	0x8D: {0x8D, "ADC A,L", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.Registers.L)
	}},
	0x8E: {0x8E, "ADC A,(HL)", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0x8F: {0x8F, "ADC A,A", 1, func(cpu *CPU) {
		cpu.ADC_A_s(cpu.Registers.A)
	}},
	0x90: {0x90, "SUB A,B", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.Registers.B)
	}},
	0x91: {0x91, "SUB A,C", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.Registers.C)
	}},
	0x92: {0x92, "SUB A,D", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.Registers.D)
	}},
	0x93: {0x93, "SUB A,E", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.Registers.E)
	}},
	0x94: {0x94, "SUB A,H", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.Registers.H)
	}},
	0x95: {0x95, "SUB A,L", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.Registers.L)
	}},
	0x96: {0x96, "SUB A,(HL)", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0x97: {0x97, "SUB A,A", 1, func(cpu *CPU) {
		cpu.SUB_s(cpu.Registers.A)
	}},
	0x98: {0x98, "SBC A,B", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.Registers.B)
	}},
	0x99: {0x99, "SBC A,C", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.Registers.C)
	}},
	0x9A: {0x9A, "SBC A,D", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.Registers.D)
	}},
	0x9B: {0x9B, "SBC A,E", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.Registers.E)
	}},
	0x9C: {0x9C, "SBC A,H", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.Registers.H)
	}},
	0x9D: {0x9D, "SBC A,L", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.Registers.L)
	}},
	0x9E: {0x9E, "SBC A,(HL)", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0x9F: {0x9F, "SBC A,A", 1, func(cpu *CPU) {
		cpu.SBC_s(cpu.Registers.A)
	}},
	0xA0: {0xA0, "AND A,B", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.Registers.B)
	}},
	0xA1: {0xA1, "AND A,C", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.Registers.C)
	}},
	0xA2: {0xA2, "AND A,D", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.Registers.D)
	}},
	0xA3: {0xA3, "AND A,E", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.Registers.E)
	}},
	0xA4: {0xA4, "AND A,H", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.Registers.H)
	}},
	0xA5: {0xA5, "AND A,L", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.Registers.L)
	}},
	0xA6: {0xA6, "AND A,(HL)", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0xA7: {0xA7, "AND A,A", 1, func(cpu *CPU) {
		cpu.AND_s(cpu.Registers.A)
	}},
	0xA8: {0xA8, "XOR A,B", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.Registers.B)
	}},
	0xA9: {0xA9, "XOR A,C", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.Registers.C)
	}},
	0xAA: {0xAA, "XOR A,D", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.Registers.D)
	}},
	0xAB: {0xAB, "XOR A,E", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.Registers.E)
	}},
	0xAC: {0xAC, "XOR A,H", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.Registers.H)
	}},
	0xAD: {0xAD, "XOR A,L", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.Registers.L)
	}},
	0xAE: {0xAE, "XOR A,(HL)", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0xAF: {0xAF, "XOR A,A", 1, func(cpu *CPU) {
		cpu.XOR_s(cpu.Registers.A)
	}},
	0xB0: {0xB0, "OR A,B", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.Registers.B)
	}},
	0xB1: {0xB1, "OR A,C", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.Registers.C)
	}},
	0xB2: {0xB2, "OR A,D", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.Registers.D)
	}},
	0xB3: {0xB3, "OR A,E", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.Registers.E)
	}},
	0xB4: {0xB4, "OR A,H", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.Registers.H)
	}},
	0xB5: {0xB5, "OR A,L", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.Registers.L)
	}},
	0xB6: {0xB6, "OR A,(HL)", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0xB7: {0xB7, "OR A,A", 1, func(cpu *CPU) {
		cpu.OR_s(cpu.Registers.A)
	}},
	0xB8: {0xB8, "CP A,B", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.Registers.B)
	}},
	0xB9: {0xB9, "CP A,C", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.Registers.C)
	}},
	0xBA: {0xBA, "CP A,D", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.Registers.D)
	}},
	0xBB: {0xBB, "CP A,E", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.Registers.E)
	}},
	0xBC: {0xBC, "CP A,H", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.Registers.H)
	}},
	0xBD: {0xBD, "CP A,L", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.Registers.L)
	}},
	0xBE: {0xBE, "CP A,(HL)", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)))
	}},
	0xBF: {0xBF, "CP A,A", 1, func(cpu *CPU) {
		cpu.CP_s(cpu.Registers.A)
	}},
	0xC0: {0xC0, "RET NZ", 1, func(cpu *CPU) {
		cpu.RET_cc(CC_NZ)
	}},
	0xC1: {0xC1, "POP BC", 1, func(cpu *CPU) {
		cpu.POP_qq(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0xC2: {0xC2, "JP NZ, a16", 3, func(cpu *CPU) {
		cpu.JP_cc_nn(CC_NZ)
	}},
	0xC3: {0xC3, "JP d16", 3, func(cpu *CPU) {
		cpu.JP_nn()
	}},
	0xC4: {0xC4, "CALL NZ, a16", 3, func(cpu *CPU) {
		cpu.CALL_cc(CC_NZ)
	}},
	0xC5: {0xC5, "PUSH BC", 1, func(cpu *CPU) {
		cpu.PUSH_qq(&cpu.Registers.B, &cpu.Registers.C)
	}},
	0xC6: {0xC6, "ADD A,D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.ADD_s(d)
	}},
	0xC7: {0xC7, "RST 00H", 1, func(cpu *CPU) {
		cpu.RST(0x00)
	}},
	0xC8: {0xC8, "RET Z", 1, func(cpu *CPU) {
		cpu.RET_cc(CC_Z)
	}},
	0xC9: {0xC9, "RET", 1, func(cpu *CPU) {
		cpu.RET()
	}},
	0xCA: {0xCA, "JP Z, a16", 3, func(cpu *CPU) {
		cpu.JP_cc_nn(CC_Z)
	}},
	// 0xCB handled in main CPU loop
	0xCC: {0xCC, "CALL Z, a16", 3, func(cpu *CPU) {
		cpu.CALL_cc(CC_Z)
	}},
	0xCD: {0xCD, "CALL addr", 3, func(cpu *CPU) {
		cpu.CALL()
	}},
	0xCE: {0xCE, "ADC A,D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.ADC_A_s(d)
	}},
	0xCF: {0xCF, "RST 08H", 1, func(cpu *CPU) {
		cpu.RST(0x08)
	}},
	0xD0: {0xD0, "RET NC", 1, func(cpu *CPU) {
		cpu.RET_cc(CC_NC)
	}},
	0xD1: {0xD1, "POP DE", 1, func(cpu *CPU) {
		cpu.POP_qq(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0xD2: {0xD2, "JP NC,a16", 3, func(cpu *CPU) {
		cpu.JP_cc_nn(CC_NC)
	}},
	0xD4: {0xD4, "CALL NC,a16", 3, func(cpu *CPU) {
		cpu.CALL_cc(CC_NC)
	}},
	0xD5: {0xD5, "PUSH DE", 1, func(cpu *CPU) {
		cpu.PUSH_qq(&cpu.Registers.D, &cpu.Registers.E)
	}},
	0xD6: {0xD6, "SUB A,D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.SUB_s(d)
	}},
	0xD7: {0xD7, "RST 10H", 1, func(cpu *CPU) {
		cpu.RST(0x10)
	}},
	0xD8: {0xD8, "RET C", 1, func(cpu *CPU) {
		cpu.RET_cc(CC_C)
	}},
	0xD9: {0xD9, "RETI", 1, func(cpu *CPU) {
		cpu.RETI()
	}},
	0xDA: {0xDA, "JP C,a16", 3, func(cpu *CPU) {
		cpu.JP_cc_nn(CC_C)
	}},
	0xDC: {0xDC, "CALL C,a16", 3, func(cpu *CPU) {
		cpu.CALL_cc(CC_C)
	}},
	0xDE: {0xDE, "SBC A,D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.SBC_s(d)
	}},
	0xDF: {0xDF, "RST 18H", 1, func(cpu *CPU) {
		cpu.RST(0x18)
	}},
	0xE0: {0xE0, "LDH (a8),A", 2, func(cpu *CPU) {
		cpu.LDH_n_A()
	}},
	0xE1: {0xE1, "POP HL", 1, func(cpu *CPU) {
		cpu.POP_qq(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0xE2: {0xE2, "LD (C),A", 1, func(cpu *CPU) {
		cpu.LD_C_A()
	}},
	0xE5: {0xC5, "PUSH HL", 1, func(cpu *CPU) {
		cpu.PUSH_qq(&cpu.Registers.H, &cpu.Registers.L)
	}},
	0xE6: {0xE6, "AND A,D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.AND_s(d)
	}},
	0xE7: {0xE7, "RST 20H", 1, func(cpu *CPU) {
		cpu.RST(0x20)
	}},
	0xE8: {0xE8, "ADD SP,r8", 2, func(cpu *CPU) {
		cpu.ADD_SP_e()
	}},
	0xE9: {0xE9, "JP (HL)", 1, func(cpu *CPU) {
		cpu.JP_HL()
	}},
	0xEA: {0xEA, "LD (a16),A", 3, func(cpu *CPU) {
		cpu.LD_nn_A()
	}},
	0xEE: {0xEE, "XOR A,D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.XOR_s(d)
	}},
	0xEF: {0xEF, "RST 28H", 1, func(cpu *CPU) {
		cpu.RST(0x28)
	}},
	0xF0: {0xF0, "LDH A,(a8)", 2, func(cpu *CPU) {
		cpu.LDH_A_n()
	}},
	0xF1: {0xF1, "POP AF", 1, func(cpu *CPU) {
		cpu.POP_AF()
	}},
	0xF2: {0xF2, "LD A,(C)", 1, func(cpu *CPU) {
		cpu.LD_A_C()
	}},
	0xF3: {0xF3, "DI", 1, func(cpu *CPU) {
		cpu.DI()
	}},
	0xF5: {0xF5, "PUSH AF", 1, func(cpu *CPU) {
		cpu.PUSH_qq(&cpu.Registers.A, &cpu.Registers.F)
	}},
	0xF6: {0xF6, "OR A,D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.OR_s(d)
	}},
	0xF7: {0xF7, "RST 30H", 1, func(cpu *CPU) {
		cpu.RST(0x30)
	}},
	0xF8: {0xF8, "LD HL,SP+r8", 2, func(cpu *CPU) {
		cpu.LD_HL_SP_e()
	}},
	0xF9: {0xF9, "LD SP,HL", 1, func(cpu *CPU) {
		cpu.LD_SP_HL()
	}},
	0xFA: {0xFA, "LD A,(a16)", 3, func(cpu *CPU) {
		cpu.LD_A_nn()
	}},
	0xFB: {0xFB, "EI", 1, func(cpu *CPU) {
		cpu.EI()
	}},
	0xFE: {0xFE, "CP D8", 2, func(cpu *CPU) {
		cpu.PC++
		d := cpu.mmu.ReadByte((cpu.PC))
		cpu.CP_s(d)
	}},
	0xFF: {0xFF, "RST 38H", 1, func(cpu *CPU) {
		cpu.RST(0x38)
	}},
}

var CBInstructions map[byte]*Instruction = map[byte]*Instruction{
	0x00: {0x00, "RLC B", 2, func(cpu *CPU) {
		cpu.RLC_r(&cpu.Registers.B)
	}},
	0x01: {0x01, "RLC C", 2, func(cpu *CPU) {
		cpu.RLC_r(&cpu.Registers.C)
	}},
	0x02: {0x02, "RLC D", 2, func(cpu *CPU) {
		cpu.RLC_r(&cpu.Registers.D)
	}},
	0x03: {0x03, "RLC E", 2, func(cpu *CPU) {
		cpu.RLC_r(&cpu.Registers.E)
	}},
	0x04: {0x04, "RLC H", 2, func(cpu *CPU) {
		cpu.RLC_r(&cpu.Registers.H)
	}},
	0x05: {0x05, "RLC L", 2, func(cpu *CPU) {
		cpu.RLC_r(&cpu.Registers.L)
	}},
	0x06: {0x06, "RLC (HL)", 2, func(cpu *CPU) {
		cpu.RLC_HL()
	}},
	0x07: {0x07, "RLC A", 2, func(cpu *CPU) {
		cpu.RLC_r(&cpu.Registers.A)
	}},
	0x08: {0x08, "RRC B", 2, func(cpu *CPU) {
		cpu.RRC_r(&cpu.Registers.B)
	}},
	0x09: {0x09, "RRC C", 2, func(cpu *CPU) {
		cpu.RRC_r(&cpu.Registers.C)
	}},
	0x0A: {0x0A, "RRC D", 2, func(cpu *CPU) {
		cpu.RRC_r(&cpu.Registers.D)
	}},
	0x0B: {0x0B, "RRC E", 2, func(cpu *CPU) {
		cpu.RRC_r(&cpu.Registers.E)
	}},
	0x0C: {0x0C, "RRC H", 2, func(cpu *CPU) {
		cpu.RRC_r(&cpu.Registers.H)
	}},
	0x0D: {0x0D, "RRC L", 2, func(cpu *CPU) {
		cpu.RRC_r(&cpu.Registers.L)
	}},
	0x0E: {0x0E, "RRC (HL)", 2, func(cpu *CPU) {
		cpu.RRC_HL()
	}},
	0x0F: {0x0F, "RRC A", 2, func(cpu *CPU) {
		cpu.RRC_r(&cpu.Registers.A)
	}},
	0x10: {0x10, "RL B", 2, func(cpu *CPU) {
		cpu.RL_r(&cpu.Registers.B)
	}},
	0x11: {0x11, "RL C", 2, func(cpu *CPU) {
		cpu.RL_r(&cpu.Registers.C)
	}},
	0x12: {0x12, "RL D", 2, func(cpu *CPU) {
		cpu.RL_r(&cpu.Registers.D)
	}},
	0x13: {0x13, "RL E", 2, func(cpu *CPU) {
		cpu.RL_r(&cpu.Registers.E)
	}},
	0x14: {0x14, "RL H", 2, func(cpu *CPU) {
		cpu.RL_r(&cpu.Registers.H)
	}},
	0x15: {0x15, "RL L", 2, func(cpu *CPU) {
		cpu.RL_r(&cpu.Registers.L)
	}},
	0x16: {0x16, "RL (HL)", 2, func(cpu *CPU) {
		cpu.RL_HL()
	}},
	0x17: {0x17, "RL A", 2, func(cpu *CPU) {
		cpu.RL_r(&cpu.Registers.A)
	}},
	0x18: {0x18, "RR B", 2, func(cpu *CPU) {
		cpu.RR_r(&cpu.Registers.B)
	}},
	0x19: {0x19, "RR C", 2, func(cpu *CPU) {
		cpu.RR_r(&cpu.Registers.C)
	}},
	0x1A: {0x1A, "RR D", 2, func(cpu *CPU) {
		cpu.RR_r(&cpu.Registers.D)
	}},
	0x1B: {0x1B, "RR E", 2, func(cpu *CPU) {
		cpu.RR_r(&cpu.Registers.E)
	}},
	0x1C: {0x1C, "RR H", 2, func(cpu *CPU) {
		cpu.RR_r(&cpu.Registers.H)
	}},
	0x1D: {0x1D, "RR L", 2, func(cpu *CPU) {
		cpu.RR_r(&cpu.Registers.L)
	}},
	0x1E: {0x1E, "RR (HL)", 2, func(cpu *CPU) {
		cpu.RR_HL()
	}},
	0x1F: {0x3E, "RR A", 2, func(cpu *CPU) {
		cpu.RR_r(&cpu.Registers.A)
	}},
	0x20: {0x20, "SLA B", 2, func(cpu *CPU) {
		cpu.SLA_r(&cpu.Registers.B)
	}},
	0x21: {0x21, "SLA C", 2, func(cpu *CPU) {
		cpu.SLA_r(&cpu.Registers.C)
	}},
	0x22: {0x22, "SLA D", 2, func(cpu *CPU) {
		cpu.SLA_r(&cpu.Registers.D)
	}},
	0x23: {0x23, "SLA E", 2, func(cpu *CPU) {
		cpu.SLA_r(&cpu.Registers.E)
	}},
	0x24: {0x24, "SLA H", 2, func(cpu *CPU) {
		cpu.SLA_r(&cpu.Registers.H)
	}},
	0x25: {0x25, "SLA L", 2, func(cpu *CPU) {
		cpu.SLA_r(&cpu.Registers.L)
	}},
	0x26: {0x26, "SLA (HL)", 2, func(cpu *CPU) {
		cpu.SLA_HL()
	}},
	0x27: {0x27, "SLA A", 2, func(cpu *CPU) {
		cpu.SLA_r(&cpu.Registers.A)
	}},
	0x28: {0x30, "SRA B", 2, func(cpu *CPU) {
		cpu.SRA_r(&cpu.Registers.B)
	}},
	0x29: {0x29, "SRA C", 2, func(cpu *CPU) {
		cpu.SRA_r(&cpu.Registers.C)
	}},
	0x2A: {0x2A, "SRA D", 2, func(cpu *CPU) {
		cpu.SRA_r(&cpu.Registers.D)
	}},
	0x2B: {0x2B, "SRA E", 2, func(cpu *CPU) {
		cpu.SRA_r(&cpu.Registers.E)
	}},
	0x2C: {0x2C, "SRA H", 2, func(cpu *CPU) {
		cpu.SRA_r(&cpu.Registers.H)
	}},
	0x2D: {0x2D, "SRA L", 2, func(cpu *CPU) {
		cpu.SRA_r(&cpu.Registers.L)
	}},
	0x2E: {0x2E, "SRA (HL)", 2, func(cpu *CPU) {
		cpu.SRA_HL()
	}},
	0x2F: {0x2F, "SRA A", 2, func(cpu *CPU) {
		cpu.SRA_r(&cpu.Registers.A)
	}},
	0x30: {0x30, "SWAP B", 2, func(cpu *CPU) {
		cpu.SWAP_r(&cpu.Registers.B)
	}},
	0x31: {0x31, "SWAP C", 2, func(cpu *CPU) {
		cpu.SWAP_r(&cpu.Registers.C)
	}},
	0x32: {0x32, "SWAP D", 2, func(cpu *CPU) {
		cpu.SWAP_r(&cpu.Registers.D)
	}},
	0x33: {0x33, "SWAP E", 2, func(cpu *CPU) {
		cpu.SWAP_r(&cpu.Registers.E)
	}},
	0x34: {0x34, "SWAP H", 2, func(cpu *CPU) {
		cpu.SWAP_r(&cpu.Registers.H)
	}},
	0x35: {0x35, "SWAP L", 2, func(cpu *CPU) {
		cpu.SWAP_r(&cpu.Registers.L)
	}},
	0x36: {0x36, "SWAP (HL)", 2, func(cpu *CPU) {
		cpu.SWAP_HL()
	}},
	0x37: {0x37, "SWAP A", 2, func(cpu *CPU) {
		cpu.SWAP_r(&cpu.Registers.A)
	}},
	0x38: {0x30, "SRL B", 2, func(cpu *CPU) {
		cpu.SRL_r(&cpu.Registers.B)
	}},
	0x39: {0x39, "SRL C", 2, func(cpu *CPU) {
		cpu.SRL_r(&cpu.Registers.C)
	}},
	0x3A: {0x3A, "SRL D", 2, func(cpu *CPU) {
		cpu.SRL_r(&cpu.Registers.D)
	}},
	0x3B: {0x3B, "SRL E", 2, func(cpu *CPU) {
		cpu.SRL_r(&cpu.Registers.E)
	}},
	0x3C: {0x3C, "SRL H", 2, func(cpu *CPU) {
		cpu.SRL_r(&cpu.Registers.H)
	}},
	0x3D: {0x3D, "SRL L", 2, func(cpu *CPU) {
		cpu.SRL_r(&cpu.Registers.L)
	}},
	0x3E: {0x3E, "SRL (HL)", 2, func(cpu *CPU) {
		cpu.SRL_HL()
	}},
	0x3F: {0x3F, "SRL A", 2, func(cpu *CPU) {
		cpu.SRL_r(&cpu.Registers.A)
	}},
	0x40: {0x40, "BIT 0,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(0, &cpu.Registers.B)
	}},
	0x41: {0x41, "BIT 0,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(0, &cpu.Registers.C)
	}},
	0x42: {0x42, "BIT 0,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(0, &cpu.Registers.D)
	}},
	0x43: {0x43, "BIT 0,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(0, &cpu.Registers.E)
	}},
	0x44: {0x44, "BIT 0,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(0, &cpu.Registers.H)
	}},
	0x45: {0x45, "BIT 0,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(0, &cpu.Registers.L)
	}},
	0x46: {0x46, "BIT 0,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(0)
	}},
	0x47: {0x47, "BIT 0,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(0, &cpu.Registers.A)
	}},
	0x48: {0x48, "BIT 1,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(1, &cpu.Registers.B)
	}},
	0x49: {0x49, "BIT 1,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(1, &cpu.Registers.C)
	}},
	0x4A: {0x4A, "BIT 1,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(1, &cpu.Registers.D)
	}},
	0x4B: {0x4B, "BIT 1,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(1, &cpu.Registers.E)
	}},
	0x4C: {0x4C, "BIT 1,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(1, &cpu.Registers.H)
	}},
	0x4D: {0x4D, "BIT 1,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(1, &cpu.Registers.L)
	}},
	0x4E: {0x4E, "BIT 1,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(1)
	}},
	0x4F: {0x4F, "BIT 1,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(1, &cpu.Registers.A)
	}},
	0x50: {0x50, "BIT 2,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(2, &cpu.Registers.B)
	}},
	0x51: {0x51, "BIT 2,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(2, &cpu.Registers.C)
	}},
	0x52: {0x52, "BIT 2,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(2, &cpu.Registers.D)
	}},
	0x53: {0x53, "BIT 2,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(2, &cpu.Registers.E)
	}},
	0x54: {0x54, "BIT 2,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(2, &cpu.Registers.H)
	}},
	0x55: {0x55, "BIT 2,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(2, &cpu.Registers.L)
	}},
	0x56: {0x56, "BIT 2,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(2)
	}},
	0x57: {0x57, "BIT 2,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(2, &cpu.Registers.A)
	}},
	0x58: {0x58, "BIT 3,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(3, &cpu.Registers.B)
	}},
	0x59: {0x59, "BIT 3,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(3, &cpu.Registers.C)
	}},
	0x5A: {0x5A, "BIT 3,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(3, &cpu.Registers.D)
	}},
	0x5B: {0x5B, "BIT 3,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(3, &cpu.Registers.E)
	}},
	0x5C: {0x5C, "BIT 3,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(3, &cpu.Registers.H)
	}},
	0x5D: {0x5D, "BIT 3,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(3, &cpu.Registers.L)
	}},
	0x5E: {0x5E, "BIT 3,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(3)
	}},
	0x5F: {0x5F, "BIT 3,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(3, &cpu.Registers.A)
	}},
	0x60: {0x60, "BIT 4,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(4, &cpu.Registers.B)
	}},
	0x61: {0x61, "BIT 4,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(4, &cpu.Registers.C)
	}},
	0x62: {0x62, "BIT 4,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(4, &cpu.Registers.D)
	}},
	0x63: {0x63, "BIT 4,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(4, &cpu.Registers.E)
	}},
	0x64: {0x64, "BIT 4,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(4, &cpu.Registers.H)
	}},
	0x65: {0x65, "BIT 4,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(4, &cpu.Registers.L)
	}},
	0x66: {0x66, "BIT 4,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(4)
	}},
	0x67: {0x67, "BIT 4,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(4, &cpu.Registers.A)
	}},
	0x68: {0x68, "BIT 5,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(5, &cpu.Registers.B)
	}},
	0x69: {0x69, "BIT 5,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(5, &cpu.Registers.C)
	}},
	0x6A: {0x6A, "BIT 45,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(5, &cpu.Registers.D)
	}},
	0x6B: {0x6B, "BIT 5,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(5, &cpu.Registers.E)
	}},
	0x6C: {0x6C, "BIT 5,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(5, &cpu.Registers.H)
	}},
	0x6D: {0x6D, "BIT 5,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(5, &cpu.Registers.L)
	}},
	0x6E: {0x6E, "BIT 5,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(5)
	}},
	0x6F: {0x6F, "BIT 5,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(5, &cpu.Registers.A)
	}},
	0x70: {0x70, "BIT 6,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(6, &cpu.Registers.B)
	}},
	0x71: {0x71, "BIT 6,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(6, &cpu.Registers.C)
	}},
	0x72: {0x72, "BIT 6,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(6, &cpu.Registers.D)
	}},
	0x73: {0x73, "BIT 6,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(6, &cpu.Registers.E)
	}},
	0x74: {0x74, "BIT 6,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(6, &cpu.Registers.H)
	}},
	0x75: {0x75, "BIT 6,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(6, &cpu.Registers.L)
	}},
	0x76: {0x76, "BIT 6,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(6)
	}},
	0x77: {0x77, "BIT 6,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(6, &cpu.Registers.A)
	}},
	0x78: {0x78, "BIT 7,B", 2, func(cpu *CPU) {
		cpu.BIT_b_r(7, &cpu.Registers.B)
	}},
	0x79: {0x79, "BIT 7,C", 2, func(cpu *CPU) {
		cpu.BIT_b_r(7, &cpu.Registers.C)
	}},
	0x7A: {0x7A, "BIT 7,D", 2, func(cpu *CPU) {
		cpu.BIT_b_r(7, &cpu.Registers.D)
	}},
	0x7B: {0x7B, "BIT 7,E", 2, func(cpu *CPU) {
		cpu.BIT_b_r(7, &cpu.Registers.E)
	}},
	0x7C: {0x7C, "BIT 7,H", 2, func(cpu *CPU) {
		cpu.BIT_b_r(7, &cpu.Registers.H)
	}},
	0x7D: {0x7D, "BIT 7,L", 2, func(cpu *CPU) {
		cpu.BIT_b_r(7, &cpu.Registers.L)
	}},
	0x7E: {0x7E, "BIT 7,(HL)", 2, func(cpu *CPU) {
		cpu.BIT_b_HL(7)
	}},
	0x7F: {0x7F, "BIT 7,A", 2, func(cpu *CPU) {
		cpu.BIT_b_r(7, &cpu.Registers.A)
	}},
	0x80: {0x80, "RES 0,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(0, &cpu.Registers.B)
	}},
	0x81: {0x81, "RES 0,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(0, &cpu.Registers.C)
	}},
	0x82: {0x82, "RES 0,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(0, &cpu.Registers.D)
	}},
	0x83: {0x83, "RES 0,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(0, &cpu.Registers.E)
	}},
	0x84: {0x84, "RES 0,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(0, &cpu.Registers.H)
	}},
	0x85: {0x85, "RES 0,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(0, &cpu.Registers.L)
	}},
	0x86: {0x86, "RES 0,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(0)
	}},
	0x87: {0x87, "RES 0,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(0, &cpu.Registers.A)
	}},
	0x88: {0x88, "RES 1,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(1, &cpu.Registers.B)
	}},
	0x89: {0x89, "RES 1,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(1, &cpu.Registers.C)
	}},
	0x8A: {0x8A, "RES 1,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(1, &cpu.Registers.D)
	}},
	0x8B: {0x8B, "RES 1,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(1, &cpu.Registers.E)
	}},
	0x8C: {0x8C, "RES 1,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(1, &cpu.Registers.H)
	}},
	0x8D: {0x8D, "RES 1,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(1, &cpu.Registers.L)
	}},
	0x8E: {0x8E, "RES 1,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(1)
	}},
	0x8F: {0x8F, "RES 1,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(1, &cpu.Registers.A)
	}},
	0x90: {0x90, "RES 2,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(2, &cpu.Registers.B)
	}},
	0x91: {0x91, "RES 2,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(2, &cpu.Registers.C)
	}},
	0x92: {0x92, "RES 2,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(2, &cpu.Registers.D)
	}},
	0x93: {0x93, "RES 2,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(2, &cpu.Registers.E)
	}},
	0x94: {0x94, "RES 2,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(2, &cpu.Registers.H)
	}},
	0x95: {0x95, "RES 2,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(2, &cpu.Registers.L)
	}},
	0x96: {0x96, "RES 2,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(2)
	}},
	0x97: {0x97, "RES 2,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(2, &cpu.Registers.A)
	}},
	0x98: {0x98, "RES 3,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(3, &cpu.Registers.B)
	}},
	0x99: {0x99, "RES 3,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(3, &cpu.Registers.C)
	}},
	0x9A: {0x9A, "RES 3,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(3, &cpu.Registers.D)
	}},
	0x9B: {0x9B, "RES 3,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(3, &cpu.Registers.E)
	}},
	0x9C: {0x9C, "RES 3,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(3, &cpu.Registers.H)
	}},
	0x9D: {0x9D, "RES 3,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(3, &cpu.Registers.L)
	}},
	0x9E: {0x9E, "RES 3,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(3)
	}},
	0x9F: {0x9F, "RES 3,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(3, &cpu.Registers.A)
	}},
	0xA0: {0xA0, "RES 4,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(4, &cpu.Registers.B)
	}},
	0xA1: {0xA1, "RES 4,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(4, &cpu.Registers.C)
	}},
	0xA2: {0xA2, "RES 4,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(4, &cpu.Registers.D)
	}},
	0xA3: {0xA3, "RES 4,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(4, &cpu.Registers.E)
	}},
	0xA4: {0xA4, "RES 4,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(4, &cpu.Registers.H)
	}},
	0xA5: {0xA5, "RES 4,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(4, &cpu.Registers.L)
	}},
	0xA6: {0xA6, "RES 4,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(4)
	}},
	0xA7: {0xA7, "RES 4,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(4, &cpu.Registers.A)
	}},
	0xA8: {0xA8, "RES 5,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(5, &cpu.Registers.B)
	}},
	0xA9: {0xA9, "RES 5,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(5, &cpu.Registers.C)
	}},
	0xAA: {0xAA, "RES 5,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(5, &cpu.Registers.D)
	}},
	0xAB: {0xAB, "RES 5,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(5, &cpu.Registers.E)
	}},
	0xAC: {0xAC, "RES 5,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(5, &cpu.Registers.H)
	}},
	0xAD: {0xAD, "RES 5,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(5, &cpu.Registers.L)
	}},
	0xAE: {0xAE, "RES 5,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(5)
	}},
	0xAF: {0xAF, "RES 5,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(5, &cpu.Registers.A)
	}},
	0xB0: {0xB0, "RES 6,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(6, &cpu.Registers.B)
	}},
	0xB1: {0xB1, "RES 6,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(6, &cpu.Registers.C)
	}},
	0xB2: {0xB2, "RES 6,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(6, &cpu.Registers.D)
	}},
	0xB3: {0xB3, "RES 6,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(6, &cpu.Registers.E)
	}},
	0xB4: {0xB4, "RES 6,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(6, &cpu.Registers.H)
	}},
	0xB5: {0xB5, "RES 6,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(6, &cpu.Registers.L)
	}},
	0xB6: {0xB6, "RES 6,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(6)
	}},
	0xB7: {0xB7, "RES 6,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(6, &cpu.Registers.A)
	}},
	0xB8: {0xB8, "RES 7,B", 2, func(cpu *CPU) {
		cpu.RES_b_r(7, &cpu.Registers.B)
	}},
	0xB9: {0xB9, "RES 7,C", 2, func(cpu *CPU) {
		cpu.RES_b_r(7, &cpu.Registers.C)
	}},
	0xBA: {0xBA, "RES 7,D", 2, func(cpu *CPU) {
		cpu.RES_b_r(7, &cpu.Registers.D)
	}},
	0xBB: {0xBB, "RES 7,E", 2, func(cpu *CPU) {
		cpu.RES_b_r(7, &cpu.Registers.E)
	}},
	0xBC: {0xBC, "RES 7,H", 2, func(cpu *CPU) {
		cpu.RES_b_r(7, &cpu.Registers.H)
	}},
	0xBD: {0xBD, "RES 7,L", 2, func(cpu *CPU) {
		cpu.RES_b_r(7, &cpu.Registers.L)
	}},
	0xBE: {0xBE, "RES 7,(HL)", 2, func(cpu *CPU) {
		cpu.RES_b_HL(7)
	}},
	0xBF: {0xBF, "RES 7,A", 2, func(cpu *CPU) {
		cpu.RES_b_r(7, &cpu.Registers.A)
	}},
	0xC0: {0xC0, "SET 0,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(0, &cpu.Registers.B)
	}},
	0xC1: {0xC1, "SET 0,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(0, &cpu.Registers.C)
	}},
	0xC2: {0xC2, "SET 0,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(0, &cpu.Registers.D)
	}},
	0xC3: {0xC3, "SET 0,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(0, &cpu.Registers.E)
	}},
	0xC4: {0xC4, "SET 0,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(0, &cpu.Registers.H)
	}},
	0xC5: {0xC5, "SET 0,L", 2, func(cpu *CPU) {
		cpu.SET_b_r(0, &cpu.Registers.L)
	}},
	0xC6: {0xC6, "SET 0,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(0)
	}},
	0xC7: {0xC7, "SET 0,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(0, &cpu.Registers.A)
	}},
	0xC8: {0xC8, "SET 1,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(1, &cpu.Registers.B)
	}},
	0xC9: {0xC9, "SET 1,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(1, &cpu.Registers.C)
	}},
	0xCA: {0xCA, "SET 1,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(1, &cpu.Registers.D)
	}},
	0xCB: {0xCB, "SET 1,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(1, &cpu.Registers.E)
	}},
	0xCC: {0xCC, "SET 1,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(1, &cpu.Registers.H)
	}},
	0xCD: {0xCD, "SET 1,L", 2, func(cpu *CPU) {
		cpu.SET_b_r(1, &cpu.Registers.L)
	}},
	0xCE: {0xCE, "SET 1,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(1)
	}},
	0xCF: {0xCF, "SET 1,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(1, &cpu.Registers.A)
	}},
	0xD0: {0xD0, "SET 2,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(2, &cpu.Registers.B)
	}},
	0xD1: {0xD1, "SET 2,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(2, &cpu.Registers.C)
	}},
	0xD2: {0xD2, "SET 2,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(2, &cpu.Registers.D)
	}},
	0xD3: {0xD3, "SET 2,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(2, &cpu.Registers.E)
	}},
	0xD4: {0xD4, "SET 2,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(2, &cpu.Registers.H)
	}},
	0xD5: {0xD5, "SET 2,L", 2, func(cpu *CPU) {
		cpu.SET_b_r(2, &cpu.Registers.L)
	}},
	0xD6: {0xD6, "SET 2,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(2)
	}},
	0xD7: {0xD7, "SET 2,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(2, &cpu.Registers.A)
	}},
	0xD8: {0xD8, "SET 3,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(3, &cpu.Registers.B)
	}},
	0xD9: {0xD9, "SET 3,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(3, &cpu.Registers.C)
	}},
	0xDA: {0xDA, "SET 3,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(3, &cpu.Registers.D)
	}},
	0xDB: {0xDB, "SET 3,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(3, &cpu.Registers.E)
	}},
	0xDC: {0xDC, "SET 3,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(3, &cpu.Registers.H)
	}},
	0xDD: {0xDD, "SET 3,L", 2, func(cpu *CPU) {
		cpu.SET_b_r(3, &cpu.Registers.L)
	}},
	0xDE: {0xDE, "SET 3,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(3)
	}},
	0xDF: {0xDF, "SET 3,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(3, &cpu.Registers.A)
	}},
	0xE0: {0xE0, "SET 4,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(4, &cpu.Registers.B)
	}},
	0xE1: {0xE1, "SET 4,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(4, &cpu.Registers.C)
	}},
	0xE2: {0xE2, "SET 4,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(4, &cpu.Registers.D)
	}},
	0xE3: {0xE3, "SET 4,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(4, &cpu.Registers.E)
	}},
	0xE4: {0xE4, "SET 4,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(4, &cpu.Registers.H)
	}},
	0xE5: {0xE5, "SET 4,L", 2, func(cpu *CPU) {
		cpu.SET_b_r(4, &cpu.Registers.L)
	}},
	0xE6: {0xE6, "SET 4,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(4)
	}},
	0xE7: {0xE7, "SET 4,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(4, &cpu.Registers.A)
	}},
	0xE8: {0xE8, "SET 5,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(5, &cpu.Registers.B)
	}},
	0xE9: {0xE9, "SET 5,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(5, &cpu.Registers.C)
	}},
	0xEA: {0xEA, "SET 5,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(5, &cpu.Registers.D)
	}},
	0xEB: {0xEB, "SET 5,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(5, &cpu.Registers.E)
	}},
	0xEC: {0xEC, "SET 5,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(5, &cpu.Registers.H)
	}},
	0xED: {0xED, "SET 5,L", 2, func(cpu *CPU) {
		cpu.SET_b_r(5, &cpu.Registers.L)
	}},
	0xEE: {0xEE, "SET 5,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(5)
	}},
	0xEF: {0xEF, "SET 5,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(5, &cpu.Registers.A)
	}},
	0xF0: {0xF0, "SET 6,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(6, &cpu.Registers.B)
	}},
	0xF1: {0xF1, "SET 6,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(6, &cpu.Registers.C)
	}},
	0xF2: {0xF2, "SET 6,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(6, &cpu.Registers.D)
	}},
	0xF3: {0xF3, "SET 6,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(6, &cpu.Registers.E)
	}},
	0xF4: {0xF4, "SET 6,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(6, &cpu.Registers.H)
	}},
	0xF5: {0xF5, "SET 6,L", 2, func(cpu *CPU) {
		cpu.SET_b_r(6, &cpu.Registers.L)
	}},
	0xF6: {0xF6, "SET 6,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(6)
	}},
	0xF7: {0xF7, "SET 6,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(6, &cpu.Registers.A)
	}},
	0xF8: {0xF8, "SET 7,B", 2, func(cpu *CPU) {
		cpu.SET_b_r(7, &cpu.Registers.B)
	}},
	0xF9: {0xF9, "SET 7,C", 2, func(cpu *CPU) {
		cpu.SET_b_r(7, &cpu.Registers.C)
	}},
	0xFA: {0xFA, "SET 7,D", 2, func(cpu *CPU) {
		cpu.SET_b_r(7, &cpu.Registers.D)
	}},
	0xFB: {0xFB, "SET 7,E", 2, func(cpu *CPU) {
		cpu.SET_b_r(7, &cpu.Registers.E)
	}},
	0xFC: {0xFC, "SET 7,H", 2, func(cpu *CPU) {
		cpu.SET_b_r(7, &cpu.Registers.H)
	}},
	0xFD: {0xFD, "SET 7,L", 2, func(cpu *CPU) {

		cpu.SET_b_r(7, &cpu.Registers.L)
	}},
	0xFE: {0xFE, "SET 7,(HL)", 2, func(cpu *CPU) {
		cpu.SET_b_HL(7)
	}},
	0xFF: {0xFF, "SET 7,A", 2, func(cpu *CPU) {
		cpu.SET_b_r(7, &cpu.Registers.A)
	}},
}

// 8-Bit Transfer/Input-Output Instructions

// LD r,r | 1 | ---- | r=r
func (cpu *CPU) LD_r_r(register_1 *byte, register_2 *byte) {
	*register_1 = *register_2
}

// LD r,n | 2 | ---- | r=n
func (cpu *CPU) LD_r_n(register *byte) {
	cpu.PC++
	*register = cpu.mmu.ReadByte(cpu.PC)
}

// LD r,(HL) | 2 | ---- | r=(HL)
func (cpu *CPU) LD_r_HL(register *byte) {
	*register = cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L))
}

// LD (HL),r | 2 | ---- | (HL)=r
func (cpu *CPU) LD_HL_r(register *byte) {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L), *register)
}

// LD (HL),n | 3 | ---- | (HL)=n
func (cpu *CPU) LD_HL_n() {
	cpu.PC++
	n := cpu.mmu.ReadByte(cpu.PC)
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.H, cpu.Registers.L), n)
}

// LD A,(BC) | 2 | ---- | A=(BC)
func (cpu *CPU) LD_A_BC() {
	cpu.Registers.A = cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.B, cpu.Registers.C))
}

// LD A,(DE) | 2 | ---- | A=(DE)
func (cpu *CPU) LD_A_DE() {
	cpu.Registers.A = cpu.mmu.ReadByte(utils.JoinBytes(cpu.Registers.D, cpu.Registers.E))
}

// LD A,(C) | 2 | ---- | A=(0xFF00+C)
func (cpu *CPU) LD_A_C() {
	cpu.Registers.A = cpu.mmu.ReadByte(uint16(0xFF00 + uint16(cpu.Registers.C)))
}

// LD (C),A | 2 | ---- | (0xFF00+C)=A
func (cpu *CPU) LD_C_A() {
	cpu.mmu.WriteByte(uint16(0xFF00+uint16(cpu.Registers.C)), cpu.Registers.A)
}

// LDH A,(n) | 3 | ---- | A=(n)
func (cpu *CPU) LDH_A_n() {
	cpu.PC++
	n := cpu.mmu.ReadByte(cpu.PC)
	cpu.Registers.A = cpu.mmu.ReadByte(uint16(0xFF00 + uint16(n)))
}

// LDH (n),A | 3 | ---- | (n)=A
func (cpu *CPU) LDH_n_A() {
	cpu.PC++
	n := cpu.mmu.ReadByte(cpu.PC)
	cpu.mmu.WriteByte(uint16(0xFF00+uint16(n)), cpu.Registers.A)
}

// LD A,(nn) | 4 | ---- | A=(nn)
func (cpu *CPU) LD_A_nn() {
	cpu.PC++
	lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	hb := cpu.mmu.ReadByte(cpu.PC)
	cpu.Registers.A = cpu.mmu.ReadByte(utils.JoinBytes(hb, lb))
}

// LD (nn),A | 4 | ---- | (nn)=A
func (cpu *CPU) LD_nn_A() {
	cpu.PC++
	lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	hb := cpu.mmu.ReadByte(cpu.PC)
	cpu.mmu.WriteByte(utils.JoinBytes(hb, lb), cpu.Registers.A)
}

// LD A,(HLI) | 2 | ---- | A=(HL) HL=HL+1
func (cpu *CPU) LD_A_HLI() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.Registers.A = cpu.mmu.ReadByte(HL)
	HL += 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
}

// LD A,(HLD) | 2 | ---- | A=(HL) HL=HL-1
func (cpu *CPU) LD_A_HLD() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.Registers.A = cpu.mmu.ReadByte(HL)
	HL -= 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
}

// LD (BC),A | 2 | ---- | (BC)=A
func (cpu *CPU) LD_BC_A() {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.B, cpu.Registers.C), cpu.Registers.A)
}

// LD (DE),A | 2 | ---- | (DE)=A
func (cpu *CPU) LD_DE_A() {
	cpu.mmu.WriteByte(utils.JoinBytes(cpu.Registers.D, cpu.Registers.E), cpu.Registers.A)
}

// LD (HLI),A | 2 | ---- | (HL)=A HL=HL+1
func (cpu *CPU) LD_HLI_A() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.mmu.WriteByte(HL, cpu.Registers.A)
	HL += 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
}

// LD (HLD),A | 2 | ---- | (HL)=A HL=HL-1
func (cpu *CPU) LD_HLD_A() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.mmu.WriteByte(HL, cpu.Registers.A)
	HL -= 1
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(HL)
}

// 16-Bit Transfer Instructions

// LD rr,nn | 3 | ---- | rr=nn
func (cpu *CPU) LD_rr_nn(r1 *byte, r2 *byte) {
	cpu.PC++
	*r2 = cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	*r1 = cpu.mmu.ReadByte(cpu.PC)
}

// LD SP,nn | 3 | ---- | SP=nn
func (cpu *CPU) LD_SP_nn() {
	cpu.PC++
	lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	hb := cpu.mmu.ReadByte(cpu.PC)
	cpu.SP = utils.JoinBytes(hb, lb)
}

// LD SP,HL | 2 | ---- | SP=HL
func (cpu *CPU) LD_SP_HL() {
	cpu.SP = utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	cpu.cycleChannel <- 1
}

// PUSH qq | 4 | ---- | (SP-1)=qqH (SP-2)=qqL SP=SP-2
func (cpu *CPU) PUSH_qq(r1 *byte, r2 *byte) {
	cpu.cycleChannel <- 1
	cpu.pushWordToStack(utils.JoinBytes(*r1, *r2))
}

// POP qq | 3 | ---- | qqL=(SP) qqH=(SP+1) SP=SP+2
func (cpu *CPU) POP_qq(r1 *byte, r2 *byte) {
	*r1, *r2 = utils.SplitBytes(cpu.popWordFromStack())
}

// POP AF | 3 | **** | qqL=(SP) qqH=(SP+1) SP=SP+2
func (cpu *CPU) POP_AF() {
	cpu.Registers.A, cpu.Registers.F = utils.SplitBytes(cpu.popWordFromStack())
	// Since F only holds 4 bits/flags, we mask to ensure only bits 4-7 are set
	cpu.Registers.F &= 0xF0
}

// LDHL SP,e | 3 | **00 | HL=SP+e
func (cpu *CPU) LD_HL_SP_e() {
	cpu.PC++
	n := cpu.mmu.ReadByte(cpu.PC)
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(cpu.addSignedByte(cpu.SP, int8(n)))
	cpu.cycleChannel <- 1
}

// LD (nn),SP | 5 | ---- | (nn)=SPL (nn+1)==SPH
func (cpu *CPU) LD_nn_SP() {
	hb, lb := utils.SplitBytes(cpu.SP)

	cpu.PC++
	nn_lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	nn_hb := cpu.mmu.ReadByte(cpu.PC)

	nn := utils.JoinBytes(nn_hb, nn_lb)
	cpu.mmu.WriteByte(nn, lb)
	cpu.mmu.WriteByte(nn+1, hb)
}

// 8-Bit Arithmetic and Logical Operation Instructions

// ADD s | 1,2 | **0* | A=A+s
func (cpu *CPU) ADD_s(s byte) {
	cpu.Registers.A = cpu.addBytes(cpu.Registers.A, s)
}

// ADC A,s | 1,2 | **0* | A=A+s+CY
func (cpu *CPU) ADC_A_s(s byte) {
	cpu.Registers.A = cpu.addBytesWithCarry(cpu.Registers.A, s)
}

// SUB s | 1,2 | **1* | A=A-s
func (cpu *CPU) SUB_s(s byte) {
	cpu.Registers.A = cpu.subBytes(cpu.Registers.A, s)
}

// SBC A,s | 1,2 | **1* | A=A-s-CY
func (cpu *CPU) SBC_s(s byte) {
	cpu.Registers.A = cpu.subBytesWithCarry(cpu.Registers.A, s)
}

// AND s | 1,2 | 010* | A=A&s
func (cpu *CPU) AND_s(s byte) {
	cpu.Registers.A = cpu.Registers.A & s

	cpu.ResetFlag(CY)
	cpu.SetFlag(H)
	cpu.ResetFlag(N)

	if cpu.Registers.A == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}
}

// OR s | 1,2 | 000* | A=A|s
func (cpu *CPU) OR_s(s byte) {
	cpu.Registers.A = cpu.Registers.A | s

	cpu.ResetFlag(CY)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if cpu.Registers.A == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}
}

// XOR s | 1,2 | 000* | A=A^s
func (cpu *CPU) XOR_s(s byte) {
	cpu.Registers.A = cpu.Registers.A ^ s

	cpu.ResetFlag(CY)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)

	if cpu.Registers.A == 0x00 {
		cpu.SetFlag(Z)
	} else {
		cpu.ResetFlag(Z)
	}
}

//CP s | 1,2 | **1* | A-s
func (cpu *CPU) CP_s(s byte) {
	cpu.subBytes(cpu.Registers.A, s)
}

// INC r | 1 | -*0* | r=r+1
func (cpu *CPU) INC_r(r *byte) {
	*r = cpu.incByte(*r)
}

// INC (HL) | 3 | -*0* | (HL)=(HL)+1
func (cpu *CPU) INC_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	cpu.mmu.WriteByte(HL, cpu.incByte(value))
}

// DEC r | 1 | -*1* | r=r-1
func (cpu *CPU) DEC_r(r *byte) {
	*r = cpu.decByte(*r)
}

// DEC (HL) | 3 | -*1* | (HL)=(HL)-1
func (cpu *CPU) DEC_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	cpu.mmu.WriteByte(HL, cpu.decByte(value))
}

// 16-Bit Arithmetic and Logical Operation Instructions

// ADD HL,rr | 2 | **0- | HL=HL+rr
func (cpu *CPU) ADD_HL_rr(r1 *byte, r2 *byte) {
	ss := utils.JoinBytes(*r1, *r2)

	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(cpu.addHL_rr(ss))
	cpu.cycleChannel <- 1
}

// ADD HL,SP | 2 | **0- | HL=HL+SP
func (cpu *CPU) ADD_HL_SP() {
	cpu.Registers.H, cpu.Registers.L = utils.SplitBytes(cpu.addHL_rr(cpu.SP))
	cpu.cycleChannel <- 1
}

// ADD SP,e | 4 | **00 | SP=SP+e
func (cpu *CPU) ADD_SP_e() {
	cpu.PC++
	e := cpu.mmu.ReadByte(cpu.PC)
	cpu.SP = cpu.addSignedByte(cpu.SP, int8(e))
	cpu.cycleChannel <- 1
}

// INC rr | 2 | ---- | rr=rr+1
func (cpu *CPU) INC_rr(r1 *byte, r2 *byte) {
	var ss uint16 = utils.JoinBytes(*r1, *r2)
	ss += 0x01
	*r1, *r2 = utils.SplitBytes(ss)
	cpu.cycleChannel <- 1
}

// DEC rr | 2 | ---- | rr=rr-1
func (cpu *CPU) DEC_rr(r1 *byte, r2 *byte) {
	var rr uint16 = utils.JoinBytes(*r1, *r2)
	rr -= 0x01
	*r1, *r2 = utils.SplitBytes(rr)
	cpu.cycleChannel <- 1
}

// INC SP | 2 | ---- | SP=SP+1
func (cpu *CPU) INC_SP() {
	cpu.SP += 1
	cpu.cycleChannel <- 1
}

// DEC SP | 2 | ---- | SP=SP-1
func (cpu *CPU) DEC_SP() {
	cpu.SP -= 1
	cpu.cycleChannel <- 1
}

// Rotate Shift Instructions

// RLCA | 1 | *000 | A<<1 A0=A7 CY=A7
func (cpu *CPU) RLCA() {
	cpu.Registers.A = cpu.rotateLeft(cpu.Registers.A)
	// The rotateLeft function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
}

// RLA | 1 | *000 | A<<1 A0=CY CY=A7
func (cpu *CPU) RLA() {
	cpu.Registers.A = cpu.rotateLeftThroughCarry(cpu.Registers.A)
	// The rotateLeftThroughCarry function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
}

// RRCA | 1 | *000 | A>>1 A7=A0 CY=A0
func (cpu *CPU) RRCA() {
	cpu.Registers.A = cpu.rotateRight(cpu.Registers.A)
	// The rotateRight function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
}

// RRA | 1 | *000 | A>>1 A7=CY CY=A0
func (cpu *CPU) RRA() {
	cpu.Registers.A = cpu.rotateRightThroughCarry(cpu.Registers.A)
	// The rotateRightThroughCarry function sets the Z flag, but it should always be reset for this instruction
	cpu.ResetFlag(Z)
}

// RLC r | 2 | *00* | r<<1 r0=r7 CY=r7
func (cpu *CPU) RLC_r(r *byte) {
	*r = cpu.rotateLeft(*r)
	cpu.cycleChannel <- 1
}

// RLC (HL) | 4 | *00* | (HL)<<1 (HL)0=(HL)7 CY=(HL)7
func (cpu *CPU) RLC_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	cpu.cycleChannel <- 1
	value = cpu.rotateLeft(value)
	cpu.mmu.WriteByte(HL, value)
}

// RL r | 2 | *00* | r<<1 r0=CY CY=r7
func (cpu *CPU) RL_r(r *byte) {
	*r = cpu.rotateLeftThroughCarry(*r)
	cpu.cycleChannel <- 1
}

// RL (HL) | 4 | *00* | (HL)<<1 (HL)0=CY CY=(HL)7
func (cpu *CPU) RL_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.rotateLeftThroughCarry(value)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, value)
}

// RRC r | 2 | *00* | r>>1 r7=r0 CY=r0
func (cpu *CPU) RRC_r(r *byte) {
	*r = cpu.rotateRight(*r)
	cpu.cycleChannel <- 1
}

// RRC (HL) | 4 | *00* | (HL)>>1 (HL)7=(HL)0 CY=(HL)0
func (cpu *CPU) RRC_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.rotateRight(value)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, value)
}

// RR r | 2 | *00* | r>>1 r7=CY CY=r0
func (cpu *CPU) RR_r(r *byte) {
	*r = cpu.rotateRightThroughCarry(*r)
	cpu.cycleChannel <- 1
}

// RR (HL) | 4 | *00* | (HL)>>1 (HL)7=CY CY=(HL)0
func (cpu *CPU) RR_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.rotateRightThroughCarry(value)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, value)
}

// SLA r | 2 | *00* | r<<1 CY=r7
func (cpu *CPU) SLA_r(r *byte) {
	*r = cpu.shiftLeftArithmetic(*r)
	cpu.cycleChannel <- 1
}

// SLA (HL) | 4 | *00* | (HL)<<1 CY=(HL)7
func (cpu *CPU) SLA_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.shiftLeftArithmetic(value)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, value)
}

// SRA r | 2 | *00* | r>>1 r7=r7 CY=r0
func (cpu *CPU) SRA_r(r *byte) {
	*r = cpu.shiftRightArithmetic(*r)
	cpu.cycleChannel <- 1
}

// SRA (HL) | 4 | *00* | (HL)>>1 (HL)7=(HL)7 CY=(HL)0
func (cpu *CPU) SRA_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.shiftRightArithmetic(value)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, value)
}

// SRL r | 2 | *00* | r>>1 CY=r0
func (cpu *CPU) SRL_r(r *byte) {
	*r = cpu.shiftRightLogical(*r)
	cpu.cycleChannel <- 1
}

// SRL (HL) | 4 | *00* | (HL)>>1 CY=(HL)0
func (cpu *CPU) SRL_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.shiftRightLogical(value)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, value)
}

// SWAP r | 2 | 000* | r=r[4:7]&r[0:3]
func (cpu *CPU) SWAP_r(r *byte) {
	*r = cpu.swapNibbles(*r)
	cpu.cycleChannel <- 1
}

// SWAP (HL) | 4 | 000* | (HL)=(HL)[4:7]&(HL)[0:3]
func (cpu *CPU) SWAP_HL() {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	value := cpu.mmu.ReadByte(HL)
	value = cpu.swapNibbles(value)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, value)
}

// Bit Operations

// BIT b,r | 2 | -10* | Z=~rb
func (cpu *CPU) BIT_b_r(bit byte, r *byte) {
	cpu.SetFlag(H)
	cpu.ResetFlag(N)

	if utils.IsBitSet(*r, bit) {
		cpu.ResetFlag(Z)
	} else {
		cpu.SetFlag(Z)
	}
	cpu.cycleChannel <- 1
}

// BIT b,(HL) | 3 | -10* | Z=^(HL)b
func (cpu *CPU) BIT_b_HL(bit byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
	var value byte = cpu.mmu.ReadByte(HL)

	cpu.SetFlag(H)
	cpu.ResetFlag(N)

	if utils.IsBitSet(value, bit) {
		cpu.ResetFlag(Z)
	} else {
		cpu.SetFlag(Z)
	}

	cpu.cycleChannel <- 1
}

// SET b,r | 2 | ---- | rb=1
func (cpu *CPU) SET_b_r(bit byte, r *byte) {
	*r = utils.SetBit(*r, bit)
	cpu.cycleChannel <- 1
}

// SET b,(HL) | 4 | ---- | (HL)b=1
func (cpu *CPU) SET_b_HL(bit byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)

	var value byte = cpu.mmu.ReadByte(HL)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, utils.SetBit(value, bit))
}

// RES b,r | 2 | ---- | rb=0
func (cpu *CPU) RES_b_r(bit byte, r *byte) {
	*r = utils.ClearBit(*r, bit)
	cpu.cycleChannel <- 1
}

// RES b,(HL) | 4 | ---- | (HL)b=0
func (cpu *CPU) RES_b_HL(bit byte) {
	HL := utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)

	var value byte = cpu.mmu.ReadByte(HL)
	cpu.cycleChannel <- 1
	cpu.mmu.WriteByte(HL, utils.ClearBit(value, bit))
}

// Jump Instructions

// JP nn | 4 | ---- | PC=nn
func (cpu *CPU) JP_nn() {
	cpu.PC++
	lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	hb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC = utils.JoinBytes(hb, lb)
	cpu.cycleChannel <- 1
}

// JP cc,nn | 4,3 | ---- | if cc true, PC=nn
func (cpu *CPU) JP_cc_nn(conditionCode int) {
	cpu.PC++
	lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	hb := cpu.mmu.ReadByte(cpu.PC)

	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		cpu.PC = utils.JoinBytes(hb, lb)
		cpu.cycleChannel <- 1
	}
}

// JR e | 3 | ---- | PC=PC+e
func (cpu *CPU) JR_e() {
	cpu.PC++
	e := cpu.mmu.ReadByte(cpu.PC)

	// e is signed, if it is more than 127 then it is negative
	if e > 127 {
		cpu.PC -= uint16(-e)
	} else {
		cpu.PC += uint16(e)
	}

	cpu.cycleChannel <- 1
}

// JR cc,e | 3/2 | ---- | if cc true, PC=PC+e
func (cpu *CPU) JR_cc_e(conditionCode int) {
	cpu.PC++
	e := cpu.mmu.ReadByte(cpu.PC)

	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		if e > 127 {
			cpu.PC -= uint16(-e)
		} else {
			cpu.PC += uint16(e)
		}

		cpu.cycleChannel <- 1
	}
}

// JP HL | 1 | ---- | PC=HL
func (cpu *CPU) JP_HL() {
	cpu.PC = utils.JoinBytes(cpu.Registers.H, cpu.Registers.L)
}

// Call/Return Instructions

// CALL nn | 6 | ---- | (SP-1)=PCH (SP-2)=PCL PC=nn SP=SP-2
func (cpu *CPU) CALL() {
	cpu.PC++
	lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	hb := cpu.mmu.ReadByte(cpu.PC)

	cpu.cycleChannel <- 1
	cpu.pushWordToStack(cpu.PC)

	cpu.PC = utils.JoinBytes(hb, lb)
}

// CALL cc,nn | 6/3 | ---- | (SP-1)=PCH (SP-2)=PCL PC=nn SP=SP-2
func (cpu *CPU) CALL_cc(conditionCode int) {
	cpu.PC++
	lb := cpu.mmu.ReadByte(cpu.PC)
	cpu.PC++
	hb := cpu.mmu.ReadByte(cpu.PC)

	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		cpu.cycleChannel <- 1
		cpu.pushWordToStack(cpu.PC)

		cpu.PC = utils.JoinBytes(hb, lb)
	}
}

// RET | 4 | ---- | PCL=(SP) PCH=(SP+1) SP=SP+2
func (cpu *CPU) RET() {
	cpu.cycleChannel <- 1
	cpu.PC = cpu.popWordFromStack()
}

// RETI | 4 | ---- | PCL=(SP) PCH=(SP+1) SP=SP+2 IME=true
func (cpu *CPU) RETI() {
	cpu.RET()
	cpu.IME = true
}

// RET cc | 5/2 | ---- | if cc true, PCL=(SP) PCH=(SP+1) SP=SP+2
func (cpu *CPU) RET_cc(conditionCode int) {
	cpu.cycleChannel <- 1
	if ((conditionCode == CC_NZ) && !cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_Z) && cpu.IsFlagSet(Z)) ||
		((conditionCode == CC_NC) && !cpu.IsFlagSet(CY)) ||
		((conditionCode == CC_C) && cpu.IsFlagSet(CY)) {

		cpu.RET()
	}
}

// RST t | 4 | ---- | (SP-1)=PCH (SP-2)=PCL SP=SP-2 PCH=0 PCL=t
func (cpu *CPU) RST(t byte) {
	cpu.pushWordToStack(cpu.PC + 1)
	cpu.PC = uint16(t)
	cpu.cycleChannel <- 1
}

// General-Purpose Arithmetic Operations and CPU Control Instructions

// DAA | 1 | z-0x | Decimal Adjust acc
func (cpu *CPU) DAA() {
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
}

// CPL | 1 | -11- | A=^A
func (cpu *CPU) CPL() {
	cpu.Registers.A = ^cpu.Registers.A
	cpu.SetFlag(N)
	cpu.SetFlag(H)
}

// CCF | 1 | ---- | CY=~CY
func (cpu *CPU) CCF() {
	if cpu.IsFlagSet(CY) {
		cpu.ResetFlag(CY)
	} else {
		cpu.SetFlag(CY)
	}

	cpu.ResetFlag(H)
	cpu.ResetFlag(N)
}

// SCF | 1 | ---- | CY=1
func (cpu *CPU) SCF() {
	cpu.SetFlag(CY)
	cpu.ResetFlag(H)
	cpu.ResetFlag(N)
}

// NOP | 1 | ---- | No operation
func (cpu *CPU) NOP() {
	cpu.PC++
}

// Halt | 1 | ---- | Halt until interrupt occurs
func (cpu *CPU) HALT() {
	cpu.Halt = true
}

// DI | 1 | ---- | Disable interrupts, IME=0
func (cpu *CPU) DI() {
	cpu.IME = false
}

// EI | 1 | ---- | Enable interrupts, IME=1
func (cpu *CPU) EI() {
	cpu.IME = true
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
