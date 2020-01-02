package ppu

import (
	"fmt"
	"testing"
)

type MockMMU struct {
	data []byte
}

func (m *MockMMU) ReadByte(addr uint16) byte {
	return m.data[addr]
}

func (m *MockMMU) WriteByte(addr uint16, value byte) {
	m.data[addr] = value
}

func TestGetTileIdentifier(t *testing.T) {
	data := make([]byte, 0x01)
	data[0x0000] = 0xFF

	mmu := &MockMMU{
		data: data,
	}

	fetcher := &Fetcher{
		tileMapAddress:  0x0000,
		tileDataAddress: 0x8000,
		SCY:             0,
		memory:          mmu,
	}

	if fetcher.getTileIdentifier() != 0xFF {
		t.Errorf("getTileIdentifier() should have returned 0xFF but returned %x", fetcher.getTileIdentifier())
	}
}

func TestGetTileDataAddress(t *testing.T) {
	mmu := &MockMMU{
		data: make([]byte, 0x0),
	}

	tileDataAddress1 := uint16(0x8000)
	tileDataAddress2 := uint16(0x8800)

	fetcher1 := &Fetcher{
		tileMapAddress:  0x0000,
		tileDataAddress: tileDataAddress1,
		SCY:             0,
		memory:          mmu,
	}

	fetcher2 := &Fetcher{
		tileMapAddress:  0x0000,
		tileDataAddress: tileDataAddress2,
		SCY:             0,
		memory:          mmu,
	}

	cases := []struct {
		Fetcher         *Fetcher
		TileIdentifier  byte
		ExpectedAddress uint16
	}{
		{fetcher1, 0, 0x8000},
		{fetcher1, 1, 0x8010},
		{fetcher1, 128, 0x8800},
		{fetcher1, 255, 0x8FF0},
		{fetcher2, 0, 0x9000},
		{fetcher2, 1, 0x9010},
		{fetcher2, 127, 0x97F0},
		{fetcher2, 128, 0x8800},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("Tile Identifier %v", tt.TileIdentifier), func(t *testing.T) {
			dataAddress := tt.Fetcher.getTileDataAddress(tt.TileIdentifier)
			if dataAddress != tt.ExpectedAddress {
				t.Errorf("getTileIdentifier() should have returned %#x but returned %#x. tileDataAddress: %#x", tt.ExpectedAddress, dataAddress, tt.Fetcher.tileDataAddress)
			}
		})
	}
}

func TestGetTileLine(t *testing.T) {
	data := make([]byte, 0x10000)
	data[0x8000] = 0xFF
	data[0x8001] = 0xFF

	mmu := &MockMMU{
		data: data,
	}

	fetcher := &Fetcher{
		tileMapAddress:  0x0000,
		tileDataAddress: 0x8000,
		SCY:             0,
		memory:          mmu,
	}

	TileLine := fetcher.fetchTileLine(0)
	for _, pixel := range TileLine {
		if pixel.ColorIdentifier != 0x3 {
			t.Errorf("getTileLine() should have returned pixels with ColourIdentifier 3 but instead returned %x", pixel.ColorIdentifier)
		}
	}
}

func TestGetNextTileLine(t *testing.T) {
	data := make([]byte, 0x10000)
	data[0x0000] = 0x00
	data[0x0001] = 0x01
	data[0x0002] = 0x02

	data[0x8000] = 0xFF
	data[0x8001] = 0xFF
	data[0x8002] = 0xFF
	data[0x8003] = 0xFF

	data[0x8010] = 0x00
	data[0x8011] = 0xFF
	data[0x8012] = 0x00
	data[0x8013] = 0xFF

	data[0x8020] = 0xFF
	data[0x8021] = 0x00
	data[0x8022] = 0xFF
	data[0x8023] = 0x00

	mmu := &MockMMU{
		data: data,
	}

	fetcher := &Fetcher{
		tileMapAddress:  0x0000,
		tileDataAddress: 0x8000,
		LY:              0,
		SCY:             0,
		memory:          mmu,
	}

	cases := []struct {
		Fetcher                  *Fetcher
		ExpectedColourIdentifier byte
	}{
		{fetcher, 0x3}, // 11
		{fetcher, 0x2}, // 10
		{fetcher, 0x1}, // 01
	}

	for _, tt := range cases {
		tileLine := tt.Fetcher.NextTileLine()
		for _, pixel := range tileLine {
			if pixel.ColorIdentifier != tt.ExpectedColourIdentifier {
				t.Errorf("getTileLine() should have returned pixels with ColourIdentifier %x but instead returned %x", tt.ExpectedColourIdentifier, pixel.ColorIdentifier)
			}
		}
	}

	fetcher.tileMapAddress = 0x0000
	fetcher.LY = 1

	for _, tt := range cases {
		tileLine := tt.Fetcher.NextTileLine()
		for _, pixel := range tileLine {
			if pixel.ColorIdentifier != tt.ExpectedColourIdentifier {
				t.Errorf("getTileLine() should have returned pixels with ColourIdentifier %x but instead returned %x", tt.ExpectedColourIdentifier, pixel.ColorIdentifier)
			}
		}
	}
}
