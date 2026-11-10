//go:build !windows

package app

// processes returns none: MyTube hears what another process plays on
// Windows alone.
func processes() []proc { return nil }

// webViewAudio returns 0: MyTube hears its web views on Windows alone.
func webViewAudio([]proc) uint32 { return 0 }
