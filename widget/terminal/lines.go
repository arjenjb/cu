package terminal

import "strings"

// This is a line as displayed in the terminal emulator, it will never be longer than its current width
type Line struct {
	runs []Run
	brk  bool // A terminal emulator line can be backed by a longer actual line, this bool indicates that a hard break is required
}

type VirtualLine struct {
	startLine, endLine int // The line numbers this virtual line spans in the actual console
	runs               []Run
}

func (v VirtualLine) AppendLine(l2 Line) VirtualLine {
	nl := VirtualLine{
		runs:      v.runs,
		startLine: v.startLine,
		endLine:   v.endLine + 1,
	}

	for _, each := range l2.runs {
		nl.AddRun(each.offset(v.Length()))
	}

	return nl
}

func (v *VirtualLine) AddRun(tail Run) {
	v.runs = appendRun(v.runs, tail)
}

func (v VirtualLine) Length() int {
	n := 0
	for _, each := range v.runs {
		n += each.length()
	}
	return n
}

func (v *VirtualLine) Write(text string, x int, style Style) {
	v.runs = writeText(v.runs, text, x, style)
}

func (v VirtualLine) Split(maxLength int) []Line {
	var ls []Line

	var current Line

	for _, each := range v.runs {
		rest := each
		for !rest.IsEmpty() {
			var head Run
			var split bool
			head, rest, split = rest.Split(maxLength)
			current.AddRun(head)

			if split {
				// compensate for the maxlength
				rest = rest.offset(-maxLength)

				ls = append(ls, current)
				current = Line{}
			}
		}
	}

	current.brk = true
	ls = append(ls, current)
	return ls
}

func (l *Line) Write(text string, x int, style Style) {
	l.runs = writeText(l.runs, text, x, style)
}

func (l *Line) String() string {
	var sb strings.Builder
	p := 0
	for _, r := range l.runs {
		for i := p; i < r.start; i++ {
			sb.WriteByte(' ')
		}
		p = r.start + len(r.text)

		sb.WriteString(r.text)
	}
	return sb.String()
}

func (l *Line) Length() int {
	n := 0
	for _, each := range l.runs {
		n += each.length()
	}
	return n
}

func (l *Line) AddRun(tail Run) {
	l.runs = appendRun(l.runs, tail)
}
