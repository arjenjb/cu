package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/arjenjb/cu"
	widget2 "github.com/arjenjb/cu/widget"
	"image/color"
	"log/slog"
	"time"
)

var checkable widget.Bool

func textInputExample(gtx layout.Context, th *cu.Theme, editor *widget.Editor) layout.Widget {
	t := material.NewTheme()
	t.Shaper = th.Shaper
	t.Face = th.Font.SansSerif.Typeface
	t.TextSize = th.TextSize
	t.Bg = color.NRGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 125,
	}

	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(widget2.TextInput(th, editor, "Hint", 200).Layout).
		Layout
}

func checkboxExample(gtx layout.Context, th *cu.Theme) layout.Widget {
	if checkable.Update(gtx) {
		slog.Info("Checkbox toggled", "checked", checkable.Value)
	}

	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(widget2.Checkbox(th, &checkable, "Check the box or this label").Layout).
		Layout
}

func main() {
	var ops op.Ops

	var progress float32 = 0.0
	var editor = new(widget.Editor)
	editor.Mask = '•'
	editor.SingleLine = true

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
			app.Size(670, 360),
		)

		th := cu.NewDefaultTheme()

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				for {
					_, ok := editor.Update(gtx)
					if !ok {
						break
					}
				}

				th.Background(gtx)

				th.M(cu.M, th.FlexColumn(cu.Gap(cu.M)).
					Rigid(buttonExample(th)).
					Rigid(th.Hr()).
					Rigid(textInputExample(gtx, th, editor)).
					Rigid(th.Hr()).
					Rigid(checkboxExample(gtx, th)).
					Rigid(th.Hr()).
					Rigid(spinnerExample(th)).
					Rigid(th.Hr()).
					Rigid(progressExample(th, progress)).
					Rigid(th.Hr()).
					Rigid(dialogExample(gtx, th)).
					Layout)(gtx)

				e.Frame(gtx.Ops)
			}

		}
	}()

	app.Main()
}
