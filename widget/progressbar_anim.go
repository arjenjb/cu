package widget

import "math"

// Cubic bezier easing functions
func easeInOutCubic1(t float64) float64 {
	// Approximation of cubic-bezier(0.65, 0.815, 0.735, 0.395)
	if t < 0.5 {
		return 4 * t * t * t
	} else {
		return 1.0 - math.Pow(-2*t+2, 3)/2.0
	}
}

type position struct {
	left, right float64
}

func calculateBar1Position(time int64) position {
	const cycleDuration = 1500
	var progress = float64(time%cycleDuration) / cycleDuration

	left := 0.0
	right := 0.0

	if progress <= 0.6 {
		var t = progress / 0.6
		var eased = easeInOutCubic1(t)
		right = eased
	} else {
		right = 1
	}

	if progress >= 0.15 {
		var t = (progress - 0.15) / 0.85
		var eased = easeInOutCubic1(t)
		left = eased
	}

	return position{
		left,
		right,
	}
}
