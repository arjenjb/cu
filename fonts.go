package cu

import (
	_ "embed"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
)

//go:embed fonts/FiraCode-Regular.ttf
var monoFontRegularData []byte

//go:embed fonts/FiraCode-Bold.ttf
var monoFontSemiBoldData []byte

// var monoFontLight text.FontFace
var monoFontRegular text.FontFace

// var monoFontMedium text.FontFace
var monoFontBold text.FontFace

func init() {
	// Initialize fonts

	face, _ := opentype.Parse(monoFontRegularData)
	monoFontRegular = text.FontFace{Font: font.Font{Typeface: "monospace", Weight: font.Normal}, Face: face}

	face, _ = opentype.Parse(monoFontSemiBoldData)
	monoFontBold = text.FontFace{Font: font.Font{Typeface: "monospace", Weight: font.Bold}, Face: face}
}
