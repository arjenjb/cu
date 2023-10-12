package widget

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

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
