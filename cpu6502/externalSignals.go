package cpu6502

func (cpu *Cpu) clock() {
	if cpu.cycles == 0 {
		cpu.opcode = cpu.read(cpu.pc)
		cpu.pc++

		additionalCycle1 := cpu.lookup[cpu.opcode].addrmode()
		additionalCycle2 := cpu.lookup[cpu.opcode].operate()

		cpu.cycles = additionalCycle1 + additionalCycle2
	}

	cpu.cycles--
}
func (cpu *Cpu) reset() {}
func (cpu *Cpu) irq()   {}
func (cpu *Cpu) nmi()   {}
