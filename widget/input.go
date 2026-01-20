package widget

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/arjenjb/cu"
)

type TextInputWidget struct {
	Editor   *widget.Editor
	theme    cu.Theme
	Hint     string
	Width    unit.Dp
	Height   unit.Dp
	FontFace font.Typeface
	TextSize unit.Sp
}

func (t TextInputWidget) Layout(gtx layout.Context) layout.Dimensions {
	mt := material.NewTheme()
	mt.Shaper = t.theme.Shaper
	mt.TextSize = t.TextSize
	mt.Face = t.FontFace

	return InputStyle{
		CornerRadius: 4,
		Editor:       t.Editor,
		Width:        t.Width,
		Height:       t.Height,
	}.Layout(gtx, material.Editor(mt, t.Editor, t.Hint).Layout)
}

type TextInputOption func(w *TextInputWidget)

func TextInput(th cu.Theme, editor *widget.Editor, options ...TextInputOption) TextInputWidget {
	w := TextInputWidget{
		theme:    th,
		Editor:   editor,
		FontFace: th.Font.SansSerif.Typeface,
		TextSize: th.TextSize,
	}

	for _, each := range options {
		each(&w)
	}

	return w
}
