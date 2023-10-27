package widget

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image"
	"image/color"
	"math"
	"time"
)

var refTime = time.Now()

type Spinner struct {
	R     int
	start time.Time
}

func (s Spinner) Layout(gtx layout.Context) layout.Dimensions {
	// The width of a segment
	w := float32(s.R) / 3.8

	// The height of a segment
	h := .62 * float32(s.R)

	// The radius of a segment
	radius := int(w / 2.0)

	// The dimensions of the full spinner
	size := image.Point{
		X: s.R * 2,
		Y: s.R * 2,
	}

	// calculate the segment rectangle, left and right position, we already know the width and the height
	l := (float32(s.R*2) - w) / 2.0
	r := l + w

	line := func(alpha float64) {
		defer clip.RRect{
			Rect: image.Rectangle{
				Min: image.Point{X: int(l)},
				Max: image.Point{X: int(r), Y: int(h)},
			},
			SE: radius,
			SW: radius,
			NW: radius,
			NE: radius,
		}.Push(gtx.Ops).Pop()

		paint.Fill(gtx.Ops, color.NRGBA{
			R: 100,
			G: 100,
			B: 100,
			A: uint8(255 - alpha*200),
		})
	}

	now := time.Now()
	d := (now.Sub(refTime)).Seconds()

	for i := 0; i < 8; i++ {
		y := (math.Sin((d*4.0)+(2.0*math.Pi)*(float64(i)/8.0)) + 1.0) / 2.0
		line(y)

		defer op.Affine(f32.Affine2D{}.Rotate(f32.Point{
			X: float32(s.R),
			Y: float32(s.R),
		}, (math.Pi/180)*-45)).Push(gtx.Ops).Pop()
	}

	op.InvalidateOp{
		At: now.Add(50 * time.Millisecond),
	}.Add(gtx.Ops)

	return layout.Dimensions{
		Size:     size,
		Baseline: 0,
	}
}
