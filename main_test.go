package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultDownloadsDir keeps the downloads beside the cache, and with a
// portable copy in its own folder.
func TestDefaultDownloadsDir(t *testing.T) {
	root := t.TempDir()
	at := func(parts ...string) string { return filepath.Join(append([]string{root}, parts...)...) }
	for _, c := range []struct {
		name, portable, cache, want string
	}{
		{"installed", "", at("Local", "MyTube"), at("Local", "MyTube", "downloads")},
		{"portable", at("portable"), at("portable", "cache"), at("portable", "downloads")},
	} {
		if got := defaultDownloadsDir(c.portable, c.cache); got != c.want {
			t.Errorf("%s: %s, want %s", c.name, got, c.want)
		}
	}
}

// TestAdoptFolder moves MyMusic's folder to MyTube's, and leaves MyTube's
// own alone.
func TestAdoptFolder(t *testing.T) {
	base := t.TempDir()
	old, dir := filepath.Join(base, oldName), filepath.Join(base, "MyTube")
	if err := os.MkdirAll(filepath.Join(old, "data"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := adoptFolder(dir, old); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "data")); err != nil {
		t.Errorf("MyMusic's folder didn't move: %v", err)
	}
	if err := os.MkdirAll(old, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := adoptFolder(dir, old); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(old); err != nil {
		t.Errorf("MyMusic's folder moved over MyTube's: %v", err)
	}
}
