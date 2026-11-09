//go:build !windows && !linux && !(darwin && (amd64 || arm64))

package localfiles

// watchFolder watches nothing here: MyMusic reads the folders as it
// starts, and when asked.
func watchFolder(string, Changes) (func(), error) { return nil, ErrCantWatch }
