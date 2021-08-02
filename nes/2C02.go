package nes

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)
type PPU struct {
	cart       *CART
	tblName    [2][14]uint8
	tblPallet  [32]int8
	tblPattern [2][4096]int8// uture fun?

	palScreen       [0x40]color.RGBA
	SprScreen       pixel.Picture
	PixelSprites	[0x40]*pixel.Sprite
	SprFrames       []pixel.Rect
	sprNameTable    [2]*imdraw.IMDraw
	sprPatternTable []*imdraw.IMDraw

	sprBatch		*pixel.Batch
	
	win *pixelgl.Window

	FramComplete bool

	scanLine int16
	cycle    int16
}

func (bus *BUS) newPpu(win *pixelgl.Window) {
	ppu := &PPU{}
	ppu.win = win
	ppu.SetupPallet()
	bus.Ppu = ppu
}

func (pu *PPU) cpuRead(addr uint16, readonly bool) uint8 {
var data uint8
	switch(addr) {
	case 0x0000: // Contrl
		break 
	case 0x0001: // Mask
		break
	case 0x0002: // Stats
		break
	case 0x0003: // OAM Adess
		break
	case 0x0004: // OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Adress
		break
	case 0x0007: // PPU Data
	break
	}
		return data
}

func (ppu *PPU) cpuWite(addr uint16, data uint8) {
	switch(addr) {
	case 0x0000: // Contrl
		break 
	case 0x0001: // Mask
		break
	case 0x0002: // Stats
		break
	case 0x0003: // OAM Adess
		break
	case 0x0004: // OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Adress
		break
	case 00007: // PPU Data
		break
		}
}

func (ppu *PPU) ppuRead(addr uint16, readonly bool) uint8 {
	var data uint8
	addr &= 0x3FFF

	if ppu.cart.ppuRead(addr, &data) {

	}
	return data

}

func (ppu *PPU) ppuWrite(addr uint16, data uint8) {
	addr &= 0x3FFF

	if ppu.cart.ppuWrite(addr, data) {

	}
}


func (ppu *PPU) ConnectCartridge(cart *CART) {
	ppu.cart = cart
}

func (ppu *PPU) Clock() {
	
	
	ppu.PixelSprites[rand.Intn(len(ppu.PixelSprites))].Draw(ppu.sprBatch, pixel.IM.Scaled(pixel.Vec{0, 0}, 2).Moved(pixel.Vec{float64(ppu.cycle)*2, float64(ppu.scanLine)*2}))
	ppu.cycle++
	if ppu.cycle >= 256  {
		ppu.cycle = 0
		ppu.scanLine++
		
		if ppu.scanLine >= 240 {
			ppu.sprBatch.Draw(ppu.win)
			ppu.sprBatch.Clear()
			ppu.scanLine = -1
			
			ppu.FramComplete = true
		}
	}	
}

