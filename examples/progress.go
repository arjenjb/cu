package main

import (
	"gioui.org/layout"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

func progressExample(th *cu.Theme) layout.Widget {
	return th.FlexColumn().
		Rigid(th.Text("Waiting for process to finish")).
		Rigid(widget.NewProgressBar(th, 0.5).Layout).
		Layout
}
