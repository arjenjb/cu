package cu

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Palette struct {
	Text          color.NRGBA
	TextSecondary color.NRGBA
	TextDisabled  color.NRGBA

	SelectionActive color.NRGBA

	// The primary highlight color
	Primary       color.NRGBA
	ControlBorder color.NRGBA
	Link          color.NRGBA
}

type Fonts struct {
	SansSerif font.Font
	Monospace font.Font
}

type Theme struct {
	Color          Palette
	TextSize       unit.Sp
	TextSizeMedium unit.Sp
	TextSizeH1     unit.Sp
	TextSizeH2     unit.Sp
	Shaper         *text.Shaper
	Font           Fonts

	LineHeightH1 unit.Sp
	LineHeightH2 unit.Sp
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
	XS Unit = Scaled(0.25)
	S  Unit = Scaled(0.5)
	M  Unit = Scaled(1)
	L  Unit = Scaled(1.33)
)

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

func (t Theme) Link(label string, button *widget.Clickable) layout.Widget {
	return func(gtx C) D {
		var hovered bool
		if button.Hovered() {
			hovered = true
		}

		return button.Layout(gtx, func(gtx C) D {
			dim := t.Text(label, TextOptions{Color: &t.Color.Link})(gtx)

			pointer.CursorPointer.Add(gtx.Ops)

			if hovered {
				p := clip.Path{}
				p.Begin(gtx.Ops)
				p.MoveTo(f32.Point{Y: float32(dim.Size.Y - 1)})
				p.LineTo(f32.Point{Y: float32(dim.Size.Y - 1), X: float32(dim.Size.X)})

				paint.FillShape(gtx.Ops,
					t.Color.Link,
					clip.Stroke{
						Path:  p.End(),
						Width: float32(1),
					}.Op())
			}

			return dim
		})
	}
}

func NewTheme(fonts []font.FontFace) Theme {
	// var colorText = color.NRGBA{R: 52, G: 65, B: 85, A: 0xff}
	var (
		colorText            = color.NRGBA{R: 0, G: 0, B: 0, A: 0xff}
		colorTextSecondary   = color.NRGBA{R: 0x6C, G: 0x70, B: 0x7E, A: 0xff}
		colorTextDisabled    = color.NRGBA{R: 0xA8, G: 0xAD, B: 0xBD, A: 0xFF}
		colorPrimary         = color.NRGBA{R: 59, G: 130, B: 246, A: 255}
		colorLink            = color.NRGBA{R: 0x31, G: 0x5F, B: 0xBD, A: 0xFF}
		colorControlBorder   = color.NRGBA{R: 0xC9, G: 0xCC, B: 0xD6, A: 0xFF}
		colorSelectionActive = color.NRGBA{R: 0xD4, G: 0xE2, B: 0xFF, A: 0xFF}
	)

	t := Theme{
		Shaper: text.NewShaper(text.WithCollection(fonts)),
		Color: Palette{
			Text:            colorText,
			TextSecondary:   colorTextSecondary,
			TextDisabled:    colorTextDisabled,
			SelectionActive: colorSelectionActive,
			Primary:         colorPrimary,
			ControlBorder:   colorControlBorder,
			Link:            colorLink,
		},
		TextSize:       unit.Sp(13.0),
		TextSizeMedium: unit.Sp(12.0),
		TextSizeH1:     unit.Sp(20.0),
		TextSizeH2:     unit.Sp(16.0),
		LineHeightH1:   unit.Sp(24),
		LineHeightH2:   unit.Sp(20),
		Font: Fonts{
			SansSerif: font.Font{Typeface: "Roboto, SF Pro Text, Segoe UI, Dejavu, sans-serif"},
			Monospace: font.Font{Typeface: "monospace"},
		},
	}

	return t
}
