package widget

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/arjenjb/cu"
	"image"
	"image/color"
)

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
	th       *cu.Theme
	Progress float32
}

func NewProgressBar(th *cu.Theme, progress float32) ProgressBar {
	return ProgressBar{th: th, Progress: progress}
}

func (p ProgressBar) Layout(gtx layout.Context) layout.Dimensions {
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

	// Render the background
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

	// Render the progress on top
	defer clip.UniformRRect(image.Rectangle{
		Min: image.Point{Y: y},
		Max: image.Point{
			X: int(float32(gtx.Constraints.Max.X) * p.Progress),
			Y: y + h,
		},
	}, r).Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, colorProgressPrimary)

	return dims
}
