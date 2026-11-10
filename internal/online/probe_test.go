package online

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// probeAgainst has Probe ask the handlers, as YouTube's and YouTube
// Music's addresses, over TLS with the test servers' certificate, for the
// test.
func probeAgainst(t *testing.T, youtube, music http.HandlerFunc) {
	t.Helper()
	yt := httptest.NewTLSServer(youtube)
	ym := httptest.NewTLSServer(music)
	t.Cleanup(yt.Close)
	t.Cleanup(ym.Close)
	urls, client := youtubeProbes, probeClient
	t.Cleanup(func() { youtubeProbes, probeClient = urls, client })
	youtubeProbes = []string{yt.URL + "/generate_204", ym.URL + "/generate_204"}
	// The servers share their test certificate.
	probeClient = yt.Client()
	probeClient.CheckRedirect = newProbeClient().CheckRedirect
}

func status(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if code == http.StatusFound {
			http.Redirect(w, r, "/portal", code)
			return
		}
		w.WriteHeader(code)
	}
}

func TestProbe(t *testing.T) {
	for _, c := range []struct {
		name           string
		youtube, music int
		online         bool
	}{
		{"YouTube answers", http.StatusNoContent, http.StatusServiceUnavailable, true},
		{"YouTube Music answers", http.StatusServiceUnavailable, http.StatusNoContent, true},
		{"a portal's page", http.StatusOK, 0, false},
		{"a portal's redirect", http.StatusFound, 0, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			music := status(c.music)
			if c.music == 0 {
				// A portal answers for every site, which TLS doesn't trust.
				music = func(w http.ResponseWriter, r *http.Request) { panic(http.ErrAbortHandler) }
			}
			probeAgainst(t, status(c.youtube), music)
			if err := Probe(context.Background()); (err == nil) != c.online {
				t.Errorf("Probe = %v", err)
			}
		})
	}
}

func TestProbeWithoutANetwork(t *testing.T) {
	probeAgainst(t, status(http.StatusNoContent), status(http.StatusNoContent))
	// Connections fail at once, as with no route out.
	tr := probeClient.Transport.(*http.Transport).Clone()
	tr.DialContext = func(context.Context, string, string) (net.Conn, error) {
		return nil, errors.New("no route to host")
	}
	probeClient.Transport = tr
	if err := Probe(context.Background()); err == nil {
		t.Error("online without connections")
	}
}

func TestProbeDistrustsCertificates(t *testing.T) {
	probeAgainst(t, status(http.StatusNoContent), status(http.StatusNoContent))
	// The default client doesn't trust the test servers, as it wouldn't
	// trust a portal answering for YouTube.
	probeClient = newProbeClient()
	if err := Probe(context.Background()); err == nil {
		t.Error("online with untrusted certificates")
	}
}

// TestLiveProbe asks YouTube. It runs with MYMUSIC_LIVE=1.
func TestLiveProbe(t *testing.T) {
	if os.Getenv("MYMUSIC_LIVE") == "" {
		t.Skip("set MYMUSIC_LIVE=1 to ask YouTube")
	}
	if err := Probe(context.Background()); err != nil {
		t.Errorf("Probe = %v", err)
	}
}
