package main

import (
	"gioui.org/layout"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

func determinateProgressExample(th cu.Theme, progress widget.LinearProgress) layout.Widget {
	return th.FlexColumn().
		Rigid(th.Text("Waiting for process to finish")).
		Rigid(widget.NewProgressBar(th, progress).Layout).
		Layout
}

func indeterminateProgressExample(th cu.Theme, progress widget.Progress) layout.Widget {
	return th.FlexColumn().
		Rigid(th.Text("Waiting for process to finish")).
		Rigid(widget.NewProgressBar(th, progress).Layout).
		Layout
}
