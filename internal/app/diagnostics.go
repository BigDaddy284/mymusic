package app

import (
	"fmt"
	"path/filepath"

	"github.com/egoist/mygo"
)

// LogsFolder returns the folder of MyTube's log and of its crash
// reports.
func (s *System) LogsFolder() string { return s.a.opts.LogsDir }

// OpenLogsFolder opens the folder of the log in the file manager.
func (s *System) OpenLogsFolder() error {
	if s.a.opts.LogsDir == "" {
		return fmt.Errorf("MyTube keeps no log here")
	}
	return mygo.Shell.OpenPath(s.a.opts.LogsDir)
}

// LastCrash returns the name of the crash report of the run before, when
// MyTube crashed then and the user has not answered yet; "" otherwise.
func (s *System) LastCrash() string {
	s.a.mu.Lock()
	defer s.a.mu.Unlock()
	if s.a.lastCrash == "" {
		return ""
	}
	return filepath.Base(s.a.lastCrash)
}

// DismissCrash forgets the crash of the run before; its report stays in
// the logs folder.
func (s *System) DismissCrash() {
	s.a.mu.Lock()
	defer s.a.mu.Unlock()
	s.a.lastCrash = ""
}
