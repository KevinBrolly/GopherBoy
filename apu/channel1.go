package apu

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
)

const (
	NR10 = 0xFF10
	NR11 = 0xFF11
	NR12 = 0xFF12
	NR13 = 0xFF13
	NR14 = 0xFF14
)

type Channel1 struct {
	Channel

	sweepTimer  byte
	sweepEnable bool
	sweepPeriod byte
	sweepNegate bool
	sweepShift  byte

	wavePatternDuty         byte
	wavePatternDutyPosition byte
}

func NewChannel1(mmu *mmu.MMU) *Channel1 {
	channel := &Channel1{}

	mmu.MapMemory(channel, NR10)
	mmu.MapMemory(channel, NR11)
	mmu.MapMemory(channel, NR12)
	mmu.MapMemory(channel, NR13)
	mmu.MapMemory(channel, NR14)

	return channel
}

func (c *Channel1) trigger() {
	c.enable = true

	if c.length == 0 {
		c.length = 64
	}

	c.timer = (2048 - int(c.frequency)) * 4

	c.volumeEnvelopeTimer = c.volumeEnvelopePeriod
	c.volume = c.volumeEnvelopeInitial

	c.frequencyShadow = c.frequency
	c.sweepTimer = c.sweepPeriod

	if c.sweepPeriod > 0 || c.sweepShift > 0 {
		c.sweepEnable = true
	} else {
		c.sweepEnable = false
	}

	if c.sweepShift != 0 {
		c.calculateSweep()
	}

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !c.DACEnable {
		c.enable = false
	}
}

func (c *Channel1) Tick(tCycles int) {
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

func (c *Channel1) sample() byte {
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

func (c *Channel1) calculateSweep() uint16 {
	var newFrequency int
	newFrequency = int(c.frequencyShadow) >> c.sweepShift

	if c.sweepNegate {
		newFrequency = -newFrequency
	}

	newFrequency += int(c.frequencyShadow)

	// Overflow Check
	// If the new frequency is over 2047, disable the channel
	if newFrequency > 2047 {
		c.enable = false
	}

	return uint16(newFrequency)
}

func (c *Channel1) TickSweep() {
	if c.sweepTimer > 0 {
		c.sweepTimer--
	}

	if (c.sweepTimer == 0) && c.sweepEnable && (c.sweepPeriod > 0) {
		newFrequency := c.calculateSweep()
		if newFrequency <= 2047 && c.sweepShift != 0 {
			c.frequencyShadow = newFrequency
			c.frequency = newFrequency

			c.calculateSweep()
		}

		c.sweepTimer = c.sweepPeriod
	}
}

func (c *Channel1) ReadByte(addr uint16) byte {
	switch {
	case addr == NR10:
		// Bit 6-4 - Sweep Time
		// Bit 3   - Sweep Increase/Decrease
		// 	0: Addition    (frequency increases)
		// 	1: Subtraction (frequency decreases)
		// Bit 2-0 - Number of sweep shift (n: 0-7)
		var value byte
		value = (c.sweepPeriod << 4)

		if c.sweepNegate {
			value = utils.SetBit(value, 3)
		}

		value = value & c.sweepShift
		return value
	case addr == NR11:
		// Bit 7-6 - Wave Pattern Duty
		value := (c.wavePatternDuty << 6)
		return value
	case addr == NR12:
		return c.volumeEnvelopeReadByte()
	case addr == NR14:
		// Bit 6   - Counter/consecutive selection
		var value byte
		if c.lengthEnable {
			value = utils.SetBit(value, 6)
		}
		return value
	}
	return 0
}

func (c *Channel1) WriteByte(addr uint16, value byte) {
	switch {
	case addr == NR10:
		// Bit 6-4 - Sweep Time
		// Bit 3   - Sweep Increase/Decrease
		// 	0: Addition    (frequency increases)
		// 	1: Subtraction (frequency decreases)
		// Bit 2-0 - Number of sweep shift (n: 0-7)
		c.sweepPeriod = (value & 0x70) >> 4
		c.sweepNegate = utils.IsBitSet(value, 3)
		c.sweepShift = value & 0x7
	case addr == NR11:
		// Bit 7-6 - Wave Pattern Duty
		// Bit 5-0 - Sound length data
		c.wavePatternDuty = (value >> 6) & 0x3
		c.length = int(value & 0x3F)
	case addr == NR12:
		c.volumeEnvelopeWriteByte(value)
		c.DACEnable = (value & 0xf8) > 0
	case addr == NR13:
		c.writeFrequencyLowerBits(value)
	case addr == NR14:
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
