package widget

import (
	"gioui.org/font"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	theme2 "github.com/arjenjb/cu/theme"
	"image/color"
)

type ButtonWidget struct {
	label     string
	clickable *widget.Clickable
	theme     *theme2.Theme
}

func (b *ButtonWidget) Layout(ctx layout.Context) layout.Dimensions {
	return b.clickable.Layout(ctx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)

		return ButtonLayoutStyle{
			CornerRadius: 4,
			Button:       b.clickable,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			colMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			}}.Add(gtx.Ops)

			return widget.Label{Alignment: text.Middle}.Layout(gtx, b.theme.Shaper, font.Font{
				Typeface: "sans-serif",
			}, 12, b.label, colMacro.Stop())
		})

	})
}

func Button(t *theme2.Theme, clickable *widget.Clickable, label string) *ButtonWidget {
	return &ButtonWidget{
		theme:     t,
		label:     label,
		clickable: clickable,
	}
}
