package widget

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
	"image"
	"image/color"
)

type CheckboxWidget struct {
	th    *cu.Theme
	label string
	state *widget.Bool
}

func (c CheckboxWidget) Layout(gtx layout.Context) layout.Dimensions {
	dim := c.th.FlexRow(cu.Gap(cu.XS), cu.Align(layout.Middle)).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			dimn := c.renderBox(gtx)
			return dimn
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			dim := c.renderLabel(gtx)
			return dim
		}).
		Layout(gtx)
	return dim
}

func (c CheckboxWidget) renderBox(gtx layout.Context) layout.Dimensions {
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	return c.state.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		checked := c.state.Value
		pressed := c.state.Pressed()

		size := image.Pt(gtx.Dp(18), gtx.Dp(18))
		outer := image.Rectangle{
			Max: size,
		}

		inset := ifElse(!pressed, gtx.Dp(2), gtx.Dp(3))
		inner := outer.Inset(inset)

		r := gtx.Dp(3)

		gtx.Constraints.Min = size
		gtx.Constraints.Max = size

		defer clip.RRect{
			Rect: outer,
			SE:   r,
			SW:   r,
			NW:   r,
			NE:   r,
		}.Push(gtx.Ops).Pop()

		clipInner := clip.RRect{
			Rect: inner,
			SE:   r,
			SW:   r,
			NW:   r,
			NE:   r,
		}.Push(gtx.Ops)

		paint.Fill(gtx.Ops, white)

		if checked {
			paint.Fill(gtx.Ops, c.th.Color.Primary)

			offset := op.Offset(image.Pt(gtx.Dp(4), gtx.Dp(4))).Push(gtx.Ops)

			p := clip.Path{}
			p.Begin(gtx.Ops)
			p.MoveTo(f32.Pt(float32(gtx.Dp(1.5)), float32(gtx.Dp(5))))
			p.LineTo(f32.Pt(float32(gtx.Dp(4)), float32(gtx.Dp(7.5))))
			p.LineTo(f32.Pt(float32(gtx.Dp(8.5)), float32(gtx.Dp(2))))
			spec := p.End()
			paint.FillShape(gtx.Ops, white, clip.Stroke{Path: spec, Width: float32(gtx.Dp(2))}.Op())

			offset.Pop()
		} else if !pressed {
			w := gtx.Dp(1)
			paint.FillShape(gtx.Ops, c.th.Color.ControlBorder,
				clip.Stroke{
					Path:  clip.UniformRRect(inner, r).Path(gtx.Ops),
					Width: float32(w),
				}.Op(),
			)
		}

		clipInner.Pop()

		// Draw the blue border if pressed
		if pressed {
			paint.FillShape(gtx.Ops, c.th.Color.Primary,
				clip.Stroke{
					Path:  clip.UniformRRect(outer.Inset(gtx.Dp(1)), gtx.Dp(4)).Path(gtx.Ops),
					Width: float32(gtx.Dp(2)),
				}.Op(),
			)
		}

		return layout.Dimensions{
			Size:     size,
			Baseline: gtx.Dp(18),
		}
	})
}

func (c CheckboxWidget) renderLabel(gtx layout.Context) layout.Dimensions {
	return c.state.Layout(gtx, c.th.Text(c.label, cu.TextOptions{Truncate: true}))
}

func Checkbox(th *cu.Theme, ch *widget.Bool, label string) CheckboxWidget {
	return CheckboxWidget{
		th:    th,
		label: label,
		state: ch,
	}
}
