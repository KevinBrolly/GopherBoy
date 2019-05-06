package apu

import (
	"bytes"
	"encoding/binary"

	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	NR50 = 0xFF24
	NR51 = 0xFF25
	NR52 = 0xFF26

	// The frame sequencer runs at 512 Hz, which is 4194304/512=8192 clock cycles
	FrameSequencerTimerRate = 8192
	Frequency               = 44100
	Samples                 = 2048
)

// When an NRxx register is read back, the last written value ORed with the following is returned:
var apuReadMask = map[uint16]byte{
	NR10: 0x80,
	NR11: 0x3F,
	NR12: 0x00,
	NR13: 0xFF,
	NR14: 0xBF,
	NR20: 0xFF,
	NR21: 0x3F,
	NR22: 0x00,
	NR23: 0xFF,
	NR24: 0xBF,
	NR30: 0x7F,
	NR31: 0xFF,
	NR32: 0x9F,
	NR33: 0xFF,
	NR34: 0xBF,
	NR40: 0xFF,
	NR41: 0xFF,
	NR42: 0x00,
	NR43: 0x00,
	NR44: 0xBF,
	NR50: 0x00,
	NR51: 0x00,
	NR52: 0x70,

	0xFF27: 0xFF,
	0xFF28: 0xFF,
	0xFF29: 0xFF,
	0xFF2A: 0xFF,
	0xFF2B: 0xFF,
	0xFF2C: 0xFF,
	0xFF2D: 0xFF,
	0xFF2E: 0xFF,
	0xFF2F: 0xFF,
}

type APU struct {
	mmu      *mmu.MMU
	channel1 *Channel1
	channel2 *Channel2
	channel3 *Channel3
	channel4 *Channel4

	sampleTimer  int
	sampleBuffer *bytes.Buffer
	sampleCount  int

	frameSequencerTimer int
	frameSequencerStep  int

	// NR50
	outputVinSO1 bool
	outputVinSO2 bool
	volumeSO1    byte
	volumeSO2    byte

	// NR51
	output4SO1 bool
	output3SO1 bool
	output2SO1 bool
	output1SO1 bool
	output4SO2 bool
	output3SO2 bool
	output2SO2 bool
	output1SO2 bool

	// NR52
	enable bool
}

func NewAPU(mmu *mmu.MMU) *APU {
	apu := &APU{
		mmu:          mmu,
		channel1:     NewChannel1(mmu),
		channel2:     NewChannel2(mmu),
		channel3:     NewChannel3(mmu),
		channel4:     NewChannel4(mmu),
		sampleBuffer: new(bytes.Buffer),
	}

	// FF24 - NR50 - Channel control / ON-OFF / Volume
	mmu.MapMemory(apu, NR50)

	// FF25 - NR51 - Selection of Sound output terminal
	mmu.MapMemory(apu, NR51)

	// FF26 - NR52 - Sound on/off
	mmu.MapMemory(apu, NR52)

	mmu.MapMemoryRange(apu, 0xFF27, 0xFF2F)

	spec := &sdl.AudioSpec{
		Freq:     Frequency,
		Format:   sdl.AUDIO_S16,
		Channels: 2,
		Samples:  Samples,
	}

	if err := sdl.OpenAudio(spec, nil); err != nil {
		panic(err)
	}

	sdl.PauseAudio(false)
	return apu
}

func (s *APU) Tick(mCycles int) {
	tCycles := mCycles * 4

	s.channel1.Tick(tCycles)
	s.channel2.Tick(tCycles)
	s.channel3.Tick(tCycles)
	s.channel4.Tick(tCycles)

	if s.frameSequencerTimer > 0 {
		s.frameSequencerTimer = s.frameSequencerTimer - tCycles
	}

	if s.frameSequencerTimer <= 0 {
		s.tickFrameSequencer()
	}

	if s.sampleTimer > 0 {
		s.sampleTimer = s.sampleTimer - tCycles
	}

	if s.sampleTimer <= 0 {
		SO2 := int(0)
		SO1 := int(0)

		if s.enable {
			channel1Sample := int(s.channel1.sample())
			channel2Sample := int(s.channel2.sample())
			channel3Sample := int(s.channel3.sample())
			channel4Sample := int(s.channel4.sample())

			if s.output4SO2 {
				SO2 += channel4Sample
			}
			if s.output3SO2 {
				SO2 += channel3Sample
			}
			if s.output2SO2 {
				SO2 += channel2Sample
			}
			if s.output1SO2 {
				SO2 += channel1Sample
			}

			if s.output4SO1 {
				SO1 += channel4Sample
			}
			if s.output3SO1 {
				SO1 += channel3Sample
			}
			if s.output2SO1 {
				SO1 += channel2Sample
			}
			if s.output1SO1 {
				SO1 += channel1Sample
			}

			SO2 = ((0xf - SO2*2) * int(s.volumeSO2+1))
			SO1 = ((0xf - SO1*2) * int(s.volumeSO1+1))

			var L = int16(SO1)
			var R = int16(SO2)

			binary.Write(s.sampleBuffer, binary.LittleEndian, L)
			binary.Write(s.sampleBuffer, binary.LittleEndian, R)
			s.sampleCount++

			if s.sampleCount == Samples {
				s.sampleCount = 0
				sdl.QueueAudio(1, s.sampleBuffer.Bytes())
				s.sampleBuffer.Reset()
			}

			// Reload sample timer
			s.sampleTimer = 4194304 / Frequency
		}
	}
}

