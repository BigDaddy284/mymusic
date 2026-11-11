package main

import (
	"context"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/lyrics"
)

// translateLanguages are the languages Settings › Playback offers to
// translate lyrics to, by their codes of ISO 639-1.
var translateLanguages = []struct{ code, name string }{
	{"af", "Afrikaans"}, {"ar", "Arabic"}, {"bg", "Bulgarian"}, {"bn", "Bengali"}, {"ca", "Catalan"},
	{"zh", "Chinese (Simplified)"}, {"zh-Hant", "Chinese (Traditional)"}, {"hr", "Croatian"}, {"cs", "Czech"},
	{"da", "Danish"}, {"nl", "Dutch"}, {"en", "English"}, {"et", "Estonian"}, {"tl", "Filipino"},
	{"fi", "Finnish"}, {"fr", "French"}, {"de", "German"}, {"el", "Greek"}, {"he", "Hebrew"},
	{"hi", "Hindi"}, {"hu", "Hungarian"}, {"id", "Indonesian"}, {"it", "Italian"}, {"ja", "Japanese"},
	{"ko", "Korean"}, {"lv", "Latvian"}, {"lt", "Lithuanian"}, {"ms", "Malay"}, {"no", "Norwegian"},
	{"fa", "Persian"}, {"pl", "Polish"}, {"pt", "Portuguese"}, {"ro", "Romanian"}, {"ru", "Russian"},
	{"sr", "Serbian"}, {"sk", "Slovak"}, {"sl", "Slovenian"}, {"es", "Spanish"}, {"sw", "Swahili"},
	{"sv", "Swedish"}, {"ta", "Tamil"}, {"te", "Telugu"}, {"th", "Thai"}, {"tr", "Turkish"},
	{"uk", "Ukrainian"}, {"ur", "Urdu"}, {"vi", "Vietnamese"},
}

// translateTip is the tip of the translate switch above the lyrics, which
// tells what it sends, as only it does.
const translateTip = "Shows each line's translation under it. While it's on, the lyrics of the track playing and its video ID go to Unison, which has Google translate them."

// translateTarget is the language lyrics translate to: the one chosen in
// Settings › Playback, else the system's.
func (s *shell) translateTarget() string {
	if s.set.LyricsTranslateTo != "" {
		return s.set.LyricsTranslateTo
	}
	return s.systemLanguage()
}

// systemLanguage is the system's language, asked for once.
func (s *shell) systemLanguage() string {
	ls := &s.lyricsView
	if ls.system == "" {
		ls.system = s.svc.System.Language()
	}
	return ls.system
}

// lyricsExtra is what the lines of lyrics show besides their words, as
// the user chose for their language: their translation.
type lyricsExtra struct {
	translate bool
	// tr is the translation of the lyrics, once it came; nil until then,
	// and while it doesn't show.
	tr *lyrics.Translation
}

// line returns the translation of line i, "" where there is none.
func (x lyricsExtra) line(i int) string {
	if x.translate && x.tr != nil {
		return x.tr.Lines[i].Text
	}
	return ""
}

// lyricsChoices shows how lyrics show, above them: synced or as plain
// text, and translated, where they can be; and, while their translation
// loads or when it fails, says so. It tells what their lines show
// besides their words.
func (s *shell) lyricsChoices(c *ui.Context, t media.Track, ly lyrics.Lyrics, where string) lyricsExtra {
	lang, to := ly.Language, s.translateTarget()
	// Lyrics of a language MyTube can't tell may be in another.
	canTranslate := lang != to
	x := lyricsExtra{translate: canTranslate && slices.Contains(s.set.LyricsTranslate, lang)}
	var tr *loaded[lyrics.Translation]
	if x.translate {
		tr = s.translation(t, ly, to)
		if tr.done && tr.err == nil && len(tr.value.Lines) == len(ly.Lines) {
			x.tr = &tr.value
		}
	}
	if ly.Synced || canTranslate {
		s.lyricsSwitches(c, ly, lang, where, canTranslate, &x)
	}
	if tr == nil {
		return x
	}
	switch {
	case !tr.done:
		s.lyricsNote(c, "Translation is on its way…")
	case tr.err != nil:
		// One quiet line; the lyrics show all the same.
		k := &s.lk
		ui.Row(c).Padding(0, 24, 8).Gap(12).AlignItems(ui.Center).Shrink(0).Children(func() {
			ui.Text(c, "Translation isn't available right now.").FontSize(k.smallSize()).TextColor(k.muted)
			if s.link(c, "lyrics-translate-again-"+where, "try again") {
				s.svc.Catalog.TranslateAgain()
				s.forget(translationKey(t, to))
			}
		})
	}
	return x
}

