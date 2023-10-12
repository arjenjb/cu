package widget

import (
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"image"
	"image/color"
)

var colorDisabled = color.NRGBA{99, 99, 99, 255}
var colorHovered = color.NRGBA{37, 99, 235, 255}
var colorFocused = color.NRGBA{29, 78, 216, 255}

var colorPrimary = color.NRGBA{59, 130, 246, 255}

type ButtonLayoutStyle struct {
	CornerRadius unit.Dp
	Button       *widget.Clickable
}

func (b ButtonLayoutStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	min := gtx.Constraints.Min

	return b.Button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				var background = colorPrimary
				rr := gtx.Dp(b.CornerRadius)
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
				switch {
				case gtx.Queue == nil:
					background = colorDisabled

				case b.Button.Pressed():
					background = colorFocused
				case b.Button.Hovered():
					background = colorHovered
				}
				paint.Fill(gtx.Ops, background)

				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = min
				gtx.Constraints.Min.X = gtx.Dp(80)
				gtx.Constraints.Min.Y = gtx.Dp(26)
				return layout.Center.Layout(gtx, w)
			}),
		)
	})
}
