package mmu

type MBC0 struct {
	mmu           *MMU
	CartridgeData []byte
}

func NewMBC0(mmu *MMU, data []byte) *MBC0 {
	mbc0 := &MBC0{
		mmu:           mmu,
		CartridgeData: data,
	}

	// Cartridge ROM range
	mmu.MapMemoryRange(mbc0, 0x0000, 0x7FFF)

	return mbc0
}

func (mbc *MBC0) ReadByte(addr uint16) byte {
	return mbc.CartridgeData[addr]
}

func (mbc *MBC0) WriteByte(addr uint16, value byte) {
	// Some games, such as Tetris do not use a MBC but still attempt to write to certain
	// MBC related addresses, for example Tetris writes the value 1 to 0x2000, possibly
	// some test code to switch the ROM bank number at a time when the game development
	// did use a MBC/did not fit into 32KB.

	// So we do not error on writes but just do nothing.
}
