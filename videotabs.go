package main

// Video tabs: videos kept along the top bar of the Video side, beside
// Random, to come back to later. A tab plays its video from where it was
// left; its × closes it. The + button keeps the video that plays, and a
// card's menu keeps any video. The tabs are kept in the settings, from one
// start to the next.

import (
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// tabVideo is the video of a tab.
func tabVideo(t settings.VideoTab) youtube.Video {
	return youtube.Video{ID: t.ID, Title: t.Title, Channel: t.Channel, ChannelID: t.ChannelID, Duration: t.Duration,
		Thumbnail: []media.Image{{URL: "https://i.ytimg.com/vi/" + t.ID + "/mqdefault.jpg", Width: 320, Height: 180}}}
}

// hasVideoTab reports whether a tab keeps the video id.
func (s *shell) hasVideoTab(id string) bool {
	return slices.ContainsFunc(s.set.VideoTabs, func(t settings.VideoTab) bool { return t.ID == id })
}

// addVideoTab keeps v in a tab, at the end, unless one keeps it already.
func (s *shell) addVideoTab(v youtube.Video) {
	if v.ID == "" || s.hasVideoTab(v.ID) {
		return
	}
	if len(s.set.VideoTabs) >= settings.MaxVideoTabs {
		s.notes = append(s.notes, notice("You have the most tabs the bar keeps: close one first"))
		return
	}
	s.set.VideoTabs = append(s.set.VideoTabs, settings.VideoTab{ID: v.ID, Title: v.Title, Channel: v.Channel, ChannelID: v.ChannelID, Duration: v.Duration})
	s.save()
}

// closeVideoTab closes the tab of the video id.
func (s *shell) closeVideoTab(id string) {
	s.set.VideoTabs = slices.DeleteFunc(slices.Clone(s.set.VideoTabs), func(t settings.VideoTab) bool { return t.ID == id })
	s.save()
}

// videoTabs is the row of tabs in the top bar, with the + button that
// keeps the video that plays. It takes the room the bar has left, and
// its tabs narrow to fit it.
func (s *shell) videoTabs(c *ui.Context) {
	k := &s.lk
	vp := &s.videos.player
	tabs := s.set.VideoTabs
	row := ui.Row(c.Key("video-tabs")).Grow(1.5).Basis(0).MinWidth(0).Gap(tabGap).AlignItems(ui.Center).Clip()
	// As many tabs as fit at their least, the rest in the menu of a button
	// after them.
	room := row.Bounds().W - 2*(34+tabGap)
	fit := len(tabs)
	if room > 0 {
		fit = min(fit, max(0, int((room+tabGap)/(tabLeast+tabGap))))
	}
	row.Children(func() {
		var list []youtube.Video
		for _, t := range tabs {
			list = append(list, tabVideo(t))
		}
		for i, t := range tabs[:fit] {
			s.videoTab(c, "video-tab-"+count(i), t, list, t.ID == vp.current.ID)
		}
		if rest := tabs[fit:]; len(rest) > 0 {
			more := s.buttonBase(c.Key("video-tabs-more")).Height(34).PaddingX(10).Radius(k.radiusSm).Shrink(0).
				Background(k.surface2).Label("More tabs").Cursor(ui.CursorPointer)
			more.Menu(func(m *ui.Menu) {
				for _, t := range rest {
					if m.Item(t.Title).Checked(t.ID == vp.current.ID).Chosen() {
						s.watch(tabVideo(t), list)
					}
				}
			})
			more.Children(func() { ui.Text(c, "+"+count(len(rest))).FontSize(k.smallSize()).SingleLine() })
		}
		cur := vp.current
		label := "Keep the video that plays in a tab"
		if cur.ID == "" {
			label = "Play a video to keep it in a tab"
		}
		add := s.iconButton(c, "video-tab-add", "plus", label, false, 18).Shrink(0).Disabled(cur.ID == "" || s.hasVideoTab(cur.ID))
		if add.Clicked() {
			s.addVideoTab(cur)
		}
	})
}

// The least width of a tab, and the gap between two.
const (
	tabLeast = 130
	tabGap   = 6
)

// videoTab is the tab of t: its thumbnail and title, and a close button.
// Clicked, it plays its video, with the tabs as up next.
func (s *shell) videoTab(c *ui.Context, key string, t settings.VideoTab, list []youtube.Video, current bool) {
	k := &s.lk
	v := tabVideo(t)
	tab := s.buttonBase(c.Key(key)).Height(34).PaddingX(4).Gap(6).Radius(k.radiusSm).AlignItems(ui.Center).Justify(ui.Start).
		Width(200).MinWidth(tabLeast).Shrink(1).Label(t.Title).Tooltip(t.Title).Cursor(ui.CursorPointer).
		ContextMenu(func(m *ui.Menu) {
			if m.Item("Play").Chosen() {
				s.watch(v, list)
			}
			if m.Item("Copy the link").Chosen() {
				copyText(youtubeWatchURL(t.ID))
				s.notes = append(s.notes, notice("Copied the link to "+t.Title))
			}
			m.Separator()
			if m.Item("Close the tab").Chosen() {
				s.closeVideoTab(t.ID)
			}
			if m.Item("Close the other tabs").Disabled(len(s.set.VideoTabs) < 2).Chosen() {
				s.set.VideoTabs = slices.DeleteFunc(slices.Clone(s.set.VideoTabs), func(x settings.VideoTab) bool { return x.ID != t.ID })
				s.save()
			}
		})
	hovered := tab.Hovered()
	switch {
	case current:
		tab.Background(k.selected)
	case hovered:
		tab.Background(k.hover)
	default:
		tab.Background(k.surface2)
	}
	if tab.Clicked() && !current {
		s.watch(v, list)
	}
	tab.Children(func() {
		ui.Box(c).Size(40, 23).Shrink(0).Radius(4).Clip().Children(func() {
			s.videoThumb(c, key+"-thumb", v, 40, 23)
			s.watchedBar(c, v, 40, 2)
		})
		ui.Text(c, t.Title).FontSize(k.captionSize()).SingleLine().Grow(1).Basis(0).MinWidth(0).TextAlign(ui.Start)
		if s.iconButton(c, key+"-close", "close", "Close the tab", false, 14).Shrink(0).Clicked() {
			s.closeVideoTab(t.ID)
		}
	})
}
