package cpu6502

// THIS FILES IS A UTIL FILE FOR DEBUGGING

// stores pointers of operations which I will be performing
type Operation func() uint8
type AddressMode func() uint8

// Instruction structure
type Instruction struct {
	name     string      // Name of the instruction
	operate  Operation   // Operation function pointer
	addrmode AddressMode // Addressing mode function pointer
	cycles   uint8       // Number of clock cycles
}
