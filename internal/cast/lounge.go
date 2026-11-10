package cast

// YouTube's lounge: the service YouTube's phone app controls the YouTube
// app of a TV through, by the TV's screen ID. MyTube binds to it as a
// remote, then sends it commands, as "setPlaylist" to play a video.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// LoungeURL is where the lounge is, replaced in tests.
var LoungeURL = "https://www.youtube.com"

// Lounge is a remote of a TV's YouTube app, bound to the lounge.
type Lounge struct {
	client *http.Client
	screen string

	mu            sync.Mutex
	token         string
	sid, gsession string
	id            string
	rid           int
}

var (
	loungeSID     = regexp.MustCompile(`"c","(.*?)","`)
	loungeSession = regexp.MustCompile(`"S","(.*?)"]`)
)

// errLoungeSession is why a command fails whose session the lounge forgot.
var errLoungeSession = errors.New("the lounge forgot the session")

// BindLounge binds to the lounge as a remote of the screen screenID.
func BindLounge(ctx context.Context, client *http.Client, screenID string) (*Lounge, error) {
	b := make([]byte, 13)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	l := &Lounge{client: client, screen: screenID, id: hex.EncodeToString(b)}
	if err := l.bind(ctx); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *Lounge) post(ctx context.Context, path string, query, form url.Values) (*http.Response, error) {
	u := LoungeURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://www.youtube.com")
	if l.token != "" {
		req.Header.Set("X-YouTube-LoungeId-Token", l.token)
	}
	return l.client.Do(req)
}

// bind asks the lounge for the screen's token, then for a session of a
// remote.
func (l *Lounge) bind(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	res, err := l.post(ctx, "/api/lounge/pairing/get_lounge_token_batch", nil, url.Values{"screen_ids": {l.screen}})
	if err != nil {
		return fmt.Errorf("YouTube's lounge: %w", err)
	}
	var tokens struct {
		Screens []struct {
			LoungeToken string `json:"loungeToken"`
		} `json:"screens"`
	}
	err = json.NewDecoder(res.Body).Decode(&tokens)
	res.Body.Close()
	if err != nil || res.StatusCode != http.StatusOK || len(tokens.Screens) == 0 || tokens.Screens[0].LoungeToken == "" {
		return fmt.Errorf("YouTube's lounge didn't know the device (HTTP %d)", res.StatusCode)
	}
	l.token = tokens.Screens[0].LoungeToken
	l.rid = 0
	res, err = l.post(ctx, "/api/lounge/bc/bind", url.Values{"RID": {"0"}, "VER": {"8"}, "CVER": {"1"}}, url.Values{
		"device": {"REMOTE_CONTROL"}, "id": {l.id}, "name": {"MyTube"}, "mdx-version": {"3"},
		"pairing_type": {"cast"}, "app": {"android-phone-13.14.55"},
	})
	if err != nil {
		return fmt.Errorf("YouTube's lounge: %w", err)
	}
	body, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	res.Body.Close()
	sid, gs := loungeSID.FindSubmatch(body), loungeSession.FindSubmatch(body)
	if res.StatusCode != http.StatusOK || sid == nil || gs == nil {
		return fmt.Errorf("YouTube's lounge didn't open a session (HTTP %d)", res.StatusCode)
	}
	l.sid, l.gsession = string(sid[1]), string(gs[1])
	l.rid = 1
	return nil
}

// Do sends the command cmd, with its arguments, binding again once when
// the lounge forgot the session.
func (l *Lounge) Do(ctx context.Context, cmd string, args map[string]string) error {
	err := l.do(ctx, cmd, args)
	if errors.Is(err, errLoungeSession) {
		if err = l.bind(ctx); err == nil {
			err = l.do(ctx, cmd, args)
		}
	}
	return err
}

func (l *Lounge) do(ctx context.Context, cmd string, args map[string]string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	form := url.Values{"count": {"1"}, "req0__sc": {cmd}}
	for k, v := range args {
		form.Set("req0_"+k, v)
	}
	q := url.Values{"SID": {l.sid}, "gsessionid": {l.gsession}, "RID": {strconv.Itoa(l.rid)}, "VER": {"8"}, "CVER": {"1"}}
	res, err := l.post(ctx, "/api/lounge/bc/bind", q, form)
	if err != nil {
		return fmt.Errorf("YouTube's lounge: %w", err)
	}
	io.Copy(io.Discard, io.LimitReader(res.Body, 1<<20))
	res.Body.Close()
	l.rid++
	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest, http.StatusNotFound, http.StatusGone:
		return errLoungeSession
	}
	return fmt.Errorf("YouTube's lounge: HTTP %d", res.StatusCode)
}

// PlayVideo has the TV play the video id from at seconds in.
func (l *Lounge) PlayVideo(ctx context.Context, id string, at float64) error {
	return l.Do(ctx, "setPlaylist", map[string]string{
		"videoId": id, "videoIds": id, "listId": "", "currentIndex": "-1",
		"currentTime": strconv.FormatFloat(max(at, 0), 'f', 0, 64), "audioOnly": "false",
	})
}

// Play, Pause and Seek control the video on the TV.
func (l *Lounge) Play(ctx context.Context) error  { return l.Do(ctx, "play", nil) }
func (l *Lounge) Pause(ctx context.Context) error { return l.Do(ctx, "pause", nil) }
func (l *Lounge) Seek(ctx context.Context, at float64) error {
	return l.Do(ctx, "seekTo", map[string]string{"newTime": strconv.FormatFloat(max(at, 0), 'f', 0, 64)})
}

// SetVolume sets the TV's volume, from 0 to 100.
func (l *Lounge) SetVolume(ctx context.Context, v float64) error {
	return l.Do(ctx, "setVolume", map[string]string{"volume": strconv.Itoa(int(min(max(v, 0), 100)))})
}
