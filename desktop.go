package main

import (
	"context"

	"github.com/egoist/mygo/ui"
)

func init() { addSettings("desktop", 10, (*shell).desktopSettings) }

// desktopRows are the rows of Settings › On the desktop that register
// from their features' files, below Discord's.
const desktopRows = "desktop › On the desktop"

func (s *shell) desktopSettings(c *ui.Context) {
	s.group(c, "On the desktop", func() {
		toggle := func(key, title, description string, v *bool, words ...string) {
			s.row(c, title, description, func() {
				if s.switchOf(c.Key(key), v).Label(title).Changed() {
					s.save()
				}
			}, words...)
		}
		toggle("tray", "Close to the tray", "Closing the window keeps the music playing; quit from the tray icon.", &s.set.CloseToTray, "minimize background")
		toggle("media-keys", "Media keys", "The keyboard's play, next and previous keys control MyTube. Keys of your own that work everywhere are in Settings › Shortcuts.", &s.set.MediaKeys)
		toggle("notifications", "Notifications", "Tell what starts playing while MyTube is in the background.", &s.set.Notifications)
		toggle("title", "Show what plays in the window title", "The window's title names the track that plays, for the taskbar, the window switcher and screen readers.", &s.set.ShowInTitle)
		toggle("discord", "Show what plays on Discord", "While Discord runs, your profile says what you listen to: you choose its lines, cover, time and button.", &s.set.DiscordPresence)
		if s.set.DiscordPresence {
			describe := func() {
				ui.Text(c, "Discord shows the name of MyTube's application. To show another, create one in the ")
				s.linkOf(c, "developer portal", "https://discord.com/developers/applications")
				ui.Text(c, " and paste its application ID.")
			}
			s.rowWith(c, "Discord application ID", describe, func() {
				in, _ := s.settingInput(c, "discord-app", s.set.DiscordAppID, func(id string) {
					s.set.DiscordAppID = id
					s.save()
				})
				in.Placeholder("MyTube's").Width(260).Label("Discord application ID")
			}, "Discord shows the name of MyTube's application. To show another, create one in the developer portal and paste its application ID.")
		}
		s.settingsOf(c, desktopRows)
		s.row(c, "Start at login", "", func() {
			// The system says, beside the main thread.
			login := fetch(s, "openAtLogin", func(context.Context) (bool, error) { return s.svc.Preferences.OpenAtLogin(), nil })
			on := login.value
			if s.switchOf(c.Key("login"), &on).Label("Start at login").Disabled(!login.done).Changed() {
				login.value = on
				go func() {
					err := s.svc.Preferences.SetOpenAtLogin(on)
					s.update(func() {
						if err != nil {
							s.notes = append(s.notes, noticeOf(err))
						}
						s.forget("openAtLogin")
					})
				}()
			}
		}, "startup launch")
		s.row(c, "Settings file", "Your settings, as settings.json in your data folder. Save it while MyTube runs, and they apply at once.", func() {
			if s.button(c, "settings-file-open", "open", false).Clicked() {
				s.openSettingsFile()
			}
			if s.button(c, "settings-file-show", "show in folder", false).Clicked() {
				s.showSettingsFile()
			}
		}, "settings.json json config")
		s.row(c, "Reset settings", "Puts back the settings of a new install; your accounts, library and add-ons stay.", func() {
			if s.button(c, "reset", "reset", false).Clicked() {
				s.resetting = true
			}
		}, "defaults")
	})
	if s.confirm(c, "reset-settings", &s.resetting, "reset settings", "Put back every setting as it was on a new install?", "reset") {
		if _, err := s.svc.Preferences.Reset(); err != nil {
			s.notes = append(s.notes, noticeOf(err))
		}
	}
}

func init() { addSettings("shortcuts", 10, (*shell).shortcutSettings) }

