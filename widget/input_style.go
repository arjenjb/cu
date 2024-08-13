package widget

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"image"
)

type InputStyle struct {
	CornerRadius unit.Dp
	Editor       *widget.Editor
	Width        unit.Dp
}

func (b InputStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		// Draws the background of the button, with a border radius and handles
		// hovering and clicking style changes
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			//gtx.Constraints.Min.X = b.Width
			//gtx.Constraints.Max.X = b.Width

			rect := image.Rectangle{Max: gtx.Constraints.Min}
			inner := rect.Inset(gtx.Dp(4))

			var background = colorNormal
			var borderColor = colorNone

			if b.Editor.ReadOnly {
				background = colorNone

			}

			rr := gtx.Dp(b.CornerRadius)

			switch {
			case b.Editor.ReadOnly:
				background = colorDisabled
			case gtx.Focused(b.Editor):
				// Draw the outline of an active button
				w := gtx.Dp(2)
				paint.FillShape(gtx.Ops, colorPrimary,
					clip.Stroke{
						Path:  clip.UniformRRect(rect.Inset(w), rr+w).Path(gtx.Ops),
						Width: float32(w),
					}.Op(),
				)
			}

			if !b.Editor.ReadOnly {
				shape := clip.UniformRRect(inner, rr)
				defer shape.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, background)
			}

			// draw the border
			borderColor = borderColorNormal
			w := gtx.Dp(1)

			paint.FillShape(gtx.Ops, borderColor,
				clip.Stroke{
					Path:  clip.UniformRRect(inner, rr).Path(gtx.Ops),
					Width: float32(w),
				}.Op(),
			)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		// Draw the text on top of the background
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(b.Width)
			gtx.Constraints.Max.X = gtx.Dp(b.Width)

			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   12,
				Right:  12,
			}.Layout(gtx, w)
		}),
	)
}
