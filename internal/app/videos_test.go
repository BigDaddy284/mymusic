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
