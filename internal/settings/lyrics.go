package settings

import (
	"regexp"
	"slices"
)

// languageRe matches the code of a language, as "ja" or "zh-Hant".
var languageRe = regexp.MustCompile(`^[a-z]{2,3}(-[A-Z][a-z]{3})?$`)

// Lyrics is how the lyrics show, as the user leaves the switches beside
// them.
type Lyrics struct {
	// LyricsPlain shows synchronized lyrics as plain text.
	LyricsPlain bool `json:"lyricsPlain"`
	// LyricsTranslate are the languages of the lyrics the user translates,
	// codes of ISO 639-1 as "ja": their translation shows under each line.
	// "" stands for lyrics whose language MyTube can't tell.
	LyricsTranslate []string `json:"lyricsTranslate"`
}

// defaultLyrics is how a new install shows the lyrics.
func defaultLyrics() Lyrics {
	var l Lyrics
	l.LyricsTranslate = []string{}
	return l
}

// clean is l without the languages that aren't one.
func (l Lyrics) clean() Lyrics {
	translate := []string{}
	for _, lang := range l.LyricsTranslate {
		if (lang == "" || languageRe.MatchString(lang)) && !slices.Contains(translate, lang) && len(translate) < 200 {
			translate = append(translate, lang)
		}
	}
	l.LyricsTranslate = translate
	return l
}

// Translation is the language lyrics translate to, as Settings › Playback
// › Lyrics sets it. It is apart from Lyrics, as settings.json has its key
// after the playback's.
type Translation struct {
	// LyricsTranslateTo is the language lyrics translate to, a code of ISO
	// 639-1; "" is the system's.
	LyricsTranslateTo string `json:"lyricsTranslateTo"`
}

// clean is t with the system's language in place of one that isn't one.
func (t Translation) clean() Translation {
	if !languageRe.MatchString(t.LyricsTranslateTo) {
		t.LyricsTranslateTo = ""
	}
	return t
}
