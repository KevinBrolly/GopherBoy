package apu

import "github.com/kevinbrolly/GopherBoy/utils"

type Channel struct {
	DACEnable bool
	enable    bool
	timer     int

	length       int
	lengthEnable bool

	frequency       uint16
	frequencyShadow uint16
}

func (c *Channel) TickLength() {
	if c.lengthEnable && c.length > 0 {
		c.length--

		if c.length == 0 {
			c.enable = false
		}
	}
}

func (c *Channel) writeFrequencyLowerBits(value byte) {
	// value = lower 8 bits of 11 bit frequency. Next 3 bits are in NRx4
	c.frequency = (c.frequency & 0x700) | uint16(value)
}

func (c *Channel) writeFrequencyHigherBits(value byte) {
	// value = Bit 2-0 - Frequency's higher 3 bits (x) (Write Only)
	frequencyHighBits := uint16(value&0x7) << 8
	c.frequency = (c.frequency & 0xFF) | uint16(frequencyHighBits)
}

type VolumeEnvelope struct {
	volume                  byte
	volumeEnvelopeInitial   byte
	volumeEnvelopeDirection bool
	volumeEnvelopePeriod    byte
	volumeEnvelopeTimer     byte
}

func (v *VolumeEnvelope) TickVolumeEnvelope() {
	if v.volumeEnvelopePeriod > 0 {
		if v.volumeEnvelopeTimer > 0 {
			v.volumeEnvelopeTimer--
		}
		if v.volumeEnvelopeTimer == 0 {
			if v.volumeEnvelopeDirection {
				if v.volume < 0xF {
					v.volume++
				}
			} else {
				if v.volume > 0 {
					v.volume--
				}
			}

			v.volumeEnvelopeTimer = v.volumeEnvelopePeriod
		}
	}
}

func (v *VolumeEnvelope) volumeEnvelopeReadByte() byte {
	// Bit 7-4 - Initial Volume of envelope (0-0Fh) (0=No Sound)
	// Bit 3   - Envelope Direction (0=Decrease, 1=Increase)
	// Bit 2-0 - Number of envelope sweep (n: 0-7)
	var value byte
	value = v.volumeEnvelopeInitial << 4
	if v.volumeEnvelopeDirection {
		value = utils.SetBit(value, 3)
	}
	value = value | v.volumeEnvelopePeriod
	return value
}

func (v *VolumeEnvelope) volumeEnvelopeWriteByte(value byte) {
	// Bit 7-4 - Initial Volume of envelope (0-0Fh) (0=No Sound)
	// Bit 3   - Envelope Direction (0=Decrease, 1=Increase)
	// Bit 2-0 - Number of envelope sweep (n: 0-7)
	v.volumeEnvelopeInitial = (value >> 4) & 0x0f
	v.volumeEnvelopeDirection = utils.IsBitSet(value, 3)
	v.volumeEnvelopePeriod = value & 0x7
}

type Square struct {
	Channel
	VolumeEnvelope

	wavePatternDuty         byte
	wavePatternDutyPosition byte
}

func (c *Square) Tick(tCycles int) {
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

func (c *Square) sample() byte {
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

func (c *Square) ReadByte(addr uint16) byte {
	var value byte

	switch addr {
	case NR11, NR21:
		// Bit 7-6 - Wave Pattern Duty
		value = (c.wavePatternDuty << 6)
	case NR12, NR22:
		value = c.volumeEnvelopeReadByte()
	case NR14, NR24:
		// Bit 6   - Counter/consecutive selection
		if c.lengthEnable {
			value = utils.SetBit(value, 6)
		}
	}

	return value
}

func (c *Square) WriteByte(addr uint16, value byte) {
	switch addr {
	case NR11, NR21:
		// Bit 7-6 - Wave Pattern Duty
		// Bit 5-0 - Sound length data
		c.wavePatternDuty = (value >> 6) & 0x3
		c.length = int(value & 0x3F)
	case NR12, NR22:
		c.volumeEnvelopeWriteByte(value)
		c.DACEnable = (value & 0xf8) > 0

		// Any time the DAC is off the channel is kept disabled
		if !c.DACEnable {
			c.enable = false
		}
	case NR13, NR23:
		c.writeFrequencyLowerBits(value)
	}
}
