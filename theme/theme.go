package theme

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

type Theme struct {
	Color    Palette
	TextSize unit.Sp
	Shaper   *text.Shaper
}

func (t Theme) FlexRow() layout.Flex {
	return layout.Flex{
		Axis: layout.Horizontal,
	}
}

func (t Theme) FlexColumn() layout.Flex {
	return layout.Flex{
		Axis: layout.Vertical,
	}
}

func (t Theme) Hr() layout.Widget {
	return func(gtx C) D {
		h := unit.Dp(2)

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

func NewTheme(fonts []font.FontFace) *Theme {
	var colorText = color.NRGBA{R: 52, G: 65, B: 85, A: 0xff}

	t := &Theme{
		Shaper: text.NewShaper(text.NoSystemFonts(), text.WithCollection(fonts)),
		Color: Palette{
			Text: colorText,
		},
		TextSize: unit.Sp(13.0),
	}

	return t
}
