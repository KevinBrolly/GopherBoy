package mmu

type MBC2 struct {
	mmu            *MMU
	CartridgeData  []byte
	ROM            []byte
	RAM            []byte
	RAMEnabled     bool
	CurrentROMBank int
}

func NewMBC2(mmu *MMU, data []byte) *MBC2 {
	mbc2 := &MBC2{
		mmu:            mmu,
		CartridgeData:  data,
		RAM:            make([]byte, 512),
		RAMEnabled:     false,
		CurrentROMBank: 1,
	}

	// Cartridge ROM range
	mmu.MapMemoryRange(mbc2, 0x0000, 0x7FFF)
	// External RAM range
	mmu.MapMemoryRange(mbc2, 0xA000, 0xA1FF)

	return mbc2
}

func (mbc *MBC2) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF:
		return mbc.CartridgeData[addr]
	case addr >= 0x4000 && addr <= 0x7FFF:
		addr := addr - 0x4000
		return mbc.CartridgeData[int(addr)+(int(mbc.CurrentROMBank)*0x4000)]
	case addr >= 0xA000 && addr <= 0xA1FF:
		if mbc.RAMEnabled {
			addr := addr - 0xA000
			return mbc.RAM[int(addr)] & 0x0F
		}
	}

	return 0
}

func (mbc *MBC2) WriteByte(addr uint16, value byte) {
	switch {
	// 0000-1FFF - RAM Enable (Write Only)
	// The least significant bit of the upper address byte must be zero to enable/disable cart RAM.
	// For example the following addresses can be used to enable/disable cart RAM: 0000-00FF, 0200-02FF, 0400-04FF, ..., 1E00-1EFF.
	// The suggested address range to use for MBC2 ram enable/disable is 0000-00FF.
	case addr >= 0x0000 && addr <= 0x1FFF:
		if addr&0x100 == 0 {
			switch value & 0xF {
			case 0xA:
				mbc.RAMEnabled = true
			case 0x0:
				mbc.RAMEnabled = false
			}
		}
	// 2000-3FFF - ROM Bank Number (Write Only)
	// Writing a value (XXXXBBBB - X = Don't cares, B = bank select bits) into 2000-3FFF area will select an appropriate ROM bank at 4000-7FFF.
	// The least significant bit of the upper address byte must be one to select a ROM bank. For example the following addresses can be used to select a ROM bank: 2100-21FF, 2300-23FF, 2500-25FF, ..., 3F00-3FFF.
	// The suggested address range to use for MBC2 rom bank selection is 2100-21FF.
	case addr >= 0x2000 && addr <= 0x3FFF:
		if addr&0x100 == 1 {
			mbc.CurrentROMBank = int(value & 0x0F)
		}
	}
}
