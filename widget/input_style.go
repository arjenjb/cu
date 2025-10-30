package widget

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type InputStyle struct {
	CornerRadius unit.Dp
	Editor       *widget.Editor
	Width        unit.Dp
	Height       unit.Dp
}

func (b InputStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		// Draws the background of the button, with a border radius and handles
		// hovering and clicking style changes
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
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
			if b.Width > 0 {
				gtx.Constraints.Min.X = gtx.Dp(b.Width)
				gtx.Constraints.Max.X = gtx.Dp(b.Width)
			} else {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
			}

			if b.Height > 0 {
				gtx.Constraints.Min.Y = gtx.Dp(b.Height)
				gtx.Constraints.Max.Y = gtx.Dp(b.Height)
			} else if !b.Editor.SingleLine {
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			}

			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   12,
				Right:  12,
			}.Layout(gtx, w)
		}),

		// Render the scrollbar
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if b.Editor.SingleLine {
				return layout.Dimensions{}
			}

			// scrollbar area
			offset := gtx.Dp(8)
			area := image.Rectangle{
				Min: image.Point{
					X: gtx.Constraints.Max.X - int(float32(offset)*1.7),
					Y: gtx.Constraints.Min.Y + offset}, // width of the scrollbar
				Max: image.Point{
					X: gtx.Constraints.Max.X - offset,
					Y: gtx.Constraints.Max.Y - offset},
			}

			// The viewport height divided by the height of the document
			viewportHeight := float32(b.Editor.Dimensions().Size.Y)
			documentHeight := float32(b.Editor.FullDimensions().Size.Y)
			l := viewportHeight / documentHeight

			if l < 1.0 {
				scrollBarHeight := viewportHeight * l
				scrollMax := documentHeight - viewportHeight
				scrollTop := float32(b.Editor.ScrollTop())
				offsetTop := (viewportHeight - scrollBarHeight) * (scrollTop / scrollMax)

				bar := area
				bar.Max.Y = area.Min.Y + int(scrollBarHeight)
				bar = bar.Add(image.Point{Y: int(offsetTop)})

				defer clip.RRect{
					Rect: bar,
					SE:   gtx.Dp(3),
					SW:   gtx.Dp(3),
					NW:   gtx.Dp(3),
					NE:   gtx.Dp(3),
				}.Push(gtx.Ops).Pop()

				paint.Fill(gtx.Ops, color.NRGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 128,
				})
			}

			return layout.Dimensions{
				Size: gtx.Constraints.Max,
			}
		}),
	)
}
