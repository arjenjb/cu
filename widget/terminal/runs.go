package terminal

// Runs should always be consecutive
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

func (r Run) length() int {
	return len(r.text)
}

// Shift the start position of this run by length
func (r Run) offset(o int) Run {
	if o == 0 {
		return r
	}

	r.start += o
	return r
}

func (r Run) Split(maxLength int) (Run, Run, bool) {
	if r.end() <= maxLength {
		return r, Run{}, false
	} else {
		l := maxLength - r.start
		head := Run{
			start: r.start,
			text:  r.text[0:l],
			style: r.style,
		}

		tail := Run{
			start: r.start + l,
			text:  r.text[l:],
			style: r.style,
		}

		return head, tail, true
	}
}

func appendRun(runs []Run, tail Run) []Run {
	if len(runs) == 0 {
		return []Run{tail}
	}

	i := len(runs) - 1
	last := runs[i]

	if last.end() == tail.start && last.style == tail.style {
		runs[i].text += tail.text
		return runs
	} else {
		return append(runs, tail)
	}
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
