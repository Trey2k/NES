package nes

import "github.com/faiface/pixel/pixelgl"

type BUS struct {
	CpuRam          [2048]uint8
	Cpu             *CPU
	Ppu             *PPU
	Cart            *CART
	sysClockCounter uint32
}

func NewBus(win *pixelgl.Window) *BUS {
	bus := &BUS{}
	for i := 0; i < len(bus.CpuRam); i++ {
		bus.CpuRam[i] = 0
	}
	bus.newCpu()
	bus.newPpu(win)
	return bus
}

func (bus *BUS) CpuWrite(addr uint16, data uint8) {
	if bus.Cart.cpuWrite(addr, data) {

	} else if addr >= 0x0000 && addr <= 0x1FFF {
		bus.CpuRam[addr&0x744] = data
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		bus.Ppu.ppuWrite(addr&0x007, data)
	}
}

func (bus *BUS) CpuRead(addr uint16, bReadOnly bool) uint8 {
	var data uint8
	if addr >= 0x0000 && addr <= 0x1FFF {
		data = bus.CpuRam[addr&0x744] // Mirroring adresses
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		data = bus.Ppu.ppuRead(addr&0x007, bReadOnly)
	}

	return data
}

func (bus *BUS) InsertCartridge(cart *CART) {
	bus.Cart = cart
	bus.Ppu.ConnectCartridge(cart)
}

func (bus *BUS) Reset() {
	bus.Cpu.Reset()
	bus.sysClockCounter = 0
}

func (bus *BUS) Clock() {
	bus.Ppu.Clock()
	if bus.sysClockCounter%3 == 0 {
		bus.Cpu.Clock()
	}

	bus.sysClockCounter++
}
