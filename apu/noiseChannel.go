package apu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type NoiseChannel struct {
	Channel
	VolumeEnvelope

	// Linear Feedback Shift Register — 15-bit
	LFSR uint16

	shiftClockFrequency byte
	counterWidth        bool
	dividingRatio       byte
}

func (c *NoiseChannel) trigger(frameSequencerStep int) {
	// Channel is enabled
	c.enable = true

	// If length counter is zero, it is set to 64.
	if c.length == 0 {
		c.length = 64
		// If a channel is triggered when the frame sequencer's
		// next step is one that doesn't clock the length counter
		// and the length counter is now enabled and length is
		// being set to 64 (256 for wave channel) because it was
		// previously zero, it is set to 63 instead (255 for wave channel).
		if frameSequencerStep%2 != 0 && c.lengthEnable {
			c.length--
		}
	}

	// Frequency timer is reloaded with period.
	c.timer = c.getDividingRatio() << c.shiftClockFrequency

	// Volume envelope timer is reloaded with period
	// The volume envelope and sweep timers treat a period of 0 as 8
	c.volumeEnvelopeTimer = c.volumeEnvelopePeriod
	if c.volumeEnvelopePeriod == 0 {
		c.volumeEnvelopeTimer = 8
	}

	// Channel volume is reloaded from NRx2.
	c.volume = c.volumeEnvelopeInitial

	// Noise channel's LFSR bits are all set to 1.
	c.LFSR = 0x7FFF

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !c.DACEnable {
		c.enable = false
	}
}

func (c *NoiseChannel) Tick(tCycles int) {
	if c.timer > 0 {
		c.timer -= tCycles
	}

	if c.timer <= 0 {
		// When clocked by the frequency timer, the low two bits (0 and 1)
		// are XORed
		bit := (c.LFSR & 0x1) ^ ((c.LFSR >> 1) & 0x1)
		// All LFSR bits are shifted right by one
		c.LFSR = c.LFSR >> 1
		// And the result of the XOR is put into the now-empty high bit
		c.LFSR = c.LFSR | (bit << 14)
		// If width mode is 1 (NR43), the XOR result is ALSO put into
		// bit 6 AFTER the shift, resulting in a 7-bit LFSR.
		if c.counterWidth {
			c.LFSR = c.LFSR &^ 0x40
			c.LFSR = c.LFSR | (bit << 6)
		}

		// Reload timer
		c.timer = c.getDividingRatio() << c.shiftClockFrequency
	}
}

func (c *NoiseChannel) sample() byte {
	if c.enable && c.DACEnable {
		// The waveform output is bit 0 of the LFSR, INVERTED
		bit := c.LFSR & 0x1

		if bit == 1 {
			return c.volume
		}
	}

	return 0
}

func (c *NoiseChannel) getDividingRatio() int {
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

func (c *NoiseChannel) ReadByte(addr uint16) byte {
	var value byte

	switch addr {
	case NR41:
		// Bit 5-0 - Sound length data
		value = byte(c.length)
	case NR42:
		value = c.volumeEnvelopeReadByte()
	case NR43:
		// Bit 7-4 - Shift Clock Frequency (s)
		// Bit 3   - Counter Step/Width (0=15 bits, 1=7 bits)
		// Bit 2-0 - Dividing Ratio of Frequencies (r)
		value = c.shiftClockFrequency << 4
		if c.counterWidth {
			value = utils.SetBit(value, 3)
		}
		value = value | c.dividingRatio
	case NR44:
		// Bit 6   - Counter/consecutive selection
		if c.lengthEnable {
			value = utils.SetBit(value, 6)
		}
	}

	return value
}

func (c *NoiseChannel) WriteByte(addr uint16, value byte) {
	switch addr {
	case NR41:
		// Bit 5-0 - Sound length data

		// Writing a byte to NRx1 loads the length counter with 64 - data
		c.length = 64 - int(value&0x3f)
	case NR42:
		c.volumeEnvelopeWriteByte(value)
		c.DACEnable = (value & 0xf8) > 0

		// Any time the DAC is off the channel is kept disabled
		if !c.DACEnable {
			c.enable = false
		}
	case NR43:
		c.shiftClockFrequency = value >> 4
		c.counterWidth = utils.IsBitSet(value, 3)
		c.dividingRatio = value & 0x7
	}
}

func (c *NoiseChannel) WriteTriggerByte(value byte, frameSequencerStep int) {
	// Bit 7   - Initial (1=Restart Sound)
	// Bit 6   - Counter/consecutive selection
	// 		  (1=Stop output when length in NR11 expires)
	c.lengthEnable = utils.IsBitSet(value, 6)

	// Make sure we trigger after lengthEnable is set
	if utils.IsBitSet(value, 7) {
		c.trigger(frameSequencerStep)
	}
}
