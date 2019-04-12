package apu

import "github.com/kevinbrolly/GopherBoy/utils"

type Channel struct {
	DACEnable bool
	enable    bool
	timer     int

	length       int
	lengthEnable bool

	volume                  byte
	volumeEnvelopeInitial   byte
	volumeEnvelopeDirection bool
	volumeEnvelopePeriod    byte
	volumeEnvelopeTimer     byte

	frequency       uint16
	frequencyShadow uint16
}

func (c *Channel) TickLength() {
	if c.length > 0 {
		c.length--
	}

	if c.length == 0 && c.lengthEnable {
		c.enable = false
	}
}

func (c *Channel) TickVolumeEnvelope() {
	if c.volumeEnvelopePeriod > 0 {
		if c.volumeEnvelopeTimer > 0 {
			c.volumeEnvelopeTimer--
		}
		if c.volumeEnvelopeTimer == 0 {
			if c.volumeEnvelopeDirection {
				if c.volume < 0xF {
					c.volume++
				}
			} else {
				if c.volume > 0 {
					c.volume--
				}
			}

			c.volumeEnvelopeTimer = c.volumeEnvelopePeriod
		}
	}
}

func (c *Channel) volumeEnvelopeReadByte() byte {
	// Bit 7-4 - Initial Volume of envelope (0-0Fh) (0=No Sound)
	// Bit 3   - Envelope Direction (0=Decrease, 1=Increase)
	// Bit 2-0 - Number of envelope sweep (n: 0-7)
	var value byte
	value = c.volumeEnvelopeInitial << 4
	if c.volumeEnvelopeDirection {
		value = utils.SetBit(value, 3)
	}
	value = value | c.volumeEnvelopePeriod
	return value
}

func (c *Channel) volumeEnvelopeWriteByte(value byte) {
	// Bit 7-4 - Initial Volume of envelope (0-0Fh) (0=No Sound)
	// Bit 3   - Envelope Direction (0=Decrease, 1=Increase)
	// Bit 2-0 - Number of envelope sweep (n: 0-7)
	c.volumeEnvelopeInitial = value >> 4
	c.volumeEnvelopeDirection = utils.IsBitSet(value, 3)
	c.volumeEnvelopePeriod = (value & 0x3)
}

func (c *Channel) writeHighFrequency(value byte) {
	// value = lower 8 bits of 11 bit frequency. Next 3 bits are in NRx4
	c.frequency = (c.frequency & 0x700) | uint16(value)
}
