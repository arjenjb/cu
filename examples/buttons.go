package main

import (
	"cu"
	"cu/widget"
	"gioui.org/layout"
	widget2 "gioui.org/widget"
)

var button1 widget2.Clickable
var button2 widget2.Clickable
var button3 widget2.Clickable

func buttonExample(th *cu.Theme) layout.Widget {

	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(widget.Button(th, &button1, "Cancel").Layout).
		Rigid(widget.Button(th, &button2, "Disabled", widget.Disabled(true)).Layout).
		Rigid(widget.Button(th, &button3, "Go", widget.Primary()).Layout).
		Layout
}
