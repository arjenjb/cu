package terminal

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"log/slog"
	"strconv"
	"strings"
	"unicode"
)

type AnsiToken interface {
}

type EmptyToken struct {
	AnsiToken
}

type NumberToken struct {
	AnsiToken
	Number int
}

type CharacterToken struct {
	AnsiToken
	Character byte
}

type EscapeSequence struct {
	Tokens []AnsiToken
}

func (t EscapeSequence) ApplyOn(s *Screen) error {
	final := t.Tokens[len(t.Tokens)-1]
	switch v := final.(type) {
	case CharacterToken:
		switch v.Character {
		case 'm':
			return t.applyGraphicRenditionOn(s)
		case 'H':
			return t.applyCursorPosition(s)
		default:
			return fmt.Errorf("cannot deal with this escape sequence: %c", v.Character)
		}
	}

	return nil
}

func (t EscapeSequence) applyGraphicRenditionOn(s *Screen) error {
	code, err := t.ExpectNumber(0, 0)
	if err != nil {
		return err
	}

	if code == 0 {
		s.Reset()
	} else if code == 1 {
		s.SetBold(true)
	} else if code == 2 {
		// Faint not implemented
		s.SetFaint(true)

	} else if code >= 30 && code <= 37 {
		s.SetForegroundColorAnsi8(code-30, false)
	} else if code == 38 {
		// RGB Code
		n, err := t.ExpectNumber(1, 0)
		if err != nil {
			return err
		}
		if n == 2 {
			// Gather the rest
			r, _ := t.ExpectNumber(2, 0)
			g, _ := t.ExpectNumber(3, 0)
			b, _ := t.ExpectNumber(4, 0)
			s.SetForegroundColor(color.NRGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255,
			})
		} else if n == 5 {
			c, _ := t.ExpectNumber(2, 0)
			s.SetForegroundColorAnsi256(c)
		}

	} else if code >= 90 && code <= 97 {
		s.SetForegroundColorAnsi8(code-90, true)
	} else {
		slog.Debug("Unsupported code", "code", code)
	}

	return nil
}

func (t EscapeSequence) applyCursorPosition(s *Screen) error {
	if len(t.Tokens) != 3 {
		return errors.New("invalid escape sequence")
	}
	line, err := t.ExpectNumber(0, 1)
	if err != nil {
		return err
	}
	column, err := t.ExpectNumber(1, 1)
	if err != nil {
		return err
	}

	s.cursor.Y = line - 1
	s.cursor.X = column - 1
	return nil
}

func (t EscapeSequence) ExpectNumber(n int, default_ int) (int, error) {
	if len(t.Tokens) < n+1 {
		return -1, errors.New("invalid token")
	}

	token := t.Tokens[n]
	switch v := token.(type) {
	case EmptyToken:
		return default_, nil
	case NumberToken:
		return v.Number, nil
	default:
		return 0, errors.New("invalid token")
	}
}

type AnsiReader struct {
	Input []byte
	p     int

	Screen *Screen
}

func (r *AnsiReader) Parse() {
	for r.HasNext() {
		if bytes.Equal(r.PeekN(1), []byte{'\n'}) {
			r.Next()
			r.Screen.WriteNewLine()

		} else if bytes.Equal(r.PeekN(1), []byte{'\r'}) {
			r.Next()
			r.Screen.WriteCarriageReturn()

		} else if bytes.Equal(r.PeekN(2), []byte{'\x1b', '['}) {
			seq := r.parseEscapeSequence()
			seq.ApplyOn(r.Screen)

		} else {
			start := r.p
			r.p++
			r.UpUntil(func(b byte) bool { return b == '\x1b' || b == '\n' || b == '\r' })

			text := string(r.Input[start:r.p])
			r.Screen.WriteCharacters(text)
		}
	}
}

func (r *AnsiReader) HasNext() bool {
	return r.p < len(r.Input)
}

func (r *AnsiReader) Next() byte {
	b := r.Input[r.p]
	r.p++
	return b
}

func (r *AnsiReader) PeekN(i int) []byte {
	if len(r.Input[r.p:]) < i {
		return []byte{}
	}
	return r.Input[r.p : r.p+i]
}

func (r *AnsiReader) UpUntil(f func(b byte) bool) string {
	start := r.p
	for r.HasNext() {
		if f(r.Input[r.p]) {
			break
		}
		r.p++
	}
	return string(r.Input[start:r.p])
}

func (r *AnsiReader) parseEscapeSequence() EscapeSequence {
	r.p += 2

	var tokens []AnsiToken
	lastWasNumber := false

	for r.HasNext() {
		b := r.Input[r.p]

		if unicode.IsDigit(rune(b)) {
			buf := strings.Builder{}

			for r.HasNext() {
				b = r.Input[r.p]
				if unicode.IsDigit(rune(b)) {
					buf.WriteByte(b)
					r.p++
				} else {
					break
				}
			}

			n, _ := strconv.Atoi(buf.String())
			tokens = append(tokens, NumberToken{Number: n})
			lastWasNumber = true
		} else if b == ';' {
			if !lastWasNumber {
				tokens = append(tokens, EmptyToken{})
			}
			lastWasNumber = false
			r.p++
		} else {
			tokens = append(tokens, CharacterToken{Character: b})
			r.p++
			break
		}
	}

	return EscapeSequence{Tokens: tokens}
}
