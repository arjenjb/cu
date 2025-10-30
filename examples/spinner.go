package main

import (
	"gioui.org/layout"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

func spinnerExample(th cu.Theme) layout.Widget {
	return th.M(cu.M, widget.Spinner{
		R: 12,
	}.Layout)
}
