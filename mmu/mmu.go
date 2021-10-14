package mmu

import (
	"sync"

	"github.com/kevinbrolly/GopherBoy/utils"
)

type Memory interface {
	ReadByte(addr uint16) byte
	WriteByte(addr uint16, b byte)
}

type MMU struct {
	locations    map[uint16]Memory
	mutex        *sync.Mutex
	cycleChannel chan int
}

func NewMMU(cycleChannel chan int) *MMU {
	mmu := &MMU{
		locations:    make(map[uint16]Memory),
		mutex:        &sync.Mutex{},
		cycleChannel: cycleChannel,
	}

	return mmu
}

func (m *MMU) MapMemory(memory Memory, addr uint16) {
	m.mutex.Lock()
	m.locations[addr] = memory
	m.mutex.Unlock()
}

func (m *MMU) MapMemoryRange(memory Memory, startAddr uint16, endAddr uint16) {
	m.mutex.Lock()
	for addr := startAddr; addr <= endAddr; addr++ {
		m.locations[addr] = memory
	}
	m.mutex.Unlock()
}

func (m *MMU) RequestInterrupt(interrupt byte) {
	var IFAddress uint16 = 0xFF0F

	IF := utils.SetBit(m.ReadByte(IFAddress), interrupt)
	m.WriteByte(IFAddress, IF)
}

func (m *MMU) ReadByte(addr uint16) byte {
	m.mutex.Lock()
	var res byte
	if l := m.locations[addr]; l != nil {
		res = l.ReadByte(addr)
	} else {
		res = 0
	}
	m.mutex.Unlock()
	m.cycleChannel <- 1
	return res
}

func (m *MMU) WriteByte(addr uint16, value byte) {
	m.mutex.Lock()
	if l := m.locations[addr]; l != nil {
		l.WriteByte(addr, value)
	}
	m.cycleChannel <- 1
	m.mutex.Unlock()
}
