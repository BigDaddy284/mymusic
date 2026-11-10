package app

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/egoist/mygo"

	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/player"
	"github.com/bigdaddy284/mymusic/internal/stream"
)

// PlaybackMode is how a track plays.
type PlaybackMode string

const (
	// PlaybackStream plays the track in MyTube's audio engine, from Url.
	PlaybackStream PlaybackMode = "stream"
)

// PlaybackSource tells the player how to play a track.
type PlaybackSource struct {
	Mode PlaybackMode `json:"mode"`
	// URL is the audio of a stream.
	URL string `json:"url,omitempty"`
	// SurroundURL is the 5.1 surround audio of a stream, when YouTube has
	// some, which plays in its place while the settings ask for 5.1.
	SurroundURL string `json:"surroundUrl,omitempty"`
	// Track is the track that plays.
	Track media.Track `json:"track"`
	// Cached is set when the audio plays from the cache, without the
	// network.
	Cached bool `json:"cached,omitempty"`
	// Live is set for the stream of a live video, as it goes on, which the
	// engine plays from near its live edge: it has no length, can't seek,
	// and isn't cached.
	Live bool `json:"live,omitempty"`
}

// Player prepares tracks for the player of the main window, and relays
// the player's state and commands between windows.
type Player struct{ a *App }

// Prepare finds the audio of a track to play in the engine, found ahead
// so that errors such as YouTube refusing it show before playback starts.
// It keeps how loud YouTube measured the track, or asks, for the engine to
// play it at an even loudness: a live stream has no measure.
// It tells the 5.1 surround audio of the track too, when YouTube has some.
// The downloads are told whether the video it found is live, as one
// doesn't download.
//
// A track downloaded plays from its download, without the network, also
// while it answers, in stereo: downloads hold no 5.1 audio.
func (p *Player) Prepare(ctx context.Context, track media.Track) (PlaybackSource, error) {
	if track.Source == media.SourceFile {
		return p.a.prepareFile(track)
	}
	if k, ok := p.a.downloads.Lookup(track.ID); ok {
		p.learnLoudness(k.Video)
		return PlaybackSource{Mode: PlaybackStream, URL: streamURL(k.Video), Track: k.Track, Cached: true}, nil
	}
	yt, err := onYouTube(track)
	if err != nil {
		return PlaybackSource{}, err
	}
	_, videoID, err := media.ParseID(yt.ID)
	if err != nil {
		return PlaybackSource{}, err
	}
	src := PlaybackSource{Mode: PlaybackStream, URL: streamURL(videoID), Track: yt}
	surroundURL := schemeURL(streamScheme) + strings.TrimPrefix(stream.SurroundPath(videoID), "/")
	offline := p.a.offline()
	// The stereo audio of its recording may be downloaded for another
	// track.
	stereo := p.a.downloads.Snapshot().HasVideo(videoID) || p.a.cache.Has(videoID)
	if cachedSurround := p.a.cache.Has(stream.SurroundKey(videoID)); cachedSurround || stereo {
		surround := p.a.settings.Get().Surround
		if cachedSurround {
			src.SurroundURL = surroundURL
		}
		// With 5.1 on, a track held in stereo plays its 5.1 audio when
		// YouTube has some, which takes asking it while the network
		// answers, for up to surroundWait: else it plays what is held,
		// and the answer, which comes on, is there for its next play.
		if surround && !cachedSurround && !offline {
			ask, cancel := context.WithTimeout(ctx, surroundWait)
			if s, err := p.a.resolveStream(ask, videoID); err == nil && s.Surround != nil {
				src.SurroundURL = surroundURL
			}
			cancel()
		}
		// It plays without the network when the audio it plays is held.
		if surround && src.SurroundURL != "" {
			src.Cached = cachedSurround
		} else {
			src.Cached = stereo
		}
		// A track cached before MyTube kept measures gets one now.
		p.learnLoudness(videoID)
		return src, nil
	}
	// A stream found before plays without the network, as the engine may
	// hold its audio already; yt-dlp needs the network to find one.
	if p.a.resolver.Fresh(videoID) {
		live := p.a.resolver.Live(videoID)
		p.a.downloads.SetLive(videoID, live)
		if live {
			return liveSource(src, videoID), nil
		}
		// Resolve gives what it found at once.
		if s, err := p.a.resolver.Resolve(ctx, videoID); err == nil && s.Surround != nil {
			src.SurroundURL = surroundURL
		}
		// Its measure was kept as it was found, else it is asked for now.
		p.learnLoudness(videoID)
		return src, nil
	}
	if offline {
		return PlaybackSource{}, player.ErrOffline
	}
	s, err := p.a.resolveStream(ctx, videoID)
	if err != nil {
		return PlaybackSource{}, err
	}
	p.a.downloads.SetLive(videoID, s.Live)
	if s.Live {
		return liveSource(src, videoID), nil
	}
	if s.Surround != nil {
		src.SurroundURL = surroundURL
	}
	if s.Measured {
		p.a.loudness.keep(videoID, s.Loudness)
	} else {
		// yt-dlp tells no loudness: YouTube's player API does, while
		// YouTube holds yt-dlp's stream back.
		p.a.loudness.learn(videoID)
	}
	return src, nil
}

