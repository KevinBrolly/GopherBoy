package cartridge

import (
	"GopherBoy/mmu"
)

type MBC1 struct {
	mmu            *mmu.MMU
	cartridge      *Cartridge
	CartridgeData  []byte
	ROM            []byte
	RAM            []byte
	RAMEnabled     bool
	CurrentROMBank int
	CurrentRAMBank int
	BankingMode    int
}

func NewMBC1(mmu *mmu.MMU, cartridge *Cartridge) *MBC1 {
	mbc1 := &MBC1{
		mmu:            mmu,
		cartridge:      cartridge,
		RAM:            make([]byte, 0x8000),
		RAMEnabled:     false,
		CurrentROMBank: 1,
		BankingMode:    ROMBankingMode,
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
		return mbc.cartridge.CartridgeData[addr]
	case addr >= 0x4000 && addr <= 0x7FFF:
		addr := addr - 0x4000
		return mbc.cartridge.CartridgeData[int(addr)+(int(mbc.CurrentROMBank)*0x4000)]
	case addr >= 0xA000 && addr <= 0xBFFF:
		addr := addr - 0xA000
		return mbc.RAM[int(addr)+mbc.CurrentRAMBank*0x2000]
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
		mbc.CurrentROMBank = int(value & 0x1F)
		// When 0 is written, the MBC translates that to Bank 1
		// This is fine because ROM Bank 0 can be always directly accessed by reading from 0000-3FFF
	case addr >= 0x4000 && addr <= 0x5FFF:
		switch mbc.BankingMode {
		case ROMBankingMode:
			mbc.CurrentROMBank = int(value & 0xE0)
			// When 0 is written, the MBC translates that to Bank 1
			// This is fine because ROM Bank 0 can be always directly accessed by reading from 0000-3FFF
		case RAMBankingMode:
			mbc.CurrentRAMBank = int(value & 0x3)
		}
	case addr >= 0x6000 && addr <= 0x7FFF:
		switch value & 0x1 {
		case 0:
			mbc.BankingMode = ROMBankingMode
		case 1:
			mbc.BankingMode = RAMBankingMode
			// only RAM Bank 0 can be used during RAM Banking Mode
			// and only ROM Banks 00-1Fh can be used during ROM Banking Mode
			// so we set the RAM bank to 0 in this case
			mbc.CurrentRAMBank = 0
		}
	case addr >= 0xA000 && addr <= 0xBFFF:
		if mbc.RAMEnabled {
			addr := addr - 0xA000
			mbc.RAM[addr+(uint16(mbc.CurrentRAMBank)*0x2000)] = value
		}
	}
}
