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

var colorDisabled = color.NRGBA{142, 142, 147, 255}
var colorPrimaryHovered = color.NRGBA{37, 99, 235, 255}
var colorPrimaryFocused = color.NRGBA{29, 78, 216, 255}
var colorPrimary = color.NRGBA{59, 130, 246, 255}
var colorNormal = color.NRGBA{255, 255, 255, 255}
var colorNone = color.NRGBA{0, 0, 0, 0}

var borderColorNormal = color.NRGBA{0xC9, 0xCC, 0xD6, 255}

type ButtonStyle struct {
	CornerRadius unit.Dp
	Disabled     bool
	Primary      bool
	Button       *widget.Clickable
}

func (b ButtonStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	min := gtx.Constraints.Min

	btn := func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,

			// Draws the background of the button, with a border radius and handles
			// hovering and clicking style changes
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				rect := image.Rectangle{Max: gtx.Constraints.Min}
				inner := rect.Inset(gtx.Dp(4))

				var background = colorNormal
				var borderColor = colorNone

				if b.Disabled {
					background = colorNone
				} else if b.Primary {
					background = colorPrimary
				}

				rr := gtx.Dp(b.CornerRadius)

				switch {
				case gtx.Queue == nil || b.Disabled:
					background = colorDisabled

				case b.Button.Pressed():
					// Draw the outline of an active button
					w := gtx.Dp(2)
					paint.FillShape(gtx.Ops, colorPrimary,
						clip.Stroke{
							Path:  clip.UniformRRect(rect.Inset(w), rr+w).Path(gtx.Ops),
							Width: float32(w),
						}.Op(),
					)
				}

				if !b.Disabled {
					shape := clip.UniformRRect(inner, rr)
					defer shape.Push(gtx.Ops).Pop()
					paint.Fill(gtx.Ops, background)
				}

				if !b.Primary {
					borderColor = borderColorNormal
					w := gtx.Dp(1)

					paint.FillShape(gtx.Ops, borderColor,
						clip.Stroke{
							Path:  clip.UniformRRect(inner, rr).Path(gtx.Ops),
							Width: float32(w),
						}.Op(),
					)
				}

				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),

			// Draw the text on top of the background
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = min
				return layout.Center.Layout(gtx, w)
			}),
		)
	}

	if b.Disabled {
		return btn(gtx)
	} else {
		return b.Button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			semantic.Button.Add(gtx.Ops)
			return btn(gtx)
		})
	}
}
