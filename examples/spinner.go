package main

import (
	"gioui.org/layout"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

func spinnerExample(th *cu.Theme) layout.Widget {
	return th.M(cu.M, widget.Spinner{
		R: 12,
	}.Layout)
}

func progressExample(th *cu.Theme) layout.Widget {
	return th.FlexColumn().
		Rigid(th.Text("Downloading smalltalk image...")).
		Rigid(widget.NewProgressBar(th, 0.5).Layout).
		//Rigid(th.Text("12 kB/s", cu.TextOptions{Color: &th.Color.TextSecondary, Size: th.TextSizeMedium})).
		Layout
}
