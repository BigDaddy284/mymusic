package main

import (
	"fmt"
	"slices"
	"testing"

	"github.com/bigdaddy284/mymusic/internal/player"
	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/lyrics"
)

// kiteLines are made-up lyrics in Japanese, and kiteTranslation their
// made-up translation.
var (
	kiteLines       = []lyrics.Line{{TimeMs: 0, Text: "紙の凧"}, {TimeMs: 60000, Text: "空へ"}}
	kiteTranslation = lyrics.Translation{From: "ja", To: "en", Lines: []lyrics.TranslatedLine{
		{Text: "Paper kite"}, {Text: "To the sky"}}}
)

// lyricsIn shows made-up lyrics in the language lang in the lyrics panel,
// with their translation into English when it is given.
func lyricsIn(t *testing.T, lang string, tr *loaded[lyrics.Translation]) (*shell, func() []string) {
	t.Helper()
	s := newTestShell(t)
	s.loads["lyrics:"+testTracks[0].ID] = &loaded[lyrics.Lyrics]{done: true, value: lyrics.Lyrics{Synced: true, Source: "Unison", Language: lang, Lines: kiteLines}}
	if tr != nil {
		s.loads["lyrics-translation:"+testTracks[0].ID+":en"] = tr
	}
	s.p.PlayTracks(testTracks, 0)
	s.settle(t, func(st player.State) bool { return st.Playing })
	tt := testerOf(t, s, 1280, 900)
	s.open("/lyrics")
	tt.Frame()
	return s, func() []string {
		s.drain()
		tt.Frame()
		return tt.Texts()
	}
}

// TestLyricsSwitches offers to translate lyrics in another language than
// the one they translate to, and nothing to show them in Latin letters.
func TestLyricsSwitches(t *testing.T) {
	for _, c := range []struct {
		lang      string
		translate bool
	}{{"ko", true}, {"ja", true}, {"es", true}, {"en", false}, {"", true}} {
		_, texts := lyricsIn(t, c.lang, nil)
		got := texts()
		if slices.Contains(got, "translate") != c.translate {
			t.Errorf("%q lyrics: translate shows: %v", c.lang, !c.translate)
		}
		if slices.Contains(got, "romaji") || slices.Contains(got, "romanize") {
			t.Errorf("%q lyrics offer Latin letters: %q", c.lang, got)
		}
		if !slices.Contains(got, kiteLines[0].Text) {
			t.Errorf("%q lyrics: the lyrics don't show: %q", c.lang, got)
		}
	}
}

// TestLyricsSwitchesInBlend names the switches above the lyrics as the
// layouts set labels: capitalised in Blend, where Tiles' are lowercase.
func TestLyricsSwitchesInBlend(t *testing.T) {
	s, texts := lyricsIn(t, "ko", nil)
	st := s.svc.Preferences.Get()
	st.Layout = settings.LayoutBlend
	if _, err := s.svc.Preferences.Set(st); err != nil {
		t.Fatal(err)
	}
	s.refresh()
	got := texts()
	for _, name := range []string{"Synced", "Translate"} {
		if !slices.Contains(got, name) {
			t.Errorf("%q doesn't show in Blend: %q", name, got)
		}
	}
	if slices.Contains(got, "translate") {
		t.Error("translate shows lowercase in Blend")
	}
}

// TestLyricsTranslated shows each line's translation under it.
func TestLyricsTranslated(t *testing.T) {
	s, texts := lyricsIn(t, "ja", &loaded[lyrics.Translation]{done: true, value: kiteTranslation})
	if got := texts(); slices.Contains(got, "Paper kite") {
		t.Errorf("translated while off: %q", got)
	}
	s.set.LyricsTranslate = []string{"ja"}
	got := texts()
	if i := slices.Index(got, "紙の凧"); i < 0 || i+1 >= len(got) || got[i+1] != "Paper kite" {
		t.Errorf("translated: %q", got)
	}
	// Translations of other lyrics of the track don't show on these.
	s.loads["lyrics-translation:"+testTracks[0].ID+":en"] = &loaded[lyrics.Translation]{done: true, value: lyrics.Translation{Lines: kiteTranslation.Lines[:1]}}
	if got := texts(); !slices.Contains(got, "紙の凧") || slices.Contains(got, "Paper kite") {
		t.Errorf("a translation of other lyrics: %q", got)
	}
}

// TestLyricsCreditTranslated keeps the credit Unison's lyrics ask for under
// them while they show translated.
func TestLyricsCreditTranslated(t *testing.T) {
	s, texts := lyricsIn(t, "ja", &loaded[lyrics.Translation]{done: true, value: kiteTranslation})
	s.loads["lyrics:"+testTracks[0].ID].(*loaded[lyrics.Lyrics]).value.SourceURL = lyrics.UnisonURL
	s.set.LyricsTranslate = []string{"ja"}
	got := texts()
	if !slices.Contains(got, "Paper kite") {
		t.Fatalf("not translated: %q", got)
	}
	if !slices.Contains(got, "lyrics from Unison ("+lyrics.UnisonURL+")") {
		t.Errorf("Unison's credit: %q", got)
	}
}

