package online

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	// The proxy of Windows' settings is in http.DefaultTransport before
	// probeClient clones it.
	_ "github.com/bigdaddy284/mymusic/internal/sysproxy"
)

// The addresses Probe asks, and how; tests change them. YouTube's answer
// 204 with no body, where yt-dlp and YouTube Music ask; a captive portal's
// page doesn't.
var (
	youtubeProbes = []string{"https://www.youtube.com/generate_204", "https://music.youtube.com/generate_204"}
	probeClient   = newProbeClient()
)

// newProbeClient returns a client that asks once, on a connection of its
// own, and takes no redirect: a portal redirects.
func newProbeClient() *http.Client {
	c := &http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		t = t.Clone()
		t.DisableKeepAlives = true
		c.Transport = t
	}
	return c
}

// Probe asks YouTube and YouTube Music whether they answer, at once, and
// returns nil as soon as one does, or why neither did.
func Probe(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	answers := make(chan error, len(youtubeProbes))
	for _, u := range youtubeProbes {
		go func() { answers <- ask(ctx, u) }()
	}
	var errs []error
	for range youtubeProbes {
		err := <-answers
		if err == nil {
			return nil
		}
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// ask asks u, and returns nil when it answers 204 No Content.
func ask(ctx context.Context, u string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	res, err := probeClient.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s answered HTTP %d", u, res.StatusCode)
	}
	return nil
}
