package main

// SponsorBlock's segments of the video that plays, which the player skips,
// mutes, or marks on the seek bar, as Settings › Videos asks, and the
// video's chapters, which the seek bar marks too, and a menu of the
// player bar jumps between.

import (
	"context"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/app"
	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// segmentState follows the segments of the video that plays.
type segmentState struct {
	// of is the video they are of; skipped the segments skipped once, by
	// their UUID, which play when the user goes back into them; muting is
	// the one the player muted the video for.
	of      string
	skipped map[string]bool
	muting  string
}

// videoSegments are SponsorBlock's segments of the video id, or none with
// SponsorBlock off.
func (s *shell) videoSegments(id string) []app.Segment {
	if id == "" || s.set.SponsorBlockOff {
		return nil
	}
	l := fetch(s, "videos:segments:"+id, func(ctx context.Context) ([]app.Segment, error) { return s.svc.Videos.Segments(ctx, id) })
	if !l.done || l.err != nil {
		return nil
	}
	return l.value
}

// videoChapters are the chapters of the video id, once the watch page told
// them.
func (s *shell) videoChapters(id string) []youtube.Chapter {
	if id == "" {
		return nil
	}
	d := s.videoDetails(id)
	if !d.done || d.err != nil {
		return nil
	}
	return d.value.Chapters
}

// videoChapterAt is the index of the chapter at t, or -1.
func videoChapterAt(chapters []youtube.Chapter, t float64) int {
	at := -1
	for i, ch := range chapters {
		if ch.Start <= t+0.25 {
			at = i
		}
	}
	return at
}

// segmentName is the name of a category of segments, as Settings shows it.
func segmentName(category string) string {
	for _, k := range settings.SegmentKinds {
		if k.ID == category {
			return k.Name
		}
	}
	return category
}

// followSegments skips the segment the video plays in, once, or mutes the
// video through it, as the settings ask: after each poll of the player.
func (s *shell) followSegments() {
	vp := &s.videos.player
	sg := &vp.segments
	id := vp.current.ID
	if sg.of != id {
		s.endSegmentMute()
		*sg = segmentState{of: id, skipped: map[string]bool{}}
	}
	st := vp.state
	if id == "" || st.Duration <= 0 || vp.settled < 1 {
		return
	}
	muteIn := ""
	for _, seg := range s.videoSegments(id) {
		inside := st.Time >= seg.Start && st.Time < seg.End-0.3
		if !inside || seg.End <= seg.Start {
			continue
		}
		action := s.set.SegmentAction(seg.Category)
		if seg.Mute && action == settings.SegmentSkip {
			action = settings.SegmentMute
		}
		switch action {
		case settings.SegmentSkip:
			if sg.skipped[seg.UUID] || !st.Playing {
				continue
			}
			sg.skipped[seg.UUID] = true
			end := min(seg.End, st.Duration)
			s.videoCommand("seek", end)
			vp.state.Time = end
			back := seg.Start
			s.actionToasts = append(s.actionToasts, actionToast{"Skipped: " + segmentName(seg.Category), "undo", func() {
				if vp.current.ID == id {
					s.videoCommand("seek", back)
					vp.state.Time = back
				}
			}})
			return
		case settings.SegmentMute:
			muteIn = seg.UUID
		}
	}
	switch {
	case muteIn != "" && sg.muting == "":
		sg.muting = muteIn
		s.localVideoCommand("mute", true)
	case muteIn == "" && sg.muting != "":
		s.endSegmentMute()
	}
}

// endSegmentMute gives the video its sound back, after a segment muted it,
// unless the user muted it.
func (s *shell) endSegmentMute() {
	vp := &s.videos.player
	if vp.segments.muting == "" {
		return
	}
	vp.segments.muting = ""
	if !s.set.VideoMuted {
		s.localVideoCommand("mute", false)
	}
}

// segmentColor is the color SponsorBlock gives a category on the seek bar.
func segmentColor(category string) ui.Color {
	switch category {
	case "sponsor":
		return ui.RGBA(0, 212, 0, 1)
	case "selfpromo":
		return ui.RGBA(255, 255, 0, 1)
	case "interaction":
		return ui.RGBA(204, 0, 255, 1)
	case "intro":
		return ui.RGBA(0, 255, 255, 1)
	case "outro":
		return ui.RGBA(2, 2, 237, 1)
	case "preview":
		return ui.RGBA(0, 143, 214, 1)
	case "filler":
		return ui.RGBA(115, 0, 255, 1)
	case "music_offtopic":
		return ui.RGBA(255, 153, 0, 1)
	case "poi_highlight":
		return ui.RGBA(255, 22, 132, 1)
	}
	return ui.RGBA(128, 128, 128, 1)
}

// seekMarks marks the segments and the chapters of the video that plays
// under the track of the seek bar, w wide, of a video duration long: it
// goes before the seek bar, in the box that holds it, so the bar paints
// over it.
func (s *shell) seekMarks(c *ui.Context, w float32, duration float64) {
	vp := &s.videos.player
	id := vp.current.ID
	if id == "" || duration <= 0 || w <= 0 {
		return
	}
	pad := c.Theme().Space(5) / 2
	track := w - 2*pad
	if track <= 0 {
		return
	}
	x := func(t float64) float32 { return pad + float32(min(max(t/duration, 0), 1))*track }
	top := c.Theme().Space(2.5) + 4
	for i, seg := range s.videoSegments(id) {
		if s.set.SegmentAction(seg.Category) == settings.SegmentOff {
			continue
		}
		left, right := x(seg.Start), x(seg.End)
		width := max(right-left, 2)
		if seg.Category == "poi_highlight" {
			left, width = left-2, 4
		}
		ui.Box(c.Key("seek-seg-"+count(i))).Absolute().Left(left).Top(top).Size(width, 3).Radius(1).Background(segmentColor(seg.Category))
	}
	for i, ch := range s.videoChapters(id) {
		if ch.Start <= 0 {
			continue
		}
		ui.Box(c.Key("seek-chapter-"+count(i))).Absolute().Left(x(ch.Start)-1).Top(top-1).Size(2, 5).Background(s.lk.muted)
	}
}

// chaptersMenu lists the chapters of the video that plays, to jump to
// one, in the player bar's menu: none shows without chapters.
func (s *shell) chaptersMenu(m *ui.Menu) {
	vp := &s.videos.player
	chapters := s.videoChapters(vp.current.ID)
	if len(chapters) == 0 {
		return
	}
	now := videoChapterAt(chapters, vp.state.Time)
	m.Submenu("Chapters", func(m *ui.Menu) {
		for i, ch := range chapters {
			if m.Item(clock(time.Duration(ch.Start*float64(time.Second))) + "  " + ch.Title).Checked(i == now).Chosen() {
				s.videoCommand("seek", ch.Start)
				vp.state.Time = ch.Start
			}
		}
	})
}
