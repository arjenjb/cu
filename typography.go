package cu

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type D = layout.Dimensions
type C = layout.Context

func textWidget(t Theme, label string, size unit.Sp, weight font.Weight, alignment text.Alignment) layout.Widget {
	f := t.Font.SansSerif
	f.Weight = weight

	return func(gtx C) D {
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: t.Color.Text}.Add(gtx.Ops)
		return widget.Label{
			Alignment: alignment,
		}.Layout(gtx, t.Shaper, f, size, label, colMacro.Stop())
	}
}

type TextOptions struct {
	Size     float32
	Bold     bool
	Centered bool
}

func (t Theme) Text(label string, opts ...TextOptions) layout.Widget {
	size := t.TextSize
	weight := font.Normal
	alignment := text.Start

	for _, opt := range opts {
		if opt.Size == 0.0 {
			size = t.TextSize
		} else {
			size = unit.Sp(opt.Size) * t.TextSize
		}

		if opt.Bold {
			weight = font.Bold
		}

		if opt.Centered {
			alignment = text.Middle
		}
	}

	return textWidget(t, label, size, weight, alignment)
}

func (t Theme) H1(label string) layout.Widget {
	return t.Mb(Scaled(1.2), textWidget(t, label, t.TextSize*1, font.Bold, text.Start))
}

func (t Theme) Paragraph(label string) layout.Widget {
	return t.Mb(M, textWidget(t, label, t.TextSize, font.Normal, text.Start))
}
