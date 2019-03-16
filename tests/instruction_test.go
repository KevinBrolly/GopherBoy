package Gameboy_test

import (
	"GoBoy/Gameboy"
	"bytes"
	"fmt"
	"testing"
)

func compareGameboy(t *testing.T, initialGameboy *Gameboy.Gameboy, expectedGameboy *Gameboy.Gameboy) {
	if initialGameboy.CPU.Registers.A != expectedGameboy.CPU.Registers.A {
		t.Errorf("A = %#x, expected: A: %#x", initialGameboy.CPU.Registers.A, expectedGameboy.CPU.Registers.A)
	}
	if initialGameboy.CPU.Registers.B != expectedGameboy.CPU.Registers.B {
		t.Errorf("B = %#x, expected: B: %#x", initialGameboy.CPU.Registers.B, expectedGameboy.CPU.Registers.B)
	}
	if initialGameboy.CPU.Registers.C != expectedGameboy.CPU.Registers.C {
		t.Errorf("C = %#x, expected: C: %#x", initialGameboy.CPU.Registers.C, expectedGameboy.CPU.Registers.C)
	}
	if initialGameboy.CPU.Registers.D != expectedGameboy.CPU.Registers.D {
		t.Errorf("D = %#x, expected: D: %#x", initialGameboy.CPU.Registers.D, expectedGameboy.CPU.Registers.D)
	}
	if initialGameboy.CPU.Registers.E != expectedGameboy.CPU.Registers.E {
		t.Errorf("E = %#x, expected: E: %#x", initialGameboy.CPU.Registers.E, expectedGameboy.CPU.Registers.E)
	}
	if initialGameboy.CPU.Registers.H != expectedGameboy.CPU.Registers.H {
		t.Errorf("H = %#x, expected: H: %#x", initialGameboy.CPU.Registers.H, expectedGameboy.CPU.Registers.H)
	}
	if initialGameboy.CPU.Registers.L != expectedGameboy.CPU.Registers.L {
		t.Errorf("L = %#x, expected: L: %#x", initialGameboy.CPU.Registers.L, expectedGameboy.CPU.Registers.L)
	}

	if initialGameboy.CPU.Registers.F != expectedGameboy.CPU.Registers.F {
		t.Errorf("Z = %t, N = %t, H = %t, CY = %t, expected: Z = %t, N = %t, H = %t, CY = %t",
			initialGameboy.CPU.IsFlagSet(Gameboy.Z),
			initialGameboy.CPU.IsFlagSet(Gameboy.N),
			initialGameboy.CPU.IsFlagSet(Gameboy.H),
			initialGameboy.CPU.IsFlagSet(Gameboy.CY),
			expectedGameboy.CPU.IsFlagSet(Gameboy.Z),
			expectedGameboy.CPU.IsFlagSet(Gameboy.N),
			expectedGameboy.CPU.IsFlagSet(Gameboy.H),
			expectedGameboy.CPU.IsFlagSet(Gameboy.CY),
		)
	}

	if initialGameboy.CPU.SP != expectedGameboy.CPU.SP {
		t.Errorf("SP = %v, expected: %v", initialGameboy.CPU.SP, expectedGameboy.CPU.SP)
	}

	if initialGameboy.CPU.PC != expectedGameboy.CPU.PC {
		t.Errorf("PC = %v, expected: %v", initialGameboy.CPU.PC, expectedGameboy.CPU.PC)
	}

	if !bytes.Equal(initialGameboy.ROM, expectedGameboy.ROM) {
		t.Errorf("MMU ROM = %v, expected: %v", initialGameboy.ROM, expectedGameboy.ROM)
	}

	if initialGameboy.WorkingRAM != expectedGameboy.WorkingRAM {
		t.Errorf("MMU WorkingRAM = %v, expected: %v", initialGameboy.WorkingRAM, expectedGameboy.WorkingRAM)
	}

	if initialGameboy.ZeroPageRAM != expectedGameboy.ZeroPageRAM {
		t.Errorf("MMU ZeroPageRAM = %v, expected: %v", initialGameboy.ZeroPageRAM, expectedGameboy.ZeroPageRAM)
	}
}

