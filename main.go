package main

import (
	"fmt"
	"time"

	"github.com/Trey2k/NES/nes"
	emulator "github.com/Trey2k/NES/nes"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type display struct {
	cpuInfo  *text.Text
	codeInfo *text.Text
	ramInfo  *text.Text
	win      *pixelgl.Window
	nes      *emulator.BUS
}

func setupNES(win *pixelgl.Window) *emulator.BUS {
	nes := nes.NewBus(win)
	cart, err := emulator.NewCartridge("assets/SMB.nes")
	if err != nil {
		panic(err)
	}
	if cart.ImageValid {
		nes.InsertCartridge(cart)
		nes.Cpu.Reset()
	} else {
		panic("BAD IMAGE")
	}

	return nes
}

func (disp *display) update() {
	disp.win.Clear(colornames.Blue)
	disp.nes.Clock()
	for !disp.nes.Ppu.FramComplete {
		disp.nes.Clock()
	}
	disp.nes.Ppu.FramComplete = false

	if disp.win.JustPressed(pixelgl.KeySpace) {
		disp.nes.Cpu.Clock()
		for !disp.nes.Cpu.Complete() {
			disp.nes.Clock()
		}
	} else if disp.win.JustPressed(pixelgl.KeyEnter) {

		for !disp.nes.Cpu.Complete() {
			disp.nes.Clock()
		}

	} else if disp.win.JustPressed(pixelgl.KeyR) {
		disp.nes.Reset()
	}
	disp.drawCpu()
	disp.drawCode()

	//disp.drawRam()
}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "NES Emulator",
		Bounds: pixel.R(0, 0, 680, 480),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	nes := setupNES(win)
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	cpuInfo := text.New(pixel.V(450, 450), basicAtlas)
	codeInfo := text.New(pixel.V(450, 340), basicAtlas)
	ramInfo := text.New(pixel.V(15, 450), basicAtlas)
	disp := &display{
		win:      win,
		nes:      nes,
		cpuInfo:  cpuInfo,
		codeInfo: codeInfo,
		ramInfo:  ramInfo,
	}
	win.Clear(colornames.Blue)
	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	for !win.Closed() {
		win.Clear(colornames.Blue)
		cpuInfo.Draw(win, pixel.IM.Scaled(cpuInfo.Orig, 1))
		codeInfo.Draw(win, pixel.IM.Scaled(codeInfo.Orig, 1))
		//ramInfo.Draw(win, pixel.IM.Scaled(ramInfo.Orig, 1))

		disp.update()
		win.Update()
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)

}
