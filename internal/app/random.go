package app

// Random, on the Video side, searches for something nobody would think to
// search for. It doesn't come back to what it searched before: it takes
// the titles of random Wikipedia articles, a batch at a time, out of more
// than six million, and searches each once, and when Wikipedia doesn't
// answer, it goes through every pairing of its own words with a way to
// film them, in a shuffled order, before it starts over.

import (
	"context"
	"math/rand/v2"
	"strings"
	"sync"
	"time"
)

// randomBatch is how many titles of random articles Random asks Wikipedia
// for at a time.
const randomBatch = 20

// randomPicks is what Random keeps between its searches.
var randomPicks struct {
	sync.Mutex
	// titles are titles of random articles not searched yet.
	titles []string
	// words are pairings of the built-in words not searched yet, in a
	// shuffled order.
	words []string
	// used are the titles Random searched in this run of the app.
	used map[string]bool
}

// RandomTopic returns something random to search for, never one it
// returned before: the title of a random article of English Wikipedia, or,
// when Wikipedia doesn't answer, a word of a list of the app's own, with a
// way to film it. fromWikipedia tells which.
func (v *Videos) RandomTopic(ctx context.Context) (topic string, fromWikipedia bool) {
	if t := nextRandomTitle(ctx); t != "" {
		return t, true
	}
	return RandomWord(), false
}

func nextRandomTitle(ctx context.Context) string {
	p := &randomPicks
	p.Lock()
	defer p.Unlock()
	if p.used == nil {
		p.used = map[string]bool{}
	}
	for tries := 0; tries < 3; tries++ {
		for len(p.titles) > 0 {
			t := p.titles[0]
			p.titles = p.titles[1:]
			if key := strings.ToLower(t); !p.used[key] {
				p.used[key] = true
				return t
			}
		}
		ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
		titles, err := randomArticles(ctx, wikipediaAPI, randomBatch)
		cancel()
		if err != nil {
			return ""
		}
		p.titles = titles
	}
	return ""
}

// RandomWord returns a pairing of a word of the built-in list with a way
// to film it, as "glacier time lapse": each of the thousands of pairings
// once, in a shuffled order, before any comes again.
func RandomWord() string {
	p := &randomPicks
	p.Lock()
	defer p.Unlock()
	if len(p.words) == 0 {
		p.words = make([]string, 0, len(randomWords)*len(randomFilms))
		for _, w := range randomWords {
			for _, f := range randomFilms {
				p.words = append(p.words, strings.TrimSpace(w+" "+f))
			}
		}
		rand.Shuffle(len(p.words), func(i, j int) { p.words[i], p.words[j] = p.words[j], p.words[i] })
	}
	w := p.words[len(p.words)-1]
	p.words = p.words[:len(p.words)-1]
	return w
}

// randomFilms are ways to film anything, paired with the built-in words.
var randomFilms = []string{
	"", "documentary", "time lapse", "how it's made", "history", "close up", "tour",
	"explained", "slow motion", "world record", "for beginners", "in the 1970s",
	"drone footage", "behind the scenes", "fails", "restoration", "vlog", "4K",
	"news report", "ASMR",
}
