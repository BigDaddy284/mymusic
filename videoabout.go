package main

// The Description panel of the Video side: what YouTube's watch page shows
// under the video that plays, its views, date and description. Its links
// work: those to YouTube's videos, channels, playlists and hashtags open
// here, its times seek the video, and the rest open in the browser.

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// videoAboutPanel shows the description of the video that plays.
func (s *shell) videoAboutPanel(c *ui.Context) {
	k := &s.lk
	vp := &s.videos.player
	v := vp.current
	s.panelHead(c, "Description", func() {
		if v.ID != "" && s.iconButton(c, "about-copy", "link", "Copy the link", false, 16).Clicked() {
			copyText(youtubeWatchURL(v.ID))
			s.notes = append(s.notes, notice("Copied the link to "+pick(v.Title != "", v.Title, "the video")))
		}
	})
	if v.ID == "" {
		s.emptyIn(c, "Nothing plays yet. Pick a video, and its description shows here.")
		return
	}
	d := s.videoDetails(v.ID)
	s.fillFromDetails(v.ID, d)
	a := d.value.About
	ui.Scroll(c.Key("about-"+v.ID)).Grow(1).MinHeight(0).Padding(0, 12, 24).Children(func() {
		title := pick(a.Title != "", a.Title, v.Title)
		ui.Text(c, title).FontSize(k.bodySize()).FontWeight(600).Padding(0, 12, 4)
		ch := v
		if ch.ChannelID == "" && a.ChannelID != "" {
			ch.Channel, ch.ChannelID = a.Channel, a.ChannelID
		}
		if ch.Channel != "" {
			ui.Box(c).Padding(0, 12, 12).Children(func() { s.channelLink(c, "about-channel", ch, k.smallSize()) })
		}
		if !d.done {
			s.emptyIn(c, "Looking for the description…")
			return
		}
		if d.err != nil {
			s.emptyIn(c, d.err.Error())
			ui.Row(c).Padding(0, 24, 8).Children(func() {
				if s.button(c, "about-retry", "try again", false).Clicked() {
					s.forget("videos:details:" + v.ID)
				}
			})
			return
		}
		// The description sits in a box, as on YouTube, its views and date
		// at its top.
		ui.Column(c.Key("about-box")).Padding(12).Gap(2).Radius(k.radiusSm).Background(k.surface2).Children(func() {
			views := pick(a.Views != "", a.Views, d.value.Views)
			ui.RichText(c).FontSize(k.smallSize()).FontWeight(600).Padding(0, 0, 8).Children(func() {
				sep := ""
				for _, f := range []string{views, a.Date} {
					if f != "" {
						ui.Text(c, sep+f)
						sep = "  "
					}
				}
				for _, r := range a.Tags {
					ui.Text(c, sep)
					s.aboutRun(c, r)
					sep = " "
				}
			})
			if len(a.Description) == 0 {
				ui.Text(c, "This video has no description.").FontSize(k.smallSize()).TextColor(k.muted)
				return
			}
			for _, line := range a.Description {
				if blankLine(line) {
					ui.Box(c).Height(10)
					continue
				}
				ui.RichText(c).FontSize(k.smallSize()).Children(func() {
					for _, r := range line {
						s.aboutRun(c, r)
					}
				})
			}
		})
	})
	if s.showsMini() {
		w := s.panelWidth(c) - 32
		ui.Box(c).Padding(12, 16, 16).Shrink(0).BorderWidth(1, 0, 0, 0).BorderColor(k.border).Children(func() {
			s.videoMini(c, "mini-about", w)
		})
	}
}

// fillFromDetails gives the video id, which plays, the title and channel
// its watch page d tells, where it lacks them: a video opened by a link
// knows only its ID, and the link's text, which YouTube cuts short.
func (s *shell) fillFromDetails(id string, d *loaded[youtube.VideoDetails]) {
	vp := &s.videos.player
	a := d.value.About
	if !d.done || d.err != nil || id == "" || vp.current.ID != id {
		return
	}
	t := vp.current.Title
	if a.Title != "" && (t == "" || strings.HasSuffix(t, "...") || strings.HasSuffix(t, "…")) {
		vp.current.Title = a.Title
	}
	if vp.current.ChannelID == "" && a.ChannelID != "" {
		vp.current.Channel, vp.current.ChannelID = a.Channel, a.ChannelID
	}
}

