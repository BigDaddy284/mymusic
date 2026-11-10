package main

// One search for both sides: a search of the Video side shows YouTube
// Music's songs too, under its Songs chip, which play on the music side;
// a search of the Music side shows the first of YouTube's videos under
// All, which play on the Video side, with a link to the rest.

import (
	"context"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/app"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// songResults shows the songs YouTube Music finds for q, as the Music
// side's search does.
func (s *shell) songResults(c *ui.Context, q string) {
	res := fetch(s, "search:"+q, func(ctx context.Context) (app.SearchResponse, error) { return s.svc.Catalog.Search(ctx, q) })
	switch {
	case !res.done:
		s.loading(c)
		return
	case res.err != nil:
		s.problem(c, res.err, func() { s.forget("search:" + q) })
		return
	}
	r := mergeResults(res.value, s.searchScope)
	ui.Column(c).Margin(0, 0, 32).Children(func() {
		s.trackList(c, "song-results", r.Tracks, trackListOptions{art: true, album: true, plays: true, empty: "No songs match."})
	})
	if s.button(c, "songs-on-music", "search music for it", false).AlignSelf(ui.Start).Clicked() {
		s.navigate(searchPath(q))
	}
}

// videoResultsShelf shows the first videos of YouTube's website for q, on
// the Music side's search, with a link to all of them on the Video side.
func (s *shell) videoResultsShelf(c *ui.Context, q string) {
	key := "videos:search:" + q
	res := fetch(s, key, func(ctx context.Context) (youtube.VideoPage, error) { return s.svc.Videos.Search(ctx, q) })
	if !res.done || res.err != nil {
		return
	}
	videos := s.shownVideos(append([]youtube.Video{}, res.value.Videos...))
	if len(videos) == 0 {
		return
	}
	k := &s.lk
	ui.Row(c).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
		s.sectionTitle(c, "On YouTube")
		b := s.buttonBase(c.Key("all-youtube-videos")).Label("Show all YouTube videos").Cursor(ui.CursorPointer)
		if b.Clicked() {
			s.videos.query = q
			s.videos.searchChip = 0
			s.navigate(videoSearchPath(q))
		}
		b.Children(func() { ui.Text(c, "show all").FontSize(k.smallSize()).TextColor(k.accent) })
	})
	// One row of cards.
	n, _ := s.videoColumns(c, 0)
	ui.Column(c).Margin(0, 0, 48).Children(func() {
		s.videoGrid(c, "music-search-videos", videos[:min(len(videos), max(n, 2))])
	})
}
