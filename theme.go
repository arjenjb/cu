package cu

import (
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"image/color"
)

type Palette struct {
	Text color.NRGBA
}

type Fonts struct {
	SansSerif font.Font
	Monospace font.Font
}

type Theme struct {
	Color    Palette
	TextSize unit.Sp
	Shaper   *text.Shaper
	Font     Fonts
}

type Flex struct {
	widget   layout.Flex
	gap      unit.Dp
	children []layout.FlexChild
}

func (f Flex) Layout(gtx layout.Context) layout.Dimensions {
	children := f.children
	if f.gap > 0 {
		var spacer layout.Widget
		if f.widget.Axis == layout.Horizontal {
			spacer = HSpacer(f.gap)
		} else {
			spacer = VSpacer(f.gap)
		}

		n := []layout.FlexChild{}
		for i, each := range children {
			if i > 0 {
				n = append(n, layout.Rigid(spacer))
			}
			n = append(n, each)
		}
		children = n
	}

	return f.widget.Layout(gtx, children...)
}

func (f Flex) Rigid(w layout.Widget) Flex {
	f.children = append(f.children, layout.Rigid(w))
	return f
}

func (f Flex) Flexed(weight float32, w layout.Widget) Flex {
	f.children = append(f.children, layout.Flexed(weight, w))
	return f
}

type CuFlexOption func(w *Flex, t Theme)

func Spacing(s layout.Spacing) func(w *Flex, t Theme) {
	return func(w *Flex, t Theme) {
		w.widget.Spacing = s
	}
}

func Align(a layout.Alignment) func(w *Flex, t Theme) {
	return func(w *Flex, t Theme) {
		w.widget.Alignment = a
	}
}

type Unit interface {
	Dp(t Theme) unit.Dp
}

type scaledUnit struct {
	scale float32
}

func (s scaledUnit) Dp(t Theme) unit.Dp {
	return unit.Dp(s.scale) * unit.Dp(16)
}

func Scaled(n float32) Unit {
	return scaledUnit{scale: n}
}

var (
	XS Unit = Scaled(0.33)
	S  Unit = Scaled(0.66)
	M  Unit = Scaled(1)
	L  Unit = Scaled(2)
)

func Gap(s Unit) func(w *Flex, t Theme) {
	return func(w *Flex, t Theme) {
		w.gap = s.Dp(t)
	}
}

func (t Theme) FlexRow(options ...CuFlexOption) Flex {
	l := Flex{
		widget: layout.Flex{
			Axis: layout.Horizontal,
		},
		gap: 0,
	}
	for _, opt := range options {
		opt(&l, t)
	}
	return l
}

func (t Theme) FlexColumn(options ...CuFlexOption) Flex {
	l := Flex{
		widget: layout.Flex{
			Axis: layout.Vertical,
		},
		gap: 0,
	}
	for _, opt := range options {
		opt(&l, t)
	}
	return l
}

func (t Theme) Hr() layout.Widget {
	return func(gtx C) D {
		h := unit.Dp(1)

		gtx.Constraints.Min.Y = gtx.Dp(h)

		var path clip.Path
		path.Begin(gtx.Ops)
		path.MoveTo(f32.Pt(0, float32(h)/2))
		path.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(h)/2))
		path.Close()

		paint.FillShape(gtx.Ops,
			color.NRGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 255 / 5,
			},
			clip.Stroke{
				Path:  path.End(),
				Width: float32(h),
			}.Op())

		return layout.Spacer{Height: h}.Layout(gtx)
	}
}

func (t Theme) Background(gtx layout.Context) {
	colorBackground := color.NRGBA{
		R: 0xF7,
		G: 0xF8,
		B: 0xFA,
		A: 255,
	}
	paint.Fill(gtx.Ops, colorBackground)
}

func NewTheme(fonts []font.FontFace) *Theme {
	var colorText = color.NRGBA{R: 52, G: 65, B: 85, A: 0xff}

	t := &Theme{
		Shaper: text.NewShaper(text.WithCollection(fonts)),
		Color: Palette{
			Text: colorText,
		},
		TextSize: unit.Sp(13.0),
		Font: Fonts{
			SansSerif: font.Font{Typeface: "Roboto, SF Pro, Segoe UI, Dejavu, sans-serif"},
			Monospace: font.Font{Typeface: "monospace"},
		},
	}

	return t
}