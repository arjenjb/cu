package widget

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/arjenjb/cu"
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

type ProgressBarStyle struct {
	Background color.NRGBA
	Foreground color.NRGBA
	Height     unit.Dp
}

var DefaultProgressBarStyle = ProgressBarStyle{
	Height:     unit.Dp(4),
	Background: colorProgressBg,
	Foreground: colorProgressPrimary,
}

type ProgressBar struct {
	th       cu.Theme
	Progress Progress
	Style    ProgressBarStyle
}

type Progress interface{ implementsProgress() }

type LinearProgress float32
type indeterminedProgressType time.Time

var IndeterminedProgress Progress

func init() {
	IndeterminedProgress = NewIndeterminedProgress()
}
func NewIndeterminedProgress() Progress {
	return indeterminedProgressType(time.Now())
}

func (indeterminedProgressType) implementsProgress() {}
func (LinearProgress) implementsProgress()           {}

func NewProgressBar(th cu.Theme, progress Progress) ProgressBar {
	return ProgressBar{th: th, Progress: progress, Style: DefaultProgressBarStyle}
}

func (p ProgressBar) Layout(gtx layout.Context) layout.Dimensions {
	h := gtx.Dp(p.Style.Height)
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

	paint.Fill(gtx.Ops, p.Style.Background)

	// Render the progress on top
	switch pr := p.Progress.(type) {
	case LinearProgress:
		defer clip.UniformRRect(image.Rectangle{
			Min: image.Point{Y: y},
			Max: image.Point{
				X: int(float32(gtx.Constraints.Max.X) * float32(pr)),
				Y: y + h,
			},
		}, r).Push(gtx.Ops).Pop()

	case indeterminedProgressType:
		elapsed := gtx.Now.Sub(time.Time(pr)).Milliseconds()
		p := calculateBar1Position(elapsed)

		defer clip.UniformRRect(image.Rectangle{
			Min: image.Point{
				X: int(p.left * float64(gtx.Constraints.Max.X)),
				Y: y},
			Max: image.Point{
				X: int(p.right * float64(gtx.Constraints.Max.X)),
				Y: y + h,
			},
		}, r).Push(gtx.Ops).Pop()

		gtx.Execute(op.InvalidateCmd{})
	}

	paint.Fill(gtx.Ops, p.Style.Foreground)

	return dims
}