// aboutRun is a run of a line of the description: a link that opens what
// it links, a time that seeks the video, or plain text.
func (s *shell) aboutRun(c *ui.Context, r media.TextRun) {
	switch {
	case r.Seek:
		if s.linkOf(c, r.Text, "").Label("Play from " + r.Text).Tooltip("Play from " + r.Text).Clicked() {
			s.seekVideo(float64(r.AtMs) / 1000)
		}
	case r.Link != "":
		tip := r.Link
		if _, ok := youtubeLinkOf(r.Link); ok {
			tip = "Open in MyTube: " + r.Link
		}
		if s.linkOf(c, r.Text, "").Tooltip(tip).Clicked() {
			s.openAboutLink(r)
		}
	default:
		ui.Text(c, r.Text)
	}
}

// seekVideo has the video that plays go to at seconds.
func (s *shell) seekVideo(at float64) {
	vp := &s.videos.player
	if vp.current.ID == "" || vp.web == nil && vp.casting == nil {
		return
	}
	s.videoCommand("seek", at)
	vp.state.Time = at
}

// openAboutLink opens the link of the run r: what is YouTube's opens here,
// and the rest in the browser.
func (s *shell) openAboutLink(r media.TextRun) {
	l, ok := youtubeLinkOf(r.Link)
	if !ok {
		if u := media.WebLink(r.Link); u != "" {
			openWebPage(s, u)
		}
		return
	}
	vp := &s.videos.player
	switch {
	case l.video != "" && l.video == vp.current.ID:
		if l.at > 0 {
			s.seekVideo(l.at)
		}
		s.navigate(watchPath(l.video))
	case l.video != "":
		s.watch(youtube.Video{ID: l.video, Title: linkTitle(r.Text)}, nil)
	case l.channel != "":
		s.openChannel(youtube.Video{ChannelID: l.channel, Channel: strings.TrimSpace(r.Text)})
	case l.playlist != "":
		s.navigate(videoPlaylistPath(l.playlist))
	case l.search != "":
		s.navigate(videoSearchPath(l.search))
	}
}

// youtubeLink is what a link to YouTube opens here: a video, from a time
// in it, a channel, a playlist, or a search, as a hashtag's.
type youtubeLink struct {
	video, channel, playlist, search string
	at                               float64
}

// youtubeLinkOf reads a link to YouTube that opens here, of youtube.com,
// its music and mobile sites, or youtu.be: ok is false for another.
func youtubeLinkOf(raw string) (youtubeLink, bool) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "http" && u.Scheme != "https" {
		return youtubeLink{}, false
	}
	host := strings.ToLower(u.Hostname())
	for _, p := range []string{"www.", "m.", "music."} {
		host = strings.TrimPrefix(host, p)
	}
	q := u.Query()
	at := linkSeconds(q.Get("t"))
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	switch {
	case host == "youtu.be" && len(parts) == 1 && parts[0] != "":
		return youtubeLink{video: parts[0], at: at}, true
	case host != "youtube.com":
		return youtubeLink{}, false
	case parts[0] == "watch" && q.Get("v") != "":
		return youtubeLink{video: q.Get("v"), at: at}, true
	case (parts[0] == "shorts" || parts[0] == "live" || parts[0] == "embed") && len(parts) > 1 && parts[1] != "":
		return youtubeLink{video: parts[1], at: at}, true
	case parts[0] == "channel" && len(parts) > 1 && strings.HasPrefix(parts[1], "UC"):
		return youtubeLink{channel: parts[1]}, true
	case parts[0] == "playlist" && q.Get("list") != "":
		return youtubeLink{playlist: q.Get("list")}, true
	case parts[0] == "hashtag" && len(parts) > 1 && parts[1] != "":
		return youtubeLink{search: "#" + parts[1]}, true
	case parts[0] == "results" && q.Get("search_query") != "":
		return youtubeLink{search: q.Get("search_query")}, true
	}
	return youtubeLink{}, false
}

// linkSeconds reads the time of a link: "90", "90s" or "1h2m3s".
func linkSeconds(t string) float64 {
	if t == "" {
		return 0
	}
	if n, err := strconv.Atoi(strings.TrimSuffix(t, "s")); err == nil {
		return float64(max(n, 0))
	}
	if d, err := time.ParseDuration(t); err == nil && d > 0 {
		return d.Seconds()
	}
	return 0
}

// linkTitle is the title a link's text gives the video it opens until
// its watch page tells it: "• Its title", as YouTube writes a link to a
// video, but none for an address.
func linkTitle(text string) string {
	t := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(text), "•·▶ "))
	if strings.Contains(t, "://") || strings.HasPrefix(t, "youtu") || strings.HasPrefix(t, "www.") {
		return ""
	}
	return t
}
