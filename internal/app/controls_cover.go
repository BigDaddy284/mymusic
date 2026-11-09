package app

// The system's media controls show the covers of the artwork store, which
// the windows show too: each cover downloads once, and shows offline once
// kept.

import (
	"context"
	"sync"
	"time"
)

const (
	// coverWait bounds the wait for a cover.
	coverWait = 30 * time.Second
	// coverRetry is how long the controls go without a cover that didn't
	// load before they ask for it again.
	coverRetry = time.Minute
)

// controlsCovers hands the system's controls the cover of the track that
// shows, from the artwork store, which it reads beside them.
type controlsCovers struct {
	// load returns the image at an address, as the artwork store does.
	load func(ctx context.Context, rawURL string) ([]byte, error)
	// changed is called once a cover asked for loaded, or didn't.
	changed func()
	now     func() time.Time

	mu sync.Mutex
	// url is the address of the cover kept: data, or none when it didn't
	// load at failed.
	url    string
	data   []byte
	failed time.Time
	// want is the address asked for last, and loading those that load.
	want    string
	loading map[string]bool
}

// newControlsCovers returns the covers of the artwork store, which call
// changed once one loaded.
func (a *App) newControlsCovers(changed func()) *controlsCovers {
	return &controlsCovers{load: a.coverBytes, changed: changed, now: time.Now}
}

// get returns the cover at rawURL, once kept. Until then it returns
// nothing, and loads it beside the caller.
func (cc *controlsCovers) get(rawURL string) []byte {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.want = rawURL
	if rawURL == "" {
		return nil
	}
	if cc.url == rawURL && (cc.data != nil || cc.now().Sub(cc.failed) < coverRetry) {
		return cc.data
	}
	if !cc.loading[rawURL] {
		if cc.loading == nil {
			cc.loading = map[string]bool{}
		}
		cc.loading[rawURL] = true
		go cc.fetch(rawURL)
	}
	return nil
}

// fetch reads the cover at rawURL from the store, which downloads it unless
// it has it, and keeps it if it's still the one asked for.
func (cc *controlsCovers) fetch(rawURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), coverWait)
	data, err := cc.load(ctx, rawURL)
	cancel()
	if err != nil || len(data) == 0 {
		data = nil // the store logs why
	}
	if cc.keep(rawURL, data) {
		cc.changed()
	}
}

// keep keeps data, the cover at rawURL, or that it didn't load, if it's
// still the one asked for, and reports whether it was.
func (cc *controlsCovers) keep(rawURL string, data []byte) bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	delete(cc.loading, rawURL)
	if cc.want != rawURL {
		return false
	}
	cc.url, cc.data, cc.failed = rawURL, data, time.Time{}
	if data == nil {
		cc.failed = cc.now()
	}
	return true
}
