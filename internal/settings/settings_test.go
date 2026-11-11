package settings

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
)

func TestDefaultIsClean(t *testing.T) {
	d := Default()
	a, _ := json.Marshal(d)
	b, _ := json.Marshal(d.Clean())
	if string(a) != string(b) {
		t.Errorf("Clean changed the defaults:\n%s\n%s", a, b)
	}
}

func TestCleanClampsAndRepairs(t *testing.T) {
	s := Settings{
		Appearance: Appearance{
			Theme:     "../../etc",
			Accent:    "red; background: url(x)",
			FontScale: 5,
		},
		Playback: Playback{
			Crossfade: 60,
			Volume:    -1,
			Equalizer: Equalizer{Gains: []float64{99, -99}},
		},
		Tools:    Tools{JSRuntime: "python"},
		Accounts: Accounts{YouTubeChannel: "1234; x"},
	}.Clean()
	if s.Theme != DefaultTheme || s.Accent != "" || s.FontScale != 1 {
		t.Errorf("appearance and playback not repaired: %+v", s)
	}
	if s.Crossfade != 12 || s.Volume != 0 || len(s.Equalizer.Gains) != len(Bands) {
		t.Errorf("ranges not clamped: %+v", s)
	}
	if s.JSRuntime != "" {
		t.Errorf("tools not repaired: %+v", s)
	}
	if s.YouTubeChannel != "" {
		t.Errorf("a channel that is not a page ID was kept: %q", s.YouTubeChannel)
	}
	if c := (Settings{Accounts: Accounts{YouTubeChannel: " 112233445566778899 "}}).Clean(); c.YouTubeChannel != "112233445566778899" {
		t.Errorf("a page ID was not kept: %q", c.YouTubeChannel)
	}
	if !Default().YouTubeHistory {
		t.Error("plays do not reach the YouTube history by default")
	}
}

func TestValidID(t *testing.T) {
	for id, want := range map[string]bool{"mix": true, "sleep-timer": true, "a1": true, "": false, "-x": false, "Mix": false, "a/b": false, "a.b": false} {
		if got := ValidID(id); got != want {
			t.Errorf("ValidID(%q) = %v", id, got)
		}
	}
}

func TestActiveTheme(t *testing.T) {
	s := Default()
	if got := s.ActiveTheme(true); got != DefaultTheme {
		t.Errorf("not following the system, dark: %q", got)
	}
	s.Theme = "sakura"
	s.FollowSystem = true
	if got := s.ActiveTheme(true); got != DefaultTheme {
		t.Errorf("dark appearance: %q", got)
	}
	if got := s.ActiveTheme(false); got != DefaultLightTheme {
		t.Errorf("light appearance: %q", got)
	}
	s.LightTheme, s.DarkTheme = "Not An ID", ""
	if c := s.Clean(); c.LightTheme != DefaultLightTheme || c.DarkTheme != DefaultTheme {
		t.Errorf("cleaned to %q and %q", c.LightTheme, c.DarkTheme)
	}
}

func TestShortcutsClean(t *testing.T) {
	s := Default()
	s.Shortcuts = map[string]string{"player.toggle": "Mod+P", "player.mute": "", "": "K", "x": strings.Repeat("A", 65)}
	c := s.Clean()
	if c.Shortcuts["player.toggle"] != "Mod+P" {
		t.Errorf("chosen keys dropped: %v", c.Shortcuts)
	}
	if v, ok := c.Shortcuts["player.mute"]; !ok || v != "" {
		t.Errorf("a removed shortcut must stay removed: %v", c.Shortcuts)
	}
	if _, ok := c.Shortcuts[""]; ok || len(c.Shortcuts) != 2 {
		t.Errorf("invalid entries kept: %v", c.Shortcuts)
	}
	if (Settings{}).Clean().Shortcuts == nil {
		t.Error("Clean must give an empty map, not nil")
	}
}

// TestGlobalShortcutsClean keeps the global shortcuts chosen, of IDs this
// version doesn't know too, and drops entries without keys or too long.
func TestGlobalShortcutsClean(t *testing.T) {
	if d := Default().GlobalShortcuts; d == nil || len(d) != 0 {
		t.Errorf("a new install has global shortcuts %v", d)
	}
	if (Settings{}).Clean().GlobalShortcuts == nil {
		t.Error("Clean must give an empty map, not nil")
	}
	s := Default()
	s.GlobalShortcuts = map[string]string{"player.toggle": "Mod+Alt+P", "plugin.later": "Mod+Alt+L", "": "Mod+Alt+K",
		"player.mute": "", "player.next": strings.Repeat("A", 65), strings.Repeat("x", 129): "Mod+Alt+X"}
	c := s.Clean()
	if c.GlobalShortcuts["player.toggle"] != "Mod+Alt+P" || c.GlobalShortcuts["plugin.later"] != "Mod+Alt+L" || len(c.GlobalShortcuts) != 2 {
		t.Errorf("cleaned to %v", c.GlobalShortcuts)
	}
	c.GlobalShortcuts["player.toggle"] = "F9"
	if s.GlobalShortcuts["player.toggle"] != "Mod+Alt+P" {
		t.Error("changing the cleaned settings changed the settings")
	}
	// Which of too many stay is up to the map's order: only how many.
	s.GlobalShortcuts = map[string]string{}
	for i := range 40 {
		s.GlobalShortcuts["command."+strings.Repeat("a", i+1)] = "Mod+Alt+P"
	}
	if n := len(s.Clean().GlobalShortcuts); n != 32 {
		t.Errorf("%d of 40 global shortcuts kept, want 32", n)
	}
}

