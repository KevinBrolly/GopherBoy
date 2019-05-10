package apu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type Square1Channel struct {
	Square

	sweepTimer  byte
	sweepEnable bool
	sweepPeriod byte
	sweepNegate bool
	sweepShift  byte
}

func (c *Square1Channel) trigger() {
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

func (c *Square1Channel) WriteTriggerByte(value byte) {
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

func (c *Square1Channel) calculateSweep() uint16 {
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

func (c *Square1Channel) TickSweep() {
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

func (c *Square1Channel) ReadSweep(addr uint16) byte {
	var value byte

	// Bit 6-4 - Sweep Time
	// Bit 3   - Sweep Increase/Decrease
	// 	0: Addition    (frequency increases)
	// 	1: Subtraction (frequency decreases)
	// Bit 2-0 - Number of sweep shift (n: 0-7)
	value = (c.sweepPeriod << 4)

	if c.sweepNegate {
		value = utils.SetBit(value, 3)
	}

	value = value | c.sweepShift

	return value
}

func (c *Square1Channel) WriteSweep(value byte) {
	// Bit 6-4 - Sweep Time
	// Bit 3   - Sweep Increase/Decrease
	// 	0: Addition    (frequency increases)
	// 	1: Subtraction (frequency decreases)
	// Bit 2-0 - Number of sweep shift (n: 0-7)
	c.sweepPeriod = (value & 0x70) >> 4
	c.sweepNegate = utils.IsBitSet(value, 3)
	c.sweepShift = value & 0x7
}
