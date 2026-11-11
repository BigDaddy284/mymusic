package addons

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func themeJSON(id string) string { return strings.Replace(validTheme, `"test-theme"`, `"`+id+`"`, 1) }

func testManager(t *testing.T) *Manager {
	t.Helper()
	builtin := fstest.MapFS{
		"themes/mix/theme.json":  {Data: []byte(themeJSON("mix"))},
		"themes/mix/preview.png": {Data: []byte("mix")},
	}
	dir := t.TempDir()
	write := func(name, data string) {
		p := filepath.Join(dir, filepath.FromSlash(name))
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("themes/ocean/theme.json", themeJSON("ocean"))
	write("themes/ocean/preview.png", "ocean")
	write("themes/mix/theme.json", themeJSON("mix"))             // takes a built-in id
	write("themes/misnamed/theme.json", themeJSON("other-name")) // folder and id differ
	write("themes/broken/theme.json", "{")
	return NewManager(dir, builtin)
}

func TestList(t *testing.T) {
	c := testManager(t).List()
	if len(c.Themes) != 2 || c.Themes[0].ID != "mix" || !c.Themes[0].Builtin || c.Themes[1].ID != "ocean" || c.Themes[1].Builtin {
		t.Errorf("themes = %+v", c.Themes)
	}
	if len(c.Problems) != 3 {
		t.Errorf("problems = %+v", c.Problems)
	}
}

func TestUninstall(t *testing.T) {
	m := testManager(t)
	if err := m.Uninstall(KindTheme, "mix"); !errors.Is(err, ErrBuiltin) {
		t.Errorf("uninstalling a built-in theme: %v", err)
	}
	if err := m.Uninstall(KindTheme, "ocean"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(m.Dir, "themes", "ocean")); !os.IsNotExist(err) {
		t.Error("the theme is still there")
	}
	if err := m.Uninstall(KindTheme, "../themes"); err == nil {
		t.Error("uninstalled outside the add-ons folder")
	}
}