// TestCleanCopies changes what Clean returned: the settings it cleaned
// stay as they were.
func TestCleanCopies(t *testing.T) {
	s := Default()
	s.ThemeOptions = map[string]map[string]string{"mix": {"texture": "dots"}}
	s.Shortcuts = map[string]string{"player.toggle": "Space"}
	c := s.Clean()
	c.Equalizer.Gains[0] = 6
	c.ThemeOptions["mix"]["texture"] = "none"
	c.Shortcuts["player.toggle"] = "K"
	if s.Equalizer.Gains[0] != 0 || s.ThemeOptions["mix"]["texture"] != "dots" || s.Shortcuts["player.toggle"] != "Space" {
		t.Errorf("changing the cleaned settings changed the settings: %v %v %v", s.Equalizer.Gains, s.ThemeOptions, s.Shortcuts)
	}
}

// TestShowInTitle names what plays in the window's title in new installs
// and in the settings of earlier versions, which don't have the key.
func TestShowInTitle(t *testing.T) {
	for file, want := range map[string]bool{`{}`: true, `{"theme":"mix"}`: true, `{"showInTitle":false}`: false, `{"showInTitle":true}`: true} {
		s := Default()
		if err := json.Unmarshal([]byte(file), &s); err != nil {
			t.Fatal(err)
		}
		if s.Clean().ShowInTitle != want {
			t.Errorf("%s shows what plays in the title %v, not %v", file, s.ShowInTitle, want)
		}
	}
}

// TestCleanPanel shows the queue for settings from before the panel, and
// keeps the user's choice.
func TestCleanPanel(t *testing.T) {
	for in, want := range map[string]string{"": "queue", "bogus": "queue", "queue": "queue", "lyrics": "lyrics", "none": "none"} {
		if got := (Settings{Appearance: Appearance{Panel: in}}).Clean().Panel; got != want {
			t.Errorf("Clean of the panel %q gives %q, not %q", in, got, want)
		}
	}
}

// TestLoudness has Even loudness on in new installs, and in the settings of
// earlier versions, which don't have it, until the user turns it off.
func TestLoudness(t *testing.T) {
	for file, want := range map[string]bool{
		`{}`:                                  true,
		`{"normalize":true}`:                  true,
		`{"loudness":false}`:                  false,
		`{"loudness":false,"normalize":true}`: false,
	} {
		s := Default()
		if err := json.Unmarshal([]byte(file), &s); err != nil {
			t.Fatal(err)
		}
		if s = s.Clean(); s.Loudness != want {
			t.Errorf("%s: Even loudness %v, not %v", file, s.Loudness, want)
		}
	}
}

// TestSurround has 5.1 surround off in new installs, and in the settings
// of earlier versions, which don't have it, until the user turns it on.
func TestSurround(t *testing.T) {
	for file, want := range map[string]bool{`{}`: false, `{"surround":false}`: false, `{"surround":true}`: true} {
		s := Default()
		if err := json.Unmarshal([]byte(file), &s); err != nil {
			t.Fatal(err)
		}
		if s = s.Clean(); s.Surround != want {
			t.Errorf("%s: 5.1 surround %v, not %v", file, s.Surround, want)
		}
	}
}

// TestZoom keeps a zoom of the steps, takes the nearest step for another
// in range, and 100% for one outside them, as settings from before it.
func TestZoom(t *testing.T) {
	for _, c := range []struct{ in, want float64 }{
		{1.25, 1.25}, {2, 2}, {0.8, 0.8}, {1.3, 1.25}, {1.2, 1.25}, {0.85, 0.8},
		{0, 1}, {0.5, 1}, {3, 1}, {math.NaN(), 1},
	} {
		s := Default()
		s.Zoom = c.in
		if got := s.Clean().Zoom; got != c.want {
			t.Errorf("a zoom of %v reads as %v, want %v", c.in, got, c.want)
		}
	}
}

// TestSkipDisliked passes over disliked songs in new installs, and in the
// settings of earlier versions, which don't have it, until the user turns
// it off.
func TestSkipDisliked(t *testing.T) {
	for file, want := range map[string]bool{`{}`: true, `{"autoplay":false}`: true, `{"skipDisliked":false}`: false, `{"skipDisliked":true}`: true} {
		s := Default()
		if err := json.Unmarshal([]byte(file), &s); err != nil {
			t.Fatal(err)
		}
		if s.Clean().SkipDisliked != want {
			t.Errorf("%s skips disliked songs %v, not %v", file, s.SkipDisliked, want)
		}
	}
}

// TestCleanLayout gives settings from before the layout setting, and
// those of a layout no more, Studio's.
func TestCleanLayout(t *testing.T) {
	for in, want := range map[Layout]Layout{"": LayoutStudio, "bogus": LayoutStudio, LayoutStudio: LayoutStudio, LayoutBlend: LayoutBlend,
		LayoutExpressive: LayoutExpressive, LayoutCrate: LayoutCrate, LayoutTiles: LayoutTiles} {
		if got := (Settings{Appearance: Appearance{Layout: in}}).Clean().Layout; got != want {
			t.Errorf("Clean of layout %q = %q, want %q", in, got, want)
		}
	}
}

// TestRenamedThemes gives settings from before MyTube its themes.
func TestRenamedThemes(t *testing.T) {
	s := Settings{Appearance: Appearance{Theme: "mymusic", LightTheme: "mymusic-light", DarkTheme: "mymusic"}}.Clean()
	if s.Theme != DefaultTheme || s.LightTheme != DefaultLightTheme || s.DarkTheme != DefaultTheme {
		t.Errorf("cleaned to %q, %q and %q", s.Theme, s.LightTheme, s.DarkTheme)
	}
}
