package dialog

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	widget2 "gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

type confirmDialog struct {
	Theme *cu.Theme
	// Title is shown in the titlebar, if left empty it will default to Message
	Title string
	// BigMessage is shown in H2 style at the top of window
	BigMessage string
	// NormalMessage is an optional additional message shown in normal font
	NormalMessage string

	AcceptLabel string
	CancelLabel string

	btnAccept *widget2.Clickable
	btnCancel *widget2.Clickable
	Accepted  bool
}

func (a confirmDialog) Layout(gtx layout.Context) layout.Dimensions {
	th := a.Theme
	th.Background(gtx)

	return th.M(cu.M,
		th.FlexColumn(cu.Gap(cu.S)).
			RigidIf(len(a.BigMessage) > 0, th.H2(a.BigMessage)).
			Rigid(th.Text(a.NormalMessage)).
			Flexed(1, th.FlexRow(cu.Align(layout.End)).
				Flexed(1, cu.HSpacer(0)).
				Rigid(widget.Button(th, a.btnCancel, a.CancelLabel).Layout).
				Rigid(widget.Button(th, a.btnAccept, a.AcceptLabel, widget.Primary()).Layout).
				Layout).
			Layout)(gtx)
}

func (a confirmDialog) Show() bool {
	w := app.NewWindow(func(metric unit.Metric, config *app.Config) {
		config.Title = a.Title
	}, app.Size(400, 140))

	var ops op.Ops

	var done = false
	for !done {
		select {
		case evt := <-w.Events():
			switch e := evt.(type) {
			case system.DestroyEvent:
				done = true

			case system.FrameEvent:
				if a.btnCancel.Clicked() {
					w.Perform(system.ActionClose)
				}
				if a.btnAccept.Clicked() {
					a.Accepted = true
					w.Perform(system.ActionClose)
				}

				gtx := layout.NewContext(&ops, e)
				a.Layout(gtx)
				e.Frame(&ops)
			}
		}
	}

	return a.Accepted
}

func NewConfirmDialog(th *cu.Theme) *confirmDialog {
	return &confirmDialog{
		Theme:       th,
		Title:       "Confirm",
		AcceptLabel: "OK",
		CancelLabel: "Cancel",
		btnAccept:   &widget2.Clickable{},
		btnCancel:   &widget2.Clickable{},
	}
}
