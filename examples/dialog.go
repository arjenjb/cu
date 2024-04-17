package main

import (
	"fmt"
	"gioui.org/layout"
	widget2 "gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/dialog"
	"github.com/arjenjb/cu/widget"
)

var btnAlertDialog = &widget2.Clickable{}
var btnConfirmDialog = &widget2.Clickable{}

func dialogExample(gtx layout.Context, th *cu.Theme) func(gtx layout.Context) layout.Dimensions {

	if btnAlertDialog.Clicked(gtx) {
		go func() {
			d := dialog.NewAlertDialog(th)
			d.Title = "Alert!!!"
			d.BigMessage = "Beware"
			d.NormalMessage = "All your base are belong to us"
			d.Show()
		}()
	}

	if btnConfirmDialog.Clicked(gtx) {
		go func() {
			d := dialog.NewConfirmDialog(th)
			d.Title = "Are you sure"
			d.NormalMessage = "Do you want to continue?"
			d.AcceptLabel = "Continue"
			if d.Show() {
				fmt.Println("User accepted")
			} else {
				fmt.Println("User canceled")
			}
		}()
	}

	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(widget.Button(th, btnAlertDialog, "Alert dialog").Layout).
		Rigid(widget.Button(th, btnConfirmDialog, "Confirm dialog").Layout).
		Layout
}
