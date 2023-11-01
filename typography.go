package cu

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"image/color"
)

type D = layout.Dimensions
type C = layout.Context

func textWidget(t Theme, label string, size unit.Sp, lineHeight unit.Sp, weight font.Weight, alignment text.Alignment, c color.NRGBA, maxLines int) layout.Widget {
	f := t.Font.SansSerif
	f.Weight = weight

	return func(gtx C) D {
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: c}.Add(gtx.Ops)
		return widget.Label{
			Alignment:  alignment,
			LineHeight: lineHeight,
			MaxLines:   maxLines,
		}.Layout(gtx, t.Shaper, f, size, label, colMacro.Stop())
	}
}

type TextOptions struct {
	Size     unit.Sp
	Bold     bool
	Centered bool
	Color    *color.NRGBA
	Truncate bool
}

func (t Theme) Text(label string, opts ...TextOptions) layout.Widget {
	size := t.TextSize
	weight := font.Normal
	alignment := text.Start
	color := t.Color.Text
	maxLines := 0

	for _, opt := range opts {
		if opt.Size != 0.0 {
			size = opt.Size
		}

		if opt.Bold {
			weight = font.Bold
		}

		if opt.Centered {
			alignment = text.Middle
		}

		if opt.Color != nil {
			color = *opt.Color
		}

		if opt.Truncate {
			maxLines = 1
		}
	}

	return textWidget(t, label, size, 0, weight, alignment, color, maxLines)
}

func (t Theme) H1(label string) layout.Widget {
	return textWidget(t, label, t.TextSizeH1, t.LineHeightH1, font.Bold, text.Start, t.Color.Text, 0)
}

func (t Theme) H2(label string) layout.Widget {
	return textWidget(t, label, t.TextSizeH2, t.LineHeightH2, font.Bold, text.Start, t.Color.Text, 0)
}

func (t Theme) Paragraph(label string) layout.Widget {
	return t.Mb(M, textWidget(t, label, t.TextSize, unit.Sp(18), font.Normal, text.Start, t.Color.Text, 0))
}
