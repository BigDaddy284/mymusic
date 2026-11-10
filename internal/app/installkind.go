package app

import (
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
)

// InstallKind tells how this copy of MyTube was installed: "installer",
// "portable" or "development".
func (s *System) InstallKind() string { return installKind() }

// installKind tells how this copy of MyTube was installed.
func installKind() string {
	if mygo.IsDev() {
		return "development"
	}
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		exe = real
	}
	return kindOf(filepath.Dir(exe), func(p string) bool {
		_, err := os.Stat(p)
		return err == nil
	})
}

// kindOf tells how a copy whose executable is in dir was installed: the
// installer leaves its uninstaller beside it.
func kindOf(dir string, exists func(string) bool) string {
	if exists(filepath.Join(dir, "Uninstall.exe")) {
		return "installer"
	}
	return "portable"
}