// surroundWait bounds how long a track cached in stereo waits, with 5.1
// on, for YouTube to tell whether it has 5.1 audio.
var surroundWait = 3 * time.Second

// liveSource is src as the live stream of a video, which the engine opens
// at its own address.
func liveSource(src PlaybackSource, videoID string) PlaybackSource {
	src.Live = true
	src.URL = schemeURL(streamScheme) + strings.TrimPrefix(stream.LivePath(videoID), "/")
	return src
}

// learnLoudness asks YouTube how loud a video MyTube holds the audio of
// is, unless the book knows, or the network is gone: the engine would wait
// for an answer that can't come.
func (p *Player) learnLoudness(videoID string) {
	if p.a.online != nil && p.a.online.Offline() {
		return
	}
	p.a.loudness.learn(videoID)
}

// StreamFailure tells why the audio of a track, the one Prepare returned,
// failed to load: in plain words, or "" when its audio arrived,
// and the engine could not play it.
func (p *Player) StreamFailure(track media.Track) string {
	_, videoID, err := media.ParseID(track.ID)
	if err != nil {
		return ""
	}
	return p.a.proxy.Failure(videoID)
}

// streamURL is the address of a video's audio, which the engine opens.
func streamURL(videoID string) string {
	return schemeURL(streamScheme) + strings.TrimPrefix(stream.Path(videoID), "/")
}

// Report tells the app what the player plays, for the tray, the main
// window's title, notifications, the system's media controls and
// Last.fm.
func (p *Player) Report(np NowPlaying) {
	p.a.mu.Lock()
	prev := p.a.nowPlaying
	p.a.nowPlaying = np
	started := p.a.startOf(np)
	tray := p.a.tray
	bar := p.a.taskbar
	p.a.mu.Unlock()
	s := p.a.settings.Get()
	p.a.scrobble(np)
	p.a.keepPosition(np)
	p.a.reportWatchtime(np, s.YouTubeHistory)
	p.a.reportControls(np)
	if bar != nil {
		bar.Show(np.plays(), np.Track != nil)
	}

	if tip := trayTip(np); tray != nil && tip != trayTip(prev) {
		tray.SetToolTip(tip)
	}
	if windowTitle(np, s.ShowInTitle) != windowTitle(prev, s.ShowInTitle) {
		p.a.retitle()
	}
	if started && s.Notifications && !p.a.inFront() {
		go p.a.notifyTrack(*np.Track)
	}
}

// startOf reports whether np tells that a track started: it sounds, no
// longer loading, and it isn't the track whose start was told last. A
// track is reported loading first, then playing, so it is the second
// report that tells it started; paused and played again, it didn't start
// again. Nothing playing at all forgets the track, which starts again
// when played again. a.mu guards it.
func (a *App) startOf(np NowPlaying) bool {
	switch {
	case np.Track == nil:
		a.started = ""
		return false
	case !np.Playing || np.Loading || np.Track.ID == a.started:
		return false
	}
	a.started = np.Track.ID
	return true
}

// command has the player do what the tray, the media keys, the system's
// media controls or a deep link ask.
func (a *App) command(cmd PlayerCommand) {
	if fn := a.commandHook.Load(); fn != nil && (*fn)(cmd) {
		return
	}
	if p := a.goPlayer.Load(); p != nil {
		a.commandGoPlayer(p, cmd)
	}
}

// OnCommand has fn see each command of the tray, the media keys, the
// system's media controls and the global shortcuts first, until cancel,
// on any goroutine: a command fn takes, returning true, doesn't reach the
// player, as the Video side's player takes them while a video plays.
func (a *App) OnCommand(fn func(PlayerCommand) bool) (cancel func()) {
	f := &fn
	a.commandHook.Store(f)
	return func() { a.commandHook.CompareAndSwap(f, nil) }
}

