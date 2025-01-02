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
func (cpu *Cpu) reset() {
	cpu.a = 0
	cpu.x = 0
	cpu.y = 0
	cpu.stkp = 0xFD
	cpu.status = 0x00 | U

	cpu.addr_abs = 0xFFFC
	lo := cpu.read(cpu.addr_abs + 0)
	hi := cpu.read(cpu.addr_abs + 1)

	cpu.pc = uint16(hi<<8) | uint16(lo)
	cpu.addr_rel = 0x0000
	cpu.addr_abs = 0x0000
	cpu.fetched = 0x00

	cpu.cycles = 0

}
func (cpu *Cpu) irq() {
	if cpu.getFlag(I) == 0 {
		cpu.write(0x0100+uint16(cpu.stkp), (uint8(cpu.pc)>>8)&0x00FF)
		cpu.stkp--
		cpu.write(0x0100+uint16(cpu.stkp), uint8(cpu.pc)&0x00FF)
		cpu.stkp--

		cpu.setFlag(B, false)
		cpu.setFlag(U, true)
		cpu.setFlag(I, true)
		cpu.write(0x0100+uint16(cpu.stkp), cpu.status)
		cpu.stkp--

		cpu.addr_abs = 0xFFFE
		lo := cpu.read(cpu.addr_abs + 0)
		hi := cpu.read(cpu.addr_abs + 1)
		cpu.pc = uint16(hi<<8) | uint16(lo)

		cpu.cycles = 7
	}
}
func (cpu *Cpu) nmi() {
	cpu.write(0x0100+uint16(cpu.stkp), (uint8(cpu.pc)>>8)&0x00FF)
	cpu.stkp--
	cpu.write(0x0100+uint16(cpu.stkp), uint8(cpu.pc)&0x00FF)
	cpu.stkp--

	cpu.setFlag(B, false)
	cpu.setFlag(U, true)
	cpu.setFlag(I, true)
	cpu.write(0x0100+uint16(cpu.stkp), cpu.status)
	cpu.stkp--

	cpu.addr_abs = 0xFFFA
	lo := cpu.read(cpu.addr_abs + 0)
	hi := cpu.read(cpu.addr_abs + 1)
	cpu.pc = uint16(hi<<8) | uint16(lo)

	cpu.cycles = 7
}
