package control

import (
	"GopherBoy/mmu"
	"GopherBoy/utils"
)

const (
	P1 = 0xFF00

	JOYPAD_INTERRUPT = 4

	SELECT_BUTTON_KEYS    = 5
	SELECT_DIRECTION_KEYS = 4

	RIGHT  = 0
	LEFT   = 1
	UP     = 2
	DOWN   = 3
	A      = 4
	B      = 5
	SELECT = 6
	START  = 7
)

type Controller struct {
	mmu             *mmu.MMU
	controllerState byte
	P1              byte
}

func NewController(mmu *mmu.MMU) *Controller {
	controller := &Controller{
		mmu:             mmu,
		controllerState: 0xFF,
		P1:              0xFF,
	}

	// P1 = 0xFF00
	mmu.MapMemoryRange(controller, P1, P1)
	return controller
}

func (c *Controller) KeyPressed(key byte) {
	// Clear the bit for the pressed key
	c.controllerState = utils.ClearBit(c.controllerState, key)

	switch key {
	case RIGHT, LEFT, UP, DOWN:
		// If the game is interested in direction keys and one was pressed trigger interrupt
		if !utils.IsBitSet(c.P1, SELECT_DIRECTION_KEYS) {
			c.mmu.RequestInterrupt(JOYPAD_INTERRUPT)
		}
	case A, B, SELECT, START:
		// If the game is interested in button keys and one was pressed trigger interrupt
		if !utils.IsBitSet(c.P1, SELECT_BUTTON_KEYS) {
			c.mmu.RequestInterrupt(JOYPAD_INTERRUPT)
		}
	}
}

func (c *Controller) KeyReleased(key byte) {
	// Set the bit for the released key
	c.controllerState = utils.SetBit(c.controllerState, key)
}

func (c *Controller) getControllerState() byte {
	// And the bits in P1 so that only P14 or P15 is not set
	p1 := c.P1 & 0xFF

	switch {
	case !utils.IsBitSet(p1, SELECT_DIRECTION_KEYS):
		return c.P1 ^ (c.controllerState & 0xf)

	case !utils.IsBitSet(p1, SELECT_BUTTON_KEYS):
		return c.P1 ^ ((c.controllerState >> 4) & 0xF)

	default:
		return c.P1
	}
}

func (c *Controller) ReadByte(addr uint16) byte {
	switch {
	case addr == P1:
		return c.getControllerState()
	}

	return 0
}

func (c *Controller) WriteByte(addr uint16, value byte) {
	switch {
	case addr == P1:
		c.P1 = value
	}
}
