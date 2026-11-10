package app

import (
	"context"
	"errors"
	"net/url"
	"runtime"

	"github.com/egoist/mygo"

	"github.com/bigdaddy284/mymusic/aac"
	"github.com/bigdaddy284/mymusic/internal/artwork"
	"github.com/bigdaddy284/mymusic/internal/audio"
	"github.com/bigdaddy284/mymusic/internal/audio/output"
	"github.com/bigdaddy284/mymusic/internal/hct"
	"github.com/bigdaddy284/mymusic/internal/romaji"
	"github.com/bigdaddy284/mymusic/internal/romanize"
	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/internal/tools"
	"github.com/bigdaddy284/mymusic/opus"
)

// Info describes the app.
type Info struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
	Arch     string `json:"arch"`
	Dev      bool   `json:"dev"`
	DataDir  string `json:"dataDir"`
	// NativeAudio is set where MyTube's own audio engine can play: on
	// 64-bit Windows.
	NativeAudio bool `json:"nativeAudio"`
}

// System is what the windows need of the desktop and the app.
type System struct{ a *App }

// Info describes the app.
func (s *System) Info() Info {
	return Info{
		Name:        mygo.App.Name(),
		Version:     s.a.version(),
		Platform:    runtime.GOOS,
		Arch:        runtime.GOARCH,
		Dev:         mygo.IsDev(),
		DataDir:     s.a.opts.DataDir,
		NativeAudio: settings.NativeAudio && output.Available() == nil,
	}
}

// version is the version of MyTube that runs.
func (a *App) version() string {
	if v := mygo.App.Version(); v != "" {
		return v
	}
	return a.opts.Version
}

// Tools describes the tools MyTube runs: yt-dlp and a JavaScript
// runtime.
func (s *System) Tools(ctx context.Context) tools.Status { return s.a.tools.Status(ctx) }

// ArtworkColor returns the accent color of a piece of artwork, as
// "#rrggbb", for themes that follow what plays.
func (s *System) ArtworkColor(ctx context.Context, url string) (string, error) {
	if c, ok := s.a.colors.get(url); ok {
		return c, nil
	}
	c, err := s.a.artwork.Color(ctx, url)
	if err != nil {
		return "", err
	}
	s.a.colors.put(url, c)
	return c, nil
}

// OpenExternal opens a web page in the default browser.
func (s *System) OpenExternal(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return errors.New("only web pages open in the browser")
	}
	return mygo.Shell.OpenExternal(u.String())
}

// Notice is the license of code MyTube includes from another project.
type Notice struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// Notices returns the licenses of the code MyTube includes from other
// projects, which its copies must reproduce.
func (s *System) Notices() []Notice {
	return []Notice{
		{Name: "libopus, which decodes Opus audio", Text: opus.License},
		{Name: "SharpJaad, which decodes AAC audio", Text: aac.License},
		{Name: "go-mp3, which decodes the MP3 files on this computer", Text: audio.MP3License},
		{Name: "jfreymuth/vorbis, which decodes the Vorbis files on this computer", Text: audio.VorbisLicense},
		{Name: "Chromium's Web Audio, which shapes the sound", Text: audio.License},
		{Name: "Material Color Utilities, which make the Material themes' colours", Text: hct.License},
		{Name: "Unihan, whose readings of Chinese characters write Chinese lyrics in pinyin", Text: romanize.UnihanLicense},
		{Name: "pinyin-data and phrase-pinyin-data, the readings of Chinese characters and words in pinyin", Text: romanize.PinyinLicense},
		{Name: "kagome, which reads the kanji of Japanese lyrics for their romaji", Text: romaji.License},
		{Name: "The IPA dictionary (mecab-ipadic), which romaji downloads the first time it's turned on", Text: romaji.DictionaryNotice},
	}
}

// CacheUsage is what the audio cache holds.
type CacheUsage struct {
	Bytes  int64 `json:"bytes"`
	Tracks int   `json:"tracks"`
}

// CacheUsage returns what the audio cache holds.
func (s *System) CacheUsage() CacheUsage {
	b, n := s.a.cache.Usage()
	return CacheUsage{Bytes: b, Tracks: n}
}

// ClearCache empties the audio cache.
func (s *System) ClearCache() error { return s.a.cache.Clear() }

// Artwork downloads the artwork the windows show, and keeps it on disk.
func (s *System) Artwork() *artwork.Store { return s.a.artStore }

// ArtworkUsage is what the artwork kept on disk holds.
type ArtworkUsage struct {
	Bytes  int64 `json:"bytes"`
	Images int   `json:"images"`
}

// ArtworkUsage returns what the artwork kept on disk holds.
func (s *System) ArtworkUsage() ArtworkUsage {
	b, n := s.a.artStore.Usage()
	return ArtworkUsage{Bytes: b, Images: n}
}

// ClearArtwork removes the artwork kept on disk.
func (s *System) ClearArtwork() error { return s.a.artStore.Clear() }

// ShowMiniPlayer shows the mini player, a small window over the others.
func (s *System) ShowMiniPlayer() { s.a.showMini() }

// ShowMainWindow shows the main window and closes the mini player.
func (s *System) ShowMainWindow() {
	s.a.showMain()
	s.a.mu.Lock()
	mini := s.a.mini
	s.a.mu.Unlock()
	if mini != nil {
		mini.Close()
	}
}

// Notify shows a desktop notification, for plugins.
func (s *System) Notify(title, body string) error {
	return mygo.NewNotification(mygo.NotificationOptions{Title: title, Body: body}).Show()
}

// Quit quits MyTube.
func (s *System) Quit() { mygo.App.Quit() }
