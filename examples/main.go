package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
	widget2 "github.com/arjenjb/cu/widget"
	"github.com/rs/zerolog/log"
	"time"
)

var checkable widget.Bool

func checkboxExample(th *cu.Theme) layout.Widget {
	if checkable.Changed() {
		if checkable.Value {
			log.Info().Str("checked", "yes").Msg("Checkbox toggled")
		} else {
			log.Info().Str("checked", "no").Msg("Checkbox toggled")
		}
	}

	return th.FlexRow(cu.Gap(cu.XS)).
		Rigid(widget2.Checkbox(th, &checkable, "Check the box or this label").Layout).
		Layout
}

func main() {
	w := app.NewWindow(app.Size(unit.Dp(670), unit.Dp(360)))
	th := cu.NewDefaultTheme()

	var ops op.Ops

	var progress float32 = 0.0

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
						Rigid(checkboxExample(th)).
						Rigid(th.Hr()).
						Rigid(spinnerExample(th)).
						Rigid(th.Hr()).
						Rigid(progressExample(th, progress)).
						Rigid(th.Hr()).
						Rigid(dialogExample(th)).
						Layout)(gtx)

					e.Frame(gtx.Ops)
				}
			}
		}
	}()

	app.Main()
}
