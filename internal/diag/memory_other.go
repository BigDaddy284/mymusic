//go:build !windows

package diag

// processMemory tells nothing of the process where the system doesn't tell
// it cheaply.
func processMemory() (resident, private uint64) { return 0, 0 }
