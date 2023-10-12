package theme

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type D = layout.Dimensions
type C = layout.Context

func textWidget(t Theme, label string, size unit.Sp, weight font.Weight) layout.Widget {
	return func(gtx C) D {
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: t.Color.Text}.Add(gtx.Ops)
		return widget.Label{}.Layout(gtx, t.Shaper, font.Font{Typeface: "sans-serif", Weight: weight}, size, label, colMacro.Stop())
	}
}

func (t Theme) H1(label string) layout.Widget {
	return t.Mb(1.2, textWidget(t, label, t.TextSize*1, font.Bold))
}

func (t Theme) Paragraph(label string) layout.Widget {
	return t.Mb(1, textWidget(t, label, t.TextSize, font.Normal))
}
