package main

import (
	"NES_Emulator/bus"
	"NES_Emulator/cpu6502"
)

func main() {
	bus.NewBus()
	cpu6502.NewCpu()
}
