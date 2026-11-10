package app

import (
	"context"

	"github.com/bigdaddy284/mymusic/internal/library"
	"github.com/bigdaddy284/mymusic/internal/media"
)

// Collection is the user's library: liked and disliked tracks and
// playlists mixing every source, and the history.
type Collection struct{ a *App }

// Liked returns the liked tracks, the last liked first.
func (c *Collection) Liked() []media.Track { return c.a.library.Liked() }

// Version counts the changes of the library, as a play or a new mix: what
// a window made of it holds until the version changes.
func (c *Collection) Version() uint64 { return c.a.library.Version() }

// LikedIDs returns the IDs of the liked tracks, for lists to show hearts.
func (c *Collection) LikedIDs() []string {
	liked := c.a.library.Liked()
	ids := make([]string, len(liked))
	for i, t := range liked {
		ids[i] = t.ID
	}
	return ids
}

// SetLiked likes or unlikes a track, and YouTube Music's on YouTube Music
// too when the user signed in. Liking a disliked track
// takes the dislike back; the player hears of that from the caller.
func (c *Collection) SetLiked(ctx context.Context, track media.Track, liked bool) error {
	if err := c.a.library.SetLiked(track, liked); err != nil {
		return err
	}
	c.a.syncYouTubeLike(ctx, track, liked)
	return nil
}

// Disliked returns the disliked tracks, the last disliked first.
func (c *Collection) Disliked() []media.Track { return c.a.library.Disliked() }

// SetDisliked dislikes a track, or takes the dislike back, in MyTube's
// library, and on YouTube Music for its tracks when the user signed in.
// The player hears of the dislike from the caller, which tells it as the
// user asks.
func (c *Collection) SetDisliked(ctx context.Context, track media.Track, disliked bool) error {
	if _, err := c.a.library.SetDisliked(track, disliked); err != nil {
		return err
	}
	c.a.syncYouTubeDislike(ctx, track, disliked)
	return nil
}

// Playlists returns the playlists made in MyTube, without their tracks.
func (c *Collection) Playlists() []media.Playlist { return c.a.library.Playlists() }

// localID returns the library's ID of a playlist made in MyTube.
func localID(id string) (string, error) {
	source, native, err := media.ParseID(id)
	if err != nil || source != media.SourceLocal {
		return "", library.ErrNoPlaylist
	}
	return native, nil
}

// CreatePlaylist makes a playlist, holding tracks.
func (c *Collection) CreatePlaylist(title, description string, tracks []media.Track) (media.Playlist, error) {
	return c.a.library.CreatePlaylist(title, description, tracks)
}

// UpdatePlaylist changes the title and the description of a playlist.
func (c *Collection) UpdatePlaylist(id, title, description string) error {
	native, err := localID(id)
	if err != nil {
		return err
	}
	return c.a.library.UpdatePlaylist(native, title, description)
}

// DeletePlaylist deletes a playlist, and the file of the cover chosen for
// it, if any.
func (c *Collection) DeletePlaylist(id string) error {
	native, err := localID(id)
	if err != nil {
		return err
	}
	cover, _ := c.a.library.PlaylistCover(native)
	if err := c.a.library.DeletePlaylist(native); err != nil {
		return err
	}
	c.a.removeMixCover(cover)
	return nil
}

// AddToPlaylist appends tracks to a playlist made in MyTube.
func (c *Collection) AddToPlaylist(ctx context.Context, id string, tracks []media.Track) error {
	native, err := localID(id)
	if err != nil {
		return err
	}
	return c.a.library.AddToPlaylist(native, tracks)
}

// RemoveFromPlaylist removes the track at index from a playlist.
func (c *Collection) RemoveFromPlaylist(id string, index int) error {
	native, err := localID(id)
	if err != nil {
		return err
	}
	return c.a.library.RemoveFromPlaylist(native, index)
}

// MoveInPlaylist moves the track at from to the position to.
func (c *Collection) MoveInPlaylist(id string, from, to int) error {
	native, err := localID(id)
	if err != nil {
		return err
	}
	return c.a.library.MoveInPlaylist(native, from, to)
}

// ImportPlaylist copies a playlist of a service into a playlist of
// MyTube, where tracks of other services can join it.
func (c *Collection) ImportPlaylist(ctx context.Context, id string) (media.Playlist, error) {
	p, err := c.a.wholePlaylist(ctx, id)
	if err != nil {
		return media.Playlist{}, err
	}
	return c.a.library.CreatePlaylist(p.Title, p.Description, p.Tracks)
}

// RecordPlay adds a track to the history, once it played long enough to
// count.
func (c *Collection) RecordPlay(track media.Track) error {
	err := c.a.library.RecordPlay(track)
	c.a.reportYouTubePlay(track)
	return err
}

// History returns the last plays, the latest first.
func (c *Collection) History(limit int) []library.Play { return c.a.library.History(limit) }

// ClearHistory forgets what the user played.
func (c *Collection) ClearHistory() error { return c.a.library.ClearHistory() }

// Stats sums up the plays of the last days, or of all time when days is 0.
func (c *Collection) Stats(days int) library.Stats { return c.a.library.Stats(days) }

// Searches returns the recent searches.
func (c *Collection) Searches() []string { return c.a.library.Searches() }

// ClearSearches forgets the recent searches.
func (c *Collection) ClearSearches() error { return c.a.library.ClearSearches() }
