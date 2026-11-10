package library

import (
	"errors"
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// fromSpotify reports whether what has the ID id, of source, came from
// Spotify, which MyTube 0.x played from and kept in the library, and
// which MyTube plays no more.
func fromSpotify(id string, source media.Source) bool {
	return source == "spotify" || strings.HasPrefix(id, "sp:")
}

// spotifyTrack reports whether t came from Spotify.
func spotifyTrack(t media.Track) bool { return fromSpotify(t.ID, t.Source) }

// holdsSpotify reports whether d holds anything that came from Spotify.
func holdsSpotify(d Data) bool {
	if slices.ContainsFunc(d.Liked, spotifyTrack) || slices.ContainsFunc(d.Disliked, spotifyTrack) ||
		slices.ContainsFunc(d.History, func(p Play) bool { return spotifyTrack(p.Track) }) ||
		slices.ContainsFunc(d.Artists, func(a FollowedArtist) bool { return fromSpotify(a.ID, a.Source) }) ||
		slices.ContainsFunc(d.Podcasts, func(p media.Podcast) bool { return fromSpotify(p.ID, p.Source) }) {
		return true
	}
	return slices.ContainsFunc(d.Playlists, func(p Playlist) bool { return slices.ContainsFunc(p.Tracks, spotifyTrack) })
}

// dropSpotify removes from d, which it changes in place, what came from
// Spotify.
func dropSpotify(d *Data) {
	track := spotifyTrack
	d.Liked = slices.DeleteFunc(d.Liked, track)
	d.Disliked = slices.DeleteFunc(d.Disliked, track)
	for i := range d.Playlists {
		d.Playlists[i].Tracks = slices.DeleteFunc(d.Playlists[i].Tracks, track)
	}
	d.History = slices.DeleteFunc(d.History, func(p Play) bool { return track(p.Track) })
	d.Artists = slices.DeleteFunc(d.Artists, func(a FollowedArtist) bool { return fromSpotify(a.ID, a.Source) })
	d.Podcasts = slices.DeleteFunc(d.Podcasts, func(p media.Podcast) bool { return fromSpotify(p.ID, p.Source) })
}

// forgetSpotify drops what the library kept of Spotify, once: the file as
// it was stays beside it, as library.json.spotify, the first time.
func (l *Library) forgetSpotify() error {
	if !holdsSpotify(l.file.Get()) {
		return nil
	}
	backup := l.file.Path() + ".spotify"
	if _, err := os.Stat(backup); errors.Is(err, fs.ErrNotExist) {
		if data, err := os.ReadFile(l.file.Path()); err == nil {
			if err := os.WriteFile(backup, data, 0o600); err != nil {
				return err
			}
		}
	}
	return l.file.Update(func(d *Data) error {
		dropSpotify(d)
		return nil
	})
}
