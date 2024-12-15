package cpu6502

import "NES_Emulator/bus"

type Cpu struct {
	// Registers
	a      uint8  // Accumulator
	x      uint8  // X Register
	y      uint8  // Y Register
	stkp   uint8  // Stack Pointer
	pc     uint16 // Program Counter
	status uint8  // Status Register

	// Internal CPU states
	fetched  uint8 // data fetched on operation to be used later
	addr_abs uint16
	addr_rel uint16
	opcode   uint8
	cycles   uint8

	bus *bus.Bus

	lookup []Instruction
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
	cpu := &Cpu{}

	// the following is a 16x16 matrix as per operation-codes
	cpu.lookup = []Instruction{
		{"BRK", cpu.BRK, cpu.IMM, 7}, {"ORA", cpu.ORA, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 3}, {"ORA", cpu.ORA, cpu.ZP0, 3}, {"ASL", cpu.ASL, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PHP", cpu.PHP, cpu.IMP, 3}, {"ORA", cpu.ORA, cpu.IMM, 2}, {"ASL", cpu.ASL, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.NOP, cpu.IMP, 4}, {"ORA", cpu.ORA, cpu.ABS, 4}, {"ASL", cpu.ASL, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BPL", cpu.BPL, cpu.REL, 2}, {"ORA", cpu.ORA, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"ORA", cpu.ORA, cpu.ZPX, 4}, {"ASL", cpu.ASL, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"CLC", cpu.CLC, cpu.IMP, 2}, {"ORA", cpu.ORA, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"ORA", cpu.ORA, cpu.ABX, 4}, {"ASL", cpu.ASL, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"JSR", cpu.JSR, cpu.ABS, 6}, {"AND", cpu.AND, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"BIT", cpu.BIT, cpu.ZP0, 3}, {"AND", cpu.AND, cpu.ZP0, 3}, {"ROL", cpu.ROL, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PLP", cpu.PLP, cpu.IMP, 4}, {"AND", cpu.AND, cpu.IMM, 2}, {"ROL", cpu.ROL, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"BIT", cpu.BIT, cpu.ABS, 4}, {"AND", cpu.AND, cpu.ABS, 4}, {"ROL", cpu.ROL, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BMI", cpu.BMI, cpu.REL, 2}, {"AND", cpu.AND, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"AND", cpu.AND, cpu.ZPX, 4}, {"ROL", cpu.ROL, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"SEC", cpu.SEC, cpu.IMP, 2}, {"AND", cpu.AND, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"AND", cpu.AND, cpu.ABX, 4}, {"ROL", cpu.ROL, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"RTI", cpu.RTI, cpu.IMP, 6}, {"EOR", cpu.EOR, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 3}, {"EOR", cpu.EOR, cpu.ZP0, 3}, {"LSR", cpu.LSR, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PHA", cpu.PHA, cpu.IMP, 3}, {"EOR", cpu.EOR, cpu.IMM, 2}, {"LSR", cpu.LSR, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"JMP", cpu.JMP, cpu.ABS, 3}, {"EOR", cpu.EOR, cpu.ABS, 4}, {"LSR", cpu.LSR, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BVC", cpu.BVC, cpu.REL, 2}, {"EOR", cpu.EOR, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"EOR", cpu.EOR, cpu.ZPX, 4}, {"LSR", cpu.LSR, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"CLI", cpu.CLI, cpu.IMP, 2}, {"EOR", cpu.EOR, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"EOR", cpu.EOR, cpu.ABX, 4}, {"LSR", cpu.LSR, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"RTS", cpu.RTS, cpu.IMP, 6}, {"ADC", cpu.ADC, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 3}, {"ADC", cpu.ADC, cpu.ZP0, 3}, {"ROR", cpu.ROR, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PLA", cpu.PLA, cpu.IMP, 4}, {"ADC", cpu.ADC, cpu.IMM, 2}, {"ROR", cpu.ROR, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"JMP", cpu.JMP, cpu.IND, 5}, {"ADC", cpu.ADC, cpu.ABS, 4}, {"ROR", cpu.ROR, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BVS", cpu.BVS, cpu.REL, 2}, {"ADC", cpu.ADC, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"ADC", cpu.ADC, cpu.ZPX, 4}, {"ROR", cpu.ROR, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"SEI", cpu.SEI, cpu.IMP, 2}, {"ADC", cpu.ADC, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"ADC", cpu.ADC, cpu.ABX, 4}, {"ROR", cpu.ROR, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"???", cpu.NOP, cpu.IMP, 2}, {"STA", cpu.STA, cpu.IZX, 6}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 6}, {"STY", cpu.STY, cpu.ZP0, 3}, {"STA", cpu.STA, cpu.ZP0, 3}, {"STX", cpu.STX, cpu.ZP0, 3}, {"???", cpu.XXX, cpu.IMP, 3}, {"DEY", cpu.DEY, cpu.IMP, 2}, {"???", cpu.NOP, cpu.IMP, 2}, {"TXA", cpu.TXA, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"STY", cpu.STY, cpu.ABS, 4}, {"STA", cpu.STA, cpu.ABS, 4}, {"STX", cpu.STX, cpu.ABS, 4}, {"???", cpu.XXX, cpu.IMP, 4},
		{"BCC", cpu.BCC, cpu.REL, 2}, {"STA", cpu.STA, cpu.IZY, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 6}, {"STY", cpu.STY, cpu.ZPX, 4}, {"STA", cpu.STA, cpu.ZPX, 4}, {"STX", cpu.STX, cpu.ZPY, 4}, {"???", cpu.XXX, cpu.IMP, 4}, {"TYA", cpu.TYA, cpu.IMP, 2}, {"STA", cpu.STA, cpu.ABY, 5}, {"TXS", cpu.TXS, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 5}, {"???", cpu.NOP, cpu.IMP, 5}, {"STA", cpu.STA, cpu.ABX, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"???", cpu.XXX, cpu.IMP, 5},
		{"LDY", cpu.LDY, cpu.IMM, 2}, {"LDA", cpu.LDA, cpu.IZX, 6}, {"LDX", cpu.LDX, cpu.IMM, 2}, {"???", cpu.XXX, cpu.IMP, 6}, {"LDY", cpu.LDY, cpu.ZP0, 3}, {"LDA", cpu.LDA, cpu.ZP0, 3}, {"LDX", cpu.LDX, cpu.ZP0, 3}, {"???", cpu.XXX, cpu.IMP, 3}, {"TAY", cpu.TAY, cpu.IMP, 2}, {"LDA", cpu.LDA, cpu.IMM, 2}, {"TAX", cpu.TAX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"LDY", cpu.LDY, cpu.ABS, 4}, {"LDA", cpu.LDA, cpu.ABS, 4}, {"LDX", cpu.LDX, cpu.ABS, 4}, {"???", cpu.XXX, cpu.IMP, 4},
		{"BCS", cpu.BCS, cpu.REL, 2}, {"LDA", cpu.LDA, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 5}, {"LDY", cpu.LDY, cpu.ZPX, 4}, {"LDA", cpu.LDA, cpu.ZPX, 4}, {"LDX", cpu.LDX, cpu.ZPY, 4}, {"???", cpu.XXX, cpu.IMP, 4}, {"CLV", cpu.CLV, cpu.IMP, 2}, {"LDA", cpu.LDA, cpu.ABY, 4}, {"TSX", cpu.TSX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 4}, {"LDY", cpu.LDY, cpu.ABX, 4}, {"LDA", cpu.LDA, cpu.ABX, 4}, {"LDX", cpu.LDX, cpu.ABY, 4}, {"???", cpu.XXX, cpu.IMP, 4},
		{"CPY", cpu.CPY, cpu.IMM, 2}, {"CMP", cpu.CMP, cpu.IZX, 6}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"CPY", cpu.CPY, cpu.ZP0, 3}, {"CMP", cpu.CMP, cpu.ZP0, 3}, {"DEC", cpu.DEC, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"INY", cpu.INY, cpu.IMP, 2}, {"CMP", cpu.CMP, cpu.IMM, 2}, {"DEX", cpu.DEX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"CPY", cpu.CPY, cpu.ABS, 4}, {"CMP", cpu.CMP, cpu.ABS, 4}, {"DEC", cpu.DEC, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BNE", cpu.BNE, cpu.REL, 2}, {"CMP", cpu.CMP, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"CMP", cpu.CMP, cpu.ZPX, 4}, {"DEC", cpu.DEC, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"CLD", cpu.CLD, cpu.IMP, 2}, {"CMP", cpu.CMP, cpu.ABY, 4}, {"NOP", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"CMP", cpu.CMP, cpu.ABX, 4}, {"DEC", cpu.DEC, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"CPX", cpu.CPX, cpu.IMM, 2}, {"SBC", cpu.SBC, cpu.IZX, 6}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"CPX", cpu.CPX, cpu.ZP0, 3}, {"SBC", cpu.SBC, cpu.ZP0, 3}, {"INC", cpu.INC, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"INX", cpu.INX, cpu.IMP, 2}, {"SBC", cpu.SBC, cpu.IMM, 2}, {"NOP", cpu.NOP, cpu.IMP, 2}, {"???", cpu.SBC, cpu.IMP, 2}, {"CPX", cpu.CPX, cpu.ABS, 4}, {"SBC", cpu.SBC, cpu.ABS, 4}, {"INC", cpu.INC, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BEQ", cpu.BEQ, cpu.REL, 2}, {"SBC", cpu.SBC, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"SBC", cpu.SBC, cpu.ZPX, 4}, {"INC", cpu.INC, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"SED", cpu.SED, cpu.IMP, 2}, {"SBC", cpu.SBC, cpu.ABY, 4}, {"NOP", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"SBC", cpu.SBC, cpu.ABX, 4}, {"INC", cpu.INC, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
	}

	return cpu
}

func (cpu *Cpu) Cleanup() {}

func (cpu *Cpu) ConnectBus(bus *bus.Bus) {
	cpu.bus = bus
}

func (cpu *Cpu) read(addr uint16) uint8 {
	return cpu.bus.Read(addr, false)
}

func (cpu *Cpu) write(addr uint16, data uint8) {
	cpu.bus.Write(addr, data)
}

func (cpu *Cpu) getFlag(flag uint8) uint8 {
	return 0x0000
}

func (cpu *Cpu) setFlag(flag uint8, value bool) {}

func (cpu *Cpu) fetch() uint8 {
	return 0x0000
}
