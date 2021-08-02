package nes

// Illigal Opcode Catch All
func (cpu *CPU) XXX() uint8 {
	return 0
}

func (cpu *CPU) ADC() uint8 {
	cpu.fetch()
	temp := uint16(cpu.A) + uint16(cpu.Fetched) + uint16(cpu.GetFlag(C))
	cpu.SetFlag(C, temp > 255) // Set carry flag if needed
	cpu.SetFlag(Z, (temp&0x00FF) == 0)
	cpu.SetFlag(N, temp&0x80 != 0)
	// If we overflowed set the overflow flag
	cpu.SetFlag(V, (^(uint16(cpu.A)^uint16(cpu.Fetched))&(uint16(cpu.A)^temp))&0x0080 != 0)
	cpu.A = uint8(temp) & 0x00FF
	return 1
}
func (cpu *CPU) AND() uint8 {
	cpu.fetch()
	cpu.A = cpu.A & cpu.Fetched
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, cpu.A&0x80 != 0)
	return 1
}
func (cpu *CPU) ASL() uint8 {
	return 0
}
func (cpu *CPU) BCC() uint8 {
	if cpu.GetFlag(C) == 0 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) BCS() uint8 {
	if cpu.GetFlag(C) == 1 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) BEQ() uint8 {
	if cpu.GetFlag(Z) == 1 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) BIT() uint8 {
	return 0
}
func (cpu *CPU) BMI() uint8 {
	if cpu.GetFlag(N) == 1 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) BNE() uint8 {
	if cpu.GetFlag(Z) == 0 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) BPL() uint8 {
	if cpu.GetFlag(N) == 0 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) BRK() uint8 {
	return 0
}
func (cpu *CPU) BVC() uint8 {
	if cpu.GetFlag(V) == 0 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) BVS() uint8 {
	if cpu.GetFlag(V) == 1 {
		cpu.Cycles++
		cpu.AddrAbs = cpu.PC + cpu.AddrRel

		if (cpu.AddrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
			cpu.Cycles++
		}
		cpu.PC = cpu.AddrAbs
	}
	return 0
}
func (cpu *CPU) CLC() uint8 {
	cpu.SetFlag(C, false)
	return 0
}
func (cpu *CPU) CLD() uint8 {
	cpu.SetFlag(D, false)
	return 0
}
func (cpu *CPU) CLI() uint8 {
	cpu.SetFlag(I, false)
	return 0
}
func (cpu *CPU) CLV() uint8 {
	cpu.SetFlag(V, false)
	return 0
}
func (cpu *CPU) CMP() uint8 {
	cpu.fetch()
	temp := uint16(cpu.A) - uint16(cpu.Fetched)
	cpu.SetFlag(C, cpu.A >= cpu.Fetched)
	cpu.SetFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.SetFlag(N, temp&0x0080 != 0)
	return 0
}
func (cpu *CPU) CPX() uint8 {
	cpu.fetch()
	temp := uint16(cpu.X) - uint16(cpu.Fetched)
	cpu.SetFlag(C, cpu.X >= cpu.Fetched)
	cpu.SetFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.SetFlag(N, temp&0x0080 != 0)
	return 0
}
func (cpu *CPU) CPY() uint8 {
	cpu.fetch()
	temp := uint16(cpu.Y) - uint16(cpu.Fetched)
	cpu.SetFlag(C, cpu.Y >= cpu.Fetched)
	cpu.SetFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.SetFlag(N, temp&0x0080 != 0)
	return 0
}
func (cpu *CPU) DEC() uint8 {
	cpu.fetch()
	temp := cpu.Fetched - 1
	cpu.write(cpu.AddrAbs, temp&0x00FF)
	cpu.SetFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.SetFlag(N, temp&0x0080 != 0)
	return 0
}
func (cpu *CPU) DEX() uint8 {
	cpu.X--
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, cpu.X&0x80 != 0)
	return 0
}
func (cpu *CPU) DEY() uint8 {
	cpu.Y--
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, cpu.Y&0x80 != 0)
	return 0
}
func (cpu *CPU) EOR() uint8 {
	cpu.fetch()
	cpu.A ^= cpu.Fetched
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, cpu.A&0x80 != 0)
	return 1
}
func (cpu *CPU) INC() uint8 {
	cpu.fetch()
	temp := cpu.Fetched + 1
	cpu.write(cpu.AddrAbs, temp&0x00FF)
	cpu.SetFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.SetFlag(N, temp&0x0080 != 0)
	return 0
}
func (cpu *CPU) INX() uint8 {
	cpu.X++
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, cpu.X&0x80 != 0)
	return 0
}
func (cpu *CPU) INY() uint8 {
	cpu.Y++
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, cpu.Y&0x80 != 0)
	return 0
}
func (cpu *CPU) JMP() uint8 {
	cpu.PC = cpu.AddrAbs
	return 0
}
func (cpu *CPU) JSR() uint8 {
	cpu.PC--

	cpu.write(uint16(0x0100)+uint16(cpu.Stkp), uint8(cpu.PC>>8&uint16(0x00FF)))
	cpu.Stkp--
	cpu.write(uint16(0x0100)+uint16(cpu.Stkp), uint8(cpu.PC&0x00FF))
	cpu.Stkp--

	cpu.PC = cpu.AddrAbs
	return 0
}
func (cpu *CPU) LDA() uint8 {
	cpu.fetch()
	cpu.A = cpu.Fetched
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, cpu.A&0x80 != 0)
	return 0
}
func (cpu *CPU) LDX() uint8 {
	cpu.fetch()
	cpu.X = cpu.Fetched
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, cpu.X&0x80 != 0)
	return 0
}
func (cpu *CPU) LDY() uint8 {
	cpu.fetch()
	cpu.Y = cpu.Fetched
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, cpu.Y&0x80 != 0)
	return 0
}
func (cpu *CPU) LSR() uint8 {
	return 0
}
func (cpu *CPU) NOP() uint8 {
	switch cpu.Opcode {
	case 0x1C:
	case 0x3C:
	case 0x5C:
	case 0x7C:
	case 0xDC:
	case 0xFC:
		return 1
	}
	return 0
}
func (cpu *CPU) ORA() uint8 {
	return 0
}
func (cpu *CPU) PHA() uint8 {
	cpu.write(0x0100+uint16(cpu.Stkp), cpu.A)
	cpu.Stkp--
	return 0
}
func (cpu *CPU) PHP() uint8 {
	cpu.Stkp++
	cpu.A = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, cpu.A&0x80 != 0)
	return 0
}
func (cpu *CPU) PLA() uint8 {
	cpu.Stkp++
	cpu.A = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, cpu.A&0x80 != 0)
	return 0
}
func (cpu *CPU) PLP() uint8 {
	cpu.Stkp++
	cpu.Status = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.SetFlag(U, true)
	return 0
}
func (cpu *CPU) ROL() uint8 {
	cpu.fetch()
	temp := uint16((cpu.Fetched << 1) | cpu.GetFlag(C))
	cpu.SetFlag(C, temp&0xFF00 != 0)
	cpu.SetFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.SetFlag(N, temp&0x0080 != 0)
	if GetFunctionName(cpu.Lookup[cpu.Opcode].Addrmode) == GetFunctionName(cpu.IMP) {
		cpu.A = uint8(temp & 0x00FF)
	} else {
		cpu.write(cpu.AddrAbs, uint8(temp&0x00FF))
	}
	return 0
}
func (cpu *CPU) ROR() uint8 {
	cpu.Stkp++
	cpu.Status = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.Status &= ^uint8(B)
	cpu.Status &= ^uint8(U)

	cpu.Stkp++
	cpu.PC = uint16(cpu.Read(0x0100 + uint16(cpu.Stkp)))
	cpu.Stkp++
	cpu.PC |= uint16(cpu.Read(0x0100+uint16(cpu.Stkp))) << 8
	return 0
}
func (cpu *CPU) RTI() uint8 {
	cpu.Stkp++
	cpu.Status = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.Status &= ^uint8(B)
	cpu.Status &= ^uint8(U)

	cpu.Stkp++
	cpu.PC = uint16(cpu.Read(0x0100 + uint16(cpu.Stkp)))
	cpu.Stkp++
	cpu.PC |= uint16(cpu.Read(0x0100+uint16(cpu.Stkp))) << 8
	return 0
}
func (cpu *CPU) RTS() uint8 {
	cpu.Stkp++
	cpu.PC = uint16(cpu.Read(0x0100 + uint16(cpu.Stkp)))
	cpu.Stkp++
	cpu.PC |= uint16(cpu.Read(0x0100+uint16(cpu.Stkp))) << 8

	cpu.PC++
	return 0
}
func (cpu *CPU) SBC() uint8 {
	cpu.fetch()
	value := uint16(cpu.Fetched) ^ 0x00FF
	temp := uint16(cpu.A) + value + uint16(cpu.GetFlag(C))
	cpu.SetFlag(C, temp > 255) // Set carry flag if needed
	cpu.SetFlag(Z, (temp&0x00FF) == 0)
	cpu.SetFlag(N, temp&0x80 != 0)
	// If we overflowed set the overflow flag
	cpu.SetFlag(V, (^(uint16(cpu.A)^uint16(cpu.Fetched))&(uint16(cpu.A)^temp))&0x0080 != 0)
	cpu.A = uint8(temp) & 0x00FF
	return 1
}
func (cpu *CPU) SEC() uint8 {
	cpu.SetFlag(C, true)
	return 0
}
func (cpu *CPU) SED() uint8 {
	cpu.SetFlag(D, true)
	return 0
}
func (cpu *CPU) SEI() uint8 {
	cpu.SetFlag(I, true)
	return 0
}
func (cpu *CPU) STA() uint8 {
	cpu.write(cpu.AddrAbs, cpu.A)
	return 0
}
func (cpu *CPU) STX() uint8 {
	cpu.write(cpu.AddrAbs, cpu.X)
	return 0
}
func (cpu *CPU) STY() uint8 {
	cpu.write(cpu.AddrAbs, cpu.Y)
	return 0
}
func (cpu *CPU) TAX() uint8 {
	cpu.X = cpu.A
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, cpu.X&0x80 != 0)
	return 0
}
func (cpu *CPU) TAY() uint8 {
	cpu.Y = cpu.A
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, cpu.Y&0x80 != 0)
	return 0
}
func (cpu *CPU) TSX() uint8 {
	cpu.X = cpu.Stkp
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, cpu.X&0x80 != 0)
	return 0
}
func (cpu *CPU) TXA() uint8 {
	cpu.A = cpu.X
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, cpu.A&0x80 != 0)
	return 0
}
func (cpu *CPU) TXS() uint8 {
	cpu.Stkp = cpu.X
	return 0
}
func (cpu *CPU) TYA() uint8 {
	cpu.A = cpu.Y
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, cpu.A&0x80 != 0)
	return 0
}
