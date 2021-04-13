package core

import (
	"github.com/enr/clui"
)

var (
	ui *clui.Clui
)

// ConfigureAppUI helps to use in library the UI created from the command.
func ConfigureAppUI(options ...func(*clui.Clui)) (*clui.Clui, error) {
	nui, err := clui.NewClui(options...)
	ui = nui
	return ui, err
}
