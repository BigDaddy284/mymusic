package app

import (
	"github.com/egoist/mygo"

	"github.com/bigdaddy284/mymusic/internal/addons"
)

// Addons lists and removes themes.
type Addons struct{ a *App }

// List returns the themes available, and those that could not
// load.
func (s *Addons) List() addons.Inventory { return s.a.addons.List() }

// Reload finds the add-ons again, after the user changed the add-ons
// folder, reads the colour file again, and tells every window.
func (s *Addons) Reload() addons.Inventory {
	inv := s.a.addons.List()
	s.a.rereadColorFile()
	notify()
	return inv
}

// Uninstall removes an add-on the user installed.
func (s *Addons) Uninstall(kind addons.Kind, id string) (addons.Inventory, error) {
	if err := s.a.addons.Uninstall(kind, id); err != nil {
		return addons.Inventory{}, err
	}
	return s.Reload(), nil
}

// OpenFolder opens the folder of the user's add-ons of a kind, where they
// can put add-ons they make.
func (s *Addons) OpenFolder(kind addons.Kind) error {
	dir, err := s.a.addons.FolderOf(kind)
	if err != nil {
		return err
	}
	return mygo.Shell.OpenPath(dir)
}
