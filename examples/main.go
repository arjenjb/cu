package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/arjenjb/cu"
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

					th.M(cu.M, th.FlexColumn(cu.Gap(cu.M)).
						Rigid(buttonExample(th)).
						Rigid(th.Hr()).
						Rigid(spinnerExample(th)).
						Rigid(th.Hr()).
						Rigid(progressExample(th)).
						Layout)(gtx)

					e.Frame(gtx.Ops)
				}
			}
		}
	}()

	app.Main()
}
