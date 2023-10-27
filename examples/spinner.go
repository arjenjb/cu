package main

import (
	"cu"
	"cu/widget"
	"gioui.org/layout"
)

func spinnerExample(th *cu.Theme) layout.Widget {
	return th.M(cu.M, widget.Spinner{
		R: 12,
	}.Layout)
}
