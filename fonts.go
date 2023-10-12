package cu

import (
	_ "embed"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
)

var th *Theme

//go:embed fonts/FiraCode-Regular.ttf
var monoFontRegularData []byte

//go:embed fonts/FiraCode-Bold.ttf
var monoFontSemiBoldData []byte

//go:embed fonts/SF-Pro-Text-Regular.otf
var textFontRegularData []byte

//go:embed fonts/SF-Pro-Text-Semibold.otf
var textFontSemiBoldData []byte

//go:embed fonts/SF-Pro-Text-Bold.otf
var textFontBoldData []byte

// var monoFontLight text.FontFace
var monoFontRegular text.FontFace

// var monoFontMedium text.FontFace
var monoFontBold text.FontFace

var textFontRegular text.FontFace
var textFontSemiBold text.FontFace
var textFontBold text.FontFace

func init() {
	// Initialize fonts
	//face, _ := opentype.Parse(monoFontLightData)
	//monoFontLight = text.FontFace{Font: font.Font{Typeface: "Roboto", Variant: "Mono", Weight: font.Light}, Face: face}

	face, _ := opentype.Parse(monoFontRegularData)
	monoFontRegular = text.FontFace{Font: font.Font{Typeface: "monospace", Weight: font.Normal}, Face: face}
	//
	//face, _ = opentype.Parse(monoFontMediumData)
	//monoFontMedium = text.FontFace{Font: font.Font{Typeface: "Roboto", Variant: "Mono", Weight: font.Medium}, Face: face}

	face, _ = opentype.Parse(monoFontSemiBoldData)
	monoFontBold = text.FontFace{Font: font.Font{Typeface: "monospace", Weight: font.Bold}, Face: face}

	face, _ = opentype.Parse(textFontRegularData)
	textFontRegular = text.FontFace{Font: font.Font{Typeface: "sans-serif", Weight: font.Normal}, Face: face}

	face, _ = opentype.Parse(textFontSemiBoldData)
	textFontSemiBold = text.FontFace{Font: font.Font{Typeface: "sans-serif", Weight: font.SemiBold}, Face: face}

	face, _ = opentype.Parse(textFontBoldData)
	textFontBold = text.FontFace{Font: font.Font{Typeface: "sans-serif", Weight: font.Bold}, Face: face}
}
