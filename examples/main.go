package main

import (
	"log/slog"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
	widget2 "github.com/arjenjb/cu/widget"
)

var checkable widget.Bool

func textInputExample(gtx layout.Context, th cu.Theme, editor *widget.Editor) layout.Widget {
	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(widget2.TextInput(th, editor, "Hint", 0, 0).Layout).
		Layout
}

func textAreaExample(gtx layout.Context, th cu.Theme, editor *widget.Editor, btn *widget.Clickable) layout.Widget {
	if btn.Clicked(gtx) {
		editor.ScrollToEnd()
	}

	for {
		_, ok := editor.Update(gtx)
		if !ok {
			break
		}

	}

	return th.FlexColumn(cu.Gap(cu.XS)).
		Flexed(1.0, widget2.TextInput(th, editor, "Hint", 0, 0).Layout).
		Rigid(widget2.Button(th, btn, "To end").Layout).
		Layout
}

func checkboxExample(gtx layout.Context, th cu.Theme) layout.Widget {
	if checkable.Update(gtx) {
		slog.Info("Checkbox toggled", "checked", checkable.Value)
	}

	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(widget2.Checkbox(th, &checkable, "Check the box or this label").Layout).
		Layout
}

var longText = `Mayor Goldie Wilson, I like the sound of that. Uh? Of course, the Enchantment Under The Sea Dance they're supposed to go to this, that's where they kiss for the first time. Damn, where is that kid. Damn. Damn damn. You're late, do you have no concept of time? No no no this sucker's electrical, but I need a nuclear reaction to generate the one point twenty-one gigawatts of electricity-

George. Lynda, first of all, I'm not your answering service. Second of all, somebody named Greg or Craig called you just a little while ago. No wait, Doc, the bruise, the bruise on your head, I know how that happened, you told me the whole story. you were standing on your toilet and you were hanging a clock, and you fell, and you hit your head on the sink, and that's when you came up with the idea for the flux capacitor, which makes time travel possible. Radiation suit, of course, cause all of the fall out from the atomic wars. This is truly amazing, a portable television studio. No wonder your president has to be an actor, he's gotta look good on television. Who's are these?

Biff. Look at the time, you've got less than 4 minutes, please hurry. Now, now, Biff, now, I never noticed any blindspot before when I would drive it. Hi, son. I still don't understand, how am I supposed to go to the dance with her, if she's already going to the dance with you. What about George?`

func main() {
	var ops op.Ops

	var progress float32 = 0.0

	var editInput = new(widget.Editor)
	editInput.Mask = '•'
	editInput.SingleLine = true

	var editArea = new(widget.Editor)
	editArea.SingleLine = false
	editArea.SetText(longText)
	editArea.ReadOnly = true

	var btnScrollDown = new(widget.Clickable)

	var indeterminedProgress = widget2.NewIndeterminedProgress()

	var w app.Window

	// Continuously update the progress bar
	go func() {
		for {
			for i := 0; i <= 1000; i++ {
				progress = 0.001 * float32(i)
				time.Sleep(10 * time.Millisecond)
				w.Invalidate()
			}
		}
	}()

	go func() {
		w.Option(
			app.Size(670, 560),
		)

		th := cu.NewDefaultTheme()

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				for {
					_, ok := editInput.Update(gtx)
					if !ok {
						break
					}
				}

				th.Background(gtx)

				th.M(cu.M, th.FlexColumn(cu.Gap(cu.M)).
					Rigid(buttonExample(th)).
					Rigid(th.Hr()).
					Rigid(textInputExample(gtx, th, editInput)).
					Rigid(th.Hr()).
					Rigid(checkboxExample(gtx, th)).
					Rigid(th.Hr()).
					Rigid(spinnerExample(th)).
					Rigid(th.Hr()).
					Rigid(determinateProgressExample(th, widget2.LinearProgress(progress))).
					Rigid(indeterminateProgressExample(th, indeterminedProgress)).
					Rigid(th.Hr()).
					Rigid(dialogExample(gtx, th)).
					Rigid(th.Hr()).
					Flexed(1.0, textAreaExample(gtx, th, editArea, btnScrollDown)).
					Layout)(gtx)

				e.Frame(gtx.Ops)
			}

		}
	}()

	app.Main()
}
