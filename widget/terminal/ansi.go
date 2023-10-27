package terminal

import (
	"gioui.org/font"
	"gioui.org/unit"
	"image/color"
)

// Normal colors
var (
	black   = color.NRGBA{A: 255}
	green   = color.NRGBA{R: 37, G: 188, B: 36, A: 255}
	blue    = color.NRGBA{R: 73, G: 46, B: 225, A: 255}
	red     = color.NRGBA{R: 170, A: 255}
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
	FgColor color.NRGBA
	BgColor color.NRGBA
	Bold    bool
}

func (s *Style) SetForegroundAnsi(i int, bright bool) {
	if !bright {

		switch i {
		case 0:
			s.FgColor = black
		case 1:
			s.FgColor = red
		case 2:
			s.FgColor = green
		case 3:
			s.FgColor = yellow
		case 4:
			s.FgColor = blue
		case 5:
			s.FgColor = magenta
		case 6:
			s.FgColor = cyan
		case 7:
			s.FgColor = white
		}
	} else {
		switch i {
		case 0:
			s.FgColor = brightBlack
		case 1:
			s.FgColor = red
		case 2:
			s.FgColor = green
		case 3:
			s.FgColor = yellow
		case 4:
			s.FgColor = blue
		case 5:
			s.FgColor = magenta
		case 6:
			s.FgColor = cyan
		case 7:
			s.FgColor = brightWhite
		}
	}
}

func (s *Style) SetForegroundColor(c color.NRGBA) {
	s.FgColor = c
}

func (s *Style) SetBackgroundColor(bgColor color.NRGBA) {
	s.BgColor = bgColor
}
