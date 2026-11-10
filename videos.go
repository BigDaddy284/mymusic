package main

// The Video side's pages: YouTube's own videos, beside the music. Signed in
// to YouTube Music, the same Google session reads the home feed, the
// subscriptions, the watch history and the playlists of YouTube's website.
// A search finds any video, and Random searches for the title of a random
// article of Wikipedia. Videos play in watch.go.

import (
	"context"
	"log"
	"net/url"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/app"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// videoHomeChips are the filters of the Video side's home page.
var videoHomeChips = []string{"All", "Subscriptions", "Continue watching"}

// videoLibraryTabs are the tabs of the Video side's library: their paths
// and labels.
var videoLibraryTabs = []struct{ id, label string }{
	{"subscriptions", "Subscriptions"},
	{"history", "History"},
	{"playlists", "Playlists"},
	{"watch-later", "Watch later"},
	{"liked", "Liked videos"},
}

// videoTrending are the topics of the Trending page, each a search of this
// week's videos: YouTube no longer has a trending page of its own.
var videoTrending = []string{"Music", "Gaming", "News", "Sports", "Comedy", "Film trailers", "Science", "Live"}

// videosState is the Video side's own state.
type videosState struct {
	// query is what the search box holds.
	query string
	// homeChip and trendingChip are the filters chosen.
	homeChip, trendingChip int
	// known are the videos the pages showed, by ID, which the watch page
	// tells of; titles the playlists' titles, by ID.
	known  map[string]youtube.Video
	titles map[string]string
	// more holds the pages after each list's first, by its key, and
	// moreToken the token of the page after those, once asked for.
	more      map[string][]youtube.Video
	moreToken map[string]string
	moreBusy  map[string]bool
	// rolling is set while Random asks Wikipedia for a topic.
	rolling bool
	// randomShown are the videos random searches showed, by video: the
	// search that showed it first, the one that may show it again.
	randomShown map[string]string
	// hidden are the videos, and the channels by videoChannelKey, the user
	// told YouTube not to recommend, left off the cards (videomenu.go).
	hidden map[string]bool
	// ratings are the account's ratings of videos the user set here, by
	// video: "like", "dislike", or "" (videorating.go).
	ratings map[string]string
	// The player (watch.go), and the preview of the card the pointer
	// rests on (preview.go).
	player  videoPlayer
	preview videoPreview
}

func videoSearchPath(q string) string    { return "/videos/search?q=" + url.QueryEscape(q) }
func videoPlaylistPath(id string) string { return "/videos/playlist/" + url.PathEscape(id) }
func watchPath(id string) string         { return "/watch/" + url.PathEscape(id) }
func youtubeWatchURL(id string) string {
	return "https://www.youtube.com/watch?v=" + url.QueryEscape(id)
}
func youtubePlaylistURL(id string) string {
	return "https://www.youtube.com/playlist?list=" + url.QueryEscape(id)
}

// videoHeader is a page's header, with the chips under it that filter it,
// if any.
func (s *shell) videoHeader(c *ui.Context, title string, chips []string, selected *int) {
	k := &s.lk
	s.pageHeader(c, title, nil, nil)
	if len(chips) == 0 {
		return
	}
	ui.Row(c).Wrap().Gap(8).Margin(0, 0, 28).Children(func() {
		for i, label := range chips {
			chip := s.buttonBase(c.Key("chip-" + count(i))).Height(34).PaddingX(14).Radius(k.round(34)).Label(label).Cursor(ui.CursorPointer)
			ink := k.text
			switch {
			case i == *selected:
				ink = k.fillAccent(chip, stateOf(chip))
			case chip.Hovered():
				chip.Background(k.surface3)
			default:
				chip.Background(k.surface2)
			}
			if chip.Clicked() {
				*selected = i
			}
			chip.Children(func() { ui.Text(c, label).FontSize(k.smallSize()).TextColor(ink).SingleLine() })
		}
	})
}

// videosHome is the Video side's home: what YouTube recommends, the
// subscriptions' latest videos, or the history, to carry on watching.
func (s *shell) videosHome(c *ui.Context) {
	vs := &s.videos
	s.pageScroll(c, func() {
		if !s.svc.Videos.SignedIn() {
			s.pageHeader(c, "Home", nil, nil)
			s.videosSignedOut(c)
			return
		}
		s.videoHeader(c, "Home", videoHomeChips, &vs.homeChip)
		switch vs.homeChip {
		case 1:
			s.videoFeed(c, "videos:subscriptions", false, "No new videos from your subscriptions.", s.svc.Videos.Subscriptions)
		case 2:
			s.videoFeed(c, "videos:history", false, "Nothing in your watch history.", s.svc.Videos.History)
		default:
			s.videoFeed(c, "videos:home", false, "YouTube recommends nothing yet: watch a few videos, or press Random.", s.svc.Videos.Home)
		}
	})
}

// videosSignedOut invites the user to sign in, which the feeds need; the
// search and Random work without.
func (s *shell) videosSignedOut(c *ui.Context) {
	k := &s.lk
	ui.Column(c).Gap(16).Margin(0, 0, 48).Children(func() {
		ui.Text(c, k.words("Watch anything")).FontSize(k.bigSize()).FontWeight(k.displayWeight()).Font(k.display).TextColor(k.accent)
		ui.Text(c, "Search YouTube for any video, or press Random. Sign in to YouTube Music to see your home feed, subscriptions, history and playlists here too: the same Google account signs in to both.").
			FontSize(k.subtitleSize()).FontWeight(k.leadWeight()).MaxWidth(640)
		ui.Row(c).Gap(12).Children(func() {
			if s.button(c, "videos-random", "random", true).Disabled(s.videos.rolling).Clicked() {
				s.randomVideoSearch()
			}
			if s.button(c, "videos-sign-in", "sign in", false).Clicked() {
				s.navigate("/settings/accounts")
			}
		})
	})
}

// randomVideoSearch asks Wikipedia for a random article and searches the
// videos for its title, or for a word of the app's own list when Wikipedia
// doesn't answer.
func (s *shell) randomVideoSearch() {
	vs := &s.videos
	if vs.rolling {
		return
	}
	vs.rolling = true
	go func() {
		topic, wiki := s.svc.Videos.RandomTopic(context.Background())
		s.update(func() {
			vs.rolling = false
			vs.query = topic
			from := "word"
			if wiki {
				from = "wikipedia"
			}
			s.navigate(videoSearchPath(topic) + "&random=" + from)
		})
	}()
}

// videoSearchPage shows the videos YouTube finds for q. A random search
// that finds none searches again for a word of the app's own list.
func (s *shell) videoSearchPage(c *ui.Context, q, random string) {
	k := &s.lk
	q = strings.TrimSpace(q)
	s.pageScroll(c, func() {
		if q == "" {
			s.pageHeader(c, "Search", nil, nil)
			ui.Column(c).Gap(16).Margin(0, 0, 48).Children(func() {
				ui.Text(c, k.words("Find any video")).FontSize(k.bigSize()).FontWeight(k.displayWeight()).Font(k.display).TextColor(k.accent)
				ui.Text(c, "Type in the search box above, or press Random to search for a random Wikipedia article.").
					FontSize(k.subtitleSize()).FontWeight(k.leadWeight()).MaxWidth(640)
				ui.Row(c).Gap(12).Children(func() {
					if s.button(c, "search-start", "start typing", true).Clicked() {
						s.searchBox.Focus()
					}
					if s.button(c, "search-random", "random", false).Disabled(s.videos.rolling).Clicked() {
						s.randomVideoSearch()
					}
				})
			})
			return
		}
		s.pageHeader(c, q, nil, nil)
		switch random {
		case "wikipedia":
			ui.RichText(c).FontSize(k.smallSize()).TextColor(k.muted).Margin(0, 0, 24).Children(func() {
				ui.Text(c, "A random Wikipedia article: ")
				s.linkOf(c, q, "https://en.wikipedia.org/wiki/"+url.PathEscape(strings.ReplaceAll(q, " ", "_")))
				ui.Text(c, ". If it finds nothing, Random picks a word from its own list instead.")
			})
		case "word":
			ui.Text(c, "Random picked a word from its own list.").FontSize(k.smallSize()).TextColor(k.muted).Margin(0, 0, 24)
		}
		key := "videos:search:" + q
		res := fetch(s, key, func(ctx context.Context) (youtube.VideoPage, error) { return s.svc.Videos.Search(ctx, q) })
		if random != "" && res.done && res.err == nil {
			res = s.freshRandomVideos(q, res)
		}
		if res.done && res.err == nil && len(res.value.Videos) == 0 && random == "wikipedia" {
			// An obscure title finds nothing: a word of the list will.
			word := randomWord()
			s.videos.query = word
			s.router.Replace(videoSearchPath(word) + "&random=word")
			return
		}
		s.videoFeedOf(c, key, true, "No videos match.", res)
	})
}

// freshRandomVideos leaves out of the videos a random search for q found
// those an earlier random search showed, as YouTube answers many an
// obscure search with the same popular videos.
func (s *shell) freshRandomVideos(q string, res *loaded[youtube.VideoPage]) *loaded[youtube.VideoPage] {
	vs := &s.videos
	if vs.randomShown == nil {
		vs.randomShown = map[string]string{}
	}
	page := res.value
	fresh := make([]youtube.Video, 0, len(page.Videos))
	for _, v := range page.Videos {
		if first, ok := vs.randomShown[v.ID]; ok && first != q {
			continue
		}
		vs.randomShown[v.ID] = q
		fresh = append(fresh, v)
	}
	if len(fresh) == len(page.Videos) {
		return res
	}
	page.Videos = fresh
	out := *res
	out.value = page
	return &out
}

// randomWord returns a word of the app's own list, replaced in tests.
var randomWord = app.RandomWord

// videoTrendingPage shows this week's videos of a topic.
func (s *shell) videoTrendingPage(c *ui.Context) {
	vs := &s.videos
	s.pageScroll(c, func() {
		s.videoHeader(c, "Trending", videoTrending, &vs.trendingChip)
		topic := videoTrending[min(max(vs.trendingChip, 0), len(videoTrending)-1)]
		key := "videos:trending:" + topic
		s.videoFeed(c, key, true, "Nothing new this week.", func(ctx context.Context) (youtube.VideoPage, error) {
			return s.svc.Videos.Recent(ctx, topic)
		})
	})
}

// videoLibraryPage is the Video side's library, at tab: one of
// videoLibraryTabs, or none, which shows the first.
func (s *shell) videoLibraryPage(c *ui.Context, tab string) {
	selected := 0
	labels := make([]string, len(videoLibraryTabs))
	for i, t := range videoLibraryTabs {
		labels[i] = t.label
		if t.id == tab {
			selected = i
		}
	}
	s.pageScroll(c, func() {
		if !s.svc.Videos.SignedIn() {
			s.pageHeader(c, "Library", nil, nil)
			s.videosSignedOut(c)
			return
		}
		before := selected
		s.pageHeader(c, "Library", labels, &selected)
		if selected != before {
			s.router.Replace("/videos/library/" + videoLibraryTabs[selected].id)
		}
		switch videoLibraryTabs[selected].id {
		case "subscriptions":
			s.videoFeed(c, "videos:subscriptions", false, "No new videos from your subscriptions.", s.svc.Videos.Subscriptions)
		case "history":
			s.videoFeed(c, "videos:history", false, "Nothing in your watch history.", s.svc.Videos.History)
		case "playlists":
			s.videoPlaylists(c)
		case "watch-later":
			s.playlistFeed(c, "WL", "Nothing saved to watch later.")
		case "liked":
			s.playlistFeed(c, "LL", "No liked videos yet.")
		}
	})
}

// videoPlaylistPage shows the videos of a playlist of YouTube's website.
func (s *shell) videoPlaylistPage(c *ui.Context, id string) {
	title := s.videos.titles[id]
	switch {
	case id == "WL":
		title = "Watch later"
	case id == "LL":
		title = "Liked videos"
	case title == "":
		title = "Playlist"
	}
	s.pageScroll(c, func() {
		s.pageHeader(c, title, nil, nil)
		ui.Row(c).Gap(12).Margin(0, 0, 32).Children(func() {
			if s.button(c, "playlist-on-youtube", "open on YouTube", false).Clicked() {
				go s.svc.System.OpenExternal(youtubePlaylistURL(id))
			}
		})
		s.playlistFeed(c, id, "This playlist has no videos.")
	})
}

// playlistFeed shows the videos of the playlist id.
func (s *shell) playlistFeed(c *ui.Context, id, empty string) {
	s.videoFeed(c, "videos:playlist:"+id, false, empty, func(ctx context.Context) (youtube.VideoPage, error) {
		return s.svc.Videos.Playlist(ctx, id)
	})
}

// videoPlaylistsLoad loads the account's playlists of YouTube's website,
// which the library and the sidebar show.
func (s *shell) videoPlaylistsLoad() *loaded[[]youtube.VideoPlaylist] {
	return fetch(s, "videos:playlists", func(ctx context.Context) ([]youtube.VideoPlaylist, error) { return s.svc.Videos.Playlists(ctx) })
}

// videoFeed shows the videos get finds, in a grid, kept by key.
func (s *shell) videoFeed(c *ui.Context, key string, search bool, empty string, get func(ctx context.Context) (youtube.VideoPage, error)) {
	s.videoFeedOf(c, key, search, empty, fetch(s, key, get))
}

// videoFeedOf shows the videos of res, and the pages after it, which a
// button below asks for.
func (s *shell) videoFeedOf(c *ui.Context, key string, search bool, empty string, res *loaded[youtube.VideoPage]) {
	vs := &s.videos
	switch {
	case !res.done:
		s.loading(c)
		return
	case res.err != nil:
		s.problem(c, res.err, func() { s.refetch(key) })
		return
	}
	videos := s.shownVideos(append(append([]youtube.Video{}, res.value.Videos...), vs.more[key]...))
	if len(videos) == 0 {
		s.emptyLine(c, empty)
		return
	}
	s.videoGrid(c, key, videos)
	token, asked := vs.moreToken[key]
	if !asked {
		token = res.value.Continue
	}
	if token == "" {
		return
	}
	label := "show more"
	if vs.moreBusy[key] {
		label = "loading…"
	}
	if s.button(c, key+"-more", label, false).Disabled(vs.moreBusy[key]).Margin(24, 0, 0).AlignSelf(ui.Start).Clicked() {
		s.moreVideos(key, search, token)
	}
}

// moreVideos asks for the page of videos after token, for the list key.
func (s *shell) moreVideos(key string, search bool, token string) {
	vs := &s.videos
	if vs.moreBusy == nil {
		vs.moreBusy, vs.more, vs.moreToken = map[string]bool{}, map[string][]youtube.Video{}, map[string]string{}
	}
	vs.moreBusy[key] = true
	go func() {
		page, err := s.svc.Videos.More(context.Background(), search, token)
		s.update(func() {
			vs.moreBusy[key] = false
			if err != nil {
				log.Printf("mytube: more videos of %s: %v", key, err)
				s.actionToasts = append(s.actionToasts, actionToast{"More videos couldn't load.", "retry", func() { s.moreVideos(key, search, token) }})
				return
			}
			vs.more[key] = append(vs.more[key], page.Videos...)
			vs.moreToken[key] = page.Continue
		})
	}()
}

// videoColumns returns how many cards of videos fit in a row of the pane,
// and their width.
func (s *shell) videoColumns(c *ui.Context) (int, float32) {
	k := &s.lk
	pane := s.paneWidth(c) - 2*k.gutter
	gap := k.gap
	const least = 240
	n := max(1, int((pane+gap)/(least+gap)))
	return n, (pane - gap*float32(n-1)) / float32(n)
}

// videoGrid lays the videos as cards, in rows that fill the pane.
func (s *shell) videoGrid(c *ui.Context, key string, videos []youtube.Video) {
	k := &s.lk
	_, w := s.videoColumns(c)
	ui.Row(c.Key(key + "-grid")).Wrap().GapX(k.gap).GapY(1.5 * k.gap).AlignItems(ui.Start).Children(func() {
		for i, v := range videos {
			s.videoCard(c, key+"-"+count(i), v, w, videos)
		}
	})
}

// videoCard is a video's thumbnail, w wide, with its length over it, its
// title, its channel, and how many saw it, when. Clicked, it plays, with
// list as up next.
func (s *shell) videoCard(c *ui.Context, key string, v youtube.Video, w float32, list []youtube.Video) {
	k := &s.lk
	h := float32(int(w*9/16 + 0.5))
	card := s.buttonBase(c.Key(key)).Width(w).Shrink(0).Column().AlignItems(ui.Start).Gap(8).Label(v.Title).Cursor(ui.CursorPointer).
		ContextMenu(func(m *ui.Menu) { s.videoMenu(m, v, list) })
	if card.Clicked() {
		s.watch(v, list)
	}
	hovered := card.Hovered()
	card.Children(func() {
		ui.Box(c).Size(w, h).Shrink(0).Radius(k.radiusSm).Clip().Children(func() {
			s.videoThumb(c, key+"-thumb", v, w, h)
			s.watchedBar(c, v, w, 4)
			s.cardPreview(c, key, v, list, hovered, w, h)
			if hovered {
				ui.Box(c).Absolute().Top(0).Left(0).Size(w, h).InsetShadow(0, 0, 0, 3, k.accent).Radius(k.radiusSm)
			}
			if v.Duration != "" {
				ui.Box(c).Absolute().Right(6).Bottom(6).PaddingX(6).PaddingY(2).Radius(4).Background(ui.RGBA(0, 0, 0, 0.8)).Children(func() {
					ui.Text(c, v.Duration).FontSize(k.captionSize()).TextColor(ui.RGBA(255, 255, 255, 1)).SingleLine()
				})
			}
		})
		ui.Row(c).Width(w).Gap(4).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
				k.labelWeight(ui.Text(c, v.Title).FontSize(k.bodySize()).SingleLine().MinWidth(0))
				if v.Channel != "" {
					ui.Text(c, v.Channel).FontSize(k.smallSize()).TextColor(k.muted).SingleLine().MinWidth(0)
				}
				if line := videoFacts(v); line != "" {
					ui.Text(c, line).FontSize(k.smallSize()).TextColor(k.muted).SingleLine().MinWidth(0)
				}
			})
			s.videoMenuButton(c, key, v, list)
		})
	})
}

