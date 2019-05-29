package apu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type Square2Channel struct {
	Square
}

func (c *Square2Channel) trigger(frameSequencerStep int) {
	c.enable = true

	if c.length == 0 {
		c.length = 64
		// If a channel is triggered when the frame sequencer's
		// next step is one that doesn't clock the length counter
		// and the length counter is now enabled and length is
		// being set to 64 (256 for wave channel) because it was
		// previously zero, it is set to 63 instead (255 for wave channel).
		if frameSequencerStep%2 == 0 && c.lengthEnable {
			c.length--
		}
	}

	c.timer = (2048 - int(c.frequency)) * 4

	// Volume envelope timer is reloaded with period
	// The volume envelope and sweep timers treat a period of 0 as 8
	c.volumeEnvelopeTimer = c.volumeEnvelopePeriod
	if c.volumeEnvelopePeriod == 0 {
		c.volumeEnvelopeTimer = 8
	}

	c.volume = c.volumeEnvelopeInitial

	c.frequencyShadow = c.frequency

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !c.DACEnable {
		c.enable = false
	}

}

func (c *Square2Channel) WriteTriggerByte(value byte, frameSequencerStep int) {
	// Bit 7   - Initial (1=Restart Sound)
	// Bit 6   - Counter/consecutive selection
	// 		  (1=Stop output when length in NR11 expires)
	// Bit 2-0 - Frequency's higher 3 bits (x)
	c.lengthEnable = utils.IsBitSet(value, 6)
	c.writeFrequencyHigherBits(value)

	// Make sure we trigger after the lengthEnable and Higher frequency bits are set
	if utils.IsBitSet(value, 7) {
		c.trigger(frameSequencerStep)
	}
}
