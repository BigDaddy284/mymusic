package settings

import (
	"slices"
)

// Addons are the plugins that run.
type Addons struct {
	// EnabledPlugins are the IDs of the plugins that run.
	EnabledPlugins []string `json:"enabledPlugins"`
}

// defaultAddons are the add-ons of a new install: no plugin runs.
func defaultAddons() Addons {
	var a Addons
	a.EnabledPlugins = []string{}
	return a
}

// clean is a with each plugin once.
func (a Addons) clean() Addons {
	plugins := []string{}
	for _, id := range a.EnabledPlugins {
		if ValidID(id) && !slices.Contains(plugins, id) {
			plugins = append(plugins, id)
		}
	}
	a.EnabledPlugins = plugins
	return a
}
