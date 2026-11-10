//go:build !windows

package localfiles

// watchFolder watches nothing here: MyTube reads the folders as it
// starts, and when asked.
func watchFolder(string, Changes) (func(), error) { return nil, ErrCantWatch }
