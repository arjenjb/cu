package main

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/rs/zerolog/log"
	"image"
	"image/color"
)

type CheckboxWidget struct {
	th    *cu.Theme
	label string
	state *widget.Bool
}

var checkable widget.Bool

func If[T any](c bool, a T, b T) T {
	if c {
		return a
	} else {
		return b
	}
}

func (c CheckboxWidget) Layout(gtx layout.Context) layout.Dimensions {
	return c.th.FlexRow(cu.Align(layout.Middle), cu.Gap(cu.XS)).
		Rigid(c.renderBox).
		Flexed(1, c.renderLabel).
		Layout(gtx)
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

		inset := If(!pressed, gtx.Dp(2), gtx.Dp(3))
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
			Baseline: 0,
		}
	})
}

func (c CheckboxWidget) renderLabel(gtx layout.Context) layout.Dimensions {
	return c.state.Layout(gtx, c.th.Text(c.label))
}

func Checkbox(th *cu.Theme, ch *widget.Bool, label string) CheckboxWidget {
	return CheckboxWidget{
		th:    th,
		label: label,
		state: ch,
	}
}

func checkboxExample(th *cu.Theme) layout.Widget {
	if checkable.Changed() {
		log.Info().Msgf("Toggled: %s", If(checkable.Value, "On", "Off"))
	}
	
	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(Checkbox(th, &checkable, "Check the box or this label").Layout).
		Layout
}

func main() {
	w := app.NewWindow(app.Size(unit.Dp(670), unit.Dp(524)))
	th := cu.NewDefaultTheme()

	var ops op.Ops

	go func() {
		for {
			select {
			case evt := <-w.Events():
				switch e := evt.(type) {
				case system.DestroyEvent:
					return

				case system.FrameEvent:
					gtx := layout.NewContext(&ops, e)

					th.Background(gtx)

					th.M(cu.M, th.FlexColumn(cu.Gap(cu.M)).
						Rigid(buttonExample(th)).
						Rigid(th.Hr()).
						Rigid(checkboxExample(th)).
						Rigid(th.Hr()).
						Rigid(spinnerExample(th)).
						Rigid(th.Hr()).
						Rigid(progressExample(th)).
						Layout)(gtx)

					e.Frame(gtx.Ops)
				}
			}
		}
	}()

	app.Main()
}
