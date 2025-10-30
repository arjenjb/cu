package cu

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
)

type Flex struct {
	widget   layout.Flex
	gap      unit.Dp
	children []layout.FlexChild
}

func (f Flex) Layout(gtx layout.Context) layout.Dimensions {
	children := f.children
	if f.gap > 0 {
		var spacer layout.Widget
		if f.widget.Axis == layout.Horizontal {
			spacer = func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{
					Size: image.Point{X: gtx.Dp(f.gap)},
				}
			}
		} else {
			spacer = func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{
					Size: image.Point{Y: gtx.Dp(f.gap)},
				}
			}
		}

		var n []layout.FlexChild
		for i, each := range children {
			if i > 0 {
				n = append(n, layout.Rigid(spacer))
			}
			n = append(n, each)
		}
		children = n
	}

	return f.widget.Layout(gtx, children...)
}

func (f Flex) Rigid(w layout.Widget) Flex {
	f.children = append(f.children, layout.Rigid(w))
	return f
}

func (f Flex) Flexed(weight float32, w layout.Widget) Flex {
	f.children = append(f.children, layout.Flexed(weight, w))
	return f
}

func (f Flex) RigidIf(condition bool, w layout.Widget) Flex {
	if condition {
		return f.Rigid(w)
	}
	return f
}

func (f Flex) FlexedIf(condition bool, weight float32, w layout.Widget) Flex {
	if condition {
		return f.Flexed(weight, w)
	}
	return f
}

type CuFlexOption func(w *Flex, t Theme)

func Spacing(s layout.Spacing) func(w *Flex, t Theme) {
	return func(w *Flex, t Theme) {
		w.widget.Spacing = s
	}
}

func Align(a layout.Alignment) func(w *Flex, t Theme) {
	return func(w *Flex, t Theme) {
		w.widget.Alignment = a
	}
}

func Gap(s Unit) func(w *Flex, t Theme) {
	return func(w *Flex, t Theme) {
		w.gap = s.Dp(t)
	}
}
