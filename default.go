package cu

import "gioui.org/font"

func NewDefaultTheme() Theme {
	fonts := []font.FontFace{
		monoFontRegular,
		monoFontBold,
	}

	return NewTheme(fonts)
}
