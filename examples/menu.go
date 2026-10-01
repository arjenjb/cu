package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	widget2 "gioui.org/widget"
	"github.com/arjenjb/cu"
	"github.com/arjenjb/cu/widget"
)

var btnMenu = &widget2.Clickable{}

var dropdownMenu = &widget.Menu{
	Items: []widget.MenuItem{
		{Label: "New file", Shortcut: "⌘N"},
		{Label: "Open…", Shortcut: "⌘O"},
		{Separator: true},
		{Label: "Save", Shortcut: "⌘S"},
		{Label: "Save as…", Disabled: true},
	},
	Dropdown: true,
}

var contextMenu = &widget.Menu{
	Items: []widget.MenuItem{
		{Label: "Rename"},
		{Label: "Duplicate"},
		{Separator: true},
		{Label: "Delete"},
	},
}

var chosenMenuItem = "nothing yet"

func menuExample(gtx layout.Context, th cu.Theme) layout.Widget {
	for {
		i, ok := dropdownMenu.Update(gtx)
		if !ok {
			break
		}
		chosenMenuItem = dropdownMenu.Items[i].Label
	}
	for {
		i, ok := contextMenu.Update(gtx)
		if !ok {
			break
		}
		chosenMenuItem = contextMenu.Items[i].Label
	}

	return th.FlexRow(cu.Gap(cu.M), cu.Align(layout.Middle)).
		// The menu covers the widget that opens it
		Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(widget.Button(th, btnMenu, "File").Layout),
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return dropdownMenu.Layout(gtx, th)
				}),
			)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(rightClickArea(th)),
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return contextMenu.Layout(gtx, th)
				}),
			)
		}).
		Rigid(th.Text("Chosen: "+chosenMenuItem, cu.TextOptions{Color: &th.Color.TextSecondary})).
		Layout
}

// rightClickArea is a bordered box for opening the context menu.
func rightClickArea(th cu.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				rect := image.Rectangle{Max: gtx.Constraints.Min}
				rr := gtx.Dp(4)

				// The clip trims the 2dp border stroke to its inner 1dp.
				shape := clip.UniformRRect(rect, rr)
				defer shape.Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, th.Color.ControlBorder,
					clip.Stroke{Path: shape.Path(gtx.Ops), Width: float32(gtx.Dp(2))}.Op())
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 6, Bottom: 6, Left: 12, Right: 12}.Layout(gtx,
					th.Text("Right click here", cu.TextOptions{Color: &th.Color.TextSecondary}))
			}),
		)
	}
}
