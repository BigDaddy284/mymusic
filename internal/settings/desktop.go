package settings

import (
	"cmp"
	"regexp"
	"strings"
)

// DefaultDiscordAppID is MyMusic's own Discord application, whose name
// Discord shows beside what plays when the user set no other.
const DefaultDiscordAppID = "1562074343691788378"

// discordIDRe matches a Discord application ID, a number.
var discordIDRe = regexp.MustCompile(`^[0-9]{15,22}$`)

// Desktop is how MyMusic fits in the desktop: Settings › Desktop.
type Desktop struct {
	// CloseToTray keeps the music playing as the window closes, from the
	// tray icon. It is off at first.
	CloseToTray bool `json:"closeToTray"`
	// Notifications tell what starts playing while MyMusic is in the
	// background.
	Notifications bool `json:"notifications"`
	// ShowInTitle names what plays in the main window's title.
	ShowInTitle bool `json:"showInTitle"`
	// KeepOnTop keeps the main window above the other windows (#350). It
	// is off at first; the mini player stays on top either way.
	KeepOnTop bool `json:"keepOnTop"`
	// MiniPanel is the panel the mini player shows beside it as it opens
	// (#371): MiniQueue, MiniLyrics or MiniOutput, or "" for none.
	MiniPanel string `json:"miniPanel"`
	// MediaKeys has the keyboard's play, next and previous keys control
	// MyMusic.
	MediaKeys bool `json:"mediaKeys"`
	// DiscordPresence shows what plays on the user's Discord profile, as
	// the Discord application of DiscordAppID, or MyMusic's own when it
	// is empty.
	DiscordPresence bool `json:"showOnDiscord"`
	// DiscordAppID is the ID of the user's own Discord application, which
	// DiscordPresence shows as; "" is MyMusic's.
	DiscordAppID string `json:"discordAppId"`
	// DiscordFirstLine is the first line of the Discord status, a template
	// whose placeholders, as {title}, name the track's (#352). Empty, the
	// line shows nothing.
	DiscordFirstLine string `json:"discordFirstLine"`
	// DiscordSecondLine is the status's second line, a template as
	// DiscordFirstLine is.
	DiscordSecondLine string `json:"discordSecondLine"`
	// DiscordCover shows the track's cover in the status, where its
	// service has it online.
	DiscordCover bool `json:"discordCover"`
	// DiscordServiceBadge shows a badge of the track's service, YouTube
	// Music, Spotify or this computer, at the cover's corner. It is off
	// at first.
	DiscordServiceBadge bool `json:"discordServiceBadge"`
	// DiscordTimeLeft shows the time played and the time left in the
	// status.
	DiscordTimeLeft bool `json:"discordTimeLeft"`
	// DiscordListenButton shows a button in the status that opens the
	// track on its service. It is off at first.
	DiscordListenButton bool `json:"discordListenButton"`
	// DiscordKeepWhilePaused keeps the status while the track is paused,
	// marked as paused and without its times, where pausing clears it at
	// first.
	DiscordKeepWhilePaused bool `json:"discordKeepWhilePaused"`
}

// defaultDesktop is how a new install fits in the desktop.
func defaultDesktop() Desktop {
	var d Desktop
	d.Notifications = true
	d.ShowInTitle = true
	d.MediaKeys = true
	d.DiscordPresence = true
	d.DiscordFirstLine, d.DiscordSecondLine = DefaultDiscordFirstLine, DefaultDiscordSecondLine
	d.DiscordCover, d.DiscordTimeLeft = true, true
	return d
}

// clean is d without a Discord application ID that isn't one, with lines of
// the Discord status that Discord takes, and with a panel of the mini
// player it has.
func (d Desktop) clean() Desktop {
	d.MiniPanel = cleanMiniPanel(d.MiniPanel)
	if d.DiscordAppID = strings.TrimSpace(d.DiscordAppID); !discordIDRe.MatchString(d.DiscordAppID) {
		d.DiscordAppID = ""
	}
	d.DiscordFirstLine, d.DiscordSecondLine = cleanDiscordLine(d.DiscordFirstLine), cleanDiscordLine(d.DiscordSecondLine)
	return d
}

// DiscordApp returns the ID of the Discord application that shows what
// plays: the user's, or MyMusic's own.
func (d Desktop) DiscordApp() string { return cmp.Or(d.DiscordAppID, DefaultDiscordAppID) }
