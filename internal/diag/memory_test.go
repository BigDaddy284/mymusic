package diag

import (
	"runtime"
	"strings"
	"testing"
)

// TestReadMemory reads the Go heap, and the process's memory where the
// system tells it: Windows its working set and its private bytes.
func TestReadMemory(t *testing.T) {
	m := ReadMemory()
	if m.HeapInUse == 0 || m.Sys < m.HeapInUse {
		t.Errorf("the Go heap: %+v", m)
	}
	if runtime.GOOS == "windows" && (m.Resident == 0 || m.Private == 0) {
		t.Errorf("no working set or private bytes: %+v", m)
	}
	if !strings.HasPrefix(m.String(), "Go heap ") || m.Process() < m.Sys {
		t.Errorf("%q, of %d bytes", m, m.Process())
	}
}

func TestMemoryString(t *testing.T) {
	m := Memory{HeapInUse: 120 << 20, Sys: 300 << 20, GCs: 45}
	if got := m.String(); got != "Go heap 120 MB in use, 300 MB from the system, 45 GCs" {
		t.Errorf("of Go alone: %q", got)
	}
	m.Resident, m.Private = 410<<20, 380<<20
	if got := m.String(); got != "Go heap 120 MB in use, 300 MB from the system, 45 GCs; resident 410 MB, private 380 MB" {
		t.Errorf("with the process's: %q", got)
	}
	if m.Process() != 410<<20 {
		t.Errorf("the process holds %d", m.Process())
	}
}
