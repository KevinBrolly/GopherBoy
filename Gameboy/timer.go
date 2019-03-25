package Gameboy

type Timer struct {
	gameboy *Gameboy

	DIV  byte // Divider
	TIMA byte // Timer Counter
	TMA  byte // Timer Modulo
	TAC  byte // Timer Controller

	dividerCounter int
	timerCounter   int
}

func NewTimer(gameboy *Gameboy) *Timer {
	timer := &Timer{gameboy: gameboy}
	timer.Reset()
	return timer
}

func (timer *Timer) Reset() {
	timer.TIMA = 0x00
	timer.TMA = 0x00
	timer.TAC = 0x05

	timer.timerCounter = 1024
	timer.dividerCounter = 0
}

func (timer *Timer) Tick(cycles byte) {
	timer.updateDividerRegister(cycles)

	if IsBitSet(timer.TAC, TIMER_STOP) {

		timer.timerCounter += int(cycles)

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
				timer.gameboy.requestInterrupt(TIMER_OVERFLOW_INTERRUPT)
			} else {
				timer.TIMA = timer.TIMA + 1
			}
		}
	}
}

func (timer *Timer) updateDividerRegister(cycles byte) {
	// Divider uses MCycles
	timer.dividerCounter += int(cycles * 4)

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
