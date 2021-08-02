package nes

type mapper000 struct {
	nPRGBanks uint8
	nCHRBanks uint8
}

func (cart *CART) newMapper000() {
	mapper := &mapper000{}
	mapper.nPRGBanks = cart.nPRGBanks
	mapper.nCHRBanks = cart.nCHRBanks
	cart.pMapper = mapper
}

func (mapper *mapper000) cpuMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		if mapper.nPRGBanks > 1 {
			*mappedAddr = uint32(addr & 0x7FFF)
			return true
		}
		*mappedAddr = uint32(addr & 0x3FFF)
		return true
	}

	return false
}
func (mapper *mapper000) cpuMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		if mapper.nPRGBanks > 1 {
			*mappedAddr = uint32(addr & 0x7FFF)
			return true
		}
		*mappedAddr = uint32(addr & 0x3FFF)
		return true
	}

	return false
}
func (mapper *mapper000) ppuMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x0000 && addr <= 0x1FFF {
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}
func (mapper *mapper000) ppuapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x0000 && addr <= 0x1FFF {
		if mapper.nCHRBanks == 0 {
			*mappedAddr = uint32(addr)
			return true
		}
	}
	return false
}