// shortcutSettings list the commands and their keys, which the user
// records by pressing them: Escape cancels, Backspace removes the
// shortcut, and a command that had the keys gives them up.
func (s *shell) shortcutSettings(c *ui.Context) {
	k := &s.lk
	groups := []string{}
	byGroup := map[string][]command{}
	palette := "the command palette"
	for _, cmd := range s.commands() {
		if _, ok := byGroup[cmd.group]; !ok {
			groups = append(groups, cmd.group)
		}
		byGroup[cmd.group] = append(byGroup[cmd.group], cmd)
		if keys := shortcutLabel(s.shortcutOf(cmd)); cmd.id == "view.palette" && keys != "" {
			palette += ", " + keys + ","
		}
	}
	// The search of the settings shows the commands it finds without the
	// section's introduction.
	if !s.find.on() {
		ui.Row(c).Gap(24).AlignItems(ui.Center).Margin(0, 0, 32).Children(func() {
			ui.Text(c, "Shortcuts of a single key work while you are not typing; "+palette+" lists every command.").
				FontSize(k.smallSize()).TextColor(k.muted).Grow(1).MinWidth(0).MaxWidth(640)
			if (len(s.set.Shortcuts) > 0 || len(s.set.GlobalShortcuts) > 0) && s.link(c, "shortcuts-reset", "reset all") {
				s.set.Shortcuts, s.set.GlobalShortcuts, s.recording = nil, nil, ""
				s.save()
			}
		})
	}
	s.globalShortcutSettings(c)
	for _, g := range groups {
		s.group(c, g, func() {
			for _, cmd := range byGroup[g] {
				current := s.shortcutOf(cmd)
				recording := s.recording == cmd.id
				// The search finds a command by its keys too.
				var keys []string
				if s.find.sifting() {
					keys = []string{shortcutWords(current)}
				}
				s.row(c, cmd.title, "", func() {
					label, ink, action := shortcutLabel(current), k.muted, "change"
					switch {
					case recording:
						label, ink, action = "press the keys · esc cancels · backspace removes", k.accent, "cancel"
					case label == "":
						label, action = "none", "set"
					}
					ui.Text(c, label).FontSize(k.smallSize()).TextColor(ink).MinWidth(120).TextAlign(ui.End).AlignSelf(ui.Center)
					b := s.button(c, "shortcut-edit-"+cmd.id, action, false).Label(action + " the shortcut of " + cmd.title)
					if recording {
						// The keys pressed name the shortcut rather than run
						// commands, until the button loses the focus.
						b.HandleInput(func(ev ui.InputEvent) bool { return s.recordKey(cmd, ev) })
						if !b.Focused() && s.recordingFocused {
							s.recording = ""
						}
						s.recordingFocused = b.Focused()
					}
					if b.Clicked() {
						s.recording, s.recordingFocused = "", false
						if !recording {
							s.recording = cmd.id
							c.Announce("Press the keys of the shortcut of " + cmd.title + ". Escape cancels, Backspace removes it.")
						}
					}
					// Reset puts the command's own keys back, which a command
					// the user gave them to gives up, as when they're recorded.
					if _, custom := s.set.Shortcuts[cmd.id]; custom && !recording &&
						s.button(c, "shortcut-reset-"+cmd.id, "reset", false).Label("reset the shortcut of "+cmd.title).Clicked() {
						s.setShortcut(cmd, cmd.shortcut)
					}
				}, keys...)
			}
		})
	}
}

// recordKey takes a key pressed as the user records cmd's shortcut; Tab
// goes on, as the focus leaves.
func (s *shell) recordKey(cmd command, ev ui.InputEvent) bool {
	if ev.Kind != ui.InputKeyDown || ev.Key == ui.KeyTab {
		return false
	}
	switch {
	case ev.Repeat:
	case ev.Key == ui.KeyEscape:
		s.recording = ""
	case ev.Key == ui.KeyBackspace || ev.Key == ui.KeyDelete:
		s.setShortcut(cmd, "")
	default:
		// A modifier alone waits for its key.
		if keys := shortcutName(ev.Mods, ev.Key); keys != "" {
			s.setShortcut(cmd, keys)
		}
	}
	return true
}

// setShortcut gives a command the keys of a shortcut, "" for none; a
// command that had them gives them up.
func (s *shell) setShortcut(cmd command, keys string) {
	s.recording = ""
	if s.set.Shortcuts == nil {
		s.set.Shortcuts = map[string]string{}
	}
	if keys != "" {
		s.giveUpGlobal(cmd.id, keys)
		for _, other := range s.commands() {
			if other.id != cmd.id && sameShortcut(s.shortcutOf(other), keys) {
				s.set.Shortcuts[other.id] = ""
				s.notes = append(s.notes, notice(other.title+" has no shortcut now"))
			}
		}
	}
	if sameShortcut(keys, cmd.shortcut) {
		delete(s.set.Shortcuts, cmd.id)
	} else {
		s.set.Shortcuts[cmd.id] = keys
	}
	s.save()
}

// buildNames tell how a copy of each kind was installed, as
// System.InstallKind says.
var buildNames = map[string]string{
	"installer": "installed with the setup", "portable": "portable", "development": "development build",
}
