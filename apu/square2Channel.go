package apu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type Square2Channel struct {
	Square
}

func (c *Square2Channel) trigger() {
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

func (c *Square2Channel) WriteTriggerByte(value byte) {
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
