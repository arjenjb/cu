module github.com/arjenjb/cu

go 1.24.0

toolchain go1.24.2

require (
	gioui.org v0.9.0
	gioui.org/x v0.9.0
	github.com/lmittmann/tint v1.1.2
	github.com/stretchr/testify v1.11.1
	golang.org/x/image v0.32.0
)

require (
	gioui.org/shader v1.0.8 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/exp/shiny v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	gioui.org v0.9.0 => github.com/ag5/gio v0.9.0-editor-patch
)