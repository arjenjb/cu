package widget

import (
	"io"
	"runtime"
	"strings"
	"sync"
	"weak"

	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/arjenjb/cu"
)

// The items of the text input menu
const (
	itemCut = iota
	itemCopy
	itemPaste
)

// inputMenu is the right click menu of a text input. TextInputWidget is
// rebuilt every frame, so the menu state is kept per editor in inputMenus.
type inputMenu struct {
	menu Menu
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

	m := &inputMenu{menu: Menu{MinWidth: 200}}
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
	hasSelection := editor.SelectedText() != ""
	m.menu.Items = []MenuItem{
		itemCut:   {Label: "Cut", Shortcut: shortcut("X"), Disabled: !hasSelection || editor.ReadOnly},
		itemCopy:  {Label: "Copy", Shortcut: shortcut("C"), Disabled: !hasSelection},
		itemPaste: {Label: "Paste", Shortcut: shortcut("V"), Disabled: editor.ReadOnly},
	}

	for {
		i, ok := m.menu.Update(gtx)
		if !ok {
			break
		}
		apply(gtx, editor, i)
	}

	// Only listen for escape while the menu is open, so it still reaches the
	// rest of the application otherwise.
	if m.menu.Active() {
		for {
			ev, ok := gtx.Event(key.Filter{Focus: editor, Name: key.NameEscape})
			if !ok {
				break
			}
			if e, ok := ev.(key.Event); ok && e.State == key.Press {
				m.menu.Dismiss()
			}
		}
	}

	dims := m.menu.Layout(gtx, th)
	if m.menu.Opened() {
		gtx.Execute(key.FocusCmd{Tag: editor})
	}
	return dims
}

func apply(gtx layout.Context, editor *widget.Editor, item int) {
	switch item {
	case itemCut:
		if text := editor.SelectedText(); text != "" && !editor.ReadOnly {
			writeClipboard(gtx, text)
			editor.Delete(1)
		}
	case itemCopy:
		if text := editor.SelectedText(); text != "" {
			writeClipboard(gtx, text)
		}
	case itemPaste:
		if !editor.ReadOnly {
			gtx.Execute(clipboard.ReadCmd{Tag: editor})
		}
	}
}

func writeClipboard(gtx layout.Context, text string) {
	gtx.Execute(clipboard.WriteCmd{
		Type: "application/text",
		Data: io.NopCloser(strings.NewReader(text)),
	})
}
