package mmu

import (
	"GopherBoy/utils"
)

type Memory interface {
	ReadByte(addr uint16) byte
	WriteByte(addr uint16, b byte)
}

type MMU struct {
	locations map[uint16]Memory
}

func NewMMU() *MMU {
	mmu := &MMU{
		locations: make(map[uint16]Memory),
	}

	return mmu
}

func (m *MMU) MapMemory(memory Memory, addr uint16) {
	m.locations[addr] = memory
}

func (m *MMU) MapMemoryRange(memory Memory, startAddr uint16, endAddr uint16) {
	for addr := startAddr; addr <= endAddr; addr++ {
		m.locations[addr] = memory
	}
}

func (m *MMU) RequestInterrupt(interrupt byte) {
	var IFAddress uint16 = 0xFF0F

	IF := utils.SetBit(m.ReadByte(IFAddress), interrupt)
	m.WriteByte(IFAddress, IF)
}

func (m *MMU) ReadByte(addr uint16) byte {
	if l := m.locations[addr]; l != nil {
		return l.ReadByte(addr)
	}
	return 0
}

func (m *MMU) WriteByte(addr uint16, value byte) {
	if l := m.locations[addr]; l != nil {
		l.WriteByte(addr, value)
	}
}
