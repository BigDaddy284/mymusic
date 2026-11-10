package main

import (
	"os"
	"strings"
	"testing"

	"github.com/grafana/sobek"
)

// spotifyPage runs the script of public/spotify.html, the page of
// MyTube's own Spotify device, without a document, in a webview whose
// EME grants the key systems of granted, or that has no EME when granted
// is nil, and returns its runtime. asked lists the key systems the page
// asked for.
func spotifyPage(t *testing.T, secure bool, granted []string) *sobek.Runtime {
	t.Helper()
	src, err := os.ReadFile("public/spotify.html")
	if err != nil {
		t.Fatal(err)
	}
	_, script, _ := strings.Cut(string(src), "<script>")
	script, _, ok := strings.Cut(script, "</script>")
	if !ok {
		t.Fatal("public/spotify.html has no script")
	}
	vm := sobek.New()
	if err := vm.Set("isSecureContext", secure); err != nil {
		t.Fatal(err)
	}
	if err := vm.Set("granted", granted); err != nil {
		t.Fatal(err)
	}
	if _, err := vm.RunString(`
		var window = globalThis;
		var asked = [];
		var navigator = granted === null ? {} : {
			requestMediaKeySystemAccess(system) {
				asked.push(system);
				return granted.includes(system)
					? Promise.resolve({ createMediaKeys: () => Promise.resolve({}) })
					: Promise.reject(new Error("NotSupportedError"));
			},
		};`); err != nil {
		t.Fatal(err)
	}
	if _, err := vm.RunScript("spotify.html", script); err != nil {
		t.Fatal(err)
	}
	return vm
}

// TestSpotifyPageDecrypts has the page of MyTube's own Spotify device
// take Widevine and PlayReady in WebView2, and nothing else: ClearKey
// doesn't play Spotify's audio (#295).
func TestSpotifyPageDecrypts(t *testing.T) {
	for _, c := range []struct {
		name    string
		secure  bool
		granted []string
		want    bool
	}{
		{"WebView2, with Widevine", true, []string{"com.widevine.alpha", "org.w3.clearkey"}, true},
		{"WebView2, with PlayReady", true, []string{"com.microsoft.playready.recommendation"}, true},
		{"ClearKey alone", true, []string{"org.w3.clearkey"}, false},
		{"no EME", true, nil, false},
		{"a page that isn't secure", false, []string{"com.widevine.alpha"}, false},
	} {
		vm := spotifyPage(t, c.secure, c.granted)
		if _, err := vm.RunString(`var supported = "pending"; mm.supported().then((v) => (supported = v), (e) => (supported = String(e)));`); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if got := vm.Get("supported").Export(); got != c.want {
			t.Errorf("%s: mm.supported() = %v, want %v (asked %v)", c.name, got, c.want, vm.Get("asked").Export())
		}
	}
	// Widevine first: WebView2 asks for nothing more.
	vm := spotifyPage(t, true, []string{"com.widevine.alpha", "com.microsoft.playready"})
	if _, err := vm.RunString(`mm.supported()`); err != nil {
		t.Fatal(err)
	}
	if got := vm.Get("asked").Export(); len(got.([]any)) != 1 || got.([]any)[0] != "com.widevine.alpha" {
		t.Errorf("WebView2 asked for %v, want Widevine alone", got)
	}
}
