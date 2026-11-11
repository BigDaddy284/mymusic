package diag

// How MyTube uses its memory: capture mode's act:memory reads it, and
// act:heap saves a profile of the heap.

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

// Memory is how MyTube uses its memory at a moment: the Go heap's, as
// the runtime tells it, and the process's, where the system tells it
// cheaply.
type Memory struct {
	// HeapInUse is the bytes of the Go heap's spans in use, and Sys those
	// Go holds from the system, the heap's and the rest.
	HeapInUse, Sys uint64
	// GCs counts the collections of the Go heap so far.
	GCs uint32
	// Resident is the process's resident size, the bytes of its pages in
	// memory, and 0 where the system doesn't tell it cheaply.
	Resident uint64
	// Private is Windows' private bytes of the process, its memory no
	// other process shares.
	Private uint64
}

// ReadMemory reads how MyTube uses its memory now.
func ReadMemory() Memory {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m := Memory{HeapInUse: ms.HeapInuse, Sys: ms.Sys, GCs: ms.NumGC}
	m.Resident, m.Private = processMemory()
	return m
}

// String writes m as a line of a log: "Go heap 120 MB in use, 300 MB from
// the system, 45 GCs; resident 410 MB", and the private bytes where the
// system tells them.
func (m Memory) String() string {
	s := fmt.Sprintf("Go heap %d MB in use, %d MB from the system, %d GCs", m.HeapInUse>>20, m.Sys>>20, m.GCs)
	var proc []string
	for _, v := range []struct {
		name string
		n    uint64
	}{{"resident", m.Resident}, {"private", m.Private}} {
		if v.n > 0 {
			proc = append(proc, fmt.Sprintf("%s %d MB", v.name, v.n>>20))
		}
	}
	if len(proc) > 0 {
		s += "; " + strings.Join(proc, ", ")
	}
	return s
}

// Process is the most the process holds, of what is known: its private
// bytes or resident size, or what Go holds from the system.
func (m Memory) Process() uint64 { return max(m.Private, m.Resident, m.Sys) }

// WriteHeapProfile collects the Go heap's garbage, so that the profile
// tells what is in use, and writes a profile of the heap to path: where
// it was allocated, by how much, and nothing of what it holds.
func WriteHeapProfile(path string) error {
	runtime.GC()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := pprof.WriteHeapProfile(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