// trayTip is the tray's tooltip for what plays. Windows cuts tooltips at
// 127 characters.
func trayTip(np NowPlaying) string { return cut(playingLine(np), 120) }

// windowTitle is the main window's title: what plays, unless the settings
// keep it to MyTube's name.
func windowTitle(np NowPlaying, show bool) string {
	if !show {
		return "MyTube"
	}
	return cut(playingLine(np), 256)
}

// playingLine names what plays, as the tray's tooltip and the main
// window's title show it: "One More Time - Daft Punk", after "Paused: "
// while it's paused, or "MyTube" with nothing to play.
func playingLine(np NowPlaying) string {
	if np.Track == nil {
		return "MyTube"
	}
	line := oneLine(np.Track.Title)
	if line == "" {
		return "MyTube"
	}
	if artists := oneLine(np.Track.ArtistNames()); artists != "" {
		line += " - " + artists
	}
	if !np.plays() {
		line = "Paused: " + line
	}
	return line
}

// oneLine is s as a title shows it, on one line: valid UTF-8, with
// control characters, line breaks among them, as
// spaces, and no runs of spaces. The marks that turn the direction of the
// text after them go, so that a title can't turn its artists around.
func oneLine(s string) string {
	s = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsControl(r):
			return ' '
		// The embeddings and overrides, then the isolates.
		case r >= '\u202a' && r <= '\u202e', r >= '\u2066' && r <= '\u2069':
			return -1
		}
		return r
	}, strings.ToValidUTF8(s, "�"))
	return strings.Join(strings.Fields(s), " ")
}

// cut cuts s to n characters at most, the last of them an ellipsis.
func cut(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}

// showNote shows a notification; tests, which run no MyGo loop for it to
// wait for, replace it.
var showNote = (*mygo.Notification).Show

// notifyTrack shows a notification of the track that started, with its
// cover once it loaded, unless another track started meanwhile
// (trackNote).
func (a *App) notifyTrack(t media.Track) {
	opts, ok := a.trackNote(t)
	if !ok {
		return
	}
	n := mygo.NewNotification(opts)
	n.OnClick(a.showMain)
	if err := showNote(n); err != nil && !errors.Is(err, mygo.ErrNotificationsDenied) {
		log.Printf("mytube: notification: %v", err)
	}
}

// notifyLiked tells that the track that plays joined the liked songs, or
// left them, while neither the main window nor the mini player is in
// front: liked from elsewhere, as by a global shortcut, nothing else
// shows it.
func (a *App) notifyLiked(t media.Track, liked bool) {
	title, body := likedNote(t, liked)
	a.notifyRated("liked", title, body)
}

// notifyDisliked tells that the track that plays was disliked, or no
// longer, while neither the main window nor the mini player is in front,
// as notifyLiked does.
func (a *App) notifyDisliked(t media.Track, disliked bool) {
	title, body := dislikedNote(t, disliked)
	a.notifyRated("disliked", title, body)
}

// notifyRated shows the notification id of a track liked or disliked from
// outside MyTube's windows, unless one that shows it is in front.
func (a *App) notifyRated(id, title, body string) {
	if !a.settings.Get().Notifications || a.inFront() {
		return
	}
	n := mygo.NewNotification(mygo.NotificationOptions{ID: id, Title: title, Body: body, Silent: true})
	n.OnClick(a.showMain)
	if err := showNote(n); err != nil && !errors.Is(err, mygo.ErrNotificationsDenied) {
		log.Printf("mytube: notification: %v", err)
	}
}

// likedNote is what the notification of a track liked, or no longer,
// says.
func likedNote(t media.Track, liked bool) (title, body string) {
	title = "Removed from liked songs"
	if liked {
		title = "Added to liked songs"
	}
	return title, noteBody(t)
}

// dislikedNote is what the notification of a track disliked, or no
// longer, says.
func dislikedNote(t media.Track, disliked bool) (title, body string) {
	title = "Removed the dislike"
	if disliked {
		title = "Disliked"
	}
	return title, noteBody(t)
}

// noteBody names a track in a notification: its title and artists.
func noteBody(t media.Track) string {
	body := t.Title
	if artists := t.ArtistNames(); artists != "" {
		body += " - " + artists
	}
	return body
}
