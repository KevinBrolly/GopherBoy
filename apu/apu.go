package apu

import (
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
	sampleBufferSize        = 1024
)

type APU struct {
	channel1 *Channel1
	channel2 *Channel2
	channel3 *Channel3
	channel4 *Channel4

	sampleTimer       int
	sampleBuffer      []byte
	sampleBufferIndex int

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
		channel1:     NewChannel1(mmu),
		channel2:     NewChannel2(mmu),
		channel3:     NewChannel3(mmu),
		channel4:     NewChannel4(mmu),
		sampleBuffer: make([]byte, sampleBufferSize),
	}

	// FF24 - NR50 - Channel control / ON-OFF / Volume
	mmu.MapMemory(apu, NR50)

	// FF25 - NR51 - Selection of Sound output terminal
	mmu.MapMemory(apu, NR51)

	// FF26 - NR52 - Sound on/off
	mmu.MapMemory(apu, NR52)

	spec := &sdl.AudioSpec{
		Freq:     44100,
		Format:   sdl.AUDIO_S16,
		Channels: 2,
		Samples:  1024,
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
		SO2 := byte(0)
		SO1 := byte(0)

		if s.enable {
			channel1Sample := s.channel1.sample()
			channel2Sample := s.channel2.sample()
			channel3Sample := s.channel3.sample()
			channel4Sample := s.channel4.sample()

			//fmt.Printf("channel1Sample: %v\n", channel1Sample)
			//fmt.Printf("channel2Sample: %v\n", channel1Sample)
			//fmt.Printf("channel3Sample: %v\n", channel1Sample)
			//fmt.Printf("channel4Sample: %v\n", channel1Sample)

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
		}

		SO2 = SO2 * (s.volumeSO2 + 1)
		SO1 = SO1 * (s.volumeSO1 + 1)

		s.sampleBuffer[s.sampleBufferIndex] = SO1
		s.sampleBufferIndex++

		s.sampleBuffer[s.sampleBufferIndex] = SO2
		s.sampleBufferIndex++

		if s.sampleBufferIndex == sampleBufferSize {
			s.sampleBufferIndex = 0

			for sdl.GetQueuedAudioSize(1) > (sampleBufferSize * 4) {
				sdl.Delay(1)
			}

			sdl.QueueAudio(1, s.sampleBuffer)
		}

		//sdl.QueueAudio(1, []byte{SO2, SO1})

		// Reload sample timer
		s.sampleTimer = 4194304 / 44100
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
	switch {
	case addr == NR50:
		var value byte
		if s.outputVinSO1 {
			value = utils.SetBit(value, 3)
		}

		if s.outputVinSO2 {
			value = utils.SetBit(value, 7)
		}

		value |= s.volumeSO1
		value |= s.volumeSO2

		return value
	case addr == NR51:
		var value byte
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
		return value
	case addr == NR52:
		var value byte
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
		return value
	}
	return 0
}

func (s *APU) WriteByte(addr uint16, value byte) {
	switch {
	case addr == NR50:
		s.outputVinSO1 = utils.IsBitSet(value, 3)
		s.outputVinSO2 = utils.IsBitSet(value, 7)
		s.volumeSO1 = (value & 0x7)
		s.volumeSO2 = (value & 0x70)
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
		s.enable = utils.IsBitSet(value, 7)
	}
}
