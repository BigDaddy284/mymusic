package app

import (
	"io/fs"

	"github.com/egoist/mygo"

	"github.com/bigdaddy284/mymusic/internal/addons"
)

// Addons lists and removes themes and plugins, and keeps what
// plugins store.
type Addons struct{ a *App }

// List returns the themes and plugins available, and those that could not
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

// AddonFiles returns the files of an add-on, for the native interface's
// plugins.
func (a *App) AddonFiles(kind addons.Kind, id string) (fs.FS, error) { return a.addons.Files(kind, id) }

// Storage returns every value a plugin keeps.
func (s *Addons) Storage(plugin string) (map[string]string, error) { return s.a.storage.All(plugin) }

// StorageSet keeps a value of a plugin.
func (s *Addons) StorageSet(plugin, key, value string) error {
	return s.a.storage.Set(plugin, key, value)
}

// StorageDelete removes a value of a plugin.
func (s *Addons) StorageDelete(plugin, key string) error { return s.a.storage.Delete(plugin, key) }
