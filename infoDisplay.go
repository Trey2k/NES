package main

import (
	"fmt"

	emulator "github.com/Trey2k/NES/nes"
	"golang.org/x/image/colornames"
)

func (disp *display) drawCpu() {
	disp.cpuInfo.Clear()
	disp.cpuInfo.Color = colornames.White
	fmt.Fprint(disp.cpuInfo, "STATUS: ")
	setColor(disp.cpuInfo, disp.nes, emulator.N)
	fmt.Fprint(disp.cpuInfo, " N ")
	setColor(disp.cpuInfo, disp.nes, emulator.V)
	fmt.Fprint(disp.cpuInfo, " V ")
	setColor(disp.cpuInfo, disp.nes, emulator.U)
	fmt.Fprint(disp.cpuInfo, " - ")
	setColor(disp.cpuInfo, disp.nes, emulator.B)
	fmt.Fprint(disp.cpuInfo, " B ")
	setColor(disp.cpuInfo, disp.nes, emulator.D)
	fmt.Fprint(disp.cpuInfo, " D ")
	setColor(disp.cpuInfo, disp.nes, emulator.I)
	fmt.Fprint(disp.cpuInfo, " I ")
	setColor(disp.cpuInfo, disp.nes, emulator.Z)
	fmt.Fprint(disp.cpuInfo, " Z ")
	setColor(disp.cpuInfo, disp.nes, emulator.C)
	fmt.Fprint(disp.cpuInfo, " C\n")
	disp.cpuInfo.Color = colornames.White
	fmt.Fprintf(disp.cpuInfo, "PC: %s\n", emulator.Hex16(disp.nes.Cpu.PC))
	fmt.Fprintf(disp.cpuInfo, "A: %s\n", emulator.Hex8(disp.nes.Cpu.A))
	fmt.Fprintf(disp.cpuInfo, "X: %s\n", emulator.Hex8(disp.nes.Cpu.X))
	fmt.Fprintf(disp.cpuInfo, "Y: %s\n", emulator.Hex8(disp.nes.Cpu.Y))
	fmt.Fprintf(disp.cpuInfo, "Stack P: %s\n", emulator.Hex8(disp.nes.Cpu.Stkp))
}

func (disp *display) drawCode() {
	disp.codeInfo.Clear()
	addr := disp.nes.Cpu.PC
	asm, addrs := disp.nes.Cpu.Disassemble(addr-10, addr+10)
	for i := 0; i < len(asm); i++ {
		if addrs[i] == addr {
			disp.codeInfo.Color = colornames.Aqua
			fmt.Fprintln(disp.codeInfo, asm[i])
			disp.codeInfo.Color = colornames.White
			continue
		}
		fmt.Fprintln(disp.codeInfo, asm[i])
	}

}

func (disp *display) drawRam() {
	disp.ramInfo.Clear()
	nAddr := uint16(0x0000)
	for i := 0; i < 16; i++ {
		fmt.Fprintf(disp.ramInfo, "$%04X:", nAddr)
		for x := 0; x < 16; x++ {
			fmt.Fprintf(disp.ramInfo, " %02X", disp.nes.CpuRead(nAddr, true))
			nAddr += 1
		}
		fmt.Fprint(disp.ramInfo, " \n")
	}
	fmt.Fprint(disp.ramInfo, " \n")
	fmt.Fprint(disp.ramInfo, " \n")

	nAddr = uint16(0x8000)
	for i := 0; i < 16; i++ {
		fmt.Fprintf(disp.ramInfo, "$%04X:", nAddr)
		for x := 0; x < 16; x++ {
			fmt.Fprintf(disp.ramInfo, " %02X", disp.nes.CpuRead(nAddr, true))
			nAddr += 1
		}
		fmt.Fprint(disp.ramInfo, " \n")
	}
}
