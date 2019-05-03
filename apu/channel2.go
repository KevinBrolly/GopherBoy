package apu

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
)

const (
	NR20 = 0xFF15
	NR21 = 0xFF16
	NR22 = 0xFF17
	NR23 = 0xFF18
	NR24 = 0xFF19
)

type Channel2 struct {
	Channel

	wavePatternDuty         byte
	wavePatternDutyPosition byte
}

func NewChannel2(mmu *mmu.MMU) *Channel2 {
	channel := &Channel2{}

	mmu.MapMemory(channel, NR20)
	mmu.MapMemory(channel, NR21)
	mmu.MapMemory(channel, NR22)
	mmu.MapMemory(channel, NR23)
	mmu.MapMemory(channel, NR24)

	return channel
}

func (c *Channel2) trigger() {
	c.enable = true

	if c.length == 0 {
		c.length = 64
	}

	c.timer = (2048 - int(c.frequency)) * 4

	c.volumeEnvelopeTimer = c.volumeEnvelopePeriod
	c.volume = c.volumeEnvelopeInitial

	c.frequencyShadow = c.frequency

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !c.DACEnable {
		c.enable = false
	}
}

func (c *Channel2) Tick(tCycles int) {
	// A square channel's frequency timer period is set to (2048-frequency)*4.
	// Four duty cycles are available, each waveform taking 8 frequency timer clocks to cycle through:
	// Duty   Waveform    Ratio
	// -------------------------
	// 0      00000001    12.5%
	// 1      10000001    25%
	// 2      10000111    50%
	// 3      01111110    75%

	if c.timer > 0 {
		c.timer -= tCycles
	}
	if c.timer <= 0 {
		// Increment position of the duty waveform
		c.wavePatternDutyPosition++
		if c.wavePatternDutyPosition == 8 {
			c.wavePatternDutyPosition = 0
		}

		// Reload timer
		c.timer = (2048 - int(c.frequency)) * 4
	}
}

func (c *Channel2) sample() byte {
	if !c.enable && !c.DACEnable {
		return 0
	}

	var pattern byte
	switch c.wavePatternDuty {
	case 0:
		pattern = 0x1 // 00000001
	case 1:
		pattern = 0x81 // 10000001
	case 2:
		pattern = 0x87 // 10000111
	case 3:
		pattern = 0x7E // 01111110
	}

	if utils.IsBitSet(pattern, (7 - c.wavePatternDutyPosition)) {
		return c.volume
	}

	return 0
}

func (c *Channel2) ReadByte(addr uint16) byte {
	var value byte

	switch {
	case addr == NR21:
		// Bit 7-6 - Wave Pattern Duty
		value = (c.wavePatternDuty << 6)
	case addr == NR22:
		value = c.volumeEnvelopeReadByte()
	case addr == NR24:
		// Bit 6   - Counter/consecutive selection
		if c.lengthEnable {
			value = utils.SetBit(value, 6)
		}
	}

	return value | apuReadMask[addr]
}

func (c *Channel2) WriteByte(addr uint16, value byte) {
	switch {
	case addr == NR21:
		// Bit 7-6 - Wave Pattern Duty
		// Bit 5-0 - Sound length data
		c.wavePatternDuty = (value >> 6) & 0x3
		c.length = int(value & 0x3F)
	case addr == NR22:
		c.volumeEnvelopeWriteByte(value)
		c.DACEnable = (value & 0xf8) > 0
	case addr == NR23:
		c.writeFrequencyLowerBits(value)
	case addr == NR24:
		// Bit 7   - Initial (1=Restart Sound)
		// Bit 6   - Counter/consecutive selection
		// 		  (1=Stop output when length in NR11 expires)
		// Bit 2-0 - Frequency's higher 3 bits (x)
		c.lengthEnable = utils.IsBitSet(value, 6)

		c.writeFrequencyHigherBits(value)

		// Make sure we trigger after the lengthEnable and Higher frequency bits are set
		if utils.IsBitSet(value, 7) {
			c.trigger()
		}
	}
}
