package nes

import (
	"encoding/binary"
	"fmt"
	"os"
)

type CART struct {
	VPRGMemory []uint8
	VCHRMemory []uint8
	nMapperID  uint8
	nPRGBanks  uint8
	nCHRBanks  uint8
	mirror     uint8
	pMapper    MAPPER
	ImageValid bool
}

type CartHeader struct {
	Name [4]byte
	PrgRomChunks,
	ChrRomChunks,
	Mapper1,
	Mapper2,
	PrgRamSize,
	TvSystem1,
	TvSystem2 byte
	_ [5]byte // unused padding
}

const (
	HORIZONTAL = iota
	VERTICAL
	ONESCREEN_LO
	ONESCREEN_HI
)

func NewCartridge(fileName string) (*CART, error) {
	header := &CartHeader{}
	cart := &CART{}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := binary.Read(file, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	if header.Mapper1&0x04 != 0 {

	}
	cart.nMapperID = ((header.Mapper2 >> 4) << 4) | (header.Mapper1 >> 4)
	if header.Mapper1&0x01 != 0 {
		cart.mirror = VERTICAL
	} else {
		cart.mirror = HORIZONTAL
	}

	var nFileType uint8 = 1

	if nFileType == 0 {

	}

	if nFileType == 1 {
		cart.nPRGBanks = header.PrgRomChunks

		cart.VPRGMemory = make([]uint8, uint16(cart.nPRGBanks)*16384)
		if err := binary.Read(file, binary.LittleEndian, cart.VPRGMemory); err != nil {
			return nil, err
		}

		cart.nCHRBanks = header.ChrRomChunks
		cart.VCHRMemory = make([]uint8, uint16(cart.nCHRBanks)*8192)
		if err := binary.Read(file, binary.LittleEndian, cart.VCHRMemory); err != nil {
			panic(err)
		}

	}

	if nFileType == 2 {

	}
	fmt.Println(cart.nMapperID)
	switch cart.nMapperID {
	case 0:
		cart.newMapper000()
	}

	cart.ImageValid = true

	return cart, nil
}

func (cart *CART) cpuRead(addr uint16, data *uint8) bool {
	var mapperAdr uint32
	if cart.pMapper.cpuMapRead(addr, &mapperAdr) {
		*data = cart.VPRGMemory[mapperAdr]
		return true
	}
	return false
}

func (cart *CART) cpuWrite(addr uint16, data uint8) bool {
	var mapperAdr uint32
	if cart.pMapper.cpuMapWrite(addr, &mapperAdr) {
		cart.VPRGMemory[mapperAdr] = data
		return true
	}
	return false
}

func (cart *CART) ppuRead(addr uint16, data *uint8) bool {
	var mapperAdr uint32
	if cart.pMapper.ppuMapRead(addr, &mapperAdr) {
		*data = cart.VCHRMemory[mapperAdr]
		return true
	}
	return false
}

func (cart *CART) ppuWrite(addr uint16, data uint8) bool {
	var mapperAdr uint32
	if cart.pMapper.ppuMapRead(addr, &mapperAdr) {
		cart.VCHRMemory[mapperAdr] = data
		return true
	}
	return false
}
