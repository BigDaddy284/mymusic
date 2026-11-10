package app

import (
	"path/filepath"
	"testing"
)

func TestLastCrash(t *testing.T) {
	a := &App{lastCrash: filepath.Join(t.TempDir(), "crash-20261019-210000.log")}
	s := &System{a}
	if got := s.LastCrash(); got != "crash-20261019-210000.log" {
		t.Errorf("LastCrash = %q", got)
	}
	s.DismissCrash()
	if got := s.LastCrash(); got != "" {
		t.Errorf("after DismissCrash, LastCrash = %q", got)
	}
}
