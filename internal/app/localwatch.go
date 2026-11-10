package app

import "github.com/bigdaddy284/mymusic/internal/localfiles"

// watchLocalFiles has the watch of the folders of music follow folders,
// those scanned: once MyTube started, so that what changes in them
// shows without a scan, or, for an app whose options give their own watch,
// as tests do, from its first scan. The watch is made the first time;
// localScans.mu is held.
func (a *App) watchLocalFiles(folders []string) {
	sc := &a.localScans
	if sc.ctx.Err() != nil || !sc.started && a.opts.WatchMusic == nil {
		return
	}
	if sc.watch == nil {
		sc.watch = localfiles.NewWatcher(a.localFiles, localfiles.WatcherOptions{
			Watch: a.opts.WatchMusic,
			// A watch that lost track of what changed, or a folder that
			// went or came back, has the folders of the settings scanned.
			Rescan: func() { a.scanLocalFiles(a.settings.Get().Clean().MusicFolders, true) },
		})
	}
	sc.watch.Set(folders)
}

// closeLocalWatch stops the watch of the folders, and waits for what it
// reads, as MyTube quits: once the scans' context ended, so that it asks
// for none.
func (a *App) closeLocalWatch() {
	a.localScans.mu.Lock()
	w := a.localScans.watch
	a.localScans.mu.Unlock()
	if w != nil {
		w.Close()
	}
}

// Unwatched tells why MyTube can't watch a folder of music for changes,
// as when it's on a network drive, so that what changes in it shows only
// after a scan: "" while it watches it, before MyTube watches the
// folders, and while the folder isn't there, as its scan tells.
func (l *LocalFiles) Unwatched(folder string) string {
	l.a.localScans.mu.Lock()
	w := l.a.localScans.watch
	l.a.localScans.mu.Unlock()
	if w == nil {
		return ""
	}
	return w.Unwatched(folder)
}
