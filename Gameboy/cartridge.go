package Gameboy

import (
	"io/ioutil"
	"log"
)

const (
	ROMBankingMode = 0
	RAMBankingMode = 1

	MBC1 = 1
	MBC2 = 2
)

type Cartridge struct {
	gameboy        *Gameboy
	CartridgeData  []byte
	ROM            []byte
	RAM            []byte
	MBC            int
	RAMEnabled     bool
	CurrentROMBank int
	CurrentRAMBank int
	BankingMode    int
}

func NewCartridge(gameboy *Gameboy) *Cartridge {
	cartridge := &Cartridge{
		gameboy:        gameboy,
		RAM:            make([]byte, 0x8000),
		RAMEnabled:     false,
		CurrentROMBank: 1,
		BankingMode:    ROMBankingMode,
	}

	return cartridge
}

func (c *Cartridge) LoadCartridgeData(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	c.CartridgeData = data
	c.SetBankMode()
}

func (c *Cartridge) SetBankMode() {
	switch c.CartridgeData[0x147] {
	case 1:
		c.MBC = MBC1
	case 2:
		c.MBC = MBC1
	case 3:
		c.MBC = MBC1
	case 4:
		c.MBC = MBC1
	case 5:
		c.MBC = MBC2
	case 6:
		c.MBC = MBC2
	}
}

func (c *Cartridge) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF:
		return c.CartridgeData[addr]
	case addr >= 0x4000 && addr <= 0x7FFF:
		addr := addr - 0x4000
		return c.CartridgeData[int(addr)+(int(c.CurrentROMBank)*0x4000)]
	case addr >= 0xA000 && addr <= 0xBFFF:
		addr := addr - 0xA000
		return c.RAM[int(addr)+c.CurrentRAMBank*0x2000]
	}

	return 0
}

func (c *Cartridge) WriteByte(addr uint16, value byte) {
	if addr <= 0x8000 {
		switch c.MBC {
		case MBC1:
			switch {
			case addr >= 0x0000 && addr <= 0x1FFF:
				switch value & 0xF {
				case 0xA:
					c.RAMEnabled = true
				case 0x0:
					c.RAMEnabled = false
				}
			case addr >= 0x2000 && addr <= 0x3FFF:
				c.CurrentROMBank = (c.CurrentROMBank & 0xE0) | (int(value & 0x1F))
				// When 0 is written, the MBC translates that to Bank 1
				// This is fine because ROM Bank 0 can be always directly accessed by reading from 0000-3FFF
				if c.CurrentROMBank == 0 {
					c.CurrentROMBank = 1
				}
			case addr >= 0x4000 && addr <= 0x5FFF:
				switch c.BankingMode {
				case ROMBankingMode:
					c.CurrentROMBank = (c.CurrentROMBank & 0x1F) | (int(value & 0xE0))
					// When 0 is written, the MBC translates that to Bank 1
					// This is fine because ROM Bank 0 can be always directly accessed by reading from 0000-3FFF
					if c.CurrentROMBank == 0 {
						c.CurrentROMBank = 1
					}
				case RAMBankingMode:
					c.CurrentRAMBank = int(value & 0x3)
				}
			case addr >= 0x6000 && addr <= 0x7FFF:
				switch value & 0x1 {
				case 0:
					c.BankingMode = ROMBankingMode
				case 1:
					c.BankingMode = RAMBankingMode
					// only RAM Bank 0 can be used during RAM Banking Mode
					// and only ROM Banks 00-1Fh can be used during ROM Banking Mode
					// so we set the RAM bank to 0 in this case
					c.CurrentRAMBank = 0
				}
			case addr >= 0xA000 && addr <= 0xBFFF:
				if c.RAMEnabled {
					addr := addr - 0xA000
					c.RAM[addr+(uint16(c.CurrentRAMBank)*0x2000)] = value
				}
			}
		case MBC2:
			switch {
			case addr < 0x2000:
				// If we are in MBC2, the least significant bit of the upper address byte is zero,
				// we can enable/disable cart RAM
				if IsBitSet(value, 4) {
					// If the lower nibble of the value being written == 0xA then enable RAM banking
					switch value & 0xF {
					case 0xA:
						c.RAMEnabled = true
					case 0x0:
						c.RAMEnabled = false
					}
				}
			}
		}
	}
}
