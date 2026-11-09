package main

import (
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
		{"installed", "", at("Local", "MyMusic"), at("Local", "MyMusic", "downloads")},
		{"portable", at("portable"), at("portable", "cache"), at("portable", "downloads")},
	} {
		if got := defaultDownloadsDir(c.portable, c.cache); got != c.want {
			t.Errorf("%s: %s, want %s", c.name, got, c.want)
		}
	}
}
