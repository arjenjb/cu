package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget/terminal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"math/rand"
	"os"
	"time"
)

type TerminalWindow struct {
	Screen         *os.File
	screen         *terminal.Screen
	quitChannel    chan interface{}
	updatedChannel chan interface{}
}

func (l TerminalWindow) Close() {
	l.quitChannel <- struct{}{}
}

func (l TerminalWindow) Open() error {
	th := cu.NewDefaultTheme()
	w := app.NewWindow(app.Size(unit.Dp(670), unit.Dp(524)))

	var ops op.Ops

	button := new(widget.Clickable)
	settings := terminal.NewConsoleSettings(terminal.MaxSize(100, 30))

	for {
		select {
		case <-l.quitChannel:
			w.Perform(system.ActionClose)

		case <-l.updatedChannel:
			w.Invalidate()

		case evt := <-w.Events():
			switch e := evt.(type) {
			case system.StageEvent:
				//alwaysOnTop()

			case system.DestroyEvent:
				return nil

			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				if button.Clicked() {
					w.Perform(system.ActionClose)
				}

				terminal.Console(th, l.screen, settings)(gtx)

				e.Frame(gtx.Ops)
			}
		}
	}
}

func NewTerminalWindow(size terminal.Point) *TerminalWindow {
	updatedChannel := make(chan interface{})
	screen := terminal.NewScreen(size, updatedChannel)

	r, w, _ := os.Pipe()

	go func() {
		_, err := io.Copy(screen, r)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	return &TerminalWindow{
		Screen:         w,
		screen:         screen,
		updatedChannel: updatedChannel,
		quitChannel:    make(chan interface{}),
	}
}

func main() {
	w := NewTerminalWindow(terminal.Point{
		X: 80,
		Y: 20,
	})
	os.Stdout = w.Screen
	os.Stderr = w.Screen

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: w.Screen, TimeFormat: time.TimeOnly})

	go func() {
		err := w.Open()
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	go func() {
		for i := 0; i < 200; i++ {
			fmt.Println(randomString(82))
		}
	}()

	print("Starting main")

	app.Main()
}

func randomString(n int) string {
	s := make([]rune, n)
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < n; i++ {
		s[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(s)
}
