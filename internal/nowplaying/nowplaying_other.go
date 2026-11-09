//go:build !(windows && (amd64 || arm64))

package nowplaying

// MyMusic reaches the media controls of 64-bit Windows alone: elsewhere,
// its own audio engine has no place in the system's media controls.
func open(func(Request), func(func())) (platform, error) { return nil, ErrUnsupported }
