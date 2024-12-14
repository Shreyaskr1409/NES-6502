package olc6502

import "NES_Emulator/bus"

type Cpu struct {
	bus *bus.Bus
}

// FLAGS6502 represents the 6502 processor status flags
const (
	C uint8 = 1 << 0 // Carry Bit
	Z uint8 = 1 << 1 // Zero
	I uint8 = 1 << 2 // Disable Interrupts
	D uint8 = 1 << 3 // Decimal Mode (unused in this implementation)
	B uint8 = 1 << 4 // Break
	U uint8 = 1 << 5 // Unused
	V uint8 = 1 << 6 // Overflow
	N uint8 = 1 << 7 // Negative
)

func NewCpu() *Cpu {
	return &Cpu{}
}

func (c *Cpu) Cleanup() {}

func (c *Cpu) ConnectBus(bus *bus.Bus) {
	c.bus = bus
}

func (c *Cpu) read(addr uint16) uint8 {
	return c.bus.Read(addr, false)
}

func (c *Cpu) write(addr uint16, data uint8) {
	c.bus.Write(addr, data)
}
