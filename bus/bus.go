package bus

import "NES_Emulator/cpu6502"

type Bus struct {
	*cpu6502.Cpu
	RAM [64 * 1024]byte
}

func NewBus() *Bus {
	bus := &Bus{
		Cpu: &cpu6502.Cpu{},
	}
	for i := range bus.RAM {
		bus.RAM[i] = 0x00
	}
	bus.Cpu.ConnectBus(bus)
	return bus
}

func (b *Bus) Cleanup() {}

func (b *Bus) Write(addr uint16, data byte) {
	if addr >= 0x0000 && addr <= 0xFFFF {
		b.RAM[addr] = data
	}
}

func (b *Bus) Read(addr uint16, bReadOnly bool) byte {
	bReadOnly = false // default value

	if addr >= 0x0000 && addr <= 0xFFFF {
		return b.RAM[addr]
	}

	return 0x00
}
