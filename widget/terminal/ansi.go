package terminal

import (
	"fmt"
	"gioui.org/font"
	"gioui.org/unit"
	"image/color"
)

// Normal colors
var (
	black   = color.NRGBA{A: 255}
	green   = color.NRGBA{R: 37, G: 188, B: 36, A: 255}
	blue    = color.NRGBA{R: 73, G: 46, B: 225, A: 255}
	red     = color.NRGBA{R: 0xED, G: 0x15, B: 0x15, A: 0xFF}
	yellow  = color.NRGBA{R: 173, G: 173, B: 39, A: 255}
	magenta = color.NRGBA{R: 211, G: 56, B: 211, A: 255}
	cyan    = color.NRGBA{R: 51, G: 187, B: 200, A: 255}
	white   = color.NRGBA{R: 223, G: 224, B: 225, A: 255}
)

// Bright colors
var (
	brightBlack = color.NRGBA{R: 129, G: 131, B: 131, A: 255}
	brightWhite = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
)

type Defaults struct {
	FgColor  color.NRGBA
	BgColor  color.NRGBA
	Font     font.Font
	BoldFont font.Font
	FontSize unit.Sp
}

type Style struct {
	// Actual state
	fgColor color.NRGBA
	BgColor color.NRGBA
	Bold    bool
	Faint   bool
}

func (s *Style) SetForegroundAnsi8(i int, bright bool) {
	if !bright {
		switch i {
		case 0:
			s.fgColor = black
		case 1:
			s.fgColor = red
		case 2:
			s.fgColor = green
		case 3:
			s.fgColor = yellow
		case 4:
			s.fgColor = blue
		case 5:
			s.fgColor = magenta
		case 6:
			s.fgColor = cyan
		case 7:
			s.fgColor = white
		}
	} else {
		switch i {
		case 0:
			s.fgColor = brightBlack
		case 1:
			s.fgColor = red
		case 2:
			s.fgColor = green
		case 3:
			s.fgColor = yellow
		case 4:
			s.fgColor = blue
		case 5:
			s.fgColor = magenta
		case 6:
			s.fgColor = cyan
		case 7:
			s.fgColor = brightWhite
		}
	}
}

func (s *Style) SetForegroundColor(c color.NRGBA) {
	s.fgColor = c
}

func (s *Style) SetBackgroundColor(bgColor color.NRGBA) {
	s.BgColor = bgColor
}

func parseHexColor(s string) (c color.NRGBA) {
	c.A = 0xff
	fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	return
}

func (s *Style) SetForegroundAnsi256(c int) {
	s.fgColor = parseHexColor(colorMap256[c])
}

func (s *Style) SetFaint(b bool) {
	s.Faint = true
}

func (s *Style) FgColor() color.NRGBA {
	c := s.fgColor
	if s.Faint {
		c.A = 128
	}

	return c
}

func (s *Style) Reset(fgColor color.NRGBA, bgColor color.NRGBA) {
	s.fgColor = fgColor
	s.BgColor = bgColor
	s.Bold = false
	s.Faint = false
}
