package youtube

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestLiveShelfKeys lists the shelves of YouTube Music's home feed with
// the keys Home names them by, which home-shelf-keys.json, hand-written,
// stands for. It runs with MYTUBE_LIVE=1, signed out, or signed in with
// the Cookie header of a music.youtube.com tab, as TestLiveSignedIn is:
//
//	MYTUBE_LIVE=1 go test ./internal/youtube -run TestLiveShelfKeys -v
//	MYTUBE_LIVE=1 MYTUBE_YT_COOKIES='SAPISID=…; …' go test ./internal/youtube -run TestLiveShelfKeys -v
//
// A shelf without a key is named by its title, which changes with the
// language and the day: the test fails when no shelf has one.
func TestLiveShelfKeys(t *testing.T) {
	if os.Getenv("MYTUBE_LIVE") == "" {
		t.Skip("set MYTUBE_LIVE=1 to read YouTube Music's home feed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	c := New()
	if header := os.Getenv("MYTUBE_YT_COOKIES"); header != "" {
		cookies, err := http.ParseCookie(header)
		if err != nil {
			t.Fatal(err)
		}
		c.SetCookies(cookies)
	}
	shelves, err := c.Home(ctx)
	if err != nil || len(shelves) == 0 {
		t.Fatalf("home: %d shelves, %v", len(shelves), err)
	}
	keyed := 0
	for _, sh := range shelves {
		t.Logf("%q: %q", sh.Title, sh.Key)
		if sh.Key != "" {
			keyed++
		}
	}
	if keyed == 0 {
		t.Errorf("none of the %d shelves has a key", len(shelves))
	}
}
