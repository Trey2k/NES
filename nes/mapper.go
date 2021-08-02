package nes

type MAPPER interface {
	cpuMapRead(addr uint16, mappedAddr *uint32) bool
	cpuMapWrite(addr uint16, mappedAddr *uint32) bool
	ppuMapRead(addr uint16, mappedAddr *uint32) bool
	ppuapWrite(addr uint16, mappedAddr *uint32) bool
}
