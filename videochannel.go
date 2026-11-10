package main

// A channel's page of the Video side: its avatar, name, handle and
// subscribers, a button to subscribe, and its latest videos or its
// playlists. A channel's name on a card or in the player bar opens it.

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

func videoChannelPath(id string) string { return "/videos/channel/" + url.PathEscape(id) }
func youtubeChannelURL(id string) string {
	return "https://www.youtube.com/channel/" + url.PathEscape(id)
}

// channelTabs are the tabs of a channel's page.
var channelTabs = []string{"Videos", "Playlists"}

// openChannel shows the page of v's channel.
func (s *shell) openChannel(v youtube.Video) {
	if v.ChannelID == "" {
		return
	}
	vs := &s.videos
	if vs.channelNames == nil {
		vs.channelNames = map[string]string{}
	}
	vs.channelNames[v.ChannelID] = v.Channel
	s.navigate(videoChannelPath(v.ChannelID))
}

// channelLink is v's channel's name, in size, which opens its page when
// clicked, when the result told its ID.
func (s *shell) channelLink(c *ui.Context, key string, v youtube.Video, size float32) {
	k := &s.lk
	if v.ChannelID == "" {
		ui.Text(c, v.Channel).FontSize(size).TextColor(k.muted).SingleLine().MinWidth(0)
		return
	}
	b := s.buttonBase(c.Key(key)).MinWidth(0).Label("Open the channel " + v.Channel).Cursor(ui.CursorPointer).AlignSelf(ui.Start)
	if b.Clicked() {
		s.openChannel(v)
	}
	ink := k.muted
	if b.Hovered() {
		ink = k.text
	}
	b.Children(func() { ui.Text(c, v.Channel).FontSize(size).TextColor(ink).SingleLine().MinWidth(0) })
}

// videoChannelPage shows the channel id, its videos or its playlists.
func (s *shell) videoChannelPage(c *ui.Context, id, tab string) {
	vs := &s.videos
	playlists := tab == "playlists"
	key := "videos:channel:" + id
	if playlists {
		key += ":playlists"
	}
	res := fetch(s, key, func(ctx context.Context) (youtube.VideoChannel, error) { return s.svc.Videos.Channel(ctx, id, playlists) })
	name := vs.channelNames[id]
	if res.done && res.err == nil && res.value.Name != "" {
		name = res.value.Name
	}
	if name == "" {
		name = "Channel"
	}
	s.pageScroll(c, func() {
		selected := 0
		if playlists {
			selected = 1
		}
		before := selected
		s.pageHeader(c, name, channelTabs, &selected)
		if selected != before {
			to := videoChannelPath(id)
			if selected == 1 {
				to += "?tab=playlists"
			}
			s.router.Replace(to)
		}
		switch {
		case !res.done:
			s.loading(c)
			return
		case res.err != nil:
			s.problem(c, res.err, func() { s.refetch(key) })
			return
		}
		ch := res.value
		s.channelHead(c, ch)
		if playlists {
			if len(ch.Playlists) == 0 {
				s.emptyLine(c, "This channel has no playlists.")
				return
			}
			s.playlistGrid(c, "channel-playlists", ch.Playlists)
			return
		}
		s.videoFeedOf(c, key, false, "This channel has no videos.", &loaded[youtube.VideoPage]{done: true, value: ch.Videos})
	})
}

// channelHead is the head of a channel's page: its avatar, what YouTube
// tells of it, and the buttons to subscribe and open it on YouTube.
func (s *shell) channelHead(c *ui.Context, ch youtube.VideoChannel) {
	k := &s.lk
	vs := &s.videos
	subscribed := ch.Subscribed
	if on, ok := vs.subscribed[ch.ID]; ok {
		subscribed = on
	}
	ui.Row(c).Gap(20).AlignItems(ui.Center).Margin(0, 0, 32).Children(func() {
		const size = 88
		frame := ui.Box(c).Size(size, size).Shrink(0).Radius(size / 2).Clip().Background(k.surface3)
		var b *ui.Bitmap
		if len(ch.Avatar) > 0 {
			b = s.art.get(c, artworkURL(ch.Avatar, int(size*s.zoomOf())), s.artPixels(size))
		}
		frame.Children(func() {
			if b != nil {
				ui.Image(c, b).Size(size, size).Fit(ui.Cover)
			}
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(10).Children(func() {
			var facts []string
			for _, x := range []string{ch.Handle, ch.Subscribers, ch.Count} {
				if x != "" {
					facts = append(facts, x)
				}
			}
			if len(facts) > 0 {
				ui.Text(c, strings.Join(facts, " · ")).FontSize(k.bodySize()).TextColor(k.muted)
			}
			ui.Row(c).Gap(12).Children(func() {
				label := "subscribe"
				if subscribed {
					label = "subscribed"
				}
				if s.button(c, "channel-subscribe", label, !subscribed).Disabled(!s.svc.Videos.SignedIn()).Clicked() {
					s.subscribeChannel(ch, !subscribed)
				}
				if s.button(c, "channel-on-youtube", "open on YouTube", false).Clicked() {
					go s.svc.System.OpenExternal(youtubeChannelURL(ch.ID))
				}
				blocked := s.channelBlocked(youtube.Video{Channel: ch.Name})
				if !blocked && s.button(c, "channel-block", "block", false).Tooltip("Keep its videos off every page").Clicked() {
					s.blockChannel(youtube.Video{Channel: ch.Name})
				}
			})
		})
	})
}

// subscribeChannel subscribes the account to ch, or unsubscribes it, and
// takes it back where YouTube refuses.
func (s *shell) subscribeChannel(ch youtube.VideoChannel, on bool) {
	vs := &s.videos
	if vs.subscribed == nil {
		vs.subscribed = map[string]bool{}
	}
	vs.subscribed[ch.ID] = on
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := s.svc.Videos.Subscribe(ctx, ch.ID, on)
		s.update(func() {
			if err != nil {
				vs.subscribed[ch.ID] = !on
				s.notes = append(s.notes, noticeOf(err))
				return
			}
			s.refetch("videos:subscriptions")
			if on {
				s.notes = append(s.notes, notice("Subscribed to "+ch.Name))
			} else {
				s.notes = append(s.notes, notice("Unsubscribed from "+ch.Name))
			}
		})
	}()
}
