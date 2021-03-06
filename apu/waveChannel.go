package apu

import (
	"github.com/kevinbrolly/GopherBoy/utils"
)

type WaveChannel struct {
	Channel

	volume   byte
	position byte
	buffer   byte

	wavePatternRAM [16]byte
}

func (c *WaveChannel) trigger(frameSequencerStep int) {
	c.enable = true

	if c.length == 0 {
		c.length = 256
		// If a channel is triggered when the frame sequencer's
		// next step is one that doesn't clock the length counter
		// and the length counter is now enabled and length is
		// being set to 64 (256 for wave channel) because it was
		// previously zero, it is set to 63 instead (255 for wave channel).
		if frameSequencerStep%2 != 0 && c.lengthEnable {
			c.length--
		}
	}

	c.timer = (2048 - int(c.frequency)) * 2

	// Wave channel's position is set to 0 but sample buffer is NOT refilled.
	c.position = 0

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !c.DACEnable {
		c.enable = false
	}
}

func (c *WaveChannel) Tick(tCycles int) {
	if c.timer > 0 {
		c.timer -= tCycles
	}
	if c.timer <= 0 {
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

		c.position++
		if c.position == 32 {
			c.position = 0
		}
	}
}

func (c *WaveChannel) sample() byte {
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

func (c *WaveChannel) ReadByte(addr uint16) byte {
	var value byte

	switch {
	case addr == NR30:
		// Bit 7 - Sound Channel 3 Off  (0=Stop, 1=Playback)
		if c.DACEnable {
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
		// If the wave channel is enabled, accessing any byte from $FF30-$FF3F is equivalent
		// to accessing the current byte selected by the waveform position.
		if c.enable {
			value = c.wavePatternRAM[c.position]
		} else {
			value = c.wavePatternRAM[addr&0xF]
		}
	}

	return value
}

func (c *WaveChannel) WriteByte(addr uint16, value byte) {
	switch {
	case addr == NR30:
		// Bit 7 - Sound Channel 3 Off  (0=Stop, 1=Playback)
		c.DACEnable = utils.IsBitSet(value, 7)

		// Any time the DAC is off the channel is kept disabled
		if !c.DACEnable {
			c.enable = false
		}
	case addr == NR31:
		// Writing a byte to NRx1 loads the wave channel length counter with 256 - data
		c.length = 256 - int(value)
	case addr == NR32:
		// Bit 6-5 - Select output level
		c.volume = (value >> 5) & 0x3
	case addr == NR33:
		c.writeFrequencyLowerBits(value)
	case addr >= 0xFF30 && addr <= 0xFF3F:
		// If the wave channel is enabled, accessing any byte from $FF30-$FF3F is equivalent
		// to accessing the current byte selected by the waveform position.
		// On the DMG accesses will only work in this manner if made within a couple of clocks
		// of the wave channel accessing wave RAM; if made at any other time, reads return $FF and writes have no effect.
		if c.enable {
			c.wavePatternRAM[c.position] = value
		} else {
			c.wavePatternRAM[addr&0xF] = value
		}
	}
}

func (c *WaveChannel) WriteTriggerByte(value byte, frameSequencerStep int) {
	// Bit 7   - Initial (1=Restart Sound)
	// Bit 6   - Counter/consecutive selection
	// 		  (1=Stop output when length in NR34 expires)
	// Bit 2-0 - Frequency's higher 3 bits (x)
	c.lengthEnable = utils.IsBitSet(value, 6)

	c.writeFrequencyHigherBits(value)

	// Make sure we trigger after the lengthEnable and Higher frequency bits are set
	if utils.IsBitSet(value, 7) {
		c.trigger(frameSequencerStep)
	}
}