// translation fetches the translation of a track's lyrics into the
// language to.
func (s *shell) translation(t media.Track, ly lyrics.Lyrics, to string) *loaded[lyrics.Translation] {
	return fetch(s, translationKey(t, to), func(ctx context.Context) (lyrics.Translation, error) {
		return s.svc.Catalog.TranslateLyrics(ctx, t, ly, to)
	})
}

// translationKey is the key translation fetches the translation of a
// track's lyrics into the language into under.
func translationKey(t media.Track, into string) string {
	return "lyrics-translation:" + t.ID + ":" + into
}

// lyricsSwitches shows the choices above lyrics in the language lang:
// synced or as plain text, when they are synced; and translated, when
// they can be. It changes x as the user chooses.
func (s *shell) lyricsSwitches(c *ui.Context, ly lyrics.Lyrics, lang, where string, canTranslate bool, x *lyricsExtra) {
	ui.Row(c).Padding(0, 24, 8).GapX(16).GapY(8).Wrap().AlignItems(ui.Center).Shrink(0).Children(func() {
		if ly.Synced {
			mode := 0
			if s.set.LyricsPlain {
				mode = 1
			}
			if s.segmented(c, "lyrics-mode-"+where, "Show lyrics as", &mode, "synced", "plain text") {
				s.set.LyricsPlain = mode == 1
				s.save()
			}
		}
		if canTranslate {
			on := x.translate
			if s.lyricsSwitch(c, "lyrics-translate-"+where, "translate", "Translate lyrics", translateTip, &on) {
				x.translate = on
				s.set.LyricsTranslate = slices.DeleteFunc(slices.Clone(s.set.LyricsTranslate), func(l string) bool { return l == lang })
				if on {
					s.set.LyricsTranslate = append(s.set.LyricsTranslate, lang)
				}
				s.save()
			}
		}
	})
}

// lyricsSwitch is a switch above the lyrics with its label, and tip,
// which tells what it does: it tells whether the user flipped it.
func (s *shell) lyricsSwitch(c *ui.Context, key, label, name, tip string, on *bool) bool {
	changed := false
	ui.Row(c).Gap(8).AlignItems(ui.Center).Shrink(0).Children(func() {
		changed = s.switchOf(c.Key(key), on).Label(name).Tooltip(tip).Changed()
		ui.Text(c, s.lk.sentence(label)).FontSize(s.lk.smallSize())
	})
	return changed
}

func (s *shell) lyricsNote(c *ui.Context, text string) {
	ui.Text(c, text).FontSize(s.lk.smallSize()).TextColor(s.lk.muted).Padding(0, 24, 8).Shrink(0)
}

// translateOptions are the choices of "Translate lyrics to", the system's
// language first, made once.
func (s *shell) translateOptions() []string {
	ls := &s.lyricsView
	if ls.languages == nil {
		system := s.systemLanguage()
		for _, l := range translateLanguages {
			if l.code == system {
				system = l.name
			}
		}
		ls.languages = []string{"System language (" + system + ")"}
		for _, l := range translateLanguages {
			ls.languages = append(ls.languages, l.name)
		}
	}
	return ls.languages
}

func init() { addSettings("playback", 60, (*shell).lyricsSettings) }

// lyricsSettings chooses the language lyrics translate to.
func (s *shell) lyricsSettings(c *ui.Context) {
	s.group(c, "Lyrics", func() {
		s.row(c, "Translate lyrics to", "The translate switch beside the lyrics shows them in this language. While it's on, the lyrics of the track playing and its video ID go to Unison, which has Google translate them; nothing is sent while it's off.", func() {
			options := s.translateOptions()
			chosen := options[0]
			for i, l := range translateLanguages {
				if l.code == s.set.LyricsTranslateTo {
					chosen = options[i+1]
				}
			}
			if s.selectOf(c.Key("translate-to"), &chosen, options).Width(220).Label("Translate lyrics to").Changed() {
				s.set.LyricsTranslateTo = ""
				if i := slices.Index(options, chosen); i > 0 {
					s.set.LyricsTranslateTo = translateLanguages[i-1].code
				}
				s.save()
			}
		})
	})
}
