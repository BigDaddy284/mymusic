package main

// The Stats page of the Video side: how long the user watched, day by
// day, and the videos and channels they watched most.

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/app"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// watchedFor writes seconds of watching as "2 hr 5 min", "12 min" or,
// under a minute, "40 sec".
func watchedFor(seconds float64) string {
	m := int(seconds / 60)
	switch {
	case m >= 60:
		return fmt.Sprintf("%d hr %d min", m/60, m%60)
	case m > 0:
		return fmt.Sprintf("%d min", m)
	}
	return fmt.Sprintf("%d sec", int(seconds))
}

// videoStatsPage shows what the user watched on the Video side.
func (s *shell) videoStatsPage(c *ui.Context) {
	k := &s.lk
	labels := make([]string, len(statsRanges))
	for i, r := range statsRanges {
		labels[i] = r.label
	}
	s.pageScroll(c, func() {
		s.pageHeader(c, "Stats", labels, &s.statsRange)
		st := s.svc.Videos.Stats(statsRanges[s.statsRange].days)
		if st.Seconds < 1 {
			ui.Text(c, k.words("Nothing watched yet")).FontSize(k.bigSize()).FontWeight(k.displayWeight()).Font(k.display).TextColor(k.accent)
			ui.Text(c, "Time counts while a video plays, and a video counts as watched after thirty seconds.").FontSize(k.subtitleSize()).FontWeight(k.leadWeight())
			return
		}
		size := k.subtitleSize()
		ui.RichText(c,
			ui.Span{Text: "You watched ", Size: size, Weight: 300},
			ui.Span{Text: plural(st.Views, "video"), Size: size, Weight: 400},
			ui.Span{Text: " for ", Size: size, Weight: 300},
			ui.Span{Text: watchedFor(st.Seconds), Size: size, Weight: 400, Color: k.accent},
			ui.Span{Text: ".", Size: size, Weight: 300},
		).MaxWidth(760).Margin(0, 0, 48)
		s.dayChart(c, st.Daily, "watched")
		ui.Row(c).Gap(48).Wrap().AlignItems(ui.Start).Margin(0, 0, 48).Children(func() {
			ui.Column(c).Grow(1).Basis(420).MinWidth(0).Children(func() {
				s.sectionTitle(c, "Top videos")
				list := make([]youtube.Video, len(st.TopVideos))
				for i, w := range st.TopVideos {
					list[i] = w.Video
				}
				for i, w := range st.TopVideos {
					s.topVideoRow(c, fmt.Sprintf("top-video-%d", i), i+1, w, list)
				}
			})
			ui.Column(c).Width(300).Children(func() {
				s.sectionTitle(c, "Top channels")
				for i, ch := range st.TopChannels {
					ui.Row(c.Key(fmt.Sprintf("top-channel-%d", i))).Gap(12).Padding(8, 0).Children(func() {
						ui.Text(c, fmt.Sprint(i+1)).FontSize(k.smallSize()).TextColor(k.muted).Width(20).FontFeatures("tnum")
						ui.Text(c, ch.Name).FontSize(k.smallSize()).SingleLine().Grow(1).MinWidth(0)
						ui.Text(c, watchedFor(ch.Seconds)).FontSize(k.captionSize()).TextColor(k.muted).FontFeatures("tnum")
					})
				}
			})
		})
	})
}

// topVideoRow is a video of the top videos, with how long it was
// watched; it plays the video, with the others after it.
func (s *shell) topVideoRow(c *ui.Context, key string, rank int, w app.VideoWatch, list []youtube.Video) {
	k := &s.lk
	v := w.Video
	row := s.buttonBase(c.Key(key)).Gap(12).Padding(8).Radius(k.radiusSm).Label(v.Title).Cursor(ui.CursorPointer)
	if row.Hovered() {
		row.Background(k.hover)
	}
	if row.Clicked() {
		s.watch(v, list)
	}
	row.Children(func() {
		ui.Text(c, fmt.Sprint(rank)).FontSize(k.smallSize()).TextColor(k.muted).Width(20).FontFeatures("tnum")
		ui.Box(c).Size(71, 40).Shrink(0).Radius(4).Clip().Children(func() { s.videoThumb(c, key+"-thumb", v, 71, 40) })
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			ui.Text(c, v.Title).FontSize(k.smallSize()).SingleLine()
			if v.Channel != "" {
				ui.Text(c, v.Channel).FontSize(k.captionSize()).TextColor(k.muted).SingleLine()
			}
		})
		ui.Text(c, watchedFor(w.Seconds)).FontSize(k.captionSize()).TextColor(k.muted).FontFeatures("tnum")
	})
}
