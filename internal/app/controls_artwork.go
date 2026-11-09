package app

import "fmt"

// ControlsArtwork tells what MyMusic hands the system's media controls of
// the cover of what plays: for capture mode's act:nowplaying.
func (a *App) ControlsArtwork() string {
	c := &a.controls
	c.mu.Lock()
	ctl, info := c.c, c.shown()
	c.mu.Unlock()
	if ctl == nil || info == nil {
		return "the system's media controls show nothing"
	}
	return fmt.Sprintf("handed %q: %d bytes", info.Artwork, len(info.Cover))
}
