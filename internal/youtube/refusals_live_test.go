package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLiveAddDuplicate makes a playlist of one song on a real account,
// adds the song again, which YouTube refuses as a duplicate, then adds it
// anyway, and deletes the playlist. It changes the account, so it runs
// only with MYMUSIC_YT_WRITE=1 too, best with an account kept for tests.
// MYMUSIC_YT_RECORD names a folder to keep YouTube's answer to the
// duplicate in, for testdata/edit-playlist-duplicate.json:
//
//	MYMUSIC_LIVE=1 MYMUSIC_YT_WRITE=1 MYMUSIC_YT_COOKIES='SAPISID=…; …' MYMUSIC_YT_RECORD=dir go test ./internal/youtube -run LiveAddDuplicate -v
func TestLiveAddDuplicate(t *testing.T) {
	header := os.Getenv("MYMUSIC_YT_COOKIES")
	if os.Getenv("MYMUSIC_LIVE") == "" || os.Getenv("MYMUSIC_YT_WRITE") == "" || header == "" {
		t.Skip("set MYMUSIC_LIVE=1, MYMUSIC_YT_WRITE=1 and MYMUSIC_YT_COOKIES to add a song twice to a playlist of a signed-in account")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	c := New()
	cookies, err := http.ParseCookie(header)
	if err != nil {
		t.Fatal(err)
	}
	c.SetCookies(cookies)
	const song = "BH-SnQ8J1VU" // the first song of testdata/playlist.json
	id, err := c.CreatePlaylist(ctx, "MyMusic live test "+time.Now().Format("2006-01-02 15:04:05"), []string{song})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := c.DeletePlaylist(context.Background(), id); err != nil {
			t.Errorf("deleting %s at the end: %v", id, err)
		}
	})
	answer, err := c.call(ctx, "browse/edit_playlist", map[string]any{"playlistId": id,
		"actions": []any{map[string]any{"action": "ACTION_ADD_VIDEO", "addedVideoId": song}}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.MarshalIndent(answer.v, "", " ")
	t.Logf("the answer to the duplicate: %s", data)
	if dir := os.Getenv("MYMUSIC_YT_RECORD"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "edit-playlist-duplicate.json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := refusedEdit(answer, errors.New("refused")); answer.get("status").str() == "STATUS_SUCCEEDED" || !errors.Is(err, ErrAlreadyInPlaylist) {
		t.Errorf("the duplicate reads as %v, status %q", err, answer.get("status").str())
	}
	if err := c.AddToPlaylistAnyway(ctx, id, song); err != nil {
		t.Errorf("adding it anyway: %v", err)
	}
}
