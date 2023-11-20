package main

import (
	"gioui.org/layout"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

func progressExample(th *cu.Theme, progress float32) layout.Widget {
	return th.FlexColumn().
		Rigid(th.Text("Waiting for process to finish")).
		Rigid(widget.NewProgressBar(th, progress).Layout).
		Layout
}
