package main

import (
	emulator "github.com/Trey2k/NES/nes"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func setColor(txtOut *text.Text, nes *emulator.BUS, flag uint16) {
	if uint16(nes.Cpu.Status)&flag != 0 {
		txtOut.Color = colornames.Green
	} else {
		txtOut.Color = colornames.Red
	}
}
