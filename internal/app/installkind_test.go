package app

import (
	"path/filepath"
	"testing"
)

func TestKindOf(t *testing.T) {
	const installer = `C:\Users\me\AppData\Local\Programs\MyTube`
	for _, c := range []struct {
		name  string
		dir   string
		files []string
		want  string
	}{
		{"installer", installer, []string{filepath.Join(installer, "Uninstall.exe")}, "installer"},
		{"portable zip", `D:\Apps\MyTube`, nil, "portable"},
	} {
		exists := func(p string) bool {
			for _, f := range c.files {
				if f == p {
					return true
				}
			}
			return false
		}
		if got := kindOf(c.dir, exists); got != c.want {
			t.Errorf("%s: kindOf = %q, want %q", c.name, got, c.want)
		}
	}
}
