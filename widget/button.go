package widget

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	. "github.com/arjenjb/cu"
	"image/color"
)

var buttonTextColorPrimary = color.NRGBA{
	R: 255,
	G: 255,
	B: 255,
	A: 255,
}
var buttonTextColor = color.NRGBA{
	A: 255,
}
var buttonTextColorDisabled = color.NRGBA{
	R: 128,
	G: 128,
	B: 128,
	A: 255,
}

type ButtonWidget struct {
	label     string
	clickable *widget.Clickable
	theme     *Theme
	options   buttonOptions
}

var tag struct{}

func (b *ButtonWidget) Layout(gtx layout.Context) layout.Dimensions {

	fontSize := b.theme.TextSize
	if b.options.big {
		fontSize = 15
	}

	gtx.Constraints.Min.X = gtx.Dp(80)
	gtx.Constraints.Min.Y = gtx.Dp(32)

	var color = buttonTextColor
	if b.options.disabled {
		color = buttonTextColorDisabled
	} else if b.options.primary {
		color = buttonTextColorPrimary
	}

	buttonText := func(gtx layout.Context) layout.Dimensions {
		return b.theme.Mh(Scaled(1), func(gtx layout.Context) layout.Dimensions {
			colMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: color}.Add(gtx.Ops)
			return widget.Label{Alignment: text.Middle}.Layout(gtx, b.theme.Shaper, b.theme.Font.SansSerif, fontSize, b.label, colMacro.Stop())

		})(gtx)
	}

	return ButtonStyle{
		CornerRadius: 4,
		Button:       b.clickable,
		Disabled:     b.options.disabled,
		Primary:      b.options.primary,
	}.Layout(gtx, buttonText)
}

type ButtonOptions func(options *buttonOptions)

func Big() func(*buttonOptions) {
	return func(o *buttonOptions) {
		o.big = true
	}
}

func Disabled(b bool) func(*buttonOptions) {
	return func(o *buttonOptions) {
		o.disabled = b
	}
}

func Primary() func(*buttonOptions) {
	return func(o *buttonOptions) {
		o.primary = true
	}
}

type buttonOptions struct {
	primary  bool
	big      bool
	disabled bool
}

func Button(t *Theme, clickable *widget.Clickable, label string, options ...ButtonOptions) *ButtonWidget {
	o := buttonOptions{
		primary:  false,
		big:      false,
		disabled: false,
	}

	for _, each := range options {
		each(&o)
	}

	return &ButtonWidget{
		theme:     t,
		label:     label,
		clickable: clickable,
		options:   o,
	}
}
