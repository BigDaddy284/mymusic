package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSponsorSegments(t *testing.T) {
	sum := sha256.Sum256([]byte("dQw4w9WgXcQ"))
	prefix := hex.EncodeToString(sum[:])[:4]
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/skipSegments/"+prefix) {
			http.NotFound(w, r)
			return
		}
		if !strings.Contains(r.URL.Query().Get("categories"), "sponsor") {
			t.Errorf("categories = %q", r.URL.Query().Get("categories"))
		}
		w.Write([]byte(`[
			{"videoID":"other","segments":[{"segment":[1,2],"category":"sponsor","actionType":"skip","UUID":"x"}]},
			{"videoID":"dQw4w9WgXcQ","segments":[
				{"segment":[90,120],"category":"outro","actionType":"skip","UUID":"b"},
				{"segment":[10,25.5],"category":"sponsor","actionType":"mute","UUID":"a"},
				{"segment":[5],"category":"intro","actionType":"skip","UUID":"bad"}
			]}]`))
	}))
	defer srv.Close()
	got, err := sponsorSegments(context.Background(), srv.URL, "dQw4w9WgXcQ")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].UUID != "a" || !got[0].Mute || got[0].End != 25.5 || got[1].Category != "outro" {
		t.Fatalf("segments = %+v", got)
	}
	none, err := sponsorSegments(context.Background(), srv.URL, "unknownvid1")
	if err != nil || len(none) != 0 {
		t.Fatalf("unknown video: %v, %v", none, err)
	}
}
