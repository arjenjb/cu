package terminal

import (
	"gioui.org/font"
	"image/color"
	"strings"
	"sync"
)

type Point struct {
	X, Y int
}

func Pt(x, y int) Point {
	return Point{X: x, Y: y}
}

type Screen struct {
	cursor *Point
	Size   Point

	mu       sync.Mutex
	lines    []Line
	style    Style
	defaults Defaults

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

	s.notifyUpdated()

	return len(p), nil
}

func (s *Screen) WriteNewLine() {
	// The current line gets a line break, or allocate one if there is none
	if s.cursor.Y == len(s.lines) {
		s.lines = append(s.lines, Line{})
	}

	s.lines[s.cursor.Y].brk = true

	s.cursor.Y++
	s.cursor.X = 0

	// Allocate the new line(s)
	for s.cursor.Y >= len(s.lines) {
		s.appendLine()
	}
}

func (s *Screen) WriteString(s2 string) error {
	_, err := s.Write([]byte(s2))
	return err
}

// The following are primitives that are used by the parser
func (s *Screen) SetForegroundColorAnsi8(c int, b bool) {
	s.style.SetForegroundAnsi8(c, b)
}

func (s *Screen) SetForegroundColor(c color.NRGBA) {
	s.style.SetForegroundColor(c)
}

func (s *Screen) Reset() {
	s.style.Reset(s.defaults.FgColor, s.defaults.BgColor)
}

func (s *Screen) SetBold(b bool) {
	s.style.Bold = b
}

func (s *Screen) WriteCharacters(text string) {
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

	if s.scrollTop >= len(s.lines) {
		s.scrollTop = len(s.lines)
	}

	from := s.scrollTop
	to := min(s.scrollTop+s.Size.Y, len(s.lines))

	return s.lines[from:to]
}

func (s *Screen) appendLine() {
	// And update the top reference
	follow := s.scrollTop == s.scrollMax()

	// Create a new line
	s.lines = append(s.lines, Line{})

	if follow {
		s.scrollTop = s.scrollMax()
	}
}

// updateWidth updates the with of the terminal and recalculates the lines in the terminal
func (s *Screen) updateWidth(width int) {
	// constrain the width
	if width == s.Size.X {
		return
	}
	var displayLines []Line

	s.mu.Lock()
	defer s.mu.Unlock()

	var scrollTopTranslated bool
	var newLineNumber int

	for _, l := range s.virtualLines() {
		// rewrite top, and scrollTop if needed
		if !scrollTopTranslated && s.scrollTop <= l.endLine {
			scrollTopTranslated = true
			s.scrollTop = newLineNumber
		}

		// Split the line
		newLines := l.Split(width)
		newLineNumber += len(newLines)
		displayLines = append(displayLines, newLines...)
	}

	s.Size.X = width
	s.lines = displayLines
}

// Reconstructs the lines as originally emitted to the terminal independent of terminal width
func (s *Screen) virtualLines() []VirtualLine {
	var result []VirtualLine

	for i := 0; i < len(s.lines); i++ {
		l := s.lines[i]

		v := VirtualLine{
			startLine: i,
			endLine:   i,
			runs:      l.runs,
		}

		// Join consecutive lines
		for !l.brk && i+1 < len(s.lines) {
			i++
			l = s.lines[i]
			v = v.AppendLine(l)
		}

		result = append(result, v)
	}

	return result
}

func (s *Screen) updateHeight(newHeight int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Is it the same height?
	if newHeight == s.Size.Y {
		return
	}

	// rewrite scrollOffset
	delta := min(s.Size.Y, len(s.lines)) - newHeight
	//s.top += delta
	s.scrollTop = max(s.scrollTop+delta, 0)
	s.Size.Y = newHeight

	s.notifyUpdated()
}

func (s *Screen) SetForegroundColorAnsi256(c int) {
	s.style.SetForegroundAnsi256(c)
}

func (s *Screen) SetFaint(b bool) {
	s.style.SetFaint(true)
}

func (s *Screen) notifyUpdated() {
	// Notify the channel that the screen has been updated
	if s.updatedChannel != nil {
		// async notify
		select {
		case s.updatedChannel <- struct{}{}:
		default:
		}
	}
}

func (s *Screen) scrollMax() int {
	return max(len(s.lines)-s.Size.Y, 0)
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
		defaults:       defaults,
		style: Style{
			fgColor: defaults.FgColor,
			BgColor: defaults.BgColor,
			Bold:    false,
		},
	}
}