// TestLyricsTranslateRemembered remembers the translate switch for the
// next song in the same language.
func TestLyricsTranslateRemembered(t *testing.T) {
	s, texts := lyricsIn(t, "ja", &loaded[lyrics.Translation]{done: true, value: kiteTranslation})
	texts()
	tt := testerOf(t, s, 1280, 900)
	tt.Frame()
	click(t, s, tt, "Translate lyrics")
	if saved := s.svc.Preferences.Get(); !slices.Equal(saved.LyricsTranslate, []string{"ja"}) {
		t.Fatalf("saved %q", saved.LyricsTranslate)
	}
	if got := texts(); !slices.Contains(got, "Paper kite") {
		t.Errorf("turned on: %q", got)
	}
	click(t, s, tt, "Translate lyrics")
	if got := s.svc.Preferences.Get().LyricsTranslate; len(got) != 0 {
		t.Errorf("turned off, translated: %q", got)
	}
}

// TestLyricsTranslationFails shows the lyrics untranslated when their
// translation fails, says so in one line, and offers to try again.
func TestLyricsTranslationFails(t *testing.T) {
	failed := &loaded[lyrics.Translation]{done: true, err: fmt.Errorf("Translation isn't available right now: %w", lyrics.ErrRateLimited)}
	s, texts := lyricsIn(t, "ja", failed)
	s.set.LyricsTranslate = []string{"ja"}
	got := texts()
	if !slices.Contains(got, "Translation isn't available right now.") || !slices.Contains(got, "try again") || !slices.Contains(got, "紙の凧") {
		t.Fatalf("a failed translation: %q", got)
	}
	tt := testerOf(t, s, 1280, 900)
	tt.Frame()
	click(t, s, tt, "try again")
	if s.loads["lyrics-translation:"+testTracks[0].ID+":en"] == failed {
		t.Error("try again didn't ask for the translation again")
	}
}

// TestTranslateTo chooses the language lyrics translate to in Settings ›
// Playback, the system's by default.
func TestTranslateTo(t *testing.T) {
	s := newTestShell(t)
	// Another language chosen, the system's still names the first choice.
	s.set.LyricsTranslateTo = "ja"
	if options := s.translateOptions(); options[0] != "System language (English)" {
		t.Errorf("first option %q", options[0])
	}
	s.set.LyricsTranslateTo = ""
	if got := s.translateTarget(); got != "en" {
		t.Errorf("translate to %q by default", got)
	}
	if options := s.translateOptions(); options[0] != "System language (English)" || len(options) != len(translateLanguages)+1 {
		t.Errorf("options = %q", options[:3])
	}
	tt := testerOf(t, s, 1280, 1800)
	show(s, tt, "/settings/playback")
	if !tt.HasText("Translate lyrics to") {
		t.Fatalf("no setting: %q", tt.Texts())
	}
	s.set.LyricsTranslateTo = "ja"
	if got := s.translateTarget(); got != "ja" {
		t.Errorf("translate to %q", got)
	}
}

// TestLyricsSettingsInBlend titles the lyrics' group of the settings as
// Blend sets titles, as written, where Tiles lowers them.
func TestLyricsSettingsInBlend(t *testing.T) {
	s := newTestShell(t)
	st := s.svc.Preferences.Get()
	st.Layout = settings.LayoutBlend
	if _, err := s.svc.Preferences.Set(st); err != nil {
		t.Fatal(err)
	}
	s.refresh()
	tt := testerOf(t, s, 1280, 1800)
	show(s, tt, "/settings/playback")
	texts := tt.Texts()
	if i := slices.Index(texts, "Translate lyrics to"); i < 1 || texts[i-1] != "Lyrics" {
		t.Errorf("the group of Translate lyrics to isn't titled Lyrics: %q", texts)
	}
}

// TestLyricsSettingsEdited applies an edit of settings.json to how lyrics
// show while MyTube runs: the language they translate to, and the
// languages they translate.
func TestLyricsSettingsEdited(t *testing.T) {
	s := fileShell(t)
	editSettings(t, s, `{"lyricsTranslateTo": "ja", "lyricsTranslate": ["ko"]}`)
	s.refresh()
	if got := s.translateTarget(); got != "ja" {
		t.Errorf("translate to %q", got)
	}
	if !slices.Equal(s.set.LyricsTranslate, []string{"ko"}) {
		t.Errorf("the window has %q", s.set.LyricsTranslate)
	}
}
