package terminal

import (
	"image"
	"sync"
	"time"

	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/arjenjb/cu"
	"golang.org/x/image/math/fixed"
)

type charSize struct {
	x, y float64
}

type ConsoleSettings struct {
	paddingX, paddingY unit.Dp

	// Tracks the last gtx.Constraint.Max to compare with the next render to check for any differences
	lastLayoutWidth, lastLayoutHeight int

	// If the constraints changed, we remember which aspects changed here
	lastLayoutChange             LayoutUpdateType
	lastLayoutChangeLock         sync.Mutex
	lastLayoutChangeTimerRunning bool
	lastLayoutChangeEvents       []LayoutChangedEvent

	charSizeCache map[charSizeCacheKey]charSize
	scrollTag     event.Tag

	constraints image.Rectangle
	leeway      unit.Dp
}

type LayoutUpdateType int

const (
	LayoutUpdateNone   LayoutUpdateType = iota
	LayoutUpdateWidth  LayoutUpdateType = 1 << 0
	LayoutUpdateHeight LayoutUpdateType = 1 << 1
)

func (s *ConsoleSettings) update(th cu.Theme, screen *Screen, gtx layout.Context) layout.Context {
	cs := s.getCharSize(
		screen.defaults.Font,
		gtx.Sp(screen.defaults.FontSize),
		th.Shaper,
	)

	if s.lastLayoutWidth != gtx.Constraints.Max.X {
		s.lastLayoutWidth = gtx.Constraints.Max.X
		s.markLayoutUpdate(LayoutUpdateWidth)
	}

	if s.lastLayoutHeight != gtx.Constraints.Max.Y {
		s.lastLayoutHeight = gtx.Constraints.Max.Y
		s.markLayoutUpdate(LayoutUpdateHeight)
	}

	// Constrain the emulator layout
	if s.constraints.Max.X > 0 {
		gtx.Constraints.Max.X = min(int(cs.x*float64(s.constraints.Max.X))+gtx.Dp(s.paddingX*2+s.leeway), gtx.Constraints.Max.X)
	}
	if s.constraints.Max.Y > 0 {
		gtx.Constraints.Max.Y = min(int(cs.y*float64(s.constraints.Max.Y))+gtx.Dp(s.paddingY*2), gtx.Constraints.Max.Y)
	}

	// Consume any events
	for _, evt := range s.Events() {
		if evt.Type&LayoutUpdateWidth > 0 {
			// Calculate the max screen size
			screenWidth := int(float64(gtx.Constraints.Max.X-gtx.Dp(s.paddingX*2+s.leeway)) / cs.x)
			screen.updateWidth(screenWidth)
		}

		if evt.Type&LayoutUpdateHeight > 0 {
			screenHeight := int(float64(gtx.Constraints.Max.Y-gtx.Dp(s.paddingY*2)) / cs.y)

			if s.constraints.Max.Y > 0 && screenHeight > s.constraints.Max.Y {
				screenHeight = s.constraints.Max.Y
			}

			screen.updateHeight(screenHeight)
		}
	}

	return gtx
}

func (s *ConsoleSettings) markLayoutUpdate(u LayoutUpdateType) {
	if u == LayoutUpdateNone {
		return
	}

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

func (s *ConsoleSettings) lastLayoutUpdateTimerCallback() {
	s.lastLayoutChangeLock.Lock()
	defer s.lastLayoutChangeLock.Unlock()

	s.lastLayoutChangeEvents = append(s.lastLayoutChangeEvents, LayoutChangedEvent{
		Type: s.lastLayoutChange,
	})

	s.lastLayoutChange = LayoutUpdateNone
	s.lastLayoutChangeTimerRunning = false
}

func (s *ConsoleSettings) Events() []LayoutChangedEvent {
	evts := s.lastLayoutChangeEvents[:]
	s.lastLayoutChangeEvents = nil

	return evts
}

type charSizeCacheKey struct {
	f font.Font
	s int
}

func (s *ConsoleSettings) getCharSize(f font.Font, sizePx int, shaper *text.Shaper) charSize {
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

	shaper.LayoutString(params, "W")
	g, _ := shaper.NextGlyph()
	var charWidth = g.Advance
	charWidthf := float64((charWidth.Mul(fixed.I(1000))).Ceil()) / 1000.0

	// This is the way the styledtext widget determines the line height
	shaper.LayoutString(params, "")
	g, _ = shaper.NextGlyph()
	charHeightf := float64(g.Ascent.Ceil() + g.Descent.Ceil())

	v := charSize{
		x: charWidthf,
		y: charHeightf,
	}
	s.charSizeCache[cacheKey] = v
	return v
}

type Option func(settings *ConsoleSettings)

func MaxSize(columns, rows int) Option {
	return func(settings *ConsoleSettings) {
		p := image.Pt(columns, rows)
		settings.constraints.Max = p
	}
}

func NewConsoleSettings(opts ...Option) *ConsoleSettings {
	offsetX := unit.Dp(10)
	offsetY := unit.Dp(6)

	s := &ConsoleSettings{
		paddingX:      offsetX,
		paddingY:      offsetY,
		charSizeCache: make(map[charSizeCacheKey]charSize),
		scrollTag:     new(bool),
		leeway:        unit.Dp(20),
	}

	for _, each := range opts {
		each(s)
	}

	return s
}
