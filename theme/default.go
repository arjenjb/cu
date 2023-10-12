package theme

import "gioui.org/font"

func NewDefaultTheme() *Theme {
	fonts := []font.FontFace{
		monoFontRegular,
		monoFontBold,
		textFontRegular,
		textFontSemiBold,
		textFontBold,
	}
	th = NewTheme(fonts)
	return th
}
