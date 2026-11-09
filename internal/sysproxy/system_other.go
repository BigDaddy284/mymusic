//go:build !windows

package sysproxy

// install does nothing off Windows: the environment's proxy stays.
func install() {}