func (s *APU) tickFrameSequencer() {
	// Length Counter ticks every 2nd step at 256 Hz
	if s.frameSequencerStep%2 == 0 {
		// Tick the channels
		s.channel1.TickLength()
		s.channel2.TickLength()
		s.channel3.TickLength()
		s.channel4.TickLength()
	}

	// Volume Envelope ticks every 7th step at 64 Hz
	if s.frameSequencerStep == 7 {
		s.channel1.TickVolumeEnvelope()
		s.channel2.TickVolumeEnvelope()
		s.channel4.TickVolumeEnvelope()
	}

	// Sweep is adjusted every 2nd and 6th step at 128 Hz
	if s.frameSequencerStep == 2 || s.frameSequencerStep == 6 {
		s.channel1.TickSweep()
	}

	// Step the sequencer
	s.frameSequencerStep += 1
	if s.frameSequencerStep == 8 {
		s.frameSequencerStep = 0
	}

	// Reload sequencer timer
	s.frameSequencerTimer = FrameSequencerTimerRate
}

func (s *APU) ReadByte(addr uint16) byte {
	var value byte

	switch {
	case addr == NR50:
		if s.outputVinSO1 {
			value = utils.SetBit(value, 3)
		}

		if s.outputVinSO2 {
			value = utils.SetBit(value, 7)
		}

		value |= s.volumeSO1
		value |= s.volumeSO2
	case addr == NR51:
		if s.output4SO2 {
			value = utils.SetBit(value, 7)
		}

		if s.output3SO2 {
			value = utils.SetBit(value, 6)
		}

		if s.output2SO2 {
			value = utils.SetBit(value, 5)
		}

		if s.output1SO2 {
			value = utils.SetBit(value, 4)
		}

		if s.output4SO1 {
			value = utils.SetBit(value, 3)
		}

		if s.output3SO1 {
			value = utils.SetBit(value, 2)
		}

		if s.output2SO1 {
			value = utils.SetBit(value, 1)
		}

		if s.output1SO1 {
			value = utils.SetBit(value, 0)
		}
	case addr == NR52:
		if s.enable {
			value = utils.SetBit(value, 7)
		}

		if s.channel4.enable {
			value = utils.SetBit(value, 3)
		}

		if s.channel3.enable {
			value = utils.SetBit(value, 2)
		}

		if s.channel2.enable {
			value = utils.SetBit(value, 1)
		}

		if s.channel1.enable {
			value = utils.SetBit(value, 0)
		}
	}

	return value | apuReadMask[addr]
}

func (s *APU) WriteByte(addr uint16, value byte) {
	switch {
	case addr == NR50:
		s.outputVinSO1 = utils.IsBitSet(value, 3)
		s.outputVinSO2 = utils.IsBitSet(value, 7)
		s.volumeSO1 = value & 0x7
		s.volumeSO2 = value & 0x70
	case addr == NR51:
		s.output4SO2 = utils.IsBitSet(value, 7)
		s.output3SO2 = utils.IsBitSet(value, 6)
		s.output2SO2 = utils.IsBitSet(value, 5)
		s.output1SO2 = utils.IsBitSet(value, 4)
		s.output4SO1 = utils.IsBitSet(value, 3)
		s.output3SO1 = utils.IsBitSet(value, 2)
		s.output2SO1 = utils.IsBitSet(value, 1)
		s.output1SO1 = utils.IsBitSet(value, 0)
	case addr == NR52:
		if utils.IsBitSet(value, 7) {
			s.enable = true
		} else {
			// NR52 controls power to the sound hardware. When powered off, all registers (NR10-NR51) are instantly written with zero and any writes to those registers are ignored while power remains off
			// (except on the DMG, where length counters are unaffected by power and can still be written while off).
			s.enable = false
			for addr := NR10; addr < NR52; addr++ {
				s.mmu.WriteByte(uint16(addr), 0x00)
			}
		}
	}
}
