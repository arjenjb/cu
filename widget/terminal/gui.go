package terminal

import (
	_ "embed"
	"gioui.org/font"
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
	"sync"
	"time"
)

var scrollTag = new(bool) // We could use &pressed for this instead.

type charSize struct {
	x, y float64
}

type consoleSettings struct {
	paddingX, paddingY unit.Dp

	// Tracks the last gtx.Constraint.Max to compare with the next render to check for any differences
	lastLayoutWidth, lastLayoutHeight int

	// If the constraints changed, we remember which aspects changed here
	lastLayoutChange             LayoutUpdateType
	lastLayoutChangeLock         sync.Mutex
	lastLayoutChangeTimerRunning bool
	lastLayoutChangeEvents       []LayoutChangedEvent

	charSizeCache map[charSizeCacheKey]charSize
}

type LayoutUpdateType int

const (
	LayoutUpdateNone   LayoutUpdateType = iota
	LayoutUpdateWidth  LayoutUpdateType = 1 << 0
	LayoutUpdateHeight LayoutUpdateType = 1 << 1
)

func (s *consoleSettings) update(th *Theme, screen *Screen, gtx layout.Context) {
	if s.lastLayoutWidth != gtx.Constraints.Max.X {
		s.lastLayoutWidth = gtx.Constraints.Max.X
		s.markLayoutUpdate(LayoutUpdateWidth)
	}

	if s.lastLayoutHeight != gtx.Constraints.Max.Y {
		s.lastLayoutHeight = gtx.Constraints.Max.Y
		s.markLayoutUpdate(LayoutUpdateHeight)
	}

	// Consume any events
	for _, evt := range s.Events() {
		if evt.Type&LayoutUpdateWidth > 0 {
			// Calculate the max screen size
			charWidth := s.getCharSize(screen.defaults.Font, gtx.Sp(screen.defaults.FontSize), th.Shaper).x
			screenWidth := int(float64(gtx.Constraints.Max.X-gtx.Dp(s.paddingX*2+20)) / charWidth)

			screen.updateWidth(screenWidth)
		}

		if evt.Type&LayoutUpdateHeight > 0 {
			charHeight := s.getCharSize(screen.defaults.Font, gtx.Sp(screen.defaults.FontSize), th.Shaper).y
			screenHeight := int(float64(gtx.Constraints.Max.Y-gtx.Dp(s.paddingY*2)) / charHeight)

			screen.updateHeight(screenHeight)
		}
	}
}

func (s *consoleSettings) markLayoutUpdate(u LayoutUpdateType) {
	s.lastLayoutChangeLock.Lock()
	defer s.lastLayoutChangeLock.Unlock()

	s.lastLayoutChange |= u

	if !s.lastLayoutChangeTimerRunning {
		s.lastLayoutChangeTimerRunning = true
		go func() {
			time.Sleep(20 * time.Millisecond)
			s.lastLayoutUpdateTimerCallback()
		}()
	}
}

type LayoutChangedEvent struct {
	Type LayoutUpdateType
}

func (s *consoleSettings) lastLayoutUpdateTimerCallback() {
	s.lastLayoutChangeLock.Lock()
	defer s.lastLayoutChangeLock.Unlock()

	s.lastLayoutChangeEvents = append(s.lastLayoutChangeEvents, LayoutChangedEvent{
		Type: s.lastLayoutChange,
	})

	s.lastLayoutChangeTimerRunning = false
	s.lastLayoutChange = LayoutUpdateNone
}

func (s *consoleSettings) Events() []LayoutChangedEvent {
	evts := s.lastLayoutChangeEvents[:]
	s.lastLayoutChangeEvents = nil

	return evts
}

type charSizeCacheKey struct {
	f font.Font
	s int
}

func (s *consoleSettings) getCharSize(f font.Font, sizePx int, shaper *text.Shaper) charSize {
	cacheKey := charSizeCacheKey{
		f: f,
		s: sizePx,
	}

	if v, found := s.charSizeCache[cacheKey]; found {
		return v
	}

	params := text.Parameters{
		Font:    f,
		PxPerEm: fixed.I(sizePx),
	}

	shaper.Layout(params, strings.NewReader("A"))
	g, _ := shaper.NextGlyph()

	var charWidth = g.Advance
	var charHeight = int(g.Y) + g.Descent.Ceil()

	// Add 20px to the X position to allow a little leeway in rendering rounding
	charWidthf := float64((charWidth.Mul(fixed.I(1000))).Ceil()) / 1000.0
	charHeightf := float64(charHeight)

	v := charSize{
		x: charWidthf,
		y: charHeightf,
	}
	s.charSizeCache[cacheKey] = v
	return v
}

func NewConsoleSettings() *consoleSettings {
	offsetX := unit.Dp(10)
	offsetY := unit.Dp(6)

	return &consoleSettings{
		paddingX:      offsetX,
		paddingY:      offsetY,
		charSizeCache: make(map[charSizeCacheKey]charSize),
	}
}

func Console(th *Theme, screen *Screen, settings *consoleSettings) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Figure out character height
		settings.update(th, screen, gtx)

		return bordered(gtx, 1, color.NRGBA{
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
						Top:    settings.paddingY,
						Right:  settings.paddingX,
						Bottom: settings.paddingY,
						Left:   settings.paddingX,
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

//func getState(tag interface{}, th *Theme, screen *Screen, gtx layout.Context) *consoleCache {
//	charWidth, charHeight := getCharSize(th.Shaper, screen.defaults.Font, gtx.Sp(screen.defaults.FontSize))
//
//	screenWidth := float64(gtx.Constraints.Max.X-gtx.Dp(offsetX*2+20)) / charWidth
//	screenHeight := float64(gtx.Constraints.Max.Y-gtx.Dp(offsetY*2+2)) / charHeight
//
//	fmt.Printf("Size: %dx%d\n", math.Floor(width), math.Floor(height))
//
//	gtx.Constraints = layout.Exact(image.Point{
//		X: charWidth.Mul(fixed.I(screen.Size.X)).Ceil() + gtx.Dp(offsetX*2+20),
//		Y: charHeight*screen.Size.Y + gtx.Dp(offsetY*2+2),
//	})
//
//}

func bordered(gtx layout.Context, width unit.Dp, c color.NRGBA, f layout.Widget) layout.Dimensions {
	defer clip.Rect{
		Max: gtx.Constraints.Max,
	}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
	return layout.UniformInset(width).Layout(gtx, f)
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
