package main

// Middle-button scrolling of the page shown, as in a web browser: hold the
// wheel down and move the pointer, and the page scrolls toward it, the
// faster the further it went; or click the wheel, move, and click again
// to stop.

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// autoScroll is the middle-button scrolling of the page shown.
type autoScroll struct {
	// on is set while the page scrolls by the pointer: anchorY is where the
	// wheel went down, y where the pointer is now, in the page's box.
	on         bool
	anchorY, y float32
	// moved is set once the pointer left the anchor's dead zone; sticky
	// once the wheel was let go there, which keeps it on until a click.
	moved, sticky bool
	// last is the time of the frame that scrolled last, and at the page it
	// scrolls.
	last time.Time
	at   string
}

const (
	// autoScrollDead is how far the pointer goes from the anchor before
	// the page scrolls.
	autoScrollDead = 12
	// autoScrollRate is how many DIPs a second the page scrolls for each
	// DIP the pointer is past the dead zone.
	autoScrollRate = 6
)

// pageInput is the input handler of the page shown: the middle button
// scrolls it, and the scrolling over it may pull it (pullInput).
func (s *shell) pageInput(ev ui.InputEvent) bool {
	a := &s.autoScroll
	switch ev.Kind {
	case ui.InputPointerDown:
		if a.on && a.sticky {
			// Any click ends the scrolling it started, and does nothing
			// else.
			*a = autoScroll{}
			return true
		}
		if ev.Button == 2 {
			*a = autoScroll{on: true, anchorY: ev.Y, y: ev.Y}
			return true
		}
	case ui.InputPointerMove:
		if a.on {
			a.y = ev.Y
			if d := a.y - a.anchorY; d > autoScrollDead || d < -autoScrollDead {
				a.moved = true
			}
		}
	case ui.InputPointerUp:
		if a.on && ev.Button == 2 {
			if a.moved || ev.Canceled {
				*a = autoScroll{}
			} else {
				a.sticky = true
			}
			return true
		}
	}
	return s.pullInput(ev)
}

// followAutoScroll scrolls the page shown, page, as far as the pointer
// asks since the frame before, and asks for the next frame while it does.
func (s *shell) followAutoScroll(c *ui.Context, page ui.Element, sc *ui.ScrollState) {
	a := &s.autoScroll
	if !a.on {
		return
	}
	// Another page ends it.
	if loc := s.router.Location(); a.at != loc {
		if a.at != "" {
			*a = autoScroll{}
			return
		}
		a.at = loc
	}
	page.Cursor(ui.CursorAllScroll)
	now := c.Now()
	dt := float32(0)
	if !a.last.IsZero() {
		dt = float32(min(now.Sub(a.last), 100*time.Millisecond).Seconds())
	}
	a.last = now
	c.After(16 * time.Millisecond)
	d := a.y - a.anchorY
	switch {
	case d > autoScrollDead:
		d -= autoScrollDead
	case d < -autoScrollDead:
		d += autoScrollDead
	default:
		return
	}
	y := s.pageY() + d*autoScrollRate*dt
	sc.Y = min(max(y, 0), sc.MaxY)
}
