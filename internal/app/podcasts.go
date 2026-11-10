package app

// Podcasts: shows and their episodes, on YouTube Music, the shows the
// user follows, and where they left the episodes.

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bigdaddy284/mymusic/internal/library"
	"github.com/bigdaddy284/mymusic/internal/media"
)

// PodcastSearch holds the podcasts and episodes YouTube Music found.
type PodcastSearch struct {
	YouTube media.PodcastResults `json:"youtube"`
	// Errors are the services' errors, for those that failed.
	Errors []string `json:"errors"`
}

// SearchPodcasts searches YouTube Music's podcasts and episodes.
func (c *Catalog) SearchPodcasts(ctx context.Context, query string) (PodcastSearch, error) {
	query = strings.TrimSpace(query)
	res := PodcastSearch{YouTube: media.EmptyPodcastResults(), Errors: []string{}}
	if query == "" {
		return res, nil
	}
	found, err := c.a.youtube.SearchPodcasts(ctx, query)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, err
	}
	res.YouTube = found
	return res, nil
}

// Podcast returns a podcast with its newest episodes, and whether the
// user saved it on its service.
func (c *Catalog) Podcast(ctx context.Context, id string) (media.Podcast, error) {
	source, native, err := media.ParseID(id)
	if err != nil {
		return media.Podcast{}, err
	}
	switch source {
	case media.SourceYouTube:
		return c.a.youtube.Podcast(ctx, native)
	}
	return media.Podcast{}, fmt.Errorf("no podcast %s", id)
}

// MoreEpisodes returns the page of a podcast's episodes that more, the
// podcast's More, reads. show is the podcast, which the episodes name.
func (c *Catalog) MoreEpisodes(ctx context.Context, show media.Ref, more string) (media.EpisodePage, error) {
	source, _, err := media.ParseID(show.ID)
	if err != nil {
		return media.EpisodePage{}, err
	}
	switch source {
	case media.SourceYouTube:
		return c.a.youtube.PodcastEpisodes(ctx, show, more)
	}
	return media.EpisodePage{}, fmt.Errorf("no podcast %s", show.ID)
}

// Podcasts returns the podcasts the user follows in MyMusic, the last
// followed first.
func (c *Collection) Podcasts() []media.Podcast { return c.a.library.Podcasts() }

// Follows reports whether the user follows the podcast id in MyMusic.
func (c *Collection) Follows(id string) bool { return c.a.library.Follows(id) }

// SetPodcastFollowed follows a podcast in MyMusic, or stops, and saves it
// to the user's library of its service, or takes it out, when they signed
// in there: that it can't is logged, as for likes, unless the user stops
// following one they followed only there, which then fails.
func (c *Collection) SetPodcastFollowed(ctx context.Context, p media.Podcast, on bool) error {
	here := c.a.library.Follows(p.ID)
	if err := c.a.library.SetFollowed(p, on); err != nil {
		return err
	}
	var err error
	if p.Source == media.SourceYouTube && c.a.youtube.SignedIn() && p.PlaylistID != "" {
		err = c.a.youtube.SetSaved(ctx, p.PlaylistID, on)
	}
	switch {
	case err != nil && !on && !here:
		return err
	case err != nil:
		log.Printf("mymusic: saving a podcast on its service: %v", err)
	}
	return nil
}

// YouTubePodcasts returns the podcasts the user saved on YouTube Music.
func (c *Collection) YouTubePodcasts(ctx context.Context) ([]media.Podcast, error) {
	return c.a.youtube.LibraryPodcasts(ctx)
}

// Position returns where the user left an episode.
func (c *Collection) Position(id string) (library.Position, bool) { return c.a.positions.Get(id) }

// InProgress returns the episodes the user left partway, the latest
// first, at most limit.
func (c *Collection) InProgress(limit int) []library.Position {
	return c.a.positions.InProgress(limit)
}

// PositionsVersion counts the changes of where episodes were left.
func (c *Collection) PositionsVersion() uint64 { return c.a.positions.Version() }

// SetPlayed marks an episode played, or not played, which forgets where
// it was left.
func (c *Collection) SetPlayed(t media.Track, played bool) { c.a.positions.SetPlayed(t, played) }

// keepPosition keeps where the episode that plays is, as the player
// reports it each second: not while it loads, at where it starts.
func (a *App) keepPosition(np NowPlaying) {
	if np.Track == nil || !np.Track.IsEpisode() || np.Live || np.Loading {
		return
	}
	a.positions.Keep(*np.Track, time.Duration(np.PositionMs)*time.Millisecond, time.Duration(np.DurationMs)*time.Millisecond)
}

// ResumeAt is where an episode plays from as it becomes current.
func (s playerServices) ResumeAt(t media.Track) time.Duration { return s.a.positions.ResumeAt(t) }
