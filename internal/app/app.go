// Package app wires MyTube together: the services its windows call, the
// player, and its place on the desktop: windows, the tray, media keys and
// notifications.
package app

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/egoist/mygo"

	"github.com/bigdaddy284/mymusic/internal/addons"
	"github.com/bigdaddy284/mymusic/internal/artwork"
	"github.com/bigdaddy284/mymusic/internal/audio"
	"github.com/bigdaddy284/mymusic/internal/cast"
	"github.com/bigdaddy284/mymusic/internal/downloads"
	"github.com/bigdaddy284/mymusic/internal/library"
	"github.com/bigdaddy284/mymusic/internal/localfiles"
	"github.com/bigdaddy284/mymusic/internal/match"
	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/nowplaying"
	"github.com/bigdaddy284/mymusic/internal/online"
	"github.com/bigdaddy284/mymusic/internal/player"
	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/internal/store"
	"github.com/bigdaddy284/mymusic/internal/stream"
	"github.com/bigdaddy284/mymusic/internal/taskbar"
	"github.com/bigdaddy284/mymusic/internal/tools"
	"github.com/bigdaddy284/mymusic/internal/youtube"
	"github.com/bigdaddy284/mymusic/lyrics"
)

// artworkLimit bounds the artwork kept on disk, in bytes.
const artworkLimit = 256 << 20

// streamScheme is the scheme of the addresses of the audio PrepareStream
// finds, which the audio engine opens, and of the picture of music videos,
// which their page fetches.
const streamScheme = "stream"

// Options configure the app.
type Options struct {
	// DataDir keeps the user's settings, library and add-ons.
	DataDir string
	// CacheDir keeps what can be fetched again: audio, yt-dlp's cache.
	CacheDir string
	// DownloadsDir keeps the downloads, unless the user chose another
	// folder: DataDir's downloads when it's empty.
	DownloadsDir string
	// WatchMusic, when set, watches a folder of music for changes in place
	// of the system, from the first scan on, not only once the app started:
	// tests set it.
	WatchMusic localfiles.WatchFolder
	// Builtin holds the add-ons MyTube ships with.
	Builtin fs.FS
	// TrayIcon is the PNG of the tray icon.
	TrayIcon []byte
	// Version is the app's version, for LRCLIB's user agent.
	Version string
	// LogsDir holds the log and the crash reports.
	LogsDir string
	// LastCrash is the crash report of the run before, when it crashed.
	LastCrash string
	// Interface opens the main window, hidden when asked, and Mini the
	// mini player: windows that MyGo draws, of the root package.
	Interface func(a *App, hidden bool) *mygo.Window
	Mini      func(a *App) *mygo.Window
	// Hotkeys registers the user's global shortcuts: mygo.GlobalShortcut
	// when nil, once the app is ready.
	Hotkeys Hotkeys
	// ProblemWait is how long a settings.json with a mistake waits to be
	// told of, a second when zero: tests shorten it, at once when
	// negative.
	ProblemWait time.Duration
	// CastBrowse, when set, finds the cast devices on the network in place
	// of mDNS: tests set it.
	CastBrowse func(ctx context.Context) ([]cast.Device, error)
	// OutputDevices, when set, lists the audio output devices in place of
	// the system: tests set it.
	OutputDevices func() ([]OutputDevice, error)
	// Controls, when set, reaches the system's media controls in place of
	// nowplaying.Open on MyGo's main thread: tests, which run no event
	// loop, set it.
	Controls func(handle func(nowplaying.Request)) (*nowplaying.Controls, error)
}

