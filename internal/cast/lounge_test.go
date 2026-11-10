package cast

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

// TestLounge binds to the lounge as a remote of a screen, sends its
// commands, and binds again when the lounge forgot the session.
func TestLounge(t *testing.T) {
	var mu sync.Mutex
	var commands []url.Values
	forget := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		form, _ := url.ParseQuery(string(b))
		mu.Lock()
		defer mu.Unlock()
		switch {
		case r.URL.Path == "/api/lounge/pairing/get_lounge_token_batch":
			if form.Get("screen_ids") != "screen-1" {
				t.Errorf("token for %q", form.Get("screen_ids"))
			}
			io.WriteString(w, `{"screens":[{"screenId":"screen-1","loungeToken":"tok"}]}`)
		case r.URL.Path == "/api/lounge/bc/bind" && r.URL.Query().Get("SID") == "":
			if r.Header.Get("X-YouTube-LoungeId-Token") != "tok" {
				t.Error("bind without the lounge token")
			}
			io.WriteString(w, `123
[[0,["c","sid-1","",8]],[1,["S","gs-1"]]]`)
		case r.URL.Path == "/api/lounge/bc/bind":
			if forget {
				forget = false
				w.WriteHeader(http.StatusGone)
				return
			}
			if r.URL.Query().Get("SID") != "sid-1" || r.URL.Query().Get("gsessionid") != "gs-1" {
				t.Errorf("command of the session %q", r.URL.RawQuery)
			}
			commands = append(commands, form)
		}
	}))
	defer srv.Close()
	old := LoungeURL
	LoungeURL = srv.URL
	defer func() { LoungeURL = old }()

	ctx := context.Background()
	l, err := BindLounge(ctx, srv.Client(), "screen-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := l.PlayVideo(ctx, "dQw4w9WgXcQ", 61.4); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	forget = true
	mu.Unlock()
	if err := l.Pause(ctx); err != nil {
		t.Fatal(err)
	}
	if len(commands) != 2 {
		t.Fatalf("%d commands", len(commands))
	}
	if c := commands[0]; c.Get("req0__sc") != "setPlaylist" || c.Get("req0_videoId") != "dQw4w9WgXcQ" || c.Get("req0_currentTime") != "61" || c.Get("count") != "1" {
		t.Errorf("setPlaylist: %v", c)
	}
	if commands[1].Get("req0__sc") != "pause" {
		t.Errorf("pause: %v", commands[1])
	}
}
