package terminal

import (
	_ "embed"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/x/styledtext"
	. "github.com/arjenjb/cu"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"strings"
)

var scrollTag = new(bool) // We could use &pressed for this instead.

func Console(th *Theme, screen *Screen) layout.Widget {
	offsetX := unit.Dp(10)
	offsetY := unit.Dp(6)

	return func(gtx layout.Context) layout.Dimensions {
		// Figure out character height
		params := text.Parameters{
			Font:    screen.defaults.Font,
			PxPerEm: fixed.I(gtx.Sp(screen.defaults.FontSize)),
		}

		th.Shaper.Layout(params, strings.NewReader("A"))
		g, _ := th.Shaper.NextGlyph()

		var charWidth = g.Advance
		var charHeight = int(g.Y) + g.Descent.Ceil()

		gtx.Constraints = layout.Exact(image.Point{
			X: charWidth.Mul(fixed.I(screen.Size.X)).Ceil() + gtx.Dp(offsetX*2+20),
			Y: charHeight*screen.Size.Y + gtx.Dp(offsetY*2+2),
		})

		return bordered(gtx, gtx.Dp(unit.Dp(1)), color.NRGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 128,
		}, func(gtx layout.Context) layout.Dimensions {
			for _, each := range gtx.Queue.Events(scrollTag) {
				switch evt := each.(type) {
				case pointer.Event:
					screen.scrollTop = min(max(
						screen.scrollTop+int(evt.Scroll.Y),
						0), screen.top)
				}
			}

			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, screen.defaults.BgColor)

			// Declare the tag.
			pointer.InputOp{
				Tag:   scrollTag,
				Types: pointer.Scroll,
				ScrollBounds: image.Rectangle{
					Min: image.Point{Y: -3},
					Max: image.Point{Y: 3}},
			}.Add(gtx.Ops)

			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    offsetY,
						Right:  offsetX,
						Bottom: offsetY,
						Left:   offsetX,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						t := styledtext.Text(th.Shaper, createSpansFrom(screen)...)
						t.WrapPolicy = styledtext.WrapGraphemes
						return t.Layout(gtx, nil)
					})
				}),

				// Render the scrollbar
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// scrollbar area
					offset := gtx.Dp(8)
					area := image.Rectangle{
						Min: image.Point{gtx.Constraints.Max.X - int(float32(offset)*1.7), offset},
						Max: image.Point{gtx.Constraints.Max.X - offset, gtx.Constraints.Max.Y - offset},
					}

					l := float32(screen.Size.Y) / float32(screen.top+screen.Size.Y)

					if l < 1.0 {
						total := float32(area.Dy())
						height := total * l

						offsetTop := (total - height) * float32(screen.scrollTop) / float32(screen.top)

						bar := area
						bar.Max.Y = area.Min.Y + int(height)
						bar = bar.Add(image.Point{Y: int(offsetTop)})

						defer clip.RRect{
							Rect: bar,
							SE:   gtx.Dp(3),
							SW:   gtx.Dp(3),
							NW:   gtx.Dp(3),
							NE:   gtx.Dp(3),
						}.Push(gtx.Ops).Pop()

						paint.Fill(gtx.Ops, color.NRGBA{
							R: 255,
							G: 255,
							B: 255,
							A: 20,
						})
					}

					return layout.Dimensions{
						Size:     image.Point{0, 0},
						Baseline: 0,
					}
				}),
			)
		})
	}
}

func bordered(gtx layout.Context, width int, c color.NRGBA, f layout.Widget) layout.Dimensions {
	defer clip.Rect{
		Max: gtx.Constraints.Max,
	}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
	return layout.UniformInset(unit.Dp(width)).Layout(gtx, f)
}

func createSpansFrom(screen *Screen) []styledtext.SpanStyle {
	var spans []styledtext.SpanStyle

	for _, line := range screen.VisibleLines() {
		x := 0
		for _, run := range line.runs {
			if run.start > x {
				spans = append(spans, styledtext.SpanStyle{Content: strings.Repeat(" ", run.start-x)})
			}

			f := screen.defaults.Font
			if run.style.Bold {
				f = screen.defaults.Font
			}

			spans = append(spans, styledtext.SpanStyle{
				Content: run.text,
				Size:    screen.defaults.FontSize,
				Color:   run.style.FgColor,
				Font:    f,
			})

			x = run.end()
		}
		spans = append(spans, styledtext.SpanStyle{Content: "\n"})
	}

	return spans
}
