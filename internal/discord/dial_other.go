//go:build !windows

package discord

import (
	"context"
	"io"
)

// dial finds no Discord app: MyMusic reaches Discord's named pipe on
// Windows alone.
func dial(context.Context) (io.ReadWriteCloser, error) { return nil, ErrNotRunning }
