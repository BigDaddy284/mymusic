package app

import (
	"strings"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// largeArtwork picks the artwork the system's media controls
// show: an https image of about 512 pixels.
func largeArtwork(images []media.Image) string {
	best := ""
	bestWidth := 0
	for _, img := range images {
		if !strings.HasPrefix(img.URL, "https://") {
			continue
		}
		if best == "" || (bestWidth < 512 && img.Width > bestWidth) || (img.Width >= 512 && img.Width < bestWidth) {
			best, bestWidth = img.URL, img.Width
		}
	}
	return best
}
