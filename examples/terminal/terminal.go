package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget/terminal"
	"github.com/lmittmann/tint"

	"io"
	"log/slog"
	"math/rand"
	"os"
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
	w := &app.Window{}

	var guiReady = make(chan any)
	var ops op.Ops

	button := new(widget.Clickable)
	settings := terminal.NewConsoleSettings(terminal.MaxSize(100, 30))

	go func() {
		w.Option(app.Size(unit.Dp(670), unit.Dp(524)))
		guiReady <- struct{}{}

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				if button.Clicked(gtx) {
					w.Perform(system.ActionClose)
				}
				terminal.Console(th, l.screen, settings)(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()

	// Wait for the GUI to be ready
	<-guiReady

	for {
		select {
		case <-l.quitChannel:
			w.Perform(system.ActionClose)

		case <-l.updatedChannel:
			w.Invalidate()
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
			slog.Error(err.Error())
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

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "15:04:05",
			NoColor:    os.Getenv("NO_COLOR") == "1",
		}),
	))

	go func() {
		err := w.Open()
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	RESET := "\u001B[0m"
	BOLD := "\u001B[1m"
	FAINT := "\u001B[2m"

	go func() {
		//for i := 0; i < 200; i++ {
		//	fmt.Println(randomString(82))
		//}

		fmt.Println("ANSI Test")
		fmt.Println("=========")
		slog.Debug("This is not very important")
		slog.Info("Information message", "key", "value")
		slog.Warn("It's getting real")
		slog.Error("Oh no!")

		fmt.Println(BOLD + "This is bold" + RESET)
		fmt.Println(FAINT + "This is bold" + RESET)
		fmt.Println("\u001b[38;2;253;182;0mRgb code" + RESET)
		fmt.Println("\u001b[38;5;63m256 color code" + RESET)

		fmt.Println("")
		fmt.Println(randomString(200))
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
