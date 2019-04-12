package apu

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
)

const (
	NR41 = 0xFF20
	NR42 = 0xFF21
	NR43 = 0xFF22
	NR44 = 0xFF23
)

type Channel4 struct {
	Channel

	// Linear Feedback Shift Register — 15-bit
	LFSR uint16

	shiftClockFrequency byte
	counterWidth        bool
	dividingRatio       byte
}

func NewChannel4(mmu *mmu.MMU) *Channel4 {
	channel := &Channel4{}

	mmu.MapMemory(channel, NR41)
	mmu.MapMemory(channel, NR42)
	mmu.MapMemory(channel, NR43)
	mmu.MapMemory(channel, NR44)

	return channel
}

func (c *Channel4) trigger() {
	// Channel is enabled
	c.enable = true

	// If length counter is zero, it is set to 64.
	if c.length == 0 {
		c.length = 64
	}

	// Frequency timer is reloaded with period.
	c.timer = 524288/c.getDividingRatio()/2 ^ (int(c.frequency) + 1)

	// Volume envelope timer is reloaded with period.
	c.volumeEnvelopeTimer = c.volumeEnvelopePeriod
	// Channel volume is reloaded from NRx2.
	c.volume = c.volumeEnvelopeInitial

	// Noise channel's LFSR bits are all set to 1.
	c.LFSR = 0x7FFF
}

func (c *Channel4) Tick(tCycles int) {
	if c.timer > 0 {
		c.timer -= tCycles
	}

	if c.timer == 0 {
		// When clocked by the frequency timer, the low two bits (0 and 1)
		// are XORed
		bit := (c.LFSR & 0x1) ^ (c.LFSR & 0x2)
		// All LFSR bits are shifted right by one
		c.LFSR = c.LFSR >> 1
		// And the result of the XOR is put into the now-empty high bit
		if bit == 1 {
			c.LFSR = c.LFSR | 0x4000
		}
		// If width mode is 1 (NR43), the XOR result is ALSO put into
		// bit 6 AFTER the shift, resulting in a 7-bit LFSR.
		if c.counterWidth {
			if bit == 1 {
				c.LFSR = c.LFSR | 0x40
			} else {
				c.LFSR = c.LFSR &^ 0x40
			}
		}

		// Reload timer
		c.timer = 524288/c.getDividingRatio()/2 ^ (int(c.frequency) + 1)
	}
}

func (c *Channel4) sample() byte {
	// The waveform output is bit 0 of the LFSR, INVERTED
	bit := c.LFSR & 0x1

	if bit == 1 {
		return c.volume
	}

	return 0
}

func (c *Channel4) getDividingRatio() int {
	// Divisor code   Divisor
	// -----------------------
	// 0             8
	// 1            16
	// 2            32
	// 3            48
	// 4            64
	// 5            80
	// 6            96
	// 7            112
	if c.dividingRatio == 0 {
		return 8
	}

	return int(c.dividingRatio * 16)
}

func (c *Channel4) ReadByte(addr uint16) byte {
	switch {
	case addr == NR41:
		// Bit 5-0 - Sound length data
		return byte(c.length)
	case addr == NR42:
		return c.volumeEnvelopeReadByte()
	case addr == NR43:
		// Bit 7-4 - Shift Clock Frequency (s)
		// Bit 3   - Counter Step/Width (0=15 bits, 1=7 bits)
		// Bit 2-0 - Dividing Ratio of Frequencies (r)
		var value byte
		value = c.shiftClockFrequency << 4
		if c.counterWidth {
			value = utils.SetBit(value, 3)
		}
		value = value | c.dividingRatio
		return value
	case addr == NR44:
		// Bit 6   - Counter/consecutive selection
		var value byte
		if c.lengthEnable {
			value = utils.SetBit(value, 6)
		}
		return value
	}
	return 0
}

func (c *Channel4) WriteByte(addr uint16, value byte) {
	switch {
	case addr == NR41:
		// Bit 5-0 - Sound length data
		c.length = int(value)
	case addr == NR42:
		c.volumeEnvelopeWriteByte(value)
	case addr == NR43:
		c.shiftClockFrequency = value >> 4
		c.counterWidth = utils.IsBitSet(value, 3)
		c.dividingRatio = value & 0x7
	case addr == NR44:
		// Bit 7   - Initial (1=Restart Sound)
		// Bit 6   - Counter/consecutive selection
		// 		  (1=Stop output when length in NR11 expires)
		if utils.IsBitSet(value, 7) {
			c.trigger()
		}
		c.lengthEnable = utils.IsBitSet(value, 6)
	}
}
