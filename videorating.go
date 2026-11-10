package main

// The views of the video that plays, and the account's like and dislike
// of it, under its title in the Video side's player bar.

import (
	"context"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// videoDetails are the watch page's details of the video id.
func (s *shell) videoDetails(id string) *loaded[youtube.VideoDetails] {
	return fetch(s, "videos:details:"+id, func(ctx context.Context) (youtube.VideoDetails, error) {
		return s.svc.Videos.Details(ctx, id)
	})
}

// videoRating is the account's rating of the video id: as the user set it
// here, or as YouTube told it.
func (s *shell) videoRating(id string) string {
	if r, ok := s.videos.ratings[id]; ok {
		return r
	}
	if d := s.videoDetails(id); d.done && d.err == nil {
		return d.value.Rating
	}
	return ""
}

// rateVideo rates the video id on the account, and takes the rating back
// where YouTube refuses it.
func (s *shell) rateVideo(id, rating string) {
	vs := &s.videos
	if vs.ratings == nil {
		vs.ratings = map[string]string{}
	}
	was := s.videoRating(id)
	vs.ratings[id] = rating
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := s.svc.Videos.Rate(ctx, id, rating); err != nil {
			s.update(func() {
				vs.ratings[id] = was
				s.notes = append(s.notes, noticeOf(err))
			})
		}
	}()
}

// videoRatingRow shows how many saw v, with the like and dislike buttons
// of the account, which signs in on YouTube.
func (s *shell) videoRatingRow(c *ui.Context, v youtube.Video) {
	k := &s.lk
	views := v.Views
	d := s.videoDetails(v.ID)
	if d.done && d.err == nil && d.value.Views != "" {
		views = d.value.Views
	}
	rating := s.videoRating(v.ID)
	ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
		if views != "" {
			ui.Text(c, views).FontSize(k.captionSize()).TextColor(k.muted).SingleLine().Margin(0, 6, 0, 0)
		}
		like, dislike := "Like", "Dislike"
		if rating == "like" {
			like = "Remove the like"
		}
		if rating == "dislike" {
			dislike = "Remove the dislike"
		}
		if s.smallRatingButton(c, "v-like", "like", like, rating == "like") {
			s.rateVideo(v.ID, map[bool]string{true: "", false: "like"}[rating == "like"])
		}
		if s.smallRatingButton(c, "v-dislike", "dislike", dislike, rating == "dislike") {
			s.rateVideo(v.ID, map[bool]string{true: "", false: "dislike"}[rating == "dislike"])
		}
	})
}

// smallRatingButton is a like or dislike button the height of a line of
// small text, to fit under the title in the player bar.
func (s *shell) smallRatingButton(c *ui.Context, key, name, label string, on bool) bool {
	k := &s.lk
	b := s.buttonBase(c.Key(key)).Size(26, 24).Radius(k.round(24)).Label(label).Tooltip(label).Cursor(ui.CursorPointer)
	ink := k.muted
	switch {
	case on:
		ink = k.accent
	case b.Hovered():
		ink = k.text
	}
	if b.Hovered() {
		b.Background(k.hover)
	}
	clicked := b.Clicked()
	b.Children(func() {
		svg := icons[name]
		if on {
			svg = solidIcons[name]
		}
		ui.Icon(c, svg).FontSize(15).TextColor(ink)
	})
	return clicked
}
