package cu

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// These are all margin utilities

// M adds margin at each side of the given widget
func (t Theme) M(u Unit, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		s := u.Dp(t)
		return layout.UniformInset(s).Layout(gtx, w)
	}
}

// Mb adds margin at the bottom
func (t Theme) Mb(u Unit, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		s := u.Dp(t)
		return layout.Inset{Bottom: s}.Layout(gtx, w)
	}
}

// Mh adds horizontal margin, at the left and right side
func (t Theme) Mh(u Unit, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		s := u.Dp(t)
		return layout.Inset{Left: s, Right: s}.Layout(gtx, w)
	}
}

// Mv adds vertical margin, at the top and bottom side
func (t Theme) Mv(u Unit, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		s := u.Dp(t)
		return layout.Inset{Bottom: s, Top: s}.Layout(gtx, w)
	}
}

func VSpacer(m unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: m}.Layout(gtx)
	}
}

func HSpacer(m unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Width: m}.Layout(gtx)
	}
}

func Centered(l layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceSides,
		}.Layout(gtx, layout.Rigid(l))
	}
}
