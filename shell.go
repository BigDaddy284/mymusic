package main

import (
	"net/url"
	"strconv"

	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/egoist/mygo/ui"
)

// view builds the main window: the rail, the top bar over the main pane
// and the side panel, and the player bar across the bottom.
func (s *shell) view(c *ui.Context) {
	s.st = s.p.State()
	s.startPage()
	s.kept = s.svc.Downloads.Snapshot()
	s.applyZoom(c)
	// A window that turns narrow or wide drops the fold the user gave
	// its rail.
	if w, _ := c.Size(); (w <= narrowWidth) != s.narrow {
		s.narrow, s.narrowFold = w <= narrowWidth, nil
	}
	s.followArtwork()
	s.keepYouTubePlaylists()
	s.keepYouTubeSongs()
	c.SetTheme(&s.lk.theme)
	if s.fullAppView(c) {
		s.toasts(c)
		s.followWatch()
		return
	}
	s.toasts(c)
	s.fileToast(c)
	s.runShortcuts(c)
	s.followOnTop()
	k := &s.lk
	root := ui.Column(c).Fill().Background(k.bg).TextColor(k.text).Font(k.body)
	s.ground(c, root)
	root.Children(func() {
		ui.Row(c).Grow(1).MinHeight(0).AlignItems(ui.Stretch).Children(func() {
			s.rail(c)
			s.railResizer(c)
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				s.topBar(c)
				ui.Row(c).Grow(1).MinHeight(0).AlignItems(ui.Stretch).Children(func() {
					s.mainPane(c)
					if s.panel != "" {
						s.panelResizer(c)
						// The panel's lyrics refresh with the page.
						s.collecting = true
						s.sidePanel(c)
						s.collecting = false
					}
				})
			})
		})
		s.playerBar(c)
		s.overlay(c)
	})
	s.palette(c)
	s.coverDialog(c)
	s.questionDialog(c)
	s.trackInfoDialog(c)
	s.promptDialog(c)
	s.failureDialog(c)
	s.sweepDialogs()
	s.followVideo()
	s.followWatch()
	s.keepShown(c)
}

// mainPane shows the page of the router's path. Pages scroll themselves,
// as pageScroll does, or as lists, which build only the rows in view.
func (s *shell) mainPane(c *ui.Context) {
	ui.Column(c.Key("main-pane")).Grow(1).MinWidth(0).AlignItems(ui.Stretch).Children(func() {
		s.router.View(c, func(r *ui.Route) { s.page(c, r) }).Grow(1).MinHeight(0)
		// The video that plays goes on, small, as the Video side's other
		// pages show (watch.go).
		s.floatingMini(c)
		// And the music that plays, as the Video side's pages show
		// (musicmini.go).
		s.musicMini(c)
	})
}

// page builds the page of a route.
func (s *shell) page(c *ui.Context, r *ui.Route) {
	s.pagePath = r.Path()
	if r.Path() == s.router.Path() {
		s.followMode(r.Path())
	}
	// The page shown, not one that slides away, is what a refresh asks
	// again for.
	s.collecting = r.Path() == s.router.Path()
	switch {
	case r.Match("/"):
		r.Title("Home")
		s.homePage(c)
	case r.Match("/search"):
		r.Title("Search")
		s.searchPage(c, r.Query("q"))
	case r.Match("/album/{id}"):
		s.albumPage(c, r.Param("id"))
	case r.Match("/artist/{id}"):
		s.artistPage(c, r.Param("id"))
	case r.Match("/playlist/{id}"):
		s.playlistPage(c, r.Param("id"))
	case r.Match("/podcast/{id}"):
		s.podcastPage(c, r.Param("id"))
	case r.Match("/library"):
		s.libraryPage(c, "")
	case r.Match("/library/{tab}"):
		s.libraryPage(c, r.Param("tab"))
	case r.Match("/now-playing"):
		s.nowPlayingPage(c)
	case r.Match("/settings"):
		s.settingsPage(c, "", r.Query("q"))
	case r.Match("/settings/{tab}"):
		s.settingsPage(c, r.Param("tab"), r.Query("q"))
	case r.Match("/explore"):
		s.explorePage(c, "")
	case r.Match("/explore/{tab}"):
		s.explorePage(c, r.Param("tab"))
	case r.Match("/mood/{params}"):
		s.moodPage(c, r.Param("params"))
	case r.Match("/videos"):
		r.Title("Home")
		s.videosHome(c)
	case r.Match("/videos/search"):
		r.Title("Search")
		s.videoSearchPage(c, r.Query("q"), r.Query("random"))
	case r.Match("/videos/trending"):
		r.Title("Trending")
		s.videoTrendingPage(c)
	case r.Match("/videos/playlist/{id}"):
		r.Title("Playlist")
		s.videoPlaylistPage(c, r.Param("id"))
	case r.Match("/videos/channel/{id}"):
		r.Title("Channel")
		s.videoChannelPage(c, r.Param("id"), r.Query("tab"))
	case r.Match("/watch/{id}"):
		r.Title("Video")
		s.watchPage(c, r.Param("id"))
	case r.Match("/stats"):
		s.statsPage(c)
	default:
		s.lost(c, "This page isn't one of MyTube's.")
	}
	s.collecting = false
}

// lost stands for a page there is none of, saying why.
func (s *shell) lost(c *ui.Context, why string) {
	k := &s.lk
	s.pageScroll(c, func() {
		s.pageHeader(c, "Lost", nil, nil)
		ui.Text(c, k.words("Nothing lives here")).FontSize(k.bigSize()).FontWeight(k.displayWeight()).Font(k.display).TextColor(k.accent)
		ui.Text(c, why).FontSize(k.subtitleSize()).FontWeight(k.leadWeight()).Margin(0, 0, 24)
		if s.button(c, "go-home", "go home", true).Clicked() {
			s.navigate("/")
		}
	})
}

// toasts shows what the player told, as toasts above the player bar.
func (s *shell) toasts(c *ui.Context) {
	for _, n := range s.notes {
		switch {
		case n.Settings:
			c.ToastAction(n.Message, "settings", func() { s.navigate("/settings/playback") })
		case n.Failure != 0:
			s.failureToast(c, n)
		case n.Error:
			c.AddToast(ui.Toast{Title: n.Message, Type: "error"})
		default:
			c.Toast(n.Message)
		}
	}
	s.notes = s.notes[:0]
	for _, a := range s.actionToasts {
		c.ToastAction(a.message, a.action, a.run)
	}
	s.actionToasts = s.actionToasts[:0]
}

// actionToast is a toast with an action, as "undo" once a track is
// disliked.
type actionToast struct {
	message, action string
	run             func()
}

// showMini shows the mini player.
func (s *shell) showMini() { go s.svc.System.ShowMiniPlayer() }

// Paths of the pages of things.

func albumPath(id string) string    { return "/album/" + url.PathEscape(id) }
func artistPath(id string) string   { return "/artist/" + url.PathEscape(id) }
func playlistPath(id string) string { return "/playlist/" + url.PathEscape(id) }
func searchPath(q string) string    { return "/search?q=" + url.QueryEscape(q) }

// itemPath is the page of a shelf's item, or "" for a track.
func itemPath(it media.Item) string {
	switch {
	case it.Album != nil:
		return albumPath(it.Album.ID)
	case it.Artist != nil:
		return artistPath(it.Artist.ID)
	case it.Playlist != nil:
		return playlistPath(it.Playlist.ID)
	case it.Podcast != nil:
		return podcastPath(it.Podcast.ID)
	}
	return ""
}

// count writes n, small, as a pivot's count.
func count(n int) string { return strconv.Itoa(n) }
