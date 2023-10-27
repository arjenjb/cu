package cu

import "gioui.org/font"

func NewDefaultTheme() *Theme {
	fonts := []font.FontFace{
		monoFontRegular,
		monoFontBold,
	}

	th = NewTheme(fonts)
	return th
}
