package widget

import (
	"image"
	"image/color"
	"io"
	"runtime"
	"strings"
	"sync"
	"weak"

	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/component"
	"github.com/arjenjb/cu"
)

const (
	menuWidth      unit.Dp = 200
	menuItemHeight unit.Dp = 26
	menuItemRadius unit.Dp = 5
	menuRadius     unit.Dp = 10
	menuPadding    unit.Dp = 5
)

// inputMenu is the right click menu of a text input. TextInputWidget is
// rebuilt every frame, so the menu state is kept per editor in inputMenus.
type inputMenu struct {
	area  component.ContextArea
	cut   widget.Clickable
	copy  widget.Clickable
	paste widget.Clickable
}

var (
	inputMenusMu sync.Mutex
	inputMenus   = map[weak.Pointer[widget.Editor]]*inputMenu{}
)

// inputMenuFor returns the menu state of the editor. The state is dropped
// once the editor is garbage collected.
func inputMenuFor(editor *widget.Editor) *inputMenu {
	key := weak.Make(editor)

	inputMenusMu.Lock()
	defer inputMenusMu.Unlock()

	if m, ok := inputMenus[key]; ok {
		return m
	}

	m := &inputMenu{}
	inputMenus[key] = m
	runtime.AddCleanup(editor, func(key weak.Pointer[widget.Editor]) {
		inputMenusMu.Lock()
		delete(inputMenus, key)
		inputMenusMu.Unlock()
	}, key)
	return m
}

// Layout lays out the area that listens for right clicks, it should cover the
// text input.
func (m *inputMenu) Layout(gtx layout.Context, th cu.Theme, editor *widget.Editor) layout.Dimensions {
	m.area.Update(gtx)
	if m.area.Activated() {
		gtx.Execute(key.FocusCmd{Tag: editor})
	}

	// Only listen for escape while the menu is open, so it still reaches the
	// rest of the application otherwise.
	if m.area.Active() {
		for {
			ev, ok := gtx.Event(key.Filter{Focus: editor, Name: key.NameEscape})
			if !ok {
				break
			}
			if e, ok := ev.(key.Event); ok && e.State == key.Press {
				m.area.Dismiss()
			}
		}
	}

	return m.area.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		hasSelection := editor.SelectedText() != ""

		if m.cut.Clicked(gtx) && hasSelection && !editor.ReadOnly {
			writeClipboard(gtx, editor.SelectedText())
			editor.Delete(1)
		}
		if m.copy.Clicked(gtx) && hasSelection {
			writeClipboard(gtx, editor.SelectedText())
		}
		if m.paste.Clicked(gtx) && !editor.ReadOnly {
			gtx.Execute(clipboard.ReadCmd{Tag: editor})
		}

		gtx.Constraints = layout.Exact(image.Point{
			X: gtx.Dp(menuWidth),
			Y: gtx.Dp(3*menuItemHeight + 2*menuPadding),
		})
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				rect := image.Rectangle{Max: gtx.Constraints.Min}
				rr := gtx.Dp(menuRadius)

				menuShadow(gtx, rect, rr)

				// The clip trims the 2dp border stroke to its inner 1dp.
				shape := clip.UniformRRect(rect, rr)
				defer shape.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, colorNormal)
				paint.FillShape(gtx.Ops, th.Color.ControlBorder,
					clip.Stroke{
						Path:  shape.Path(gtx.Ops),
						Width: float32(gtx.Dp(2)),
					}.Op(),
				)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(menuPadding).Layout(gtx,
					th.FlexColumn().
						Rigid(menuItem(th, &m.cut, "Cut", shortcut("X"), hasSelection && !editor.ReadOnly)).
						Rigid(menuItem(th, &m.copy, "Copy", shortcut("C"), hasSelection)).
						Rigid(menuItem(th, &m.paste, "Paste", shortcut("V"), !editor.ReadOnly)).
						Layout)
			}),
		)
	})
}

// menuShadow draws a soft drop shadow below rect by stacking faint rounded
// rectangles that grow outwards.
func menuShadow(gtx layout.Context, rect image.Rectangle, rr int) {
	const layers = 12
	offset := image.Pt(0, gtx.Dp(3))

	for i := layers; i > 0; i-- {
		spread := gtx.Dp(unit.Dp(i))
		r := rect.Add(offset).Inset(-spread)
		paint.FillShape(gtx.Ops, color.NRGBA{A: 4},
			clip.UniformRRect(r, rr+spread).Op(gtx.Ops))
	}
}

// shortcut formats the platform's shortcut for the given key.
func shortcut(key string) string {
	if runtime.GOOS == "darwin" {
		return "⌘" + key
	}
	return "Ctrl+" + key
}

func menuItem(th cu.Theme, button *widget.Clickable, label, shortcut string, enabled bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		gtx.Constraints.Min.Y = gtx.Dp(menuItemHeight)
		gtx.Constraints.Max.Y = gtx.Constraints.Min.Y

		item := func(gtx layout.Context) layout.Dimensions {
			if enabled && button.Hovered() {
				shape := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(menuItemRadius))
				paint.FillShape(gtx.Ops, th.Color.SelectionActive, shape.Op(gtx.Ops))
			}

			labelColor := ifElse(enabled, th.Color.Text, th.Color.TextDisabled)
			shortcutColor := ifElse(enabled, th.Color.TextSecondary, th.Color.TextDisabled)

			// Align on the baseline, the shortcut symbols can come from a
			// fallback font with different line metrics than the label.
			row := th.FlexRow(cu.Align(layout.Baseline)).
				Flexed(1, th.Text(label, cu.TextOptions{Color: &labelColor})).
				Rigid(th.Text(shortcut, cu.TextOptions{Color: &shortcutColor})).
				Layout

			// layout.W and the flex both keep the minimum height, which would
			// stretch the texts and draw them at the top, so reset it to center
			// the row vertically.
			return layout.Inset{Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = image.Pt(gtx.Constraints.Max.X, 0)
					return row(gtx)
				})
			})
		}

		if !enabled {
			return item(gtx)
		}
		return button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			pointer.CursorPointer.Add(gtx.Ops)
			return item(gtx)
		})
	}
}

func writeClipboard(gtx layout.Context, text string) {
	gtx.Execute(clipboard.WriteCmd{
		Type: "application/text",
		Data: io.NopCloser(strings.NewReader(text)),
	})
}
