package terminal

import (
	"gioui.org/font"
	"github.com/ag5/go-commons/c5"
	"image/color"
	"math"
	"strings"
	"sync"
)

type Point struct {
	X, Y int
}

type Run struct {
	start int
	text  string
	style Style
}

func (r Run) end() int {
	return r.start + len(r.text)
}

func (r Run) CopyTo(x int) Run {
	if x > r.end() {
		return r
	}

	return Run{
		start: r.start,
		text:  r.text[:x-r.start],
		style: r.style,
	}
}

func (r Run) CopyFrom(i int) Run {
	return Run{
		start: i,
		text:  r.text[i-r.start:],
		style: r.style,
	}
}

func (r Run) IsEmpty() bool {
	return len(r.text) == 0
}

type Line struct {
	width int
	runs  []Run
}

func (l Line) EffectiveLines() int {
	return int(math.Floor(float64(c5.Sum(c5.Map(l.runs, func(e Run) int {
		return len(e.text)
	}))) / float64(l.width)))
}

func writeText(in []Run, text string, x int, style Style) []Run {
	// We want to write the given text to the screen at position x. First add the runs preceding the position of the
	// new text. A run that overlaps with the position of the new text we only copy that part just before it.
	var out []Run
	var i = 0

	// Check preceding
	for i < len(in) {
		r := in[i]
		if r.start < x && r.end() < x {
			out = append(out, r)
			i++
		} else if r.start < x {
			out = append(out, r.CopyTo(x))
			break
		} else {
			break
		}
	}

	out = append(out, Run{
		start: x,
		text:  text,
		style: style,
	})

	// Add remaining runs
	for i < len(in) {
		r := in[i]
		i++

		end := x + len(text)
		if r.end() <= end {
			continue
		}

		if r.start >= end {
			out = append(out, r)
		} else {
			out = append(out, r.CopyFrom(end))
		}
	}

	return out
}

func (l *Line) Write(text string, x int, style Style) {
	l.runs = writeText(l.runs, text, x, style)
}

func (l *Line) String() string {
	var sb strings.Builder
	for _, r := range l.runs {
		sb.WriteString(r.text)
	}
	return sb.String()
}

type Screen struct {
	cursor *Point
	Size   Point

	mu       sync.Mutex
	lines    []Line
	style    Style
	defaults Defaults

	// Tracks the line that is currently at the top of the screen, independent of scrolling
	top int

	// This follows the top variable, but can be overridden by scrolling the window
	scrollTop int

	updatedChannel chan interface{}
}

func (s *Screen) Lines() []Line {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lines[:]
}

func (s *Screen) Write(p []byte) (n int, err error) {
	// Scan the string
	r := &AnsiReader{
		Input:  p,
		p:      0,
		Screen: s,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r.Parse()

	// Notify the channel that the screen has been updated
	if s.updatedChannel != nil {
		// async notify
		go func() { s.updatedChannel <- struct{}{} }()
	}

	return len(p), nil
}

func (s *Screen) WriteNewLine() {
	s.cursor.Y++
	s.cursor.X = 0

	// Allocate the new line
	if s.cursor.Y > len(s.lines) {
		s.appendLine()
	}
}

func (s *Screen) WriteString(s2 string) error {
	_, err := s.Write([]byte(s2))
	return err
}

// The following are primitives that are used by the parser

func (s *Screen) SetForegroundColor(c int, b bool) {
	s.style.SetForegroundAnsi(c, b)
}

func (s *Screen) ResetColors() {
	s.style.SetForegroundColor(s.defaults.FgColor)
	s.style.SetBackgroundColor(s.defaults.BgColor)
}

func (s *Screen) SetBold(b bool) {
	s.style.Bold = b
}

func (s *Screen) WriteCharacters(text string) {
	// Ensure we have a line to write to

	// Split into chunks
	for len(text) > 0 {
		x, y := s.cursor.X, s.cursor.Y
		for y >= len(s.lines) {
			s.appendLine()
		}

		rest := s.Size.X - x
		e := min(len(text), rest)

		s.lines[y].Write(text[0:e], x, s.style)
		s.cursor.X += len(text)

		text = text[e:]
		// If there is any more text to display, wrap to the next line
		if len(text) > 0 {
			s.cursor.Y++
			s.cursor.X = 0
		}
	}

}

func (s *Screen) CursorUp(i int) {
	s.cursor.Y -= i
}

func (s *Screen) CursorRight(i int) {
	s.cursor.X += i
}

func (s *Screen) Buffer() string {
	buffer := strings.Builder{}
	for _, each := range s.lines {
		buffer.WriteString(each.String())
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (s *Screen) WriteCarriageReturn() {
	s.cursor.X = 0
}

func (s *Screen) VisibleLines() []Line {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lines[s.scrollTop:min(s.scrollTop+s.Size.Y, len(s.lines))]
}

func (s *Screen) appendLine() {
	s.lines = append(s.lines, Line{
		width: s.Size.X,
	})
	dy := max((len(s.lines)-s.top)-s.Size.Y, 0)

	if s.scrollTop == s.top {
		s.scrollTop += dy
	}
	s.top += dy
}

func min(i int, i2 int) int {
	if i < i2 {
		return i
	}
	return i2
}

func NewScreen(size Point, updatedChannel chan interface{}) *Screen {
	// background color
	defaults := Defaults{
		FgColor:  white,
		BgColor:  color.NRGBA{29, 29, 29, 255},
		Font:     font.Font{Typeface: "monospace"},
		BoldFont: font.Font{Typeface: "monospace", Weight: font.Bold},
		FontSize: 12,
	}

	return &Screen{
		updatedChannel: updatedChannel,
		cursor:         &Point{},
		Size:           size,
		lines:          nil,
		top:            0,
		defaults:       defaults,
		style: Style{
			FgColor: defaults.FgColor,
			BgColor: defaults.BgColor,
			Bold:    false,
		},
	}
}
