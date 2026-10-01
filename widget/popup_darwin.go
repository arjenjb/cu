package widget

/*
#cgo CFLAGS: -fobjc-arc
#cgo LDFLAGS: -framework AppKit -framework CoreGraphics

#include <stdint.h>

void cu_popupAttach(uintptr_t view);
void cu_popupShow(uintptr_t view, uintptr_t id, double x, double y, double flipRight, double flipBottom, double width, double height, double radius);
void cu_popupClose(uintptr_t view);
void cu_popupCleanup(uintptr_t id);
void cu_mouseLocation(double *x, double *y);
*/
import "C"

import (
	"sync"
	"sync/atomic"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
)

const popupSupported = true

var (
	// Only one popup is open at a time, a click while one is open dismisses
	// it.
	popupActive atomic.Bool
	popupNextID atomic.Uintptr

	// The text shaper of the theme belongs to the window of the menu, and
	// shapers can't be shared between goroutines, so the popups get their own.
	popupShaper = sync.OnceValue(func() *text.Shaper { return text.NewShaper() })
)

func popupOpen() bool {
	return popupActive.Load()
}

// openPopup opens a menu in a window next to the mouse. The index of the
// chosen item is passed to act, from the goroutine of the popup.
func openPopup(th cu.Theme, req popupRequest, act func(int)) {
	if !popupActive.CompareAndSwap(false, true) {
		return
	}

	// Take the mouse position now, the popup window takes a moment to appear.
	var mx, my C.double
	C.cu_mouseLocation(&mx, &my)
	mouse := f32.Pt(float32(mx), float32(my))

	// Place a context menu at the mouse, and flip it to the other side of the
	// mouse when it doesn't fit on the screen.
	p := popupPlacement{pos: mouse, flip: mouse}

	if req.dropdown {
		// The area is found on the screen from the position of the click
		// within it. Screen coordinates are in points, Gio's in pixels.
		origin := mouse.Sub(req.pointer.Div(req.scale))
		size := f32.Pt(float32(req.area.X), float32(req.area.Y)).Div(req.scale)
		gap := float32(menuDropdownGap)

		// Open below the area, aligned to its left edge. Flip it to above the
		// area, or to its right edge, when it doesn't fit on the screen.
		p.pos = f32.Pt(origin.X, origin.Y+size.Y+gap)
		p.flip = f32.Pt(origin.X+size.X, origin.Y-gap)
	}

	th.Shaper = popupShaper()
	go runPopup(th, req, p, act)
}

// popupPlacement is the top left of the popup on the screen, and the right and
// bottom edges it is moved to when it doesn't fit, in points.
type popupPlacement struct {
	pos, flip f32.Point
}

func runPopup(th cu.Theme, req popupRequest, p popupPlacement, act func(int)) {
	id := popupNextID.Add(1)

	w := new(app.Window)
	w.Option(
		app.Size(req.width, req.height),
		app.MinSize(req.width, req.height),
		app.MaxSize(req.width, req.height),
		app.Decorated(false),
	)
	defer func() {
		C.cu_popupCleanup(C.uintptr_t(id))
		popupActive.Store(false)
	}()

	var (
		ops     op.Ops
		buttons = make([]widget.Clickable, len(req.items))
		view    uintptr
		shown   bool
	)

	for {
		switch e := w.Event().(type) {
		case app.AppKitViewEvent:
			// The window is about to be shown, make it a popup while it is
			// still invisible.
			if e.Valid() && view == 0 {
				view = e.View
				w.Run(func() { C.cu_popupAttach(C.uintptr_t(view)) })
			}

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			for i := range buttons {
				if buttons[i].Clicked(gtx) && req.items[i].selectable() {
					act(i)
					w.Run(func() { C.cu_popupClose(C.uintptr_t(view)) })
					break
				}
			}

			layoutMenu(gtx, th, req.items, buttons, req.width, false)
			e.Frame(gtx.Ops)

			// Place and show the window once its contents are drawn
			if !shown && view != 0 {
				shown = true
				w.Run(func() {
					C.cu_popupShow(C.uintptr_t(view), C.uintptr_t(id),
						C.double(p.pos.X), C.double(p.pos.Y), C.double(p.flip.X), C.double(p.flip.Y),
						C.double(req.width), C.double(req.height), C.double(menuRadius))
				})
			}

		case app.DestroyEvent:
			return
		}
	}
}