func (ppu *PPU) SetupPallet() {
	ppu.palScreen[0x00] = color.RGBA{84, 84, 84, 255}
	ppu.palScreen[0x01] = color.RGBA{0, 30, 116, 255}
	ppu.palScreen[0x02] = color.RGBA{8, 16, 144, 255}
	ppu.palScreen[0x03] = color.RGBA{48, 0, 136, 255}
	ppu.palScreen[0x04] = color.RGBA{68, 0, 100, 255}
	ppu.palScreen[0x05] = color.RGBA{92, 0, 48, 255}
	ppu.palScreen[0x06] = color.RGBA{84, 4, 0, 255}
	ppu.palScreen[0x07] = color.RGBA{60, 24, 0, 255}
	ppu.palScreen[0x08] = color.RGBA{32, 42, 0, 255}
	ppu.palScreen[0x09] = color.RGBA{8, 58, 0, 255}
	ppu.palScreen[0x0A] = color.RGBA{0, 64, 0, 255}
	ppu.palScreen[0x0B] = color.RGBA{0, 60, 0, 255}
	ppu.palScreen[0x0C] = color.RGBA{0, 50, 60, 255}
	ppu.palScreen[0x0D] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x0E] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x0F] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x10] = color.RGBA{152, 150, 152, 255}
	ppu.palScreen[0x11] = color.RGBA{8, 76, 196, 255}
	ppu.palScreen[0x12] = color.RGBA{48, 50, 236, 255}
	ppu.palScreen[0x13] = color.RGBA{92, 30, 228, 255}
	ppu.palScreen[0x14] = color.RGBA{136, 20, 176, 255}
	ppu.palScreen[0x15] = color.RGBA{160, 20, 100, 255}
	ppu.palScreen[0x16] = color.RGBA{152, 34, 32, 255}
	ppu.palScreen[0x17] = color.RGBA{120, 60, 0, 255}
	ppu.palScreen[0x18] = color.RGBA{84, 90, 0, 255}
	ppu.palScreen[0x19] = color.RGBA{40, 114, 0, 255}
	ppu.palScreen[0x1A] = color.RGBA{8, 124, 0, 255}
	ppu.palScreen[0x1B] = color.RGBA{0, 118, 40, 255}
	ppu.palScreen[0x1C] = color.RGBA{0, 102, 120, 255}
	ppu.palScreen[0x1D] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x1E] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x1F] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x20] = color.RGBA{236, 238, 236, 255}
	ppu.palScreen[0x21] = color.RGBA{76, 154, 236, 255}
	ppu.palScreen[0x22] = color.RGBA{120, 124, 236, 255}
	ppu.palScreen[0x23] = color.RGBA{176, 98, 236, 255}
	ppu.palScreen[0x24] = color.RGBA{228, 84, 236, 255}
	ppu.palScreen[0x25] = color.RGBA{236, 88, 180, 255}
	ppu.palScreen[0x26] = color.RGBA{236, 106, 100, 255}
	ppu.palScreen[0x27] = color.RGBA{212, 136, 32, 255}
	ppu.palScreen[0x28] = color.RGBA{160, 170, 0, 255}
	ppu.palScreen[0x29] = color.RGBA{116, 196, 0, 255}
	ppu.palScreen[0x2A] = color.RGBA{76, 208, 32, 255}
	ppu.palScreen[0x2B] = color.RGBA{56, 204, 108, 255}
	ppu.palScreen[0x2C] = color.RGBA{56, 180, 204, 255}
	ppu.palScreen[0x2D] = color.RGBA{60, 60, 60, 255}
	ppu.palScreen[0x2E] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x2F] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x30] = color.RGBA{236, 238, 236, 255}
	ppu.palScreen[0x31] = color.RGBA{168, 204, 236, 255}
	ppu.palScreen[0x32] = color.RGBA{188, 188, 236, 255}
	ppu.palScreen[0x33] = color.RGBA{212, 178, 236, 255}
	ppu.palScreen[0x34] = color.RGBA{236, 174, 236, 255}
	ppu.palScreen[0x35] = color.RGBA{236, 174, 212, 255}
	ppu.palScreen[0x36] = color.RGBA{236, 180, 176, 255}
	ppu.palScreen[0x37] = color.RGBA{228, 196, 144, 255}
	ppu.palScreen[0x38] = color.RGBA{204, 210, 120, 255}
	ppu.palScreen[0x39] = color.RGBA{180, 222, 120, 255}
	ppu.palScreen[0x3A] = color.RGBA{168, 226, 144, 255}
	ppu.palScreen[0x3B] = color.RGBA{152, 226, 180, 255}
	ppu.palScreen[0x3C] = color.RGBA{160, 214, 228, 255}
	ppu.palScreen[0x3D] = color.RGBA{160, 162, 160, 255}
	ppu.palScreen[0x3E] = color.RGBA{0, 0, 0, 255}
	ppu.palScreen[0x3F] = color.RGBA{0, 0, 0, 255}

	upLeft := image.Point{0, 0}
	lowRight := image.Point{len(ppu.palScreen), 1}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for i := 0; i < len(ppu.palScreen); i++{
		img.Set(i, 0, ppu.palScreen[i])
		ppu.SprFrames = append(ppu.SprFrames, pixel.R(float64(i), 0, float64(i+1), 1))
	}
	ppu.SprScreen = pixel.PictureDataFromImage(img)

	ppu.sprBatch = pixel.NewBatch(&pixel.TrianglesData{}, ppu.SprScreen)

	for i:= 0; i < len(ppu.palScreen); i++{
		ppu.PixelSprites[i] = pixel.NewSprite(ppu.SprScreen, ppu.SprFrames[i])
	}
	
	
}

