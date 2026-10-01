package widget

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
	"time"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/component"
	"github.com/arjenjb/cu"
)

// PopupMenus opens menus in a separate window, so they can extend beyond the
// bounds of the application window. This is experimental and only supported
// on macOS, other platforms always draw menus inside the window.
var PopupMenus = false

const (
	menuItemHeight      unit.Dp = 26
	menuItemRadius      unit.Dp = 5
	menuItemInset       unit.Dp = 10
	menuRadius          unit.Dp = 10
	menuPadding         unit.Dp = 5
	menuSeparatorHeight unit.Dp = 9
	menuShortcutGap     unit.Dp = 24

	// The space between a dropdown menu and the area that opens it
	menuDropdownGap unit.Dp = 4
)

var menuBorderColor = color.NRGBA{0xAF, 0xAF, 0xAF, 0xFF}

// MenuItem is an entry of a Menu.
type MenuItem struct {
	Label string
	// Shortcut is shown right aligned next to the label
	Shortcut string
	Disabled bool
	// Separator draws a line instead of an item
	Separator bool
}

func (it MenuItem) selectable() bool {
	return !it.Disabled && !it.Separator
}

// Menu is a menu that opens from an area. Its Layout should cover the widget
// that opens it, for example in a layout.Expanded on top of it.
type Menu struct {
	Items []MenuItem
	// Dropdown opens the menu below the area on a primary click. Otherwise it
	// is a context menu that opens at the pointer on a secondary click.
	Dropdown bool
	// MinWidth is the minimum width of the menu, it is otherwise as wide as
	// its items.
	MinWidth unit.Dp

	area    component.ContextArea
	buttons []widget.Clickable
	opened  bool

	// Items chosen in a popup window, which runs on its own goroutine
	pendingMu sync.Mutex
	pending   []int
}

// Update returns the index of the item that was chosen, if any.
func (m *Menu) Update(gtx layout.Context) (int, bool) {
	m.pendingMu.Lock()
	if len(m.pending) > 0 {
		i := m.pending[0]
		m.pending = m.pending[1:]
		m.pendingMu.Unlock()
		return i, true
	}
	m.pendingMu.Unlock()

	for i := range m.buttons {
		if m.buttons[i].Clicked(gtx) && i < len(m.Items) && m.Items[i].selectable() {
			m.area.Dismiss()
			return i, true
		}
	}
	return 0, false
}

// Opened reports whether the menu was opened since the last call to Opened.
func (m *Menu) Opened() bool {
	opened := m.opened
	m.opened = false
	return opened
}

// Active reports whether the menu is open inside the window. A menu in a
// popup window closes by itself.
func (m *Menu) Active() bool {
	return m.area.Active()
}

// Dismiss closes the menu if it is open inside the window.
func (m *Menu) Dismiss() {
	m.area.Dismiss()
}

func (m *Menu) Layout(gtx layout.Context, th cu.Theme) layout.Dimensions {
	if len(m.buttons) != len(m.Items) {
		m.buttons = make([]widget.Clickable, len(m.Items))
	}

	if PopupMenus && popupSupported {
		return m.layoutPopup(gtx, th)
	}

	if m.Dropdown {
		m.area.Activation = pointer.ButtonPrimary
		m.area.AbsolutePosition = true
	}
	m.area.Update(gtx)
	if m.area.Activated() {
		m.opened = true
	}

	areaHeight := gtx.Constraints.Min.Y
	return m.area.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = image.Point{}
		if m.Dropdown {
			defer op.Offset(image.Pt(0, areaHeight+gtx.Dp(menuDropdownGap))).Push(gtx.Ops).Pop()
		}
		return layoutMenu(gtx, th, m.Items, m.buttons, m.width(gtx, th), true)
	})
}

// layoutPopup opens the menu in a separate window.
func (m *Menu) layoutPopup(gtx layout.Context, th cu.Theme) layout.Dimensions {
	button := ifElse(m.Dropdown, pointer.ButtonPrimary, pointer.ButtonSecondary)

	for {
		ev, ok := gtx.Event(pointer.Filter{Target: m, Kinds: pointer.Press})
		if !ok {
			break
		}
		if e, ok := ev.(pointer.Event); ok && e.Buttons.Contain(button) {
			m.opened = true
			openPopup(th, popupRequest{
				items:    append([]MenuItem(nil), m.Items...),
				width:    m.width(gtx, th),
				height:   menuHeight(m.Items),
				dropdown: m.Dropdown,
				pointer:  e.Position,
				area:     gtx.Constraints.Min,
				scale:    gtx.Metric.PxPerDp,
			}, func(i int) {
				m.pendingMu.Lock()
				m.pending = append(m.pending, i)
				m.pendingMu.Unlock()
			})
		}
	}

	// The popup runs on its own goroutine and can't redraw this window, so
	// poll for the chosen item while it is open.
	if popupOpen() {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 30)})
	}

	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, m)

	return layout.Dimensions{Size: gtx.Constraints.Min}
}

