package cartridge

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
)

const (
	ROMBankingMode = 0
	RAMBankingMode = 1
)

type MBC1 struct {
	mmu            *mmu.MMU
	CartridgeData  []byte
	RAM            []byte
	RAMEnabled     bool
	CurrentROMBank int
	CurrentRAMBank int
	BankingMode    int
}

func NewMBC1(mmu *mmu.MMU, data []byte) *MBC1 {
	mbc1 := &MBC1{
		mmu:           mmu,
		CartridgeData: data,
		RAM:           make([]byte, 0x8000),
		RAMEnabled:    false,
		BankingMode:   ROMBankingMode,
		// The ROM Bank Number defaults to 01
		CurrentROMBank: 0x01,
	}

	// Cartridge ROM range
	mmu.MapMemoryRange(mbc1, 0x0000, 0x7FFF)
	// External RAM range
	mmu.MapMemoryRange(mbc1, 0xA000, 0xBFFF)

	return mbc1
}

func (mbc *MBC1) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF:
		return mbc.CartridgeData[addr]
	case addr >= 0x4000 && addr <= 0x7FFF:
		addr := addr - 0x4000
		return mbc.CartridgeData[int(addr)+(int(mbc.CurrentROMBank)*0x4000)]
	case addr >= 0xA000 && addr <= 0xBFFF:
		if mbc.RAMEnabled {
			addr := addr - 0xA000
			return mbc.RAM[int(addr)+mbc.CurrentRAMBank*0x2000]
		}
	}

	return 0
}

func (mbc *MBC1) WriteByte(addr uint16, value byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x1FFF:
		switch value & 0xF {
		case 0xA:
			mbc.RAMEnabled = true
		case 0x0:
			mbc.RAMEnabled = false
		}
	case addr >= 0x2000 && addr <= 0x3FFF:
		mbc.CurrentROMBank = (mbc.CurrentROMBank & 0xE0) | int(value&0x1F)
	case addr >= 0x4000 && addr <= 0x5FFF:
		switch mbc.BankingMode {
		case ROMBankingMode:
			mbc.CurrentROMBank = (mbc.CurrentROMBank & 0x1F) | int(value<<5&0xE0)
		case RAMBankingMode:
			mbc.CurrentRAMBank = int(value & 0x3)
		}
	case addr >= 0x6000 && addr <= 0x7FFF:
		switch value & 0x1 {
		case ROMBankingMode:
			mbc.BankingMode = ROMBankingMode
			// only RAM Bank 0 can be used during RAM Banking Mode
			mbc.CurrentRAMBank = 0
		case RAMBankingMode:
			mbc.BankingMode = RAMBankingMode
			// only ROM Banks 00-1Fh can be used during ROM Banking Mode
			mbc.CurrentROMBank = mbc.CurrentROMBank & 0x1F
		}
	case addr >= 0xA000 && addr <= 0xBFFF:
		if mbc.RAMEnabled {
			addr := addr - 0xA000
			mbc.RAM[addr+(uint16(mbc.CurrentRAMBank)*0x2000)] = value
		}
	}

	if mbc.CurrentRAMBank&0x1F == 0 {
		mbc.CurrentRAMBank++
	}
}