// youtubeRed is the red of YouTube's bar of what was watched.
var youtubeRed = ui.RGBA(255, 0, 51, 1)

// watchedShare is how much of v was watched, from 0 to 1: here, as it was
// left, or, when it wasn't, on any device, as YouTube tells.
func (s *shell) watchedShare(v youtube.Video) float64 {
	if share := s.svc.Videos.WatchedShare(v.ID); share > 0 {
		return share
	}
	return min(float64(v.Watched)/100, 1)
}

// watchedBar is the red bar along the bottom of v's thumbnail, w wide and
// h high, that tells how much of it was watched, as YouTube's does.
func (s *shell) watchedBar(c *ui.Context, v youtube.Video, w, h float32) {
	share := s.watchedShare(v)
	if share <= 0 {
		return
	}
	ui.Box(c).Absolute().Left(0).Bottom(0).Size(w, h).Background(ui.RGBA(255, 255, 255, 0.35)).Children(func() {
		ui.Box(c).Size(float32(int(float64(w)*share+0.5)), h).Background(youtubeRed)
	})
}

// videoFacts tells how many saw a video, and when it came out.
func videoFacts(v youtube.Video) string {
	var parts []string
	for _, p := range []string{v.Views, v.Published} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, " · ")
}

// videoThumb is a video's thumbnail, filling w × h.
func (s *shell) videoThumb(c *ui.Context, key string, v youtube.Video, w, h float32) {
	frame := ui.Box(c.Key(key)).Size(w, h).Shrink(0)
	var b *ui.Bitmap
	if s.inView(c, frame) {
		b = s.art.get(c, artworkURL(v.Thumbnail, int(w*s.zoomOf())), s.artPixels(w))
	}
	if b == nil {
		frame.Background(s.lk.surface3)
	}
	frame.Children(func() {
		if b != nil {
			ui.Image(c, b).Size(w, h).Fit(ui.Cover)
		}
	})
}

