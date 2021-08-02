package nes

type CPU struct {
	bus *BUS
	// Status Register
	Status uint8
	// Accumulator Register
	A uint8
	// X Register
	X uint8
	// Y Register
	Y uint8
	// Stack Pointer
	Stkp uint8
	// Data Fetched from memory
	Fetched uint8
	// The current Opcode
	Opcode uint8
	// Number of Cycles left for current operation
	Cycles uint8
	// Program Counter
	PC uint16

	AddrAbs uint16
	// Realitive Address
	AddrRel uint16

	Lookup [256]SInstruction
}

type SInstruction struct {
	// The Name of the instruction
	Name string
	// The current operation
	Operate func() uint8
	// The current address mode
	Addrmode func() uint8

	// The number of clock Cycles the instruction requires to execute
	Cycles uint8
}

const (
	// Carry Bit
	C = (1 << 0)
	// Zero
	Z = (1 << 1)
	// Disable Interrupts
	I = (1 << 2)
	// Decimal Mode (Not used in the NES 6502)
	D = (1 << 3)
	// Break
	B = (1 << 4)
	// Unused
	U = (1 << 5)
	// Overflow
	V = (1 << 6)
	// Negative
	N = (1 << 7)
)

func (bus *BUS) newCpu() {
	cpu := &CPU{}
	cpu.bus = bus
	bus.Cpu = cpu
	cpu.initLookup()
}

func (cpu *CPU) Read(addr uint16) uint8 {
	// Default read only flag to false
	return cpu.bus.CpuRead(addr, false)
}

func (cpu *CPU) write(addr uint16, data uint8) {
	cpu.bus.CpuWrite(addr, data)
}

func (cpu *CPU) GetFlag(flag uint8) uint8 {
	if (cpu.Status & flag) > 0 {
		return 1
	} else {
		return 0
	}
}

func (cpu *CPU) SetFlag(flag uint8, v bool) {
	if v {
		cpu.Status |= flag
	} else {
		cpu.Status &= ^flag
	}
}

func (cpu *CPU) Complete() bool {
	return cpu.Cycles == 0
}

func (cpu *CPU) Clock() {
	if cpu.Cycles == 0 {
		cpu.Opcode = cpu.Read(cpu.PC)
		cpu.SetFlag(U, true)
		cpu.PC++

		cpu.Cycles = cpu.Lookup[cpu.Opcode].Cycles
		addrCycles := cpu.Lookup[cpu.Opcode].Addrmode()
		operCycles := cpu.Lookup[cpu.Opcode].Operate()
		cpu.Cycles += (addrCycles & operCycles)
	}

	cpu.Cycles--
}

// Reset Signal
func (cpu *CPU) Reset() {
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.Stkp = 0xFD
	cpu.Status = 0x00 | U

	cpu.AddrAbs = 0xFFFC
	lo := uint16(cpu.Read(cpu.AddrAbs + 0))
	hi := uint16(cpu.Read(cpu.AddrAbs + 1))

	cpu.PC = (hi << 8) | lo

	cpu.AddrRel = 0x0000
	cpu.AddrAbs = 0x0000
	cpu.Fetched = 0x00

	cpu.Cycles = 8

}

// Inturrupt Request Signal
func (cpu *CPU) irq() {
	if cpu.GetFlag(I) == 0 {
		cpu.write(0x0100+uint16(cpu.Stkp), uint8(cpu.PC>>8&0x00FF))
		cpu.Stkp--
		cpu.write(0x0100+uint16(cpu.Stkp), uint8(cpu.PC&0x00FF))
		cpu.Stkp--

		cpu.SetFlag(B, false)
		cpu.SetFlag(U, true)
		cpu.SetFlag(I, true)
		cpu.write(0x0100+uint16(cpu.Stkp), cpu.Status)
		cpu.Stkp--

		cpu.AddrAbs = 0xFFFE
		lo := uint16(cpu.Read(cpu.AddrAbs + 0))
		hi := uint16(cpu.Read(cpu.AddrAbs + 1))
		cpu.PC = (hi << 8) | lo

		cpu.Cycles = 7
	}
}

// Non Maskable Inturrupt Request Signal
func (cpu *CPU) nmi() {
	cpu.write(0x0100+uint16(cpu.Stkp), uint8(cpu.PC>>8&0x00FF))
	cpu.Stkp--
	cpu.write(0x0100+uint16(cpu.Stkp), uint8(cpu.PC&0x00FF))
	cpu.Stkp--

	cpu.SetFlag(B, false)
	cpu.SetFlag(U, true)
	cpu.SetFlag(I, true)
	cpu.write(0x0100+uint16(cpu.Stkp), cpu.Status)
	cpu.Stkp--

	cpu.AddrAbs = 0xFFFA
	lo := uint16(cpu.Read(cpu.AddrAbs + 0))
	hi := uint16(cpu.Read(cpu.AddrAbs + 1))
	cpu.PC = (hi << 8) | lo

	cpu.Cycles = 8
}

func (cpu *CPU) fetch() uint8 {
	if GetFunctionName(cpu.Lookup[cpu.Opcode].Addrmode) != GetFunctionName(cpu.IMP) {
		cpu.Fetched = cpu.Read(cpu.AddrAbs)
	}
	return cpu.Fetched
}

func (cpu *CPU) initLookup() {
	// This is really really really GROSS
	cpu.Lookup =
		[256]SInstruction{
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
}

func (cpu *CPU) Disassemble(nStart uint16, nStop uint16) ([]string, []uint16) {
	addr := nStart
	mapLines := make([]string, nStop-nStart+1)
	mapAddr := make([]uint16, nStop-nStart+1)
	index := 0
	for uint16(index) <= nStop-nStart {
		sInst := "$" + Hex16(addr) + ": "
		mapAddr[index] = addr

		opcode := cpu.bus.CpuRead(addr, true)
		addr++

		sInst += cpu.Lookup[opcode].Name + " "

		if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.IMP) {
			sInst += " {IMP}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.IMM) {
			value := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "#$" + Hex8(value) + " {IMM}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.ZP0) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "$" + Hex8(lo) + " {ZP0}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.ZPX) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "$" + Hex8(lo) + ", X {ZPX}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.ZPY) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "$" + Hex8(lo) + ", Y {ZPY}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.IZX) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "($" + Hex8(lo) + ", X) {IZX}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.IZY) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "($" + Hex8(lo) + "), Y {IZY}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.ABS) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			hi := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "$" + Hex16((uint16(hi)<<8)|uint16(lo)) + " {ABS}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.ABX) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			hi := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "$" + Hex16((uint16(hi)<<8)|uint16(lo)) + ", X {ABX}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.ABY) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			hi := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "$" + Hex16((uint16(hi)<<8)|uint16(lo)) + ", Y {ABY}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.IND) {
			lo := cpu.bus.CpuRead(addr, true)
			addr++
			hi := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "($" + Hex16((uint16(hi)<<8)|uint16(lo)) + ") {IND}"
		} else if GetFunctionName(cpu.Lookup[opcode].Addrmode) == GetFunctionName(cpu.REL) {
			value := cpu.bus.CpuRead(addr, true)
			addr++
			sInst += "$" + Hex8(value) + " [$" + Hex16(uint16(addr)+uint16(value)) + "] {REL}"
		}

		mapLines[index] = sInst

		index++
	}

	return mapLines, mapAddr
}
