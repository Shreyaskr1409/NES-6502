package cpu6502

// Addressing Modes

func (cpu *Cpu) IMP() uint8 {
	// IMPLIED
	// there is no data as part of the instruction
	// but the operation happens on accumulator
	cpu.fetched = cpu.a
	return 0
}

func (cpu *Cpu) IMM() uint8 {
	// IMMEDIATE
	cpu.pc++
	cpu.addr_abs = cpu.pc
	return 0
}

func (cpu *Cpu) ZP0() uint8 {
	// ZERO-PAGE
	// memo addr is 16 bits 0xFF55 for e.g.
	// FF represents the page
	// 55 represents the offset into the page
	// the address space is 256 pages of 256 bytes
	// here first two digits will be 00
	cpu.addr_abs = uint16(cpu.read(cpu.pc))
	cpu.pc++
	cpu.addr_abs &= 0x00FF // taking offset from first page
	return 0
}

func (cpu *Cpu) ZPX() uint8 {
	// ZERO-PAGE X-REGISTER-OFFSET
	cpu.addr_abs = uint16(cpu.read(cpu.pc) + cpu.x)
	cpu.pc++
	cpu.addr_abs &= 0x00FF // taking offset from first page
	return 0
}

func (cpu *Cpu) ZPY() uint8 {
	// ZERO-PAGE Y-REGISTER-OFFSET
	cpu.addr_abs = uint16(cpu.read(cpu.pc) + cpu.y)
	cpu.pc++
	cpu.addr_abs &= 0x00FF // taking offset from first page
	return 0
}

func (cpu *Cpu) ABS() uint8 {
	lo := cpu.read(cpu.pc)
	cpu.pc++
	hi := cpu.read(cpu.pc)
	cpu.pc++

	cpu.addr_abs = uint16((hi << 8) | lo)
	return 0
}

func (cpu *Cpu) ABX() uint8 {
	lo := cpu.read(cpu.pc)
	cpu.pc++
	hi := cpu.read(cpu.pc)
	cpu.pc++

	cpu.addr_abs = uint16((hi << 8) | lo)
	cpu.addr_abs += uint16(cpu.x)

	if (cpu.addr_abs & 0xFF00) != uint16(hi<<8) {
		// checks if page is changed after offset
		// if yes then 1 more cycle is requested
		return 1
	} else {
		return 0
	}
}

func (cpu *Cpu) ABY() uint8 {
	lo := cpu.read(cpu.pc)
	cpu.pc++
	hi := cpu.read(cpu.pc)
	cpu.pc++

	cpu.addr_abs = uint16((hi << 8) | lo)
	cpu.addr_abs += uint16(cpu.y)

	if (cpu.addr_abs & 0xFF00) != uint16(hi<<8) {
		// checks if page is changed after offset
		// if yes then 1 more cycle is requested
		return 1
	} else {
		return 0
	}
}

func (cpu *Cpu) IND() uint8 {
	ptr_lo := cpu.read(cpu.pc)
	cpu.pc++
	ptr_hi := cpu.read(cpu.pc)
	cpu.pc++

	ptr := (ptr_hi << 8) | ptr_lo

	if ptr_lo == 0x00FF {
		cpu.addr_abs = uint16((cpu.read(uint16(ptr)&0xFF00) << 8) | cpu.read(uint16(ptr+0)))
	} else {
		cpu.addr_abs = uint16((cpu.read(uint16(ptr+1)) << 8) | cpu.read(uint16(ptr+0)))
	}

	return 0
}

func (cpu *Cpu) IZX() uint8 {
	t := cpu.read(cpu.pc)
	cpu.pc++

	lo := cpu.read(uint16(t+cpu.x) & 0x00FF)
	hi := cpu.read(uint16(t+cpu.x+1) & 0x00FF)

	cpu.addr_abs = uint16((hi << 8) | lo)

	return 0
}

func (cpu *Cpu) IZY() uint8 {
	t := cpu.read(cpu.pc)
	cpu.pc++

	lo := cpu.read(uint16(t) & 0x00FF)
	hi := cpu.read(uint16(t+1) & 0x00FF)

	cpu.addr_abs = uint16((hi << 8) | lo)
	cpu.addr_abs += uint16(cpu.y)

	if (cpu.addr_abs & 0xFF00) != uint16(hi<<8) {
		// checks if page is changed after offset
		// if yes then 1 more cycle is requested
		return 1
	} else {
		return 0
	}
}

func (cpu *Cpu) REL() uint8 {
	cpu.addr_rel = uint16(cpu.read(cpu.pc))
	cpu.pc++

	if (cpu.addr_rel & 0x80) != 0 {
		cpu.addr_rel |= 0xFF00
	}

	return 0
}
