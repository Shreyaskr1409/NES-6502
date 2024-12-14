package main

import (
	"NES_Emulator/bus"
	"NES_Emulator/olc6502"
)

func main() {
	bus.NewBus()
	olc6502.NewCpu()
}
