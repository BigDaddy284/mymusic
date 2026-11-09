package settings

import (
	"slices"

	"github.com/bigdaddy284/mymusic/internal/volume"
)

// cleanVolumeCurve is c, or def for a curve MyMusic doesn't know, as in
// settings from before the curve.
func cleanVolumeCurve(c, def volume.Curve) volume.Curve {
	if !slices.Contains(volume.Curves, c) {
		return def
	}
	return c
}
