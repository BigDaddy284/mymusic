package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRandomArticle reads the title of the random article Wikipedia's API
// answers, and fails on an answer with none.
func TestRandomArticle(t *testing.T) {
	answer := `{"query":{"random":[{"id":1,"ns":0,"title":"Leptodactylus vastus"}]}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("list") != "random" || r.URL.Query().Get("rnnamespace") != "0" {
			t.Errorf("asked %s", r.URL)
		}
		w.Write([]byte(answer))
	}))
	defer srv.Close()
	title, err := randomArticle(context.Background(), srv.URL)
	if err != nil || title != "Leptodactylus vastus" {
		t.Fatalf("the random article is %q, %v", title, err)
	}
	answer = `{"query":{"random":[]}}`
	if _, err := randomArticle(context.Background(), srv.URL); err == nil {
		t.Error("an answer without an article is taken")
	}
}

// TestRandomWord picks a word of the list.
func TestRandomWord(t *testing.T) {
	if w := RandomWord(); w == "" {
		t.Error("the random word is empty")
	}
}

// TestRandomWordsDontRepeat goes through every pairing of the built-in
// words once before any comes again.
func TestRandomWordsDontRepeat(t *testing.T) {
	randomPicks.words = nil
	n := len(randomWords) * len(randomFilms)
	seen := map[string]bool{}
	for range n {
		w := RandomWord()
		if seen[w] {
			t.Fatalf("%q came again after %d", w, len(seen))
		}
		seen[w] = true
	}
}

// TestRandomTitlesDontRepeat searches each random article's title once,
// though Wikipedia answers it again.
func TestRandomTitlesDontRepeat(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Write([]byte(`{"query":{"random":[{"title":"Frogs"},{"title":"Toads"},{"title":"frogs"}]}}`))
	}))
	defer srv.Close()
	was := wikipediaAPI
	wikipediaAPI = srv.URL
	t.Cleanup(func() { wikipediaAPI = was })
	var v Videos
	a, _ := v.RandomTopic(context.Background())
	b, _ := v.RandomTopic(context.Background())
	c, wiki := v.RandomTopic(context.Background())
	if a != "Frogs" || b != "Toads" || wiki {
		t.Errorf("Random searched %q, %q, then %q (from Wikipedia: %v)", a, b, c, wiki)
	}
	if calls > 4 {
		t.Errorf("asked Wikipedia %d times", calls)
	}
}
