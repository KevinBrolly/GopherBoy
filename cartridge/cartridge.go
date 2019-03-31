package cartridge

import (
	"GopherBoy/mmu"

	"io/ioutil"
	"log"
)

const (
	ROMBankingMode = 0
	RAMBankingMode = 1
)

type MBC interface {
	ReadByte(addr uint16) byte
	WriteByte(addr uint16, b byte)
}

type Cartridge struct {
	mmu           *mmu.MMU
	CartridgeData []byte
	MBC           MBC
}

func NewCartridge(mmu *mmu.MMU) *Cartridge {
	cartridge := &Cartridge{
		mmu: mmu,
	}

	return cartridge
}

func (c *Cartridge) LoadCartridgeData(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	c.CartridgeData = data
	c.SetMBC()
}

func (c *Cartridge) SetMBC() {
	switch c.CartridgeData[0x147] {
	case 1:
		c.MBC = NewMBC1(c.mmu, c)
	case 2:
		c.MBC = NewMBC1(c.mmu, c)
	case 3:
		c.MBC = NewMBC1(c.mmu, c)
	case 4:
		c.MBC = NewMBC1(c.mmu, c)
	}
}
