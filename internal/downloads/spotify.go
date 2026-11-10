package downloads

import (
	"slices"
	"strings"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// fromSpotify reports whether what has the ID id, of source, came from
// Spotify, which MyTube 0.x downloaded the recordings of, and which
// MyTube plays no more.
func fromSpotify(id string, source media.Source) bool {
	return source == "spotify" || strings.HasPrefix(id, "sp:")
}

// forgetSpotify turns the Spotify tracks an index of MyTube 0.x holds
// into the tracks of YouTube Music they downloaded, which play on, and
// drops the Spotify albums and playlists downloaded whole, with the
// tracks only they kept: the audio no track needs then goes, at the
// next check.
func (x *index) forgetSpotify() {
	changed := false
	for id, c := range x.Collections {
		if !fromSpotify(c.ID, c.Source) && !strings.Contains(c.ID, ":sp:") {
			continue
		}
		changed = true
		delete(x.Collections, id)
		x.dropArt(c.Art)
		for _, e := range x.Tracks {
			e.Pins = slices.DeleteFunc(e.Pins, func(p string) bool { return p == id })
		}
	}
	renamed := map[string]string{}
	for id, e := range x.Tracks {
		rec := e.Recording
		e.Recording = nil
		if !fromSpotify(id, e.Track.Source) && len(e.Pins) > 0 {
			continue
		}
		changed = true
		switch {
		case len(e.Pins) > 0 && rec != nil && rec.Source == media.SourceYouTube:
			delete(x.Tracks, id)
			renamed[id] = rec.ID
			if kept, ok := x.Tracks[rec.ID]; ok {
				for _, p := range e.Pins {
					if !kept.pinned(p) {
						kept.Pins = append(kept.Pins, p)
					}
				}
				x.dropArt(e.Art)
				continue
			}
			e.Track = *rec
			x.Tracks[rec.ID] = e
		default:
			delete(x.Tracks, id)
			x.dropArt(e.Art)
		}
	}
	if !changed {
		return
	}
	for _, c := range x.Collections {
		var tracks []string
		for _, id := range c.Tracks {
			if to, ok := renamed[id]; ok {
				id = to
			}
			if _, ok := x.Tracks[id]; ok && !slices.Contains(tracks, id) {
				tracks = append(tracks, id)
			}
		}
		c.Tracks = tracks
	}
	wanted := map[string]bool{}
	for _, e := range x.Tracks {
		wanted[e.video()] = true
	}
	for video, f := range x.Files {
		if !wanted[video] {
			delete(x.Files, video)
			x.Pending = append(x.Pending, f.Path)
		}
	}
}

// dropArt has the cover p deleted, unless a track or a collection keeps
// it.
func (x *index) dropArt(p string) {
	if p == "" {
		return
	}
	for _, e := range x.Tracks {
		if e.Art == p {
			return
		}
	}
	for _, c := range x.Collections {
		if c.Art == p {
			return
		}
	}
	for u, a := range x.Art {
		if a == p {
			delete(x.Art, u)
		}
	}
	x.Pending = append(x.Pending, p)
}
