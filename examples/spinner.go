package main

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
	"image"
	"image/color"
)

func spinnerExample(th *cu.Theme) layout.Widget {
	return th.M(cu.M, widget.Spinner{
		R: 12,
	}.Layout)
}

var colorProgressPrimary = color.NRGBA{
	R: 0x35,
	G: 0x74,
	B: 0xF0,
	A: 0xFF,
}
var colorProgressBg = color.NRGBA{
	R: 0xDF,
	G: 0xE1,
	B: 0xE5,
	A: 0xFF,
}

type ProgressBar struct {
	Total    int
	Progress int
}

type ProgressWidget struct {
	Bar ProgressBar
}

func (p ProgressWidget) Layout(gtx layout.Context) layout.Dimensions {
	h := gtx.Dp(unit.Dp(4))
	r := h / 2

	dims := layout.Dimensions{
		Size: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Dp(unit.Dp(16)),
		},
		Baseline: 0,
	}

	y := (dims.Size.Y / 2) - (h / 2)

	defer clip.UniformRRect(image.Rectangle{
		Min: image.Point{
			Y: y,
		},
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: y + h,
		},
	}, r).Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, colorProgressBg)

	defer clip.UniformRRect(image.Rectangle{
		Min: image.Point{Y: y},
		Max: image.Point{
			X: int(float32(gtx.Constraints.Max.X) * 0.33),
			Y: y + h,
		},
	}, r).Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, colorProgressPrimary)

	return dims

}

func progressExample(th *cu.Theme) layout.Widget {
	return th.FlexColumn().
		Rigid(th.Text("Downloading smalltalk image...")).
		Rigid(ProgressWidget{}.Layout).
		//Rigid(th.Text("12 kB/s", cu.TextOptions{Color: &th.Color.TextSecondary, Size: th.TextSizeMedium})).
		Layout
}