func setOperand(gameboy *Gameboy.Gameboy, offset uint16, value byte) {
	gameboy.WriteByte(0x100+offset, value)
}

func resetFlags(gameboy *Gameboy.Gameboy) {
	gameboy.CPU.ResetFlag(Gameboy.Z)
	gameboy.CPU.ResetFlag(Gameboy.H)
	gameboy.CPU.ResetFlag(Gameboy.N)
	gameboy.CPU.ResetFlag(Gameboy.CY)
}

func NewTestGameboy() *Gameboy.Gameboy {
	gameboy := Gameboy.NewGameboy()
	// Set the PC to 100, to avoid the boot rom
	gameboy.CPU.PC = 0x100
	// Reset the flags to make setting up tests easier
	resetFlags(gameboy)
	return gameboy
}

func TestLD_r_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0
	initialGameboy.CPU.Registers.B = 0x40

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x40
	expectedGameboy.CPU.Registers.B = 0x40

	initialGameboy.CPU.LD_r_r(&initialGameboy.CPU.Registers.A, &initialGameboy.CPU.Registers.B)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_r_n(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.B = 0
	setOperand(initialGameboy, 1, 0x24)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.B = 0x24
	setOperand(expectedGameboy, 1, 0x24)

	initialGameboy.CPU.LD_r_n(&initialGameboy.CPU.Registers.B)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_r_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x5C)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0x5C
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x5C)

	initialGameboy.CPU.LD_r_HL(&initialGameboy.CPU.Registers.H)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_HL_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x3C
	initialGameboy.CPU.Registers.H = 0x8A
	initialGameboy.CPU.Registers.L = 0xC5

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x3C
	expectedGameboy.CPU.Registers.H = 0x8A
	expectedGameboy.CPU.Registers.L = 0xC5
	expectedGameboy.WriteByte(0x8AC5, 0x3C)

	initialGameboy.CPU.LD_HL_r(&initialGameboy.CPU.Registers.A)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_HL_n(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0x8A
	initialGameboy.CPU.Registers.L = 0xC5
	setOperand(initialGameboy, 1, 0x80)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0x8A
	expectedGameboy.CPU.Registers.L = 0xC5
	expectedGameboy.WriteByte(0x8AC5, 0x80)
	setOperand(expectedGameboy, 1, 0x80)

	initialGameboy.CPU.LD_HL_n()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_A_BC(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.B = 0xC0
	initialGameboy.CPU.Registers.C = 0x00
	initialGameboy.WriteByte(0xC000, 0x2F)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x2F
	expectedGameboy.CPU.Registers.B = 0xC0
	expectedGameboy.CPU.Registers.C = 0x00
	expectedGameboy.WriteByte(0xC000, 0x2F)

	initialGameboy.CPU.LD_A_BC()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_A_DE(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.D = 0xC0
	initialGameboy.CPU.Registers.E = 0x00
	initialGameboy.WriteByte(0xC000, 0x5F)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x5F
	expectedGameboy.CPU.Registers.D = 0xC0
	expectedGameboy.CPU.Registers.E = 0x00
	expectedGameboy.WriteByte(0xC000, 0x5F)

	initialGameboy.CPU.LD_A_DE()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_A_C(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x00
	initialGameboy.CPU.Registers.C = 0x9F
	initialGameboy.WriteByte(0xFF9F, 0x80)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x80
	expectedGameboy.CPU.Registers.C = 0x9F
	expectedGameboy.WriteByte(0xFF9F, 0x80)

	initialGameboy.CPU.LD_A_C()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_C_A(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x80
	initialGameboy.CPU.Registers.C = 0x9F

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x80
	expectedGameboy.CPU.Registers.C = 0x9F
	expectedGameboy.WriteByte(0xFF9F, 0x80)

	initialGameboy.CPU.LD_C_A()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLDH_A_n(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.WriteByte(0xFF42, 0x80)
	setOperand(initialGameboy, 1, 0x42)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x80
	expectedGameboy.WriteByte(0xFF42, 0x80)
	setOperand(expectedGameboy, 1, 0x42)

	initialGameboy.CPU.LDH_A_n()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLDH_n_A(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x80
	setOperand(initialGameboy, 1, 0x42)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x80
	expectedGameboy.WriteByte(0xFF42, 0x80)
	setOperand(expectedGameboy, 1, 0x42)

	initialGameboy.CPU.LDH_n_A()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_A_nn(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x00
	initialGameboy.WriteByte(0x8000, 0x80)
	setOperand(initialGameboy, 1, 0x00)
	setOperand(initialGameboy, 2, 0x80)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x80
	expectedGameboy.WriteByte(0x8000, 0x80)
	setOperand(expectedGameboy, 1, 0x00)
	setOperand(expectedGameboy, 2, 0x80)

	initialGameboy.CPU.LD_A_nn()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_nn_A(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x80
	setOperand(initialGameboy, 1, 0x00)
	setOperand(initialGameboy, 2, 0x80)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x80
	expectedGameboy.WriteByte(0x8000, 0x80)
	setOperand(expectedGameboy, 1, 0x00)
	setOperand(expectedGameboy, 2, 0x80)

	initialGameboy.CPU.LD_nn_A()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_A_HLI(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0x01
	initialGameboy.CPU.Registers.L = 0xFF
	initialGameboy.WriteByte(0x01FF, 0x56)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x56
	expectedGameboy.CPU.Registers.H = 0x02
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0x01FF, 0x56)

	initialGameboy.CPU.LD_A_HLI()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_A_HLD(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0x8A
	initialGameboy.CPU.Registers.L = 0x5C
	initialGameboy.WriteByte(0x8A5C, 0x3C)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x3C
	expectedGameboy.CPU.Registers.H = 0x8A
	expectedGameboy.CPU.Registers.L = 0x5B
	expectedGameboy.WriteByte(0x8A5C, 0x3C)

	initialGameboy.CPU.LD_A_HLD()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_BC_A(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x3F
	initialGameboy.CPU.Registers.B = 0x20
	initialGameboy.CPU.Registers.C = 0x5F

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x3F
	expectedGameboy.CPU.Registers.B = 0x20
	expectedGameboy.CPU.Registers.C = 0x5F
	expectedGameboy.WriteByte(0x205F, 0x3F)

	initialGameboy.CPU.LD_BC_A()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_DE_A(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x3F
	initialGameboy.CPU.Registers.D = 0x20
	initialGameboy.CPU.Registers.E = 0x5C

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x3F
	expectedGameboy.CPU.Registers.D = 0x20
	expectedGameboy.CPU.Registers.E = 0x5C
	expectedGameboy.WriteByte(0x205C, 0x3F)

	initialGameboy.CPU.LD_DE_A()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_HLI_A(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x56
	initialGameboy.CPU.Registers.H = 0xFF
	initialGameboy.CPU.Registers.L = 0xFF

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x56
	expectedGameboy.CPU.Registers.H = 0x00
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xFFFF, 0x56)

	initialGameboy.CPU.LD_HLI_A()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_HLD_A(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x05
	initialGameboy.CPU.Registers.H = 0x40
	initialGameboy.CPU.Registers.L = 0x00

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x05
	expectedGameboy.CPU.Registers.H = 0x3F
	expectedGameboy.CPU.Registers.L = 0xFF
	expectedGameboy.WriteByte(0x4000, 0x05)

	initialGameboy.CPU.LD_HLD_A()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_rr_nn(t *testing.T) {
	initialGameboy := NewTestGameboy()
	setOperand(initialGameboy, 1, 0x5B)
	setOperand(initialGameboy, 2, 0x3A)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0x3A
	expectedGameboy.CPU.Registers.L = 0x5B
	setOperand(expectedGameboy, 1, 0x5B)
	setOperand(expectedGameboy, 2, 0x3A)

	initialGameboy.CPU.LD_rr_nn(&initialGameboy.CPU.Registers.H, &initialGameboy.CPU.Registers.L)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_SP_nn(t *testing.T) {
	initialGameboy := NewTestGameboy()
	setOperand(initialGameboy, 1, 0x5B)
	setOperand(initialGameboy, 2, 0x3A)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.SP = 0x3A5B
	setOperand(expectedGameboy, 1, 0x5B)
	setOperand(expectedGameboy, 2, 0x3A)

	initialGameboy.CPU.LD_SP_nn()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_SP_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xFF
	initialGameboy.CPU.Registers.L = 0x80

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xFF
	expectedGameboy.CPU.Registers.L = 0x80
	expectedGameboy.CPU.SP = 0xFF80

	initialGameboy.CPU.LD_SP_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestPUSH_qq(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.B = 0xFF
	initialGameboy.CPU.Registers.C = 0xFC
	initialGameboy.CPU.SP = 0xFFFE

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.B = 0xFF
	expectedGameboy.CPU.Registers.C = 0xFC
	expectedGameboy.WriteByte(0xFFFE-1, 0xFF)
	expectedGameboy.WriteByte(0xFFFE-2, 0xFC)
	expectedGameboy.CPU.SP = 0xFFFC

	initialGameboy.CPU.PUSH_qq(&initialGameboy.CPU.Registers.B, &initialGameboy.CPU.Registers.C)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestPOP_qq(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.SP = 0xFFFC
	initialGameboy.WriteByte(0xFFFC, 0x5F)
	initialGameboy.WriteByte(0xFFFD, 0x3C)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.B = 0x3C
	expectedGameboy.CPU.Registers.C = 0x5F
	expectedGameboy.WriteByte(0xFFFC, 0x5F)
	expectedGameboy.WriteByte(0xFFFD, 0x3C)
	expectedGameboy.CPU.SP = 0xFFFE

	initialGameboy.CPU.POP_qq(&initialGameboy.CPU.Registers.B, &initialGameboy.CPU.Registers.C)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_HL_SP_e(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.SP = 0xFFF8
	setOperand(initialGameboy, 1, 0x02)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xFF
	expectedGameboy.CPU.Registers.L = 0xFA
	expectedGameboy.CPU.SP = 0xFFF8
	setOperand(expectedGameboy, 1, 0x02)

	initialGameboy.CPU.LD_HL_SP_e()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestLD_nn_SP(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.SP = 0xFFF8
	setOperand(initialGameboy, 1, 0x00)
	setOperand(initialGameboy, 2, 0xC1)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.SP = 0xFFF8
	expectedGameboy.WriteByte(0xC100, 0xF8)
	expectedGameboy.WriteByte(0xC101, 0xFF)
	setOperand(expectedGameboy, 1, 0x00)
	setOperand(expectedGameboy, 2, 0xC1)

	initialGameboy.CPU.LD_nn_SP()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestADD_s(t *testing.T) {
	t.Run("ADD A, r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3A
		initialGameboy.CPU.Registers.B = 0xC6

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x00
		expectedGameboy.CPU.Registers.B = 0xC6

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.ADD_s(initialGameboy.CPU.Registers.B, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("ADD A, n", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3C
		setOperand(initialGameboy, 1, 0xFF)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x3B
		setOperand(expectedGameboy, 1, 0xFF)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.ADD_s(initialGameboy.CPU.GetByteOffset(1), 2)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("ADD A, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3C
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x12)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x4E
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x12)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.ADD_s(
			initialGameboy.ReadByte(
				Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L),
			),
			2,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestADC_s(t *testing.T) {
	t.Run("ADC A r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0xE1
		initialGameboy.CPU.Registers.E = 0x0F

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0xF1
		expectedGameboy.CPU.Registers.E = 0x0F

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.ADC_A_s(initialGameboy.CPU.Registers.E, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("ADC A s", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0xE1
		setOperand(initialGameboy, 1, 0x3B)

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x1D
		setOperand(expectedGameboy, 1, 0x3B)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.ADC_A_s(initialGameboy.CPU.GetByteOffset(1), 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("ADC A (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0xE1
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x1E)

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x00
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x1E)

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.ADC_A_s(
			initialGameboy.ReadByte(Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L)),
			1,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestSUB_s(t *testing.T) {
	t.Run("SUB A, r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3E
		initialGameboy.CPU.Registers.E = 0x3E

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x00
		expectedGameboy.CPU.Registers.E = 0x3E

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.SUB_s(initialGameboy.CPU.Registers.E, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("SUB A, n", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3E
		setOperand(initialGameboy, 1, 0x0F)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x2F
		setOperand(expectedGameboy, 1, 0x0F)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.SUB_s(initialGameboy.CPU.GetByteOffset(1), 2)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("SUB A, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3E
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x40)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0xFE
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x40)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.SUB_s(
			initialGameboy.ReadByte(
				Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L),
			),
			2,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestSBC_s(t *testing.T) {
	t.Run("SBC A, r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3B
		initialGameboy.CPU.Registers.H = 0x2A

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x10
		expectedGameboy.CPU.Registers.H = 0x2A

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.SBC_s(initialGameboy.CPU.Registers.H, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("SBC A, s", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3B
		setOperand(initialGameboy, 1, 0x3A)

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x00
		setOperand(expectedGameboy, 1, 0x3A)

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.SBC_s(initialGameboy.CPU.GetByteOffset(1), 2)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("SBC A, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3B
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x4F)

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0xEB
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x4F)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.SBC_s(
			initialGameboy.ReadByte(Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L)),
			2,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestAND_s(t *testing.T) {
	t.Run("AND A, r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x5A
		initialGameboy.CPU.Registers.L = 0x3F

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x1A
		expectedGameboy.CPU.Registers.L = 0x3F

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.AND_s(initialGameboy.CPU.Registers.L, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("AND A, s", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x5A
		setOperand(initialGameboy, 1, 0x38)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x18
		setOperand(expectedGameboy, 1, 0x38)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.AND_s(initialGameboy.CPU.GetByteOffset(1), 2)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("AND A, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x5A
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x00
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x00)

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.AND_s(
			initialGameboy.ReadByte(Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L)),
			2,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestOR_s(t *testing.T) {
	t.Run("OR A, r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x5A

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x5A

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.OR_s(initialGameboy.CPU.Registers.A, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("OR A, s", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x5A
		setOperand(initialGameboy, 1, 0x03)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x5B
		setOperand(expectedGameboy, 1, 0x03)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.OR_s(initialGameboy.CPU.GetByteOffset(1), 2)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("OR A, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x5A
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x0F)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x5F
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x0F)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.OR_s(
			initialGameboy.ReadByte(Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L)),
			2,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestXOR_s(t *testing.T) {
	t.Run("XOR A, r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0xFF

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x00

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.XOR_s(initialGameboy.CPU.Registers.A, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("XOR A, s", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0xFF
		setOperand(initialGameboy, 1, 0x0F)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0xF0
		setOperand(expectedGameboy, 1, 0x0F)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.XOR_s(initialGameboy.CPU.GetByteOffset(1), 2)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("XOR A, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0xFF
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x8A)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x75
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x8A)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.XOR_s(
			initialGameboy.ReadByte(Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L)),
			2,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestCP_s(t *testing.T) {
	t.Run("CP A, r", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3C
		initialGameboy.CPU.Registers.B = 0x2F

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x3C
		expectedGameboy.CPU.Registers.B = 0x2F

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.CP_s(initialGameboy.CPU.Registers.B, 1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("CP A, s", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3C
		setOperand(initialGameboy, 1, 0x3C)

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x3C
		setOperand(expectedGameboy, 1, 0x3C)

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.CP_s(initialGameboy.CPU.GetByteOffset(1), 2)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("CP A, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x3C
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0x40)

		initialGameboy.CPU.SetFlag(Gameboy.CY)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x3C
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0x40)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.ResetFlag(Gameboy.H)
		expectedGameboy.CPU.SetFlag(Gameboy.N)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.CP_s(
			initialGameboy.ReadByte(Gameboy.JoinBytes(initialGameboy.CPU.Registers.H, initialGameboy.CPU.Registers.L)),
			2,
		)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestINC_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0xFF

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x00

	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.SetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.INC_r(&initialGameboy.CPU.Registers.A)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestINC_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()

	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x50)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x51)

	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.INC_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestDEC_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.L = 0x01

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.L = 0x00

	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.SetFlag(Gameboy.N)

	initialGameboy.CPU.DEC_r(&initialGameboy.CPU.Registers.L)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestDEC_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()

	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x00)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0xFF)

	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.SetFlag(Gameboy.H)
	expectedGameboy.CPU.SetFlag(Gameboy.N)

	initialGameboy.CPU.DEC_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestADD_HL_rr(t *testing.T) {
	t.Run("ADD HL, BC", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.H = 0x8A
		initialGameboy.CPU.Registers.L = 0x23
		initialGameboy.CPU.Registers.B = 0x06
		initialGameboy.CPU.Registers.C = 0x05

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.H = 0x90
		expectedGameboy.CPU.Registers.L = 0x28
		expectedGameboy.CPU.Registers.B = 0x06
		expectedGameboy.CPU.Registers.C = 0x05

		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)

		initialGameboy.CPU.ADD_HL_rr(&initialGameboy.CPU.Registers.B, &initialGameboy.CPU.Registers.C)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("ADD HL, HL", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.H = 0x8A
		initialGameboy.CPU.Registers.L = 0x23

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.H = 0x14
		expectedGameboy.CPU.Registers.L = 0x46

		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)
		expectedGameboy.CPU.SetFlag(Gameboy.CY)

		initialGameboy.CPU.ADD_HL_rr(&initialGameboy.CPU.Registers.H, &initialGameboy.CPU.Registers.L)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestSP_e(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.SP = 0xFFF8
	setOperand(initialGameboy, 1, 0x02)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.SP = 0xFFFA
	setOperand(expectedGameboy, 1, 0x02)

	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)
	expectedGameboy.CPU.ResetFlag(Gameboy.CY)

	initialGameboy.CPU.ADD_SP_e()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestINC_rr(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.D = 0x23
	initialGameboy.CPU.Registers.E = 0x5F

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.D = 0x23
	expectedGameboy.CPU.Registers.E = 0x60

	initialGameboy.CPU.INC_rr(&initialGameboy.CPU.Registers.D, &initialGameboy.CPU.Registers.E)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestDEC_rr(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.D = 0x23
	initialGameboy.CPU.Registers.E = 0x5F

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.D = 0x23
	expectedGameboy.CPU.Registers.E = 0x5E

	initialGameboy.CPU.DEC_rr(&initialGameboy.CPU.Registers.D, &initialGameboy.CPU.Registers.E)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestINC_SP(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.SP = 0xFFFC

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.SP = 0xFFFD

	initialGameboy.CPU.INC_SP()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestDEC_SP(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.SP = 0xFFFC

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.SP = 0xFFFB

	initialGameboy.CPU.DEC_SP()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRLCA(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x85
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x0A

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RLCA()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRLA(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x95
	initialGameboy.CPU.SetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x2B

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RLA()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRRCA(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x3B
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x9D

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RRCA()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRRA(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x81
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x40

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RRA()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRLC_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.B = 0x85
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.B = 0x0B

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RLC_r(&initialGameboy.CPU.Registers.B)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRLC_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x00)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x00)

	expectedGameboy.CPU.ResetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RLC_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRL_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.L = 0x80
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.L = 0x00

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RL_r(&initialGameboy.CPU.Registers.L)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRL_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x11)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x22)

	expectedGameboy.CPU.ResetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RL_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRRC_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.C = 0x01
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.C = 0x80

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RRC_r(&initialGameboy.CPU.Registers.C)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRRC_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x00)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x00)

	expectedGameboy.CPU.ResetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RRC_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRR_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x01
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x00

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RR_r(&initialGameboy.CPU.Registers.A)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRR_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x8A)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x45)

	expectedGameboy.CPU.ResetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.RR_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSLA_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.D = 0x80
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.D = 0x00

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SLA_r(&initialGameboy.CPU.Registers.D)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSLA_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0xFF)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0xFE)

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SLA_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSRA_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x8A
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.D = 0xC5

	expectedGameboy.CPU.ResetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SRA_r(&initialGameboy.CPU.Registers.D)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSRA_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x01)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x00)

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SRA_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSRL_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x01
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x00

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SRL_r(&initialGameboy.CPU.Registers.A)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSRL_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0xFF)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x7F)

	expectedGameboy.CPU.SetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SRL_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSWAP_r(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.A = 0x00

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.A = 0x00

	expectedGameboy.CPU.ResetFlag(Gameboy.CY)
	expectedGameboy.CPU.SetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SWAP_r(&initialGameboy.CPU.Registers.A)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestSWAP_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0xF0)
	initialGameboy.CPU.ResetFlag(Gameboy.CY)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x0F)

	expectedGameboy.CPU.ResetFlag(Gameboy.CY)
	expectedGameboy.CPU.ResetFlag(Gameboy.Z)
	expectedGameboy.CPU.ResetFlag(Gameboy.H)
	expectedGameboy.CPU.ResetFlag(Gameboy.N)

	initialGameboy.CPU.SWAP_HL()

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestBIT_b_r(t *testing.T) {
	t.Run("BIT 7, A", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x80

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x80

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)

		initialGameboy.CPU.BIT_b_r(7, &initialGameboy.CPU.Registers.A)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("BIT 4, L", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.L = 0xEF

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.L = 0xEF

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)

		initialGameboy.CPU.BIT_b_r(4, &initialGameboy.CPU.Registers.L)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestBIT_b_HL(t *testing.T) {
	t.Run("BIT 0, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0xFE)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0xFE)

		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)

		initialGameboy.CPU.BIT_b_HL(0)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("BIT 1, (HL)", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.H = 0xC0
		initialGameboy.CPU.Registers.L = 0x00
		initialGameboy.WriteByte(0xC000, 0xFE)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.H = 0xC0
		expectedGameboy.CPU.Registers.L = 0x00
		expectedGameboy.WriteByte(0xC000, 0xFE)

		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.SetFlag(Gameboy.H)
		expectedGameboy.CPU.ResetFlag(Gameboy.N)

		initialGameboy.CPU.BIT_b_HL(1)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestSET_b_r(t *testing.T) {
	t.Run("SET 3, A", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x80

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x88

		initialGameboy.CPU.SET_b_r(3, &initialGameboy.CPU.Registers.A)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("SET 7, L", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.L = 0x3B

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.L = 0xBB

		initialGameboy.CPU.SET_b_r(7, &initialGameboy.CPU.Registers.L)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestSET_b_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0x00)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0x08)

	initialGameboy.CPU.SET_b_HL(3)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestRES_b_r(t *testing.T) {
	t.Run("RES 7, A", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.A = 0x80

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.A = 0x00

		initialGameboy.CPU.RES_b_r(7, &initialGameboy.CPU.Registers.A)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("RES 1, L", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.Registers.L = 0x3B

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.Registers.L = 0x39

		initialGameboy.CPU.RES_b_r(1, &initialGameboy.CPU.Registers.L)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestRES_b_HL(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.Registers.H = 0xC0
	initialGameboy.CPU.Registers.L = 0x00
	initialGameboy.WriteByte(0xC000, 0xFF)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.Registers.H = 0xC0
	expectedGameboy.CPU.Registers.L = 0x00
	expectedGameboy.WriteByte(0xC000, 0xF7)

	initialGameboy.CPU.RES_b_HL(3)

	compareGameboy(t, initialGameboy, expectedGameboy)
}

// Jump Instructions

func TestJP_nn(t *testing.T) {
	initialGameboy := NewTestGameboy()
	initialGameboy.CPU.PC = 0x100
	setOperand(initialGameboy, 2, 0x80)
	setOperand(initialGameboy, 1, 0x00)

	expectedGameboy := NewTestGameboy()
	expectedGameboy.CPU.PC = 0x8000
	setOperand(expectedGameboy, 2, 0x80)
	setOperand(expectedGameboy, 1, 0x00)

	initialGameboy.CPU.JP_nn()
	fmt.Printf("VALUE: %#x\n\n", initialGameboy.CPU.PC)
	compareGameboy(t, initialGameboy, expectedGameboy)
}

func TestJP_cc_nn(t *testing.T) {
	t.Run("JP NZ, Z=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0xC000
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_NZ)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JP NZ, Z=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0x100
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_NZ)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JP Z, Z=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0x100
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_Z)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JP Z, Z=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0xC000
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_Z)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JP NC, CY=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0xC000
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_NC)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JP NC, CY=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0x100
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_NC)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JP C, CY=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0x100
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_C)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JP C, CY=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 2, 0xC0)
		setOperand(initialGameboy, 1, 0x00)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0xC000
		setOperand(expectedGameboy, 2, 0xC0)
		setOperand(expectedGameboy, 1, 0x00)

		initialGameboy.CPU.JP_cc_nn(Gameboy.CC_C)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

func TestJR_cc_e(t *testing.T) {
	t.Run("JR NZ, Z=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0x80)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0x180
		setOperand(expectedGameboy, 1, 0x80)

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_NZ)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JR NZ, Z=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0x80)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0x100
		setOperand(expectedGameboy, 1, 0x80)

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_NZ)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JR Z, Z=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0xB2) // -50 as "unsigned" byte

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0x100
		setOperand(expectedGameboy, 1, 0xB2) // -50 as "unsigned" byte

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_Z)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JR Z, Z=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.Z)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0xB8) // -50 as "unsigned" byte

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.Z)
		expectedGameboy.CPU.PC = 0xC8
		setOperand(expectedGameboy, 1, 0xB8) // -50 as "unsigned" byte

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_Z)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JR NC, CY=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0xB2) // -50 as "unsigned" byte

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0xB0
		setOperand(expectedGameboy, 1, 0xB2) // -50 as "unsigned" byte

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_NC)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JR NC, CY=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0xB2) // -50 as "unsigned" byte

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0x100
		setOperand(expectedGameboy, 1, 0xB2) // -50 as "unsigned" byte

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_NC)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JR C, CY=False", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.ResetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0x80)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.ResetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0x180
		setOperand(expectedGameboy, 1, 0x80)

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_C)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})

	t.Run("JR C, CY=True", func(t *testing.T) {
		initialGameboy := NewTestGameboy()
		initialGameboy.CPU.SetFlag(Gameboy.CY)
		initialGameboy.CPU.PC = 0x100
		setOperand(initialGameboy, 1, 0x80)

		expectedGameboy := NewTestGameboy()
		expectedGameboy.CPU.SetFlag(Gameboy.CY)
		expectedGameboy.CPU.PC = 0x180
		setOperand(expectedGameboy, 1, 0x80)

		initialGameboy.CPU.JR_cc_e(Gameboy.CC_C)

		compareGameboy(t, initialGameboy, expectedGameboy)
	})
}

// // Call and Return Instructions

// func TestRET(t *testing.T) {
//     initialGameboy := NewTestGameboy()
//     initialGameboy.MMU.WriteByte(0x9000, 0x03)
//     initialGameboy.MMU.WriteByte(0x9000 + 1, 0x80)
//     initialGameboy.CPU.SP = 0x9000

//     expectedGameboy := NewTestGameboy()
//     expectedGameboy.CPU.PC = 0x8003
//     expectedGameboy.CPU.SP = 0x9002

//     CPU.RET(initialGameboy.CPU)

//     compareGameboy(t, initialGameboy, expectedGameboy)
// }
