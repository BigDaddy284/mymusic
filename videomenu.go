package main

// The menu of a video card, from its ⋮ button or a right click, as on
// YouTube: add it to Up next, save it to Watch Later or a playlist, copy
// its link, tell YouTube not to recommend it or its channel, or report it
// on YouTube.

import (
	"context"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// videoMenu builds the menu of v, a card of list.
func (s *shell) videoMenu(m *ui.Menu, v youtube.Video, list []youtube.Video) {
	if m.Item("Play").Chosen() {
		s.watch(v, list)
	}
	if m.Item("Add to queue").Chosen() {
		s.queueVideo(v)
	}
	m.Separator()
	signed := s.svc.Videos.SignedIn()
	if m.Item("Save to Watch Later").Disabled(!signed).Chosen() {
		s.saveVideo(v, youtube.WatchLater, "Watch Later")
	}
	m.Submenu("Save to playlist", func(m *ui.Menu) {
		lists := s.videoPlaylistsLoad()
		switch {
		case !signed:
			m.Item("Sign in to YouTube Music to save videos").Disabled(true)
		case !lists.done:
			m.Item("Loading your playlists…").Disabled(true)
		case lists.err != nil || len(lists.value) == 0:
			m.Item("No playlists of yours").Disabled(true)
		}
		if !signed || !lists.done || lists.err != nil {
			return
		}
		for _, p := range lists.value {
			if p.ID == "" || p.ID == youtube.WatchLater || p.ID == "VL"+youtube.WatchLater {
				continue
			}
			if m.Item(p.Title).Chosen() {
				s.saveVideo(v, p.ID, p.Title)
			}
		}
	}).Disabled(!signed)
	if m.Item("Share: copy the link").Chosen() {
		copyText(youtubeWatchURL(v.ID))
		s.notes = append(s.notes, notice("Copied the link to "+v.Title))
	}
	m.Separator()
	if m.Item("Not interested").Disabled(v.Feedback.NotInterested == "").
		Tooltip("YouTube shows this for videos of your home feed").Chosen() {
		s.videoFeedback(v, v.Feedback.NotInterested, false)
	}
	if m.Item("Don't recommend channel").Disabled(v.Feedback.NoChannel == "").
		Tooltip("YouTube shows this for videos of your home feed").Chosen() {
		s.videoFeedback(v, v.Feedback.NoChannel, true)
	}
	if m.Item("Report on YouTube").Chosen() {
		go s.svc.System.OpenExternal(youtubeWatchURL(v.ID))
	}
}

// videoMenuButton is the ⋮ button by a card's title that opens its menu.
func (s *shell) videoMenuButton(c *ui.Context, key string, v youtube.Video, list []youtube.Video) {
	s.iconButton(c, key+"-menu", "moreVertical", "More actions", false, 18).Shrink(0).
		Menu(func(m *ui.Menu) { s.videoMenu(m, v, list) })
}

// queueVideo puts v next in Up next, or plays it when nothing plays.
func (s *shell) queueVideo(v youtube.Video) {
	vp := &s.videos.player
	if vp.current.ID == "" {
		s.watch(v, []youtube.Video{v})
		return
	}
	vs := &s.videos
	if vs.known == nil {
		vs.known = map[string]youtube.Video{}
	}
	vs.known[v.ID] = v
	list := slices.DeleteFunc(slices.Clone(vp.list), func(x youtube.Video) bool { return x.ID == v.ID && x.ID != vp.current.ID })
	at := slices.IndexFunc(list, func(x youtube.Video) bool { return x.ID == vp.current.ID })
	if at < 0 {
		list = append([]youtube.Video{vp.current}, list...)
		at = 0
	}
	vp.list = slices.Insert(list, at+1, v)
	s.notes = append(s.notes, notice("Playing "+v.Title+" next"))
}

// saveVideo saves v to the account's playlist id, called name.
func (s *shell) saveVideo(v youtube.Video, id, name string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := s.svc.Videos.Save(ctx, id, v.ID)
		s.update(func() {
			if err != nil {
				s.notes = append(s.notes, noticeOf(err))
				return
			}
			s.notes = append(s.notes, notice("Saved to "+name))
			if id != youtube.WatchLater {
				s.refetch("videos:playlist:" + id)
			}
		})
	}()
}

// videoFeedback tells YouTube not to recommend v, or its channel, and
// takes it, or all the channel's videos, off the cards at once.
func (s *shell) videoFeedback(v youtube.Video, token string, channel bool) {
	vs := &s.videos
	if vs.hidden == nil {
		vs.hidden = map[string]bool{}
	}
	key := v.ID
	if channel {
		key = videoChannelKey(v)
	}
	vs.hidden[key] = true
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := s.svc.Videos.Feedback(ctx, token)
		s.update(func() {
			if err != nil {
				delete(vs.hidden, key)
				s.notes = append(s.notes, noticeOf(err))
				return
			}
			if channel {
				s.notes = append(s.notes, notice("YouTube won't recommend "+v.Channel))
			} else {
				s.notes = append(s.notes, notice("YouTube won't recommend that video"))
			}
		})
	}()
}

// videoChannelKey is the key of v's channel among the hidden videos.
func videoChannelKey(v youtube.Video) string {
	if v.ChannelID != "" {
		return "channel:" + v.ChannelID
	}
	return "channel:" + v.Channel
}

// shownVideos are the videos of list the user didn't ask YouTube not to
// recommend.
func (s *shell) shownVideos(list []youtube.Video) []youtube.Video {
	hidden := s.videos.hidden
	if len(hidden) == 0 {
		return list
	}
	return slices.DeleteFunc(list, func(v youtube.Video) bool {
		return hidden[v.ID] || v.Channel != "" && hidden[videoChannelKey(v)]
	})
}
