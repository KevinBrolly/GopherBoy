package cartridge

import (
	"fmt"
	"testing"
)

func TestReadCartridge(t *testing.T) {
	// Test we can read a byte

	data := make([]byte, 0x8000)
	data[0] = 0xFF

	mmu := NewMMU()
	mbc1 := NewMBC1(mmu, data)

	if mbc1.ReadByte(0) != 0xFF {
		t.Errorf("ReadByte() should have read a value but didnt")
	}
}

func TestReadROMBanks(t *testing.T) {
	// Test we can read a byte from all the ROM banks

	// 125 banks total for mbc1
	data := make([]byte, 0x4000*125)

	for i := 0; i < 125; i++ {
		data[i*0x4000] = 0xFF
	}

	// Setup special values for Bank 0x0, 0x20, 0x40, and 0x60.
	// Any attempt to address these ROM Banks will select Bank 0x01, 0x21, 0x41, and 0x61 instead.
	data[0x4000*0x00] = 0
	data[0x4000*0x01] = 0xFF

	data[0x4000*0x20] = 0
	data[0x4000*0x21] = 0xFF

	data[0x4000*0x40] = 0
	data[0x4000*0x41] = 0xFF

	data[0x4000*0x60] = 0
	data[0x4000*0x61] = 0xFF

	mmu := NewMMU()
	mbc1 := NewMBC1(mmu, data)

	for i := 0; i < 125; i++ {
		mbc1.CurrentROMBank = i
		if mbc1.ReadByte(0x4000) != 0xFF {
			t.Errorf("ReadByte(%v) with CurrentROMBank = %v should have read a value but didnt", 0x4000, mbc1.CurrentROMBank)
		}
	}
}

func TestReadRAMBanks(t *testing.T) {
	// Test we can read a byte from all the RAM banks

	mmu := NewMMU()
	mbc1 := NewMBC1(mmu, make([]byte, 0))

	// Fill the RAM with test data
	for i := 0; i < len(mbc1.RAM); i++ {
		switch {
		// Bank 0
		case i >= 0x0000 && i <= 0x1FFF:
			mbc1.RAM[i] = 1
		// Bank 1
		case i >= 0x2000 && i <= 0x3FFF:
			mbc1.RAM[i] = 2
		// Bank 2
		case i >= 0x4000 && i <= 0x5FFF:
			mbc1.RAM[i] = 3
		// Bank 3
		case i >= 0x6000 && i <= 0x7FFF:
			mbc1.RAM[i] = 4
		}
	}

	// Disable RAM
	mbc1.RAMEnabled = false

	// With RAM disabled reads should return 0
	if mbc1.ReadByte(0xA000) != 0 {
		t.Errorf("Was able to read from RAM when it was disabled")
	}

	// re-enable RAM
	mbc1.RAMEnabled = true
	for i := 0; i < 4; i++ {
		// Read from Bank 0 should return 1, Bank 1 should return 2 etc...
		mbc1.CurrentRAMBank = i
		if mbc1.ReadByte(0xA000) != byte(i+1) {
			t.Errorf("Got wrong value from RAM bank %v", mbc1.CurrentRAMBank)
		}
	}
}

func TestRAMEnable(t *testing.T) {
	// RAM can be enabled by writing 0x0A to any address between 0x0000 - 0x1FFF
	// RAM can be disabled by writing 0x00 to the same addresses

	mmu := NewMMU()
	mbc1 := NewMBC1(mmu, make([]byte, 0))

	// Start with RAM disabled
	mbc1.RAMEnabled = false

	// Loop through all possible address values and test they work
	var i uint16
	for i = 0x0000; i <= 0x1FFF; i++ {
		mbc1.WriteByte(i, 0x0A)
		if mbc1.RAMEnabled != true {
			t.Errorf("RAM should have been enabled on write to %v with value %v but it was not", i, 0x0A)
		}

		mbc1.WriteByte(i, 0x00)
		if mbc1.RAMEnabled != false {
			t.Errorf("RAM should have been disabled on write to %v with value %v but it was not", i, 0x00)
		}
	}
}

func TestROMBankNumberSelect(t *testing.T) {
	// The ROM bank number can be changed by writing to two seperate address spaces
	// Writing to 0x2000-0x3FFF selects the lower 5 bits of the ROM Bank Number (in range 0x01-0x1F)
	// Writing to 0x4000-0x5FFF selects the upper two bits (Bit 5-6) of the ROM Bank number, depending on the current ROM/RAM Mode.

	mmu := NewMMU()
	mbc1 := NewMBC1(mmu, make([]byte, 0))

	mbc1.BankingMode = ROMBankingMode

	cases := []struct {
		Address         uint16
		Value           byte
		ExpectedROMBank int
	}{
		{0x2000, 0x10, 0x10},
		{0x3000, 0x0F, 0x0F},
		{0x3FFF, 0x1F, 0x1F},
		{0x4000, 0x01, 0x3F},
		{0x5FFF, 0x03, 0x7F},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%#x", tt.Address), func(t *testing.T) {
			mbc1.WriteByte(tt.Address, tt.Value)
			if mbc1.CurrentROMBank != tt.ExpectedROMBank {
				t.Errorf("CurrentROMBank should have been set to %08b but was set to %08b", tt.ExpectedROMBank, mbc1.CurrentROMBank)
			}
		})
	}
}
