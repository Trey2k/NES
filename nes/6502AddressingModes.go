package nes

func (cpu *CPU) IMP() uint8 {
	cpu.Fetched = cpu.A
	return 0
}

func (cpu *CPU) IMM() uint8 {
	cpu.AddrAbs = cpu.PC
	cpu.PC++
	return 0
}

func (cpu *CPU) ZP0() uint8 {
	cpu.AddrAbs = uint16(cpu.Read(cpu.PC))
	cpu.PC++
	cpu.AddrAbs &= 0x00ff
	return 0
}

func (cpu *CPU) ZPX() uint8 {
	cpu.AddrAbs = uint16(cpu.Read(cpu.PC) + cpu.X)
	cpu.PC++
	cpu.AddrAbs &= 0x00ff
	return 0
}

func (cpu *CPU) ZPY() uint8 {
	cpu.AddrAbs = uint16(cpu.Read(cpu.PC) + cpu.Y)
	cpu.PC++
	cpu.AddrAbs &= 0x00ff
	return 0
}

func (cpu *CPU) ABS() uint8 {
	lo := cpu.Read(cpu.PC)
	cpu.PC++
	hi := cpu.Read(cpu.PC)
	cpu.PC++

	cpu.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	return 0
}

func (cpu *CPU) ABX() uint8 {
	lo := cpu.Read(cpu.PC)
	cpu.PC++
	hi := cpu.Read(cpu.PC)
	cpu.PC++

	cpu.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	cpu.AddrAbs += uint16(cpu.X)

	if (cpu.AddrAbs & 0xFF00) != (uint16(hi) << 8) {
		return 1
	}
	return 0
}

func (cpu *CPU) ABY() uint8 {
	lo := cpu.Read(cpu.PC)
	cpu.PC++
	hi := cpu.Read(cpu.PC)
	cpu.PC++

	cpu.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	cpu.AddrAbs += uint16(cpu.Y)

	if (cpu.AddrAbs & 0xFF00) != (uint16(hi) << 8) {
		return 1
	}
	return 0
}

func (cpu *CPU) IND() uint8 {
	ptr_lo := uint(cpu.Read(cpu.PC))
	cpu.PC++
	ptr_hi := uint(cpu.Read(cpu.PC))
	cpu.PC++

	ptr := uint16((ptr_hi << 8) | ptr_lo)

	if ptr_lo == 0x00FF { // Simulate page boundary hardware bug
		cpu.AddrAbs = (uint16(cpu.Read(ptr&0xFF00)) << 8) | uint16(cpu.Read(ptr+0))
		return 0
	} // Otherwise behave normally

	cpu.AddrAbs = (uint16(cpu.Read(ptr+1)) << 8) | uint16(cpu.Read(ptr+0))
	return 0
}

func (cpu *CPU) IZX() uint8 {
	t := uint16(cpu.Read(cpu.PC))
	cpu.PC++

	lo := uint16(cpu.Read(t + uint16(cpu.X)&0x00FF))
	hi := uint16(cpu.Read(t + uint16(cpu.X+1)&0x00FF))

	cpu.AddrAbs = (hi << 8) | lo

	return 0
}

func (cpu *CPU) IZY() uint8 {
	t := uint16(cpu.Read(cpu.PC))
	cpu.PC++

	lo := uint16(cpu.Read(t & 0x00FF))
	hi := uint16(cpu.Read((t + 1) & 0x00FF))

	cpu.AddrAbs = (hi << 8) | lo
	cpu.AddrAbs += uint16(cpu.Y)

	if (cpu.AddrAbs & 0xFF00) != (hi << 8) {
		return 1
	}
	return 0
}

func (cpu *CPU) REL() uint8 {
	cpu.AddrRel = uint16(cpu.Read(cpu.PC))
	cpu.PC++
	if cpu.AddrRel&0x88 != 0 {
		cpu.AddrRel |= 0xFF00
	}

	return 0
}
