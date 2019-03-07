// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ncs"
)

var (
	hues = []string{
		"Y",
		"Y10R",
		"Y20R",
		"Y30R",
		"Y40R",
		"Y50R",
		"Y60R",
		"Y70R",
		"Y80R",
		"Y90R",
		"R",
		"R10B",
		"R20B",
		"R30B",
		"R40B",
		"R50B",
		"R60B",
		"R70B",
		"R80B",
		"R90B",
		"B",
		"B10G",
		"B20G",
		"B30G",
		"B40G",
		"B50G",
		"B60G",
		"B70G",
		"B80G",
		"B90G",
		"G",
		"G10Y",
		"G20Y",
		"G30Y",
		"G40Y",
		"G50Y",
		"G60Y",
		"G70Y",
		"G80Y",
		"G90Y",
	}
)

const (
	screenWidth  = 80 * 9
	screenHeight = 80 * 9

	boxWidth  = 80
	boxHeight = 80
)

func drawColorBox(screen *ebiten.Image, c ncs.Color, x, y int) {
	ebitenutil.DrawRect(screen, float64(x), float64(y), boxWidth, boxWidth, c)
}

func drawColorBoxText(screen *ebiten.Image, c ncs.Color, x, y int) {
	gb, _, _ := bitmapfont.Gothic12r.GlyphBounds('M')
	lineHeight := bitmapfont.Gothic12r.Metrics().Height.Ceil()

	tx := x + 16
	ty := y + 16 + -gb.Min.Y.Ceil()
	text.Draw(screen, colorHex(c), bitmapfont.Gothic12r, tx+1, ty+lineHeight+1, color.RGBA{0, 0, 0, 0x80})
	text.Draw(screen, c.String(), bitmapfont.Gothic12r, tx+1, ty+1, color.RGBA{0, 0, 0, 0x80})
	text.Draw(screen, colorHex(c), bitmapfont.Gothic12r, tx, ty+lineHeight, color.White)
	text.Draw(screen, c.String(), bitmapfont.Gothic12r, tx, ty, color.White)
}

func uint8Hex(b uint8) string {
	// fmt.Sprintf is slow on browsers. Use faster implementation here.
	str := ""
	b1 := b >> 4
	b2 := b & 0xf
	if b1 <= 9 {
		str += string(b1 + '0')
	} else {
		str += string(b1 - 0xa + 'A')
	}
	if b2 <= 9 {
		str += string(b2 + '0')
	} else {
		str += string(b2 - 0xa + 'A')
	}
	return str
}

func colorHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return "#" + uint8Hex(uint8(r>>8)) + uint8Hex(uint8(g>>8)) + uint8Hex(uint8(b>>8))
}

func adjustColor(c ncs.Color, blackness int, chromaticness int, hue int) ncs.Color {
	const unit = 10

	if blackness != 0 {
		if c.Blackness == 99 {
			c.Blackness = 100
		}
		c.Blackness += blackness * unit
		if c.Blackness >= 100 {
			c.Blackness = 99
		}
		if c.Blackness < 0 {
			c.Blackness = 0
		}
		if c.Blackness > 100-c.Chromaticness {
			if c.Chromaticness == 99 {
				c.Blackness = 0
			} else {
				c.Blackness = 100 - c.Chromaticness
			}
		}
	}

	if chromaticness != 0 {
		if c.Chromaticness == 99 {
			c.Chromaticness = 100
		}
		c.Chromaticness += chromaticness * unit
		if c.Chromaticness >= 100 {
			c.Chromaticness = 99
		}
		if c.Chromaticness < 0 {
			c.Chromaticness = 0
		}
		if c.Chromaticness > 100-c.Blackness {
			if c.Blackness == 99 {
				c.Chromaticness = 0
			} else {
				c.Chromaticness = 100 - c.Blackness
			}
		}
	}

	if hue != 0 {
		c.Hue += hue * unit
		for c.Hue < 0 {
			c.Hue += 400
		}
		for c.Hue >= 400 {
			c.Hue -= 400
		}
	}

	return c
}

type state struct {
	color ncs.Color
}

func (s *state) update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.color = adjustColor(s.color, 1, 0, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		s.color = adjustColor(s.color, -1, 0, 0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.color = adjustColor(s.color, 0, 1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		s.color = adjustColor(s.color, 0, -1, 0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		s.color = adjustColor(s.color, 0, 0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		s.color = adjustColor(s.color, 0, 0, -1)
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	for j := -4; j <= 4; j++ {
		for i := -4; i <= 4; i++ {
			c := adjustColor(s.color, j, 0, i)
			x := (screenWidth-boxWidth)/2 + i*boxWidth
			y := (screenHeight-boxHeight)/2 + j*boxHeight
			drawColorBox(screen, c, x, y)
		}
	}

	for j := -4; j <= 4; j++ {
		for i := -4; i <= 4; i++ {
			c := adjustColor(s.color, j, 0, i)
			x := (screenWidth-boxWidth)/2 + i*boxWidth
			y := (screenHeight-boxHeight)/2 + j*boxHeight
			drawColorBoxText(screen, c, x, y)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%0.2f", ebiten.CurrentFPS()))

	return nil
}

func main() {
	c, err := ncs.Parse("1050-R90B")
	if err != nil {
		panic(err)
	}
	s := &state{
		color: c,
	}

	if err := ebiten.Run(s.update, screenWidth, screenHeight, 1, "NCS Palette"); err != nil {
		panic(err)
	}
}
