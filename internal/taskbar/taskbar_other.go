//go:build !(windows && (amd64 || arm64))

package taskbar

// 32-bit Windows has no buttons on its taskbar's thumbnails here.
type bar struct{}

func attach(uintptr, func(Button)) (*bar, error) { return nil, ErrUnsupported }

func (*bar) show(state) {}
