package dialog

import (
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	widget2 "gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

type inputDialog struct {
	Theme *cu.Theme
	// Title is shown in the titlebar, if left empty it will default to Message
	Title string
	// BigMessage is shown in H2 style at the top of window
	BigMessage string
	// NormalMessage is an optional additional message shown in normal font
	NormalMessage string

	Hint string

	AcceptLabel string
	CancelLabel string

	btnAccept *widget2.Clickable
	btnCancel *widget2.Clickable
	Accepted  bool

	Editor *widget2.Editor
}

func (a inputDialog) Layout(gtx layout.Context) layout.Dimensions {
	th := a.Theme
	th.Background(gtx)

	return th.M(cu.M,
		th.FlexColumn(cu.Gap(cu.S)).
			RigidIf(len(a.BigMessage) > 0, th.H2(a.BigMessage)).
			Rigid(th.Text(a.NormalMessage)).
			Rigid(cu.VSpacer(12)).
			Rigid(
				th.FlexRow(cu.Align(layout.Middle)).
					Flexed(1, cu.HSpacer(0)).
					Rigid(
						widget.TextInput(th, a.Editor, a.Hint, 250).Layout,
					).
					Flexed(1, cu.HSpacer(0)).
					Layout,
			).
			Rigid(cu.VSpacer(12)).
			Flexed(1, th.FlexRow(cu.Align(layout.End)).
				Flexed(1, cu.HSpacer(0)).
				Rigid(widget.Button(th, a.btnCancel, a.CancelLabel).Layout).
				Rigid(widget.Button(th, a.btnAccept, a.AcceptLabel, widget.Primary()).Layout).
				Layout).
			Layout)(gtx)
}

func (a inputDialog) Show() (bool, string) {
	w := app.Window{}
	w.Option(
		app.Title(a.Title),
		app.Size(400, 140),
	)

	var ops op.Ops

	var done = false
	var first = true

	for !done {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			done = true

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if a.btnCancel.Clicked(gtx) {
				w.Perform(system.ActionClose)
			}
			if a.btnAccept.Clicked(gtx) {
				a.Accepted = true
				w.Perform(system.ActionClose)
			}

			if first {
				gtx.Execute(key.FocusCmd{Tag: a.Editor})
				first = false
			}

			a.Layout(gtx)
			e.Frame(&ops)
		}
	}

	return a.Accepted, a.Editor.Text()
}

func NewInputDialog(th *cu.Theme) *inputDialog {
	editor := new(widget2.Editor)
	editor.SingleLine = true

	return &inputDialog{
		Theme:       th,
		Title:       "Confirm",
		AcceptLabel: "OK",
		CancelLabel: "Cancel",
		btnAccept:   &widget2.Clickable{},
		btnCancel:   &widget2.Clickable{},
		Editor:      editor,
	}
}
