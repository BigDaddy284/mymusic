// Package brand holds MyTube's brand files: the app's icon, which
// mygo.config.ts names, and the tray's, which the app embeds, both sound
// bars, and the Music mark, a music note, which the sidebar shows.
// make_icon.py draws them all.
// The app builds them in, so they stay out of resources/, whose every file
// `mygo build` ships beside the app (#303).
//
// discord/ holds the badges of the Discord status, which tools/discordbadges
// draws: Discord reads them from GitHub, and the app doesn't build them in.
package brand

import _ "embed"

// Mark is the Music mark, the note in its gradient, which the sidebar
// shows beside the section's name.
//
//go:embed mark.svg
var Mark []byte

// VideoMark is the Video mark, a record with a play triangle in the same
// gradient: the back of the Music mark, which the sidebar turns over to
// switch between Music and Video.
//
//go:embed videomark.svg
var VideoMark []byte
