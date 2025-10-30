package widget

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/arjenjb/cu"
)

type TextInputWidget struct {
	Editor *widget.Editor
	Hint   string
	theme  cu.Theme
	Width  unit.Dp
	Height unit.Dp
}

func (t TextInputWidget) Layout(gtx layout.Context) layout.Dimensions {
	mt := material.NewTheme()
	mt.Shaper = t.theme.Shaper
	mt.Face = t.theme.Font.SansSerif.Typeface
	mt.TextSize = t.theme.TextSize

	return InputStyle{
		CornerRadius: 4,
		Editor:       t.Editor,
		Width:        t.Width,
		Height:       t.Height,
	}.Layout(gtx, material.Editor(mt, t.Editor, t.Hint).Layout)
}

func TextInput(th cu.Theme, editor *widget.Editor, hint string, width unit.Dp, height unit.Dp) TextInputWidget {
	return TextInputWidget{
		theme:  th,
		Editor: editor,
		Hint:   hint,
		Width:  width,
		Height: height,
	}
}
