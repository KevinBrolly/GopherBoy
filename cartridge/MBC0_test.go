package cartridge

import (
	"testing"
)

func TestWriteByte(t *testing.T) {
	// Writes should silently fail

	mmu := NewMMU()
	mbc0 := NewMBC0(mmu, make([]byte, 10))

	mbc0.WriteByte(0, 1)

	if mbc0.ReadByte(0) != 0 {
		t.Errorf("WriteByte() should have failed to write but didnt")
	}
}

func TestReadByte(t *testing.T) {
	// Test we can read a byte

	data := make([]byte, 1)
	data[0] = 0xFF

	mmu := NewMMU()
	mbc0 := NewMBC0(mmu, data)

	if mbc0.ReadByte(0) != 0xFF {
		t.Errorf("ReadByte() should have read a value but didnt")
	}
}