// popupRequest describes a menu to open in a popup window.
type popupRequest struct {
	items         []MenuItem
	width, height unit.Dp
	dropdown      bool

	// The position of the click that opens the menu within the area, and the
	// size of the area, in pixels.
	pointer f32.Point
	area    image.Point
	scale   float32
}

// width measures the width of the menu from its items.
func (m *Menu) width(gtx layout.Context, th cu.Theme) unit.Dp {
	macro := op.Record(gtx.Ops)
	defer macro.Stop()

	gtx.Constraints = layout.Constraints{Max: image.Pt(math.MaxInt32, math.MaxInt32)}
	content := 0
	for _, it := range m.Items {
		if it.Separator {
			continue
		}
		w := th.Text(it.Label)(gtx).Size.X
		if it.Shortcut != "" {
			w += gtx.Dp(menuShortcutGap) + th.Text(it.Shortcut)(gtx).Size.X
		}
		content = max(content, w)
	}

	width := unit.Dp(math.Ceil(float64(float32(content)/gtx.Metric.PxPerDp))) + 2*(menuItemInset+menuPadding)
	return max(width, m.MinWidth)
}

func menuHeight(items []MenuItem) unit.Dp {
	h := 2 * menuPadding
	for _, it := range items {
		h += ifElse(it.Separator, menuSeparatorHeight, menuItemHeight)
	}
	return h
}

// layoutMenu draws the menu panel. The shadow and border are left out when the
// menu is shown in its own window, which gets them from the platform.
func layoutMenu(gtx layout.Context, th cu.Theme, items []MenuItem, buttons []widget.Clickable, width unit.Dp, frame bool) layout.Dimensions {
	gtx.Constraints = layout.Exact(image.Point{X: gtx.Dp(width), Y: gtx.Dp(menuHeight(items))})

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rect := image.Rectangle{Max: gtx.Constraints.Min}
			rr := gtx.Dp(menuRadius)

			if frame {
				menuShadow(gtx, rect, rr)
			}

			// The clip trims the 2dp border stroke to its inner 1dp.
			shape := clip.UniformRRect(rect, rr)
			defer shape.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, colorNormal)
			if frame {
				paint.FillShape(gtx.Ops, menuBorderColor,
					clip.Stroke{
						Path:  shape.Path(gtx.Ops),
						Width: float32(gtx.Dp(2)),
					}.Op(),
				)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			col := th.FlexColumn()
			for i, it := range items {
				if it.Separator {
					col = col.Rigid(menuSeparator(th))
				} else {
					col = col.Rigid(menuItem(th, &buttons[i], it))
				}
			}
			return layout.UniformInset(menuPadding).Layout(gtx, col.Layout)
		}),
	)
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

func menuSeparator(th cu.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := image.Pt(gtx.Constraints.Max.X, gtx.Dp(menuSeparatorHeight))
		inset := gtx.Dp(menuItemInset)
		line := image.Rect(inset, size.Y/2, size.X-inset, size.Y/2+gtx.Dp(1))
		paint.FillShape(gtx.Ops, th.Color.ControlBorder, clip.Rect(line).Op())
		return layout.Dimensions{Size: size}
	}
}

func menuItem(th cu.Theme, button *widget.Clickable, it MenuItem) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		gtx.Constraints.Min.Y = gtx.Dp(menuItemHeight)
		gtx.Constraints.Max.Y = gtx.Constraints.Min.Y

		enabled := !it.Disabled

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
				Flexed(1, th.Text(it.Label, cu.TextOptions{Color: &labelColor}))
			if it.Shortcut != "" {
				row = row.Rigid(th.Text(it.Shortcut, cu.TextOptions{Color: &shortcutColor}))
			}

			// layout.W and the flex both keep the minimum height, which would
			// stretch the texts and draw them at the top, so reset it to center
			// the row vertically.
			return layout.Inset{Left: menuItemInset, Right: menuItemInset}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = image.Pt(gtx.Constraints.Max.X, 0)
					return row.Layout(gtx)
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
