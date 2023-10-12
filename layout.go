package cu

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func (t Theme) Mb(h float32, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(8 * h)}.Layout(gtx, w)
	}
}

func (t Theme) Mv(h float32, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		s := unit.Dp(8 * h)
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