// App is MyTube.
type App struct {
	opts Options
	// lastCrash is the crash report of the run before, until the user
	// answers about it.
	lastCrash string

	settings   *store.File[settings.Settings]
	library    *library.Library
	positions  *library.Positions
	videoStats *videoStats
	videoCast  videoCast
	// commandHook sees the player's commands first (OnCommand).
	commandHook  atomic.Pointer[func(PlayerCommand) bool]
	youtube      *youtube.Client
	lastfmStore  lastfmStore
	youtubeStore *youtubeStore
	tools        *tools.Manager
	resolver     *stream.Resolver
	proxy        *stream.Proxy
	online       *online.Watch
	cache        *stream.Cache
	downloads    *downloads.Manager
	// localFiles is the index of the music on the computer, and
	// localScans its scans.
	localFiles   *localfiles.Index
	localScans   localScans
	loudness     *loudnessBook
	matcher      *match.Matcher
	labels       *labelCache
	lyrics       *lyrics.Client
	unison       *lyrics.Unison
	translations *translations
	addons       *addons.Manager
	lyricsMemo   *memo[lyrics.Lyrics]
	related      *memo[[]media.Shelf]
	artwork      *artwork.Client
	artStore     *artwork.Store
	colors       *memo[string]
	// colorFile follows the colour file the settings name.
	colorFile colorFile
	// episodeDetails keeps the details of the episodes the player showed.
	episodeDetails *memo[media.EpisodeDetails]
	// releases keeps what the artists the user follows put out lately.
	releases releaseBook

	// watchtime reports the listening to YouTube Music's episodes to the
	// history there.
	watchtime watchtime
	// videoWatch reports the watching of the Video side's videos to the
	// history of the account on YouTube (videowatch.go).
	videoWatch watchtime
	// native is MyTube's own audio engine, made at its first use.
	native     *audio.Engine
	nativeOnce sync.Once
	// controls shows the engine's track in the system's media controls.
	controls controls
	// goPlayer is the player, once the main window opened, and queue its
	// queue as kept between runs.
	goPlayer atomic.Pointer[player.Player]
	queue    *store.File[player.Saved]
	// cast is the cast devices found, and the one the player plays on.
	cast castState
	// output is the audio output device the engine plays on.
	output outputState
	// videoPage is the page of the main window that plays the picture of
	// music videos, while it is open.
	videoPage atomic.Pointer[mygo.Page]
	// miniOpen is set while the mini player is open, which the main
	// window's frames ask without taking mu.
	miniOpen atomic.Bool
	// onNavigate shows the pages of the main window.
	onNavigate *func(target string)
	// readingProfiles is set while the accounts' profiles are read in the
	// background.
	readingProfiles atomic.Bool
	// youtubeChecking is set while a new YouTube session waits for YouTube
	// Music to take it: the cookies YouTube changes meanwhile aren't kept.
	youtubeChecking atomic.Bool
	// updating is held while an update installs.
	updating sync.Mutex
	// quitting is set once MyTube starts to quit: its windows close then,
	// even with close to the tray on.
	quitting atomic.Bool
	// global holds the user's global shortcuts as registered (#80).
	global globalShortcuts
	// metered is set while the system tells the connection is metered.
	metered atomic.Bool

	mu         sync.Mutex
	lastfm     *scrobbles
	youtubeMe  *youtube.Account
	nowPlaying NowPlaying
	// started is the ID of the track whose start the reports told last
	// (startOf).
	started    string
	main, mini *mygo.Window
	// retitling is set while the main window's title waits to change.
	retitling bool
	tray      *mygo.Tray
	// taskbar is the buttons on the thumbnail of the main window's
	// taskbar button, on Windows.
	taskbar   *taskbar.Bar
	mediaKeys bool
	// profileErrs is why the accounts' profiles could not be read, until
	// they are, or the user signs in or out.
	profileErrs profileErrors
	// flushers save what the windows hold back, as MyTube quits.
	flushers []*func()
	// stopWatch stops the watch of settings.json, as MyTube quits.
	stopWatch context.CancelFunc

	// fileMu guards file, what the windows are told of settings.json;
	// watcher is its watch, while one runs.
	fileMu  sync.Mutex
	file    fileState
	watcher atomic.Pointer[mygo.FileWatcher]
}

// New opens what the app keeps in opts.DataDir and opts.CacheDir.
func New(opts Options) (*App, error) {
	a := &App{opts: opts, lastCrash: opts.LastCrash}
	var err error
	if a.settings, err = store.OpenWith(filepath.Join(opts.DataDir, "settings.json"), settings.Default(), settings.Decode); err != nil {
		return nil, fmt.Errorf("opening settings: %w", err)
	}
	if err := a.settings.Aside(); err != nil {
		log.Printf("mytube: settings.json didn't read (%v): kept it as settings.json.broken, and started with the settings of a new install", err)
	}
	if a.library, err = library.Open(filepath.Join(opts.DataDir, "library.json")); err != nil {
		return nil, fmt.Errorf("opening the library: %w", err)
	}
	if a.positions, err = library.OpenPositions(filepath.Join(opts.DataDir, "positions.json")); err != nil {
		return nil, fmt.Errorf("opening where episodes were left: %w", err)
	}
	if a.videoStats, err = openVideoStats(filepath.Join(opts.DataDir, "video-stats.json")); err != nil {
		return nil, fmt.Errorf("opening the video stats: %w", err)
	}
	// The update preferences of the versions that looked for updates.
	removeFile(filepath.Join(opts.DataDir, "updater.json"))
	// The recordings MyTube 0.x found on YouTube Music for Spotify
	// tracks.
	removeFile(filepath.Join(opts.CacheDir, "matches.json"))
	loudness, err := store.Open(filepath.Join(opts.CacheDir, "loudness.json"), map[string]float64{})
	if err != nil {
		return nil, fmt.Errorf("opening loudness: %w", err)
	}
	aiLabels, err := store.Open(filepath.Join(opts.CacheDir, "ai-labels.json"), map[string]int64{})
	if err != nil {
		return nil, fmt.Errorf("opening AI labels: %w", err)
	}
	s := a.settings.Get().Clean()

	a.youtube = youtube.New()
	a.lastfmStore = lastfmStore{newSecretStore(opts.DataDir, lastfmSecret, "lastfm.json", "lastfm.dpapi")}
	a.openLastfm()
	a.tools = tools.New(filepath.Join(opts.DataDir, "bin"))
	a.resolver = stream.NewResolver(a.tools, filepath.Join(opts.CacheDir, "yt-dlp"))
	a.resolver.Told = a.surroundTold
	// MYTUBE_YTDLP_ONLY has yt-dlp find every stream, as before MyTube
	// asked YouTube's player API: to time one against the other, or to
	// tell which fails. The API still tells how loud tracks are.
	direct := stream.NewDirect()
	if os.Getenv("MYTUBE_YTDLP_ONLY") == "" {
		a.resolver.Direct = direct
	} else {
		log.Print("mytube: yt-dlp alone finds the streams (MYTUBE_YTDLP_ONLY)")
	}
	a.loudness = newLoudnessBook(direct.Loudness, loudness)
	a.cache = stream.NewCache(filepath.Join(opts.CacheDir, "audio"), 0)
	a.proxy = stream.NewProxy(a.resolver)
	// The player asks whether the network answers, which audio coming
	// tells too.
	a.online = online.New(online.Probe)
	a.proxy.Seen = a.online.Seen
	a.youtubeStore = newYouTubeStore(opts.DataDir)
	a.openYouTube()
	a.proxy.Cache = a.cache
	a.matcher = match.New(a.youtube)
	a.labels = newLabelCache(a.youtube.MadeWithAI, aiLabels, time.Now)
	userAgent := "MyTube/" + opts.Version + " (https://github.com/bigdaddy284/mymusic)"
	a.lyrics = lyrics.New(userAgent)
	a.unison = lyrics.NewUnison(userAgent)
	a.translations = newTranslations(filepath.Join(opts.CacheDir, "lyrics-translations"))
	a.lyricsMemo = newMemo[lyrics.Lyrics](200)
	a.related = newMemo[[]media.Shelf](100)
	a.episodeDetails = newMemo[media.EpisodeDetails](100)
	a.artStore = artwork.NewStore(filepath.Join(opts.CacheDir, "artwork"), artworkLimit)
	a.artwork = artwork.New()
	a.artwork.Store = a.artStore
	a.colors = newMemo[string](500)
	a.addons = addons.NewManager(filepath.Join(opts.DataDir, "addons"), opts.Builtin)
	a.openDownloads(s.DownloadsFolder)
	if err := a.openLocalFiles(); err != nil {
		return nil, fmt.Errorf("opening the index of the music on this computer: %w", err)
	}
	a.apply(s)
	return a, nil
}

