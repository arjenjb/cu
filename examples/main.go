package main

import (
	"cu"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {
	w := app.NewWindow(app.Size(unit.Dp(670), unit.Dp(524)))
	th := cu.NewDefaultTheme()

	var ops op.Ops

	go func() {
		for {
			select {
			case evt := <-w.Events():
				switch e := evt.(type) {
				case system.DestroyEvent:
					return

				case system.FrameEvent:
					gtx := layout.NewContext(&ops, e)

					th.Background(gtx)

					th.FlexColumn(cu.Gap(cu.M)).
						Rigid(buttonExample(th)).
						Rigid(spinnerExample(th)).
						Layout(gtx)

					e.Frame(gtx.Ops)
				}
			}
		}
	}()

	app.Main()
}
