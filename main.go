// Command mymusic is a desktop music player that brings Spotify and
// YouTube Music together in one customizable interface.
package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo"

	"github.com/bigdaddy284/mymusic/builtin"
	"github.com/bigdaddy284/mymusic/internal/app"
	"github.com/bigdaddy284/mymusic/internal/diag"
	"github.com/bigdaddy284/mymusic/internal/sysproxy"
)

//go:embed brand/tray.png
var trayIcon []byte

// allowAudio lets pages play sound no click started: the hidden page of
// MyMusic's own Spotify device gets none. WebView2 takes the arguments of
// its browser from this variable, which must be set before its first page.
func allowAudio() {
	const name, arg = "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "--autoplay-policy=no-user-gesture-required"
	if args := os.Getenv(name); !strings.Contains(args, "--autoplay-policy") {
		os.Setenv(name, strings.TrimSpace(args+" "+arg))
	}
}

// version is the version of builds made without the mygo CLI, which links
// the version of mygo.config.ts into the app.
const version = "0.5.10"

func main() {
	// MYMUSIC_DATA_DIR keeps everything in one folder: a portable copy,
	// or a test run apart from the user's own data.
	portable := os.Getenv("MYMUSIC_DATA_DIR")
	if dir := portable; dir != "" {
		mygo.App.SetPath(mygo.PathUserData, filepath.Join(dir, "data"))
		mygo.App.SetPath(mygo.PathCache, filepath.Join(dir, "cache"))
		mygo.App.SetPath(mygo.PathLogs, filepath.Join(dir, "logs"))
	}
	if !mygo.App.RequestSingleInstanceLock() {
		return // the running instance takes over
	}
	allowAudio()
	// The log and the crash reports, once this is the instance that runs.
	logsDir, err := mygo.App.Path(mygo.PathLogs)
	if err != nil {
		log.Fatal(err)
	}
	running := mygo.App.Version()
	if running == "" {
		running = version
	}
	logs, lastCrash, err := diag.Start(logsDir, running)
	if err != nil {
		log.Printf("mymusic: writing the log: %v", err)
		logsDir = ""
	} else {
		defer logs.Close()
	}
	// The proxy of Windows' settings, followed from now on, and logged.
	sysproxy.Install()
	dataDir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		log.Fatal(err)
	}
	cacheDir, err := mygo.App.Path(mygo.PathCache)
	if err != nil {
		log.Fatal(err)
	}
	a, err := app.New(app.Options{
		DataDir:      dataDir,
		CacheDir:     cacheDir,
		DownloadsDir: defaultDownloadsDir(portable, cacheDir),
		Builtin:      builtin.FS,
		TrayIcon:     trayIcon,
		Version:      version,
		LogsDir:      logsDir,
		LastCrash:    lastCrash,
		Interface:    Open,
		Mini:         OpenMini,
	})
	if err != nil {
		log.Fatal(err)
	}
	a.Bind()
	mygo.App.OnSecondInstance(a.SecondInstance)
	mygo.App.OnActivate(a.Activate)
	mygo.App.OnOpenURL(a.HandleURL)
	mygo.App.WhenReady(a.Start)
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

// defaultDownloadsDir is the folder the downloads are in, unless the user
// chose another: %LocalAppData%\MyMusic\downloads, beside the cache, in
// the folder of the app's data that Windows keeps and doesn't roam. A
// portable copy keeps them in its own folder.
func defaultDownloadsDir(portable, cacheDir string) string {
	if portable != "" {
		return filepath.Join(portable, "downloads")
	}
	return filepath.Join(cacheDir, "downloads")
}