// apply makes settings take effect in the parts of the app they configure.
func (a *App) apply(s settings.Settings) {
	a.tools.Configure(s.YTDLPPath, s.JSRuntime)
	a.resolver.SetQuality(stream.Quality(s.Quality))
	a.cache.SetLimit(int64(s.CacheSizeMB) << 20)
	a.downloads.SetMetered(a.metered.Load(), s.DownloadOnMetered)
	a.followLocalFiles(s)
	if mygo.App.IsReady() {
		a.applyMediaKeys(s.MediaKeys)
	}
	// Whether they can register yet, hotkeys tells, which tests give
	// their own. The settings' save tells the windows of the problems.
	a.applyGlobalShortcuts(false)
	if p := a.goPlayer.Load(); p != nil {
		p.SettingsChanged(s)
		a.applyEngine(s)
	}
	if a.youtube.SignedIn() && a.youtube.Channel() != s.YouTubeChannel {
		a.youtube.SetChannel(s.YouTubeChannel)
	}
	a.followColorFile(s)
	a.retitle()
}

// Bind registers the stream scheme, at which the page of the picture of
// music videos fetches it. It must run before mygo.App.Run.
func (a *App) Bind() {
	a.handleVideo()
}

// needWebView tells why what, which takes a window showing a web page,
// can't be done, or nil when it can: MyTube's own windows need no
// webview, so it runs on systems without one, where opening such a window
// would panic.
func needWebView(what string) error {
	if webViewAvailable() {
		return nil
	}
	return fmt.Errorf("%s needs Microsoft's WebView2 Runtime, which this computer doesn't have: get it from https://go.microsoft.com/fwlink/p/?LinkId=2124703", what)
}

// webViewAvailable is mygo.WebViewAvailable, asked on the main thread.
func webViewAvailable() bool {
	var available bool
	mygo.RunOnMain(func() { available = mygo.WebViewAvailable() })
	return available
}

// schemeURL returns the base URL of a custom scheme, as WebView2 serves
// custom schemes, at http://<scheme>.localhost: the stream addresses
// PrepareStream gives keep that form.
func schemeURL(scheme string) string {
	return "http://" + scheme + ".localhost/"
}

// memo remembers the results of lookups, up to a number of keys.
type memo[T any] struct {
	mu    sync.Mutex
	max   int
	order []string
	items map[string]T
}

func newMemo[T any](max int) *memo[T] {
	return &memo[T]{max: max, items: map[string]T{}}
}

func (m *memo[T]) get(key string) (T, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.items[key]
	return v, ok
}

func (m *memo[T]) put(key string, v T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[key]; !ok {
		m.order = append(m.order, key)
		if len(m.order) > m.max {
			delete(m.items, m.order[0])
			m.order = m.order[1:]
		}
	}
	m.items[key] = v
}
