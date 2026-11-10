// Package match finds the song of a music video on YouTube Music, of
// which YouTube Music makes a radio of songs, and cleans the titles of
// tracks of what names a release rather than a song.
package match

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"

	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// ErrNoMatch is returned when YouTube Music has no song close enough.
var ErrNoMatch = errors.New("no matching song found on YouTube Music")

// Searcher searches YouTube Music.
type Searcher interface {
	Search(ctx context.Context, query string, filter youtube.Filter) ([]media.Item, error)
}

// Matcher finds songs on YouTube Music.
type Matcher struct {
	YouTube Searcher
}

// New returns a matcher searching yt.
func New(yt Searcher) *Matcher {
	return &Matcher{YouTube: yt}
}

// noise matches the parts of titles that name a release rather than the
// song: "(Remastered 2011)", "- Radio Edit", "[Official Video]".
var noise = regexp.MustCompile(`(?i)\s*([(\[][^)\]]*(remaster|version|edit|mono|stereo|deluxe|official|video|audio|lyrics?|visuali[sz]er|feat\.?|ft\.?|with )[^)\]]*[)\]]|\s-\s.*(remaster|version|edit|mono|stereo)\b.*$)`)

// CleanTitle removes from a title what names the release rather than the
// song, such as "(Remastered 2011)" or "[Official Video]".
func CleanTitle(s string) string { return noise.ReplaceAllString(s, "") }

// normalize lowercases s, strips its accents and punctuation, and
// collapses its spaces.
func normalize(s string) string {
	var b strings.Builder
	space := false
	for _, r := range norm.NFKD.String(strings.ToLower(s)) {
		switch {
		case unicode.Is(unicode.Mn, r): // accents
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(r)
			space = false
		case !space && b.Len() > 0:
			b.WriteByte(' ')
			space = true
		}
	}
	return strings.TrimSpace(b.String())
}

func words(s string) []string { return strings.Fields(normalize(s)) }
