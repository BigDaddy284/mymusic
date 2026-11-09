package localfiles

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestPathsNotUTF8 keeps the path of a file whose name isn't UTF-8 in the
// index file, so that the index
// read again, as after a restart, names the same file: in base64, while
// a path that is UTF-8 reads as text.
func TestPathsNotUTF8(t *testing.T) {
	ti := newIndex(t)
	text := filepath.Join(ti.dir, "Café", "Café Nights.mp3") // decomposed
	raw := filepath.Join(ti.dir, "caf\xe9", "\xff01.mp3")
	ti.mu.Lock()
	for _, path := range []string{text, raw} {
		ti.entries[path] = Entry{Path: path, Size: 12, Modified: 34}
	}
	ti.mu.Unlock()
	if err := ti.save(); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(filepath.Dir(ti.dir), "local-files.json")
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Café Nights.mp3") || strings.Contains(string(data), `�`) {
		t.Errorf("the index file holds %s", data)
	}
	x := ti.open(t, file, filepath.Join(filepath.Dir(ti.dir), "cache", "local-art"))
	for _, path := range []string{text, raw} {
		e, ok := x.Entry(path)
		if !ok || e.Path != path || e.Size != 12 {
			t.Errorf("%q read back as %+v, %v", path, e, ok)
		}
		if got, ok := PathOf(TrackOf(e).ID); !ok || got != path {
			t.Errorf("the track of %q names %q", path, got)
		}
	}
}

// TestIndexOfBefore reads an index file whose paths are all text, as the
// versions before wrote them, a path that wasn't UTF-8 with replacement
// characters in place of its bytes. The file is of the index's own
// version: one of another is read again from the files.
func TestIndexOfBefore(t *testing.T) {
	ti := newIndex(t)
	file := filepath.Join(filepath.Dir(ti.dir), "local-files.json")
	a := filepath.Join(ti.dir, "a.mp3")
	b := filepath.Join(ti.dir, "caf�.mp3")
	quoted := func(p string) string { return strings.ReplaceAll(p, `\`, `\\`) }
	old := `{"version": ` + strconv.Itoa(version) + `, "files": [{"path": "` + quoted(a) + `", "size": 3, "modified": 4, "title": "One"},
		{"path": "` + quoted(b) + `", "size": 5, "modified": 6}]}`
	if err := os.WriteFile(file, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	x := ti.open(t, file, filepath.Join(filepath.Dir(ti.dir), "cache", "local-art"))
	if e, ok := x.Entry(a); !ok || e.Title != "One" {
		t.Errorf("a: %+v, %v", e, ok)
	}
	if e, ok := x.Entry(b); !ok || e.Size != 5 {
		t.Errorf("b: %+v, %v", e, ok)
	}
}
