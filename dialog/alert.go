package dialog

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	widget2 "gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

type AlertDialog struct {
	Theme cu.Theme
	// Title is shown in the titlebar, if left empty it will default to Message
	Title string
	// BigMessage is shown in H2 style at the top of window
	BigMessage string
	// NormalMessage is an optional additional message shown in normal font
	NormalMessage string

	button *widget2.Clickable
}

func (a AlertDialog) Layout(gtx layout.Context) layout.Dimensions {
	th := a.Theme
	th.Background(gtx)

	return th.M(cu.M,
		th.FlexColumn(cu.Gap(cu.S)).
			RigidIf(len(a.BigMessage) > 0, th.H2(a.BigMessage)).
			Rigid(th.Text(a.NormalMessage)).
			Flexed(1, th.FlexRow(cu.Align(layout.End)).Flexed(1, cu.HSpacer(0)).Rigid(widget.Button(th, a.button, "Ok", widget.Primary()).Layout).Layout).
			Layout)(gtx)
}

func (a AlertDialog) Show() {
	w := app.Window{}
	w.Option(
		app.Title(a.Title),
		app.Size(400, 140),
	)

	var ops op.Ops

	var done = false
	for !done {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			done = true

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if a.button.Clicked(gtx) {
				w.Perform(system.ActionClose)
			}

			a.Layout(gtx)
			e.Frame(&ops)
		}
	}
}

func NewAlertDialog(th cu.Theme) *AlertDialog {
	return &AlertDialog{
		Theme:  th,
		Title:  "Message",
		button: &widget2.Clickable{},
	}
}
