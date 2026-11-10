package app

// The Video page: videos of YouTube's own website, read with the session
// YouTube Music signs in with, and played in YouTube's embedded player,
// public/watch.html. Its Random button searches for the title of a random
// article of Wikipedia.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/egoist/mygo/netproxy"

	"github.com/bigdaddy284/mymusic/internal/audio/output"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// Videos reads the videos of YouTube's website.
type Videos struct{ a *App }

// SignedIn reports whether the session that signs in to YouTube Music is
// there, which signs in to YouTube's website too.
func (v *Videos) SignedIn() bool { return v.a.youtube.SignedIn() && !v.a.youtube.Rejected() }

// Search searches YouTube's website for videos.
func (v *Videos) Search(ctx context.Context, query string) (youtube.VideoPage, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return youtube.VideoPage{Videos: []youtube.Video{}}, nil
	}
	page, err := v.a.youtube.SearchVideos(ctx, query)
	return page, v.a.offlineErr(ctx, err)
}

// Recent searches YouTube's website for videos of this week.
func (v *Videos) Recent(ctx context.Context, query string) (youtube.VideoPage, error) {
	page, err := v.a.youtube.SearchRecentVideos(ctx, query)
	return page, v.a.offlineErr(ctx, err)
}

// Home returns the videos YouTube's home page recommends.
func (v *Videos) Home(ctx context.Context) (youtube.VideoPage, error) {
	page, err := v.a.youtube.HomeVideos(ctx)
	return page, v.a.offlineErr(ctx, err)
}

// More returns the page of videos after one, by its token: of a search,
// or of another page.
func (v *Videos) More(ctx context.Context, search bool, token string) (youtube.VideoPage, error) {
	endpoint := "browse"
	if search {
		endpoint = "search"
	}
	return v.a.youtube.MoreVideos(ctx, endpoint, token)
}

// Subscriptions returns the latest videos of the channels the account
// subscribes to.
func (v *Videos) Subscriptions(ctx context.Context) (youtube.VideoPage, error) {
	return v.a.youtube.Subscriptions(ctx)
}

// History returns the videos the account watched, the latest first.
func (v *Videos) History(ctx context.Context) (youtube.VideoPage, error) {
	return v.a.youtube.WatchHistory(ctx)
}

// Playlist returns the videos of a playlist.
func (v *Videos) Playlist(ctx context.Context, id string) (youtube.VideoPage, error) {
	return v.a.youtube.PlaylistVideos(ctx, id)
}

// Playlists returns the account's playlists, with Watch later and Liked
// videos first, which every account has.
func (v *Videos) Playlists(ctx context.Context) ([]youtube.VideoPlaylist, error) {
	lists, err := v.a.youtube.VideoPlaylists(ctx)
	if err != nil {
		return nil, err
	}
	out := []youtube.VideoPlaylist{
		{ID: "WL", Title: "Watch later"},
		{ID: "LL", Title: "Liked videos"},
	}
	for _, p := range lists {
		if p.ID != "WL" && p.ID != "LL" {
			out = append(out, p)
		}
	}
	return out, nil
}

// RandomTopic returns something random to search for: the title of a
// random article of English Wikipedia, or, when Wikipedia doesn't answer,
// a word of a list of the app's own. fromWikipedia tells which.
func (v *Videos) RandomTopic(ctx context.Context) (topic string, fromWikipedia bool) {
	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	if title, err := randomArticle(ctx, wikipediaAPI); err == nil {
		return title, true
	}
	return RandomWord(), false
}

// wikipediaAPI is the API of English Wikipedia, replaced in tests.
var wikipediaAPI = "https://en.wikipedia.org/w/api.php"

var wikipediaClient = &http.Client{Transport: netproxy.Transport(), Timeout: 10 * time.Second}

// randomArticle asks Wikipedia's API at api for the title of a random
// article: one of its main namespace, which holds the articles alone.
func randomArticle(ctx context.Context, api string) (string, error) {
	q := url.Values{
		"action":      {"query"},
		"list":        {"random"},
		"rnnamespace": {"0"},
		"rnlimit":     {"1"},
		"format":      {"json"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api+"?"+q.Encode(), nil)
	if err != nil {
		return "", err
	}
	// Wikipedia asks apps to name themselves.
	req.Header.Set("User-Agent", "MyTubeApp/1.0 (Windows desktop app; random video search)")
	res, err := wikipediaClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wikipedia: HTTP %d", res.StatusCode)
	}
	var body struct {
		Query struct {
			Random []struct {
				Title string `json:"title"`
			} `json:"random"`
		} `json:"query"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("wikipedia: %w", err)
	}
	if len(body.Query.Random) == 0 || strings.TrimSpace(body.Query.Random[0].Title) == "" {
		return "", errors.New("wikipedia: no random article")
	}
	return strings.TrimSpace(body.Query.Random[0].Title), nil
}

// RandomWord returns a word of the built-in list.
func RandomWord() string { return randomWords[rand.IntN(len(randomWords))] }

// randomWords are what the Random button searches for when Wikipedia
// doesn't answer: words that find plenty of videos.
var randomWords = []string{
	"aurora", "avalanche", "bakery", "balloon", "bamboo", "beekeeping", "blacksmith", "bonsai",
	"cactus", "calligraphy", "canyon", "carnival", "castle", "cathedral", "chess", "circus",
	"clockwork", "comet", "coral reef", "crossword", "desert", "dinosaur", "diving", "dragonfly",
	"drumline", "eclipse", "elephant", "espresso", "falconry", "ferris wheel", "fireworks", "fjord",
	"flamenco", "fossil", "geyser", "glacier", "glassblowing", "gondola", "harbor", "hedgehog",
	"hot air balloon", "hummingbird", "iceberg", "jellyfish", "juggling", "kayak", "kite", "lantern",
	"lighthouse", "lightning", "locomotive", "magnet", "mammoth", "marathon", "meteor", "monsoon",
	"mosaic", "mountain pass", "nebula", "night market", "octopus", "origami", "orchestra", "owl",
	"paper plane", "parrot", "pinball", "pottery", "puffin", "pyramid", "quasar", "rainforest",
	"robot", "rollercoaster", "sailing", "samurai", "sandcastle", "skyscraper", "snowboard", "space station",
	"submarine", "sundial", "sushi", "tapestry", "telescope", "thunderstorm", "tide pool", "tornado",
	"treehouse", "tugboat", "typewriter", "volcano", "waterfall", "whale", "windmill", "woodworking",
	"yo-yo", "zeppelin", "zebra", "zen garden",
}

// RouteOutput sends the sound of the web views' processes, WebView2's
// browser process pid and those it started, to the output device id, or,
// for "", to the system's default device.
func (v *Videos) RouteOutput(pid int, id string) error {
	if pid <= 0 {
		return errors.New("the video player has no process yet")
	}
	return output.RouteProcesses(output.ProcessTree(uint32(pid)), id)
}
