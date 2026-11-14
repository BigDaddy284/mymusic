package main

import (
	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// homeFilters are the home page's pivots, and the sources they keep. With
// YouTube Music the only service, Everything is it and the user's own.
var homeFilters = []struct {
	label  string
	source media.Source
}{{"Everything", ""}, {"Yours", media.SourceLocal}}

// networkOnline tells whether the network answers, as last checked; tests
// answer for it.
var networkOnline = func(s *shell) bool { return s.svc.Downloads.Online() }

// homePage shows the user's own tiles, then the shelves of every service,
// in the order the user arranged them (homelayout.go), or, while they
// edit it, the list they arrange them in (homeedit.go).
func (s *shell) homePage(c *ui.Context) {
	k := &s.lk
	s.pageScroll(c, func() {
		if s.home.editing {
			s.pageHeader(c, "Home", nil, nil)
			s.homeEditor(c)
			return
		}
		labels := make([]string, len(homeFilters))
		for i, f := range homeFilters {
			labels[i] = f.label
		}
		// Edit home goes beside the pivots.
		s.afterPivots = func() {
			if s.pivotButton(c, "home-edit", "Edit home") {
				s.home.editing = true
			}
		}
		s.pageHeader(c, "Home", labels, &s.homeFilter)
		s.afterPivots = nil
		s.toolProblems(c)

		// The track played last, as the library has it now.
		var last *media.Track
		if recent := s.svc.Collection.History(1); len(recent) > 0 {
			last = &recent[0].Track
		} else if cur, ok := s.st.Current(); ok {
			last = &cur.Track
		}
		s.tileGrid(c, func(columns int) {
			title, caption, label := "Start here", k.words("Search for something to play"), "Search for music"
			var art []media.Image
			span := tileLarge
			if last != nil {
				title, caption, label = "Jump back in", trackLine(*last), "Play "+last.Title+" again"
				art = last.Artwork
				// As a card, its artwork sits below the title: one row
				// holds them, where two left a gap between.
				if k.cards && len(art) > 0 {
					span = tileWide
				}
			}
			if s.flatTile(c, "tile-last", 0, span, columns, title, caption, label, icon("search"), art) {
				if last != nil {
					s.p.PlayTracks([]media.Track{*last}, 0)
				} else {
					s.navigate("/search")
				}
			}
			// The video that plays goes on beside it, as the Music side
			// shows: Liked songs and the mixes are in the sidebar.
			s.homeVideo(c, columns)
		})
		// The user's own shelves show at once, then those of each service
		// once they came: one that still loads holds back those after it,
		// which would move once it came (homePlaces).
		source := homeFilters[s.homeFilter].source
		layout := s.set.HomeShelves
		shown, failed := 0, false
		for _, p := range homePlaces(layout, s.homeGroups(layout), source) {
			g := p.group
			switch {
			case p.shelf != nil:
				shown += s.showHomeShelf(c, p.shelf, source)
			case !g.done:
				s.loading(c)
				return
			case g.added:
				// A shelf the user added that couldn't come says why where
				// it shows.
				key := g.key
				ui.Column(c.Key("home-added-"+key)).Margin(0, 0, k.sections()).Children(func() {
					s.sectionTitle(c, g.title)
					s.problem(c, g.err, func() { s.forget(key) })
				})
			default:
				key := g.key
				s.problem(c, g.err, func() { s.forget(key) })
				// Offline, the downloads still play: they're offered once,
				// under the first service that failed.
				if !failed && !networkOnline(s) && s.kept.Status().Tracks > 0 {
					ui.Row(c).Margin(16, 0, 0).Children(func() {
						if s.button(c, "home-downloads", "open downloads", true).Clicked() {
							s.navigate("/library/downloads")
						}
					})
				}
				failed = true
			}
		}
		if shown == 0 && !failed {
			switch source {
			case media.SourceLocal:
				s.emptyLine(c, "Play some tracks and make a mix: what you collect shows here.")
			default:
				s.emptyLine(c, "Nothing here yet.")
			}
		}
	})
}

// showHomeShelf shows the shelf sh of Home, of the pivot of source, and
// returns how many items it shows.
func (s *shell) showHomeShelf(c *ui.Context, sh *homeShelf, source media.Source) int {
	if sh.build != nil {
		return sh.build(c, sh.ui, source)
	}
	s.itemShelf(c, sh.ui, sh.shelf)
	return len(sh.shelf.Items)
}

// localShelves are the home page's shelves of what the user played and
// made, built again only once the library changed.
func (s *shell) localShelves() []media.Shelf {
	if v := s.svc.Collection.Version(); !s.localBuilt || v != s.localVersion {
		s.local, s.localVersion, s.localBuilt = s.svc.Catalog.LocalShelves(), v, true
	}
	return s.local
}

// pageScroll is a page that scrolls, inside the pane's gutter. The page
// the router shows keeps where it scrolled for the top bar and the rows
// pinned under it (showPage), or, in capture mode, has it scrolled
// (act:scroll): not the one it leaves, which shows beside it while it
// slides in. The page shown may be pulled down past its top to refresh it,
// which opens a strip above its content (pullStrip): the content is keyed,
// so that what it holds keeps its state as the strip comes and goes.
func (s *shell) pageScroll(c *ui.Context, content func()) {
	page := ui.Scroll(c).Grow(1).MinWidth(0).MinHeight(0)
	var track *ui.ScrollState
	if s.pagePath == s.router.Path() && !s.here.built {
		track = s.showPage(page)
	}
	page.Children(func() {
		if track != nil {
			s.pullStrip(c)
		}
		col := ui.Column(c.Key("page-content")).Padding(0, s.lk.gutter, 64)
		if track != nil {
			s.takeContent(page, col)
		}
		col.Children(content)
		if track != nil && s.here.room > 0 {
			col.MinHeight(s.here.room)
		}
	})
	// Once the content took where the page is: it may have moved it.
	if track != nil {
		page.TrackScroll(track).HandleInput(s.pageInput)
		s.followAutoScroll(c, page, track)
	}
}
