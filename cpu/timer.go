package cpu

import (
	"github.com/kevinbrolly/GopherBoy/mmu"
	"github.com/kevinbrolly/GopherBoy/utils"
)

//  Timer and Divider Registers
const (
	DIV  = 0xFF04
	TIMA = 0xFF05
	TMA  = 0xFF06
	TAC  = 0xFF07

	TIMER_STOP = 2
)

type Timer struct {
	mmu *mmu.MMU

	DIV  byte // Divider
	TIMA byte // Timer Counter
	TMA  byte // Timer Modulo
	TAC  byte // Timer Controller

	dividerCounter int
	timerCounter   int
	cycleChannel   chan int
}

func NewTimer(mmu *mmu.MMU, cycleChannel chan int) *Timer {
	timer := &Timer{
		mmu:          mmu,
		cycleChannel: cycleChannel,
	}

	// 0xFF04 - DIV - Divider Register
	// 0xFF05 - TIMA - Timer counter
	// 0xFF06 - TMA - Timer Modulo
	// 0xFF07 - TAC - Timer Control
	mmu.MapMemoryRange(timer, 0xFF04, 0xFF07)

	timer.Reset()

	go timer.Tick()

	return timer
}

func (timer *Timer) Reset() {
	timer.TIMA = 0x00
	timer.TMA = 0x00
	timer.TAC = 0x05

	timer.timerCounter = 1024
	timer.dividerCounter = 0
}

func (timer *Timer) Tick() {
	for cycle := range timer.cycleChannel {
		timer.updateDividerRegister(cycle)

		if utils.IsBitSet(timer.TAC, TIMER_STOP) {

			timer.timerCounter += cycle

			var threshold int
			frequency := timer.getClockFrequency()

			switch frequency {
			case 0:
				threshold = 256 // frequency 4096
			case 1:
				threshold = 4 // frequency 262144
			case 2:
				threshold = 16 // frequency 65536
			case 3:
				threshold = 64 // frequency 16382
			}

			// https://github.com/fishberg/feo-boy/blob/master/src/bus/timer.rs#L56
			// This is the source of a common timer bug and the reason blargg's instr_timing
			// test was failing with "Failure #255".
			//
			// Some instructions will go over the threshold in one instruction so rather than
			// reset the timerCounter to the current threshold we just subtract the threshold
			// to "reset" it, this way any cycles left over will still be counted/available
			// in timerCounter
			for timer.timerCounter >= threshold {
				timer.timerCounter -= threshold

				// If timer is about to overflow
				if timer.TIMA == 255 {
					timer.TIMA = timer.TMA
					timer.mmu.RequestInterrupt(TIMER_OVERFLOW_INTERRUPT)
				} else {
					timer.TIMA = timer.TIMA + 1
				}
			}
		}
	}
}

func (timer *Timer) updateDividerRegister(cycles int) {
	// Divider uses MCycles
	timer.dividerCounter += cycles

	if timer.dividerCounter >= 255 {
		timer.dividerCounter = 0
		timer.DIV++
	}
}

func (timer *Timer) getClockFrequency() byte {
	return timer.TAC & 0x03
}

func (timer *Timer) ReadByte(addr uint16) byte {
	switch {
	// Timer
	case addr == DIV:
		return timer.DIV
	case addr == TIMA:
		return timer.TIMA
	case addr == TMA:
		return timer.TMA
	case addr == TAC:
		return timer.TAC
	}

	return 0
}

func (timer *Timer) WriteByte(addr uint16, value byte) {
	switch {
	// Timer
	case addr == DIV: // Divider
		// Writing any value to DIV resets it to 0
		timer.DIV = 0
	case addr == TIMA: // Timer Counter
		timer.TIMA = value
	case addr == TMA: // Timer Modulo
		timer.TMA = value
	case addr == TAC:
		currentfreq := timer.getClockFrequency()
		timer.TAC = value

		newfreq := timer.getClockFrequency()

		if currentfreq != newfreq {
			switch newfreq {
			case 0:
				timer.timerCounter = 1024 // frequency 4096
			case 1:
				timer.timerCounter = 16 // frequency 262144
			case 2:
				timer.timerCounter = 64 // frequency 65536
			case 3:
				timer.timerCounter = 256 // frequency 16382
			}
		}
	}
}
