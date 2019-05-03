package apu

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
)

const (
	NR30 = 0xFF1A
	NR31 = 0xFF1B
	NR32 = 0xFF1C
	NR33 = 0xFF1D
	NR34 = 0xFF1E
)

type Channel3 struct {
	Channel

	position byte
	buffer   byte

	wavePatternRAM [16]byte
}

func NewChannel3(mmu *mmu.MMU) *Channel3 {
	channel := &Channel3{}

	mmu.MapMemory(channel, NR30)
	mmu.MapMemory(channel, NR31)
	mmu.MapMemory(channel, NR32)
	mmu.MapMemory(channel, NR33)
	mmu.MapMemory(channel, NR34)

	// wavePatternRAM
	mmu.MapMemoryRange(channel, 0xFF30, 0xFF3F)

	return channel
}

func (c *Channel3) trigger() {
	c.enable = true

	if c.length == 0 {
		c.length = 256
	}

	c.timer = (2048 - int(c.frequency)) * 2

	// Wave channel's position is set to 0 but sample buffer is NOT refilled.
	c.position = 0

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !c.DACEnable {
		c.enable = false
	}
}

func (c *Channel3) Tick(tCycles int) {
	if c.timer > 0 {
		c.timer -= tCycles
	}
	if c.timer <= 0 {
		c.position++
		if c.position == 32 {
			c.position = 0
		}

		// Fill the sample buffer
		// wavePatternRAM is 16 bytes, position is length 32
		// position / 2 = wavePatternRAM index
		wavePatternByte := c.wavePatternRAM[c.position/2]
		// wavePatternByte holds 2 4-bit samples
		// if position is even the low nibble is used, if odd the high nibble is used
		if c.position%2 == 0 {
			c.buffer = wavePatternByte & 0xF
		} else {
			c.buffer = wavePatternByte >> 4
		}
		// Reload timer
		c.timer = (2048 - int(c.frequency)) * 2
	}
}

func (c *Channel3) sample() byte {
	if c.enable && c.DACEnable {
		// The DAC receives the current value from the upper/lower nibble of the
		// sample buffer, shifted right by the volume control.
		// Code   Shift   Volume
		// -----------------------
		// 0      4         0% (silent)
		// 1      0       100%
		// 2      1        50%
		// 3      2        25%
		var shift byte
		switch c.volume {
		case 0:
			shift = 4
		case 1:
			shift = 0
		case 2:
			shift = 1
		case 3:
			shift = 2
		}
		return c.buffer >> shift
	}

	return 0

}

func (c *Channel3) ReadByte(addr uint16) byte {
	var value byte

	switch {
	case addr == NR30:
		// Bit 7 - Sound Channel 3 Off  (0=Stop, 1=Playback)
		if c.enable {
			value = utils.SetBit(value, 7)
		}
	case addr == NR32:
		value = c.volume << 5
	case addr == NR34:
		// Bit 6   - Counter/consecutive selection
		if c.lengthEnable {
			value = utils.SetBit(value, 6)
		}
	case addr >= 0xFF30 && addr <= 0xFF3F:
		value = c.wavePatternRAM[addr&0xF]
	}

	return value | apuReadMask[addr]
}

func (c *Channel3) WriteByte(addr uint16, value byte) {
	switch {
	case addr == NR30:
		// Bit 7 - Sound Channel 3 Off  (0=Stop, 1=Playback)
		c.enable = utils.IsBitSet(value, 7)
		c.DACEnable = utils.IsBitSet(value, 7)
	case addr == NR31:
		c.length = 256 - int(value)
	case addr == NR32:
		// Bit 6-5 - Select output level
		c.volume = (value >> 5) & 0x3
	case addr == NR33:
		c.writeFrequencyLowerBits(value)
	case addr == NR34:
		// Bit 7   - Initial (1=Restart Sound)
		// Bit 6   - Counter/consecutive selection
		// 		  (1=Stop output when length in NR34 expires)
		// Bit 2-0 - Frequency's higher 3 bits (x)
		c.lengthEnable = utils.IsBitSet(value, 6)

		c.writeFrequencyHigherBits(value)

		// Make sure we trigger after the lengthEnable and Higher frequency bits are set
		if utils.IsBitSet(value, 7) {
			c.trigger()
		}
	case addr >= 0xFF30 && addr <= 0xFF3F:
		c.wavePatternRAM[addr&0xF] = value
	}
}
