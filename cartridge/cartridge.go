package cartridge

import (
	"GopherBoy/mmu"
	"io/ioutil"
	"log"
)

type Cartridge struct {
	data []byte
	MBC  mmu.Memory
}

func NewCartridge(filename string) (cartridge *Cartridge) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	cartridge = &Cartridge{
		data: data,
	}

	switch cartridge.Type {
	case 0:
		cartridge.MBC = NewMBC0(cartridge.MMU, data)
	case 1:
		cartridge.MBC = NewMBC1(cartridge.MMU, data)
	case 2:
		cartridge.MBC = NewMBC1(cartridge.MMU, data)
	case 3:
		cartridge.MBC = NewMBC1(cartridge.MMU, data)
	case 4:
		cartridge.MBC = NewMBC1(cartridge.MMU, data)
	}

	return cartridge
}

func (c *Cartridge) Type() byte {
	return c.data[0x147]
}

func (c *Cartridge) ROMSize() byte {
	return c.data[0x148]
}

func (c *Cartridge) RAMSize() byte {
	return c.data[0x149]
}

func (c *Cartridge) DestinationCode() byte {
	return c.data[0x14A]
}

func (c *Cartridge) OldLicenseeCode() byte {
	return c.data[0x14B]
}