// videoPlaylists shows the account's playlists of YouTube's website.
func (s *shell) videoPlaylists(c *ui.Context) {
	k := &s.lk
	res := s.videoPlaylistsLoad()
	switch {
	case !res.done:
		s.loading(c)
		return
	case res.err != nil:
		s.problem(c, res.err, func() { s.refetch("videos:playlists") })
		return
	}
	_, w := s.videoColumns(c)
	h := float32(int(w*9/16 + 0.5))
	ui.Row(c.Key("video-playlists")).Wrap().GapX(k.gap).GapY(1.5 * k.gap).AlignItems(ui.Start).Children(func() {
		for i, p := range res.value {
			card := s.buttonBase(c.Key("video-playlist-" + count(i))).Width(w).Shrink(0).Column().AlignItems(ui.Start).Gap(8).
				Label(p.Title).Cursor(ui.CursorPointer)
			if card.Clicked() {
				if s.videos.titles == nil {
					s.videos.titles = map[string]string{}
				}
				s.videos.titles[p.ID] = p.Title
				s.navigate(videoPlaylistPath(p.ID))
			}
			hovered := card.Hovered()
			card.Children(func() {
				frame := ui.Box(c).Size(w, h).Shrink(0).Radius(k.radiusSm).Clip()
				var b *ui.Bitmap
				if len(p.Thumbnail) > 0 && s.inView(c, frame) {
					b = s.art.get(c, artworkURL(p.Thumbnail, int(w*s.zoomOf())), s.artPixels(w))
				}
				if b == nil {
					frame.Background(k.surface3)
				}
				frame.Children(func() {
					if b != nil {
						ui.Image(c, b).Size(w, h).Fit(ui.Cover)
					} else {
						glyph := "video"
						if p.ID == "LL" {
							glyph = "heart"
						}
						ui.Box(c).Size(w, h).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
							ui.Icon(c, icon(glyph)).FontSize(36).TextColor(k.muted)
						})
					}
					if hovered {
						ui.Box(c).Absolute().Top(0).Left(0).Size(w, h).InsetShadow(0, 0, 0, 3, k.accent).Radius(k.radiusSm)
					}
				})
				ui.Column(c).Width(w).Gap(2).Children(func() {
					k.labelWeight(ui.Text(c, p.Title).FontSize(k.bodySize()).SingleLine())
					var parts []string
					for _, x := range []string{p.Owner, p.Count} {
						if x != "" {
							parts = append(parts, x)
						}
					}
					if len(parts) > 0 {
						ui.Text(c, strings.Join(parts, " · ")).FontSize(k.smallSize()).TextColor(k.muted).SingleLine()
					}
				})
			})
		}
	})
}

