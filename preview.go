package main

// The previews of the Video side's cards: the pointer resting on a card
// plays its video in place of the thumbnail, muted, in one web view the
// cards share (public/preview.html), whose speaker button turns the sound
// on. A click on it opens the video; the pointer leaving ends it.

import (
	"net/url"
	"strconv"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

const (
	// previewAfter is how long the pointer rests on a card before its
	// video plays.
	previewAfter = 500 * time.Millisecond
	// previewGrace keeps the preview while the pointer goes from the card
	// to the web view over it, which the page tells late.
	previewGrace = 450 * time.Millisecond
	// previewPoll is how often the window asks the preview's page.
	previewPoll = 200 * time.Millisecond
)

// videoPreview is the preview of the card the pointer rests on.
type videoPreview struct {
	web    *mygo.WebView
	making bool
	stop   chan struct{}
	// key is the card's, video its video and list the list it's in;
	// since is when the pointer came on it, and seen when it, or the
	// page, last told the pointer was there.
	key   string
	video youtube.Video
	list  []youtube.Video
	since time.Time
	seen  time.Time
	// loaded is the video the page was told to play, and playing the one
	// it plays, which shows over the thumbnail.
	loaded, playing string
}

// previewState is what the preview's page tells.
type previewState struct {
	Hover   bool   `json:"hover"`
	Open    bool   `json:"open"`
	Playing bool   `json:"playing"`
	ID      string `json:"id"`
}

// previewsHeld tells whether the cards don't preview: while a video plays.
func (s *shell) previewsHeld() bool {
	vp := &s.videos.player
	return vp.current.ID != "" && (vp.state.Playing || vp.casting != nil)
}

// cardPreview plays the card's video over its thumbnail, w × h, once the
// pointer rested on it, hovered telling whether it's on the card now.
func (s *shell) cardPreview(c *ui.Context, key string, v youtube.Video, list []youtube.Video, hovered bool, w, h float32) {
	pv := &s.videos.preview
	now := c.Now()
	if hovered {
		if pv.key != key {
			pv.key, pv.video, pv.list, pv.since, pv.playing = key, v, list, now, ""
		}
		pv.seen = now
	}
	if pv.key != key || s.previewsHeld() {
		return
	}
	if !hovered {
		if left := previewGrace - now.Sub(pv.seen); left > 0 {
			c.After(left)
		} else {
			return
		}
	}
	if wait := previewAfter - now.Sub(pv.since); wait > 0 {
		c.After(wait)
		return
	}
	s.loadPreview(v.ID)
	if pv.web != nil && pv.playing == v.ID {
		ui.WebView(c.Key(key+"-preview"), pv.web).Absolute().Top(0).Left(0).Size(w, h).Label("Preview: " + v.Title)
	}
}

// loadPreview has the preview's page play the video id, making the web
// view the first time.
func (s *shell) loadPreview(id string) {
	pv := &s.videos.preview
	if pv.loaded == id || pv.making || s.win == nil || s.videoProblem() != "" {
		return
	}
	if pv.web == nil {
		pv.making = true
		win := s.win
		go func() {
			web, err := win.NewWebView(mygo.WebViewOptions{URL: "/preview.html?v=" + url.QueryEscape(id), BackgroundColor: "#000000"})
			s.update(func() {
				pv.making = false
				if err != nil {
					return
				}
				pv.web, pv.loaded = web, id
				pv.stop = make(chan struct{})
				go s.pollPreview(web, pv.stop)
			})
		}()
		return
	}
	pv.loaded = id
	web := pv.web
	go web.Page().Eval("window.mtPreviewLoad && mtPreviewLoad(" + jsString(id) + ")")
}

// pollPreview asks the preview's page what the pointer does, until stop.
func (s *shell) pollPreview(web *mygo.WebView, stop chan struct{}) {
	t := time.NewTicker(previewPoll)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
		}
		if web.IsDestroyed() {
			return
		}
		st, err := mygo.EvalAs[previewState](web.Page(), "return window.mtPreviewState ? mtPreviewState() : {}")
		if err != nil {
			continue
		}
		s.update(func() {
			pv := &s.videos.preview
			if pv.web != web || pv.key == "" {
				return
			}
			if st.Playing && st.ID == pv.loaded {
				pv.playing = st.ID
			}
			if st.Hover {
				pv.seen = time.Now()
			}
			if st.Open {
				v, list := pv.video, pv.list
				s.endPreview()
				s.watch(v, list)
			}
		})
	}
}

// followPreview ends the preview once the pointer left its card, as the
// frame built ends: no card took it.
func (s *shell) followPreview() {
	pv := &s.videos.preview
	if pv.key == "" {
		return
	}
	if time.Since(pv.seen) > previewGrace || s.previewsHeld() || !s.mode.video {
		s.endPreview()
	}
}

// endPreview stops the preview.
func (s *shell) endPreview() {
	pv := &s.videos.preview
	pv.key, pv.loaded, pv.playing, pv.list = "", "", "", nil
	if web := pv.web; web != nil {
		go web.Page().Eval("window.mtPreviewStop && mtPreviewStop()")
	}
}

// jsString writes s as a string of JavaScript.
func jsString(s string) string { return strconv.Quote(s) }
