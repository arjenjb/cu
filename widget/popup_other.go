//go:build !darwin

package widget

import "github.com/arjenjb/cu"

const popupSupported = false

func openPopup(th cu.Theme, req popupRequest, act func(int)) {}

func popupOpen() bool { return false }