// videoSearchField is the top bar's search box on the Video side, which
// searches YouTube's videos on Enter, and the Random button beside it.
func (s *shell) videoSearchField(c *ui.Context) {
	k := &s.lk
	vs := &s.videos
	box := ui.Row(c).Height(36).MaxWidth(520).Grow(1).Shrink(1).MinWidth(160).Gap(8).PaddingX(12).
		Radius(k.round(36)).Background(k.bg).Cursor(ui.CursorText)
	if box.Clicked() {
		s.searchBox.Focus()
	}
	focused := s.searchBox.Focused(c)
	ink := k.muted
	switch {
	case focused:
		box.Border(2, k.accent)
		ink = k.text
	case k.tonal:
		box.Background(k.surface2)
	default:
		box.Border(2, k.outline)
	}
	box.Children(func() {
		ui.Icon(c, icon("search")).FontSize(18).TextColor(ink)
		in := ui.TextInputBase(c.Key("video-search"), &vs.query).Placeholder(k.words("Search YouTube videos")).Grow(1).MinWidth(0).
			Bind(&s.searchBox).Label("Search YouTube videos").FontSize(k.bodySize())
		if vs.query != "" {
			in.OnShortcut(0, ui.KeyEscape, func() { vs.query = "" })
		}
		if in.Submitted() {
			if q := strings.TrimSpace(vs.query); q != "" {
				s.navigate(videoSearchPath(q))
			}
		}
	})
	label := "random"
	if vs.rolling {
		label = "finding…"
	}
	// Random is the height of the search box beside it.
	if s.button(c, "top-random", label, true).Height(36).PaddingX(20).Radius(s.lk.round(36)).AlignSelf(ui.Center).Disabled(vs.rolling).Tooltip("Search for a random Wikipedia article").Clicked() {
		s.randomVideoSearch()
	}
}
